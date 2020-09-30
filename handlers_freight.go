package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Freight by product for Zunka.
func freightsZunkaHandlerV2(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	body, err := ioutil.ReadAll(req.Body)
	if checkError(err) {
		http.Error(w, "Error reading body: %v", http.StatusInternalServerError)
		return
	}
	// log.Printf("body: %s", string(body))
	productsIn := zunkaProducts{}

	err = json.Unmarshal(body, &productsIn)
	if checkError(err) {
		http.Error(w, "Error unmarshalling body: %v", http.StatusInternalServerError)
		return
	}
	// log.Printf("[debug] products zunka: %+v", productsIn)

	// Products list for each dealer location.
	dealerProductsMap := make(map[string][]zunkaProduct)
	for _, product := range productsIn.Products {
		if product.Dealer == "Aldo" || product.Dealer == "Allnations" {
			// Invalid lenght.
			if product.Length == 0 {
				http.Error(w, fmt.Sprintf("Invalid product [%v] length [%v].", product.ID, product.Length), http.StatusBadRequest)
				return
			}
			// Invalid width.
			if product.Width == 0 {
				http.Error(w, fmt.Sprintf("Invalid product [%v] width [%v].", product.ID, product.Width), http.StatusBadRequest)
				return
			}
			// Invalid height.
			if product.Height == 0 {
				http.Error(w, fmt.Sprintf("Invalid product [%v] height [%v].", product.ID, product.Height), http.StatusBadRequest)
				return
			}
			// Invalid weight.
			if product.Weight == 0 {
				http.Error(w, fmt.Sprintf("Invalid product [%v] weight [%v].", product.ID, product.Weight), http.StatusBadRequest)
				return
			}
			// Invalid price.
			if product.Price < 1.0 || product.Price > 1000000.0 {
				http.Error(w, fmt.Sprintf("Invalid product [%v] price [%v].", product.ID, product.Price), http.StatusBadRequest)
				return
			}

			dealer := strings.ToLower(product.Dealer) + "_" + strings.ToLower(product.StockLocation)
			dealerToZunkaProducts, ok := dealerProductsMap[dealer]
			if ok {
				dealerToZunkaProducts = append(dealerToZunkaProducts, product)
			} else {
				dealerProductsMap[dealer] = []zunkaProduct{product}
			}
		}
	}

	// Create packs.
	// Zunka to client.
	zunkaToClientPack, err := createPackV2(CEP_ZUNKA, productsIn.CepDestiny, productsIn.Products)
	// log.Printf("Zunka pack: %+v\n\n", zunkaToClientPack)
	if checkError(err) {
		http.Error(w, "Could not create package. %v", http.StatusInternalServerError)
	}
	// Dealer to zunka.
	dealerPacks := []pack{}
	for _, dealerToZunkaProducts := range dealerProductsMap {
		dealer := strings.ToLower(dealerToZunkaProducts[0].Dealer) + "_" + strings.ToLower(dealerToZunkaProducts[0].StockLocation)
		p, err := createPackV2(getCEPByDealerLocation(dealer), CEP_ZUNKA, dealerToZunkaProducts)
		// log.Printf("Dealer pack: %+v\n\n", p)
		if checkError(err) {
			http.Error(w, "Could not create package. %v", http.StatusInternalServerError)
		}
		dealerPacks = append(dealerPacks, p)
	}
	// Number of pakcs come from dealers, one for each.
	dealerPacksCount := len(dealerPacks)

	chanFreightS := [](chan *freightsOk){}

	// Zunka correios.
	chanFreight := make(chan *freightsOk)
	chanFreightS = append(chanFreightS, chanFreight)
	go getCorreiosFreightByPack(chanFreight, &zunkaToClientPack)

	// Zunka motoboy.
	chanFreight = make(chan *freightsOk)
	chanFreightS = append(chanFreightS, chanFreight)
	go getMotoboyFreightByCEP(chanFreight, zunkaToClientPack.CEPDestiny)

	// Zunka region.
	chanFreight = make(chan *freightsOk)
	chanFreightS = append(chanFreightS, chanFreight)
	go getFreightRegionByCEPAndWeight(chanFreight, zunkaToClientPack.CEPDestiny, zunkaToClientPack.Weight)

	// Dealer
	for i := range dealerPacks {
		// Correios
		chanDealer := make(chan *freightsOk)
		chanFreightS = append(chanFreightS, chanDealer)
		// dealerPacks[i].Dealer = fmt.Sprintf("%v", i)
		// log.Printf("pack: %+v", &dealerPacks[i])
		go getCorreiosFreightByPack(chanDealer, &dealerPacks[i])

		// // Table
		// chanDealer = make(chan *freightsOk)
		// chanFreightS = append(chanFreightS, chanDealer)
		// go getDealerFreightByDealerAndLocationAndWeight(chanDealer, dealerPacks[i].
	}

	// Dealer region.
	// todo.

	frZunkaCorreiosS := []*freight{}
	frZunkaRegionS := []*freight{}
	frZunkaMotoboyS := []*freight{}

	type dealerFreights struct {
		dealerCount int
		freight
	}

	// Sum freight by service code.
	dealerFrCorreiosSum := make(map[string]*dealerFreights)
	dealerFrRegionSum := make(map[string]*dealerFreights)
	for _, c := range chanFreightS {
		frsOk := <-c
		if frsOk.Ok {
			// log.Printf("\nfrsOk: %+v\n", frsOk)
			for _, fr := range frsOk.Freights {
				// log.Printf("fr: %+v\n", fr)
				switch frsOk.CEPOrigin {
				// Zunka to clients.
				case CEP_ZUNKA:
					switch fr.Carrier {
					case "Correios":
						frZunkaCorreiosS = append(frZunkaCorreiosS, fr)
					case "Motoboy":
						frZunkaMotoboyS = append(frZunkaMotoboyS, fr)
					default:
						// Transportadora (tabela).
						frZunkaRegionS = append(frZunkaRegionS, fr)
					}
				// Dealers to zunka.
				default:
					var frSumMap map[string]*dealerFreights
					switch fr.Carrier {
					case "Correios":
						frSumMap = dealerFrCorreiosSum
					case "Region":
						frSumMap = dealerFrRegionSum
					}
					frSum, ok := frSumMap[fr.ServiceCode]
					if ok {
						frSum.freight.Price += fr.Price
						if fr.Deadline > frSum.freight.Deadline {
							frSum.freight.Deadline = fr.Deadline
						}
						frSum.dealerCount++
						// log.Printf("frSum: %+v", frSum)
						// log.Printf("dealerFrCorreiosSum: %+v", dealerFrCorreiosSum)
					} else {
						frSumMap[fr.ServiceCode] = &dealerFreights{
							dealerCount: 1,
							freight: freight{
								Carrier:     fr.Carrier,
								ServiceCode: fr.ServiceCode,
								ServiceDesc: fr.ServiceDesc,
								Deadline:    fr.Deadline,
								Price:       fr.Price,
							},
						}
						// log.Printf("dealerFrCorreiosSum: %+v", dealerFrCorreiosSum)
					}
				}
			}
		}
	}
	// log.Printf("dealerFrCorreiosSum: %+v", dealerFrCorreiosSum)

	// Freights with all necessaries leg added.
	var frsOut []*freight
	// Some pack come from dealer.
	if dealerPacksCount > 0 {
		// log.Printf("\nPacks from Dealer to zunka")
		// Get only valid dealer correios. Maybe some service code not received for all dealer package.
		temp := map[string]*dealerFreights{}
		for key, val := range dealerFrCorreiosSum {
			if val.dealerCount == dealerPacksCount {
				temp[key] = val
			}
		}
		dealerFrCorreiosSum = temp
		// log.Printf("dealerFrCorreiosSum: %+v", dealerFrCorreiosSum)

		// Correios.
		// Add the two legs.
		for _, frZunka := range frZunkaCorreiosS {
			frDealer, ok := dealerFrCorreiosSum[frZunka.ServiceCode]
			// If have the same service code for the two legs.
			if ok {
				frsOut = append(frsOut, &freight{
					Carrier:     frZunka.Carrier,
					ServiceCode: frZunka.ServiceCode,
					ServiceDesc: frZunka.ServiceDesc,
					Price:       frZunka.Price + frDealer.Price,
					Deadline:    frZunka.Deadline + frDealer.Deadline,
				})
			}
		}
		// Region.
		// Only if no correios result.
		if len(frsOut) == 0 {
			// Sum min with min, max with max.
			frZunkaMin := freight{
				Deadline: 1000,
			}
			frZunkaMax := freight{
				Deadline: 0,
			}
			frDealerMin := freight{
				Deadline: 1000,
			}
			frDealerMax := freight{
				Deadline: 0,
			}
			// First leg correios.
			for _, frZunka := range frZunkaCorreiosS {
				if frZunka.Deadline > frZunkaMax.Deadline {
					frZunkaMax = *frZunka
				}
				if frZunka.Deadline < frZunkaMin.Deadline {
					frZunkaMin = *frZunka
				}
			}
			// Not have first leg correios.
			if frZunkaMin.Carrier == "" {
				for _, frZunka := range frZunkaRegionS {
					if frZunka.Deadline > frZunkaMax.Deadline {
						frZunkaMax = *frZunka
					}
					if frZunka.Deadline < frZunkaMin.Deadline {
						frZunkaMin = *frZunka
					}
				}
			}
			// Second leg correios.
			for _, frDealer := range dealerFrCorreiosSum {
				if frDealer.Deadline > frDealerMax.Deadline {
					frDealerMax = frDealer.freight
				}
				if frDealer.Deadline < frDealerMin.Deadline {
					frDealerMin = frDealer.freight
				}
			}
			// Not have second leg correios.
			if frDealerMin.Carrier == "" {
				// todo.
				// for _, frDealer := range frDealerRegionS {
				// if frDealer.Deadline > frDealerMax.Deadline {
				// frDealerMax = *frDealer
				// }
				// if frDealer.Deadline < frDealerMin.Deadline {
				// frDealerMin = *frDealer
				// }
				// }
			}
			// Get the two legs.
			if frZunkaMin.Carrier != "" && frDealerMin.Carrier != "" {
				carrier := "Transportadora 1"
				if frZunkaMin.Carrier == frDealerMin.Carrier {
					carrier = frZunkaMin.Carrier
				}
				frsOut = append(frsOut, &freight{
					Carrier:     carrier,
					ServiceCode: frZunkaMin.ServiceCode,
					ServiceDesc: frZunkaMin.ServiceDesc,
					Price:       frZunkaMin.Price + frDealerMin.Price,
					Deadline:    frZunkaMin.Deadline + frDealerMin.Deadline,
				})
				carrier = "Transportadora 2"
				if frZunkaMax.Carrier == frDealerMax.Carrier {
					carrier = frZunkaMax.Carrier
				}
				frsOut = append(frsOut, &freight{
					Carrier:     carrier,
					ServiceCode: frZunkaMax.ServiceCode,
					ServiceDesc: frZunkaMax.ServiceDesc,
					Price:       frZunkaMax.Price + frDealerMax.Price,
					Deadline:    frZunkaMax.Deadline + frDealerMax.Deadline,
				})
			}
		}
	} else {
		// All product on zunka stock, nothing coming from dealers.
		frsOut = frZunkaCorreiosS
		// log.Printf("frsOut: %+v\n", frsOut)
		// log.Printf("frsOut[0]: %+v\n\n", frsOut[0])
		// log.Printf("frsOk.Freights[0]: %+v\n\n", frsOk.Freights[0])

		// Add region freights if no correios freights.
		if len(frsOut) == 0 {
			for _, fr := range frZunkaRegionS {
				// log.Printf("frZunkaRegion: %+v\n", fr)
				frsOut = append(frsOut, fr)
			}
		}
		// Add motoboy freights.
		for _, fr := range frZunkaMotoboyS {
			// log.Printf("frZunkaMotoboy: %+v\n", fr)
			frsOut = append(frsOut, fr)
		}
	}

	frsJson, err := json.Marshal(frsOut)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// log.Printf("frsJson: %+v\n\n", string(frsJson))
	w.Write(frsJson)
}

// Freight for Zunka.
func freightsZunkaHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	body, err := ioutil.ReadAll(req.Body)
	if checkError(err) {
		http.Error(w, "Error reading body: %v", http.StatusInternalServerError)
		return
	}
	// log.Printf("body: %s", string(body))
	p := pack{}
	err = json.Unmarshal(body, &p)
	if checkError(err) {
		http.Error(w, "Error unmarshalling body: %v", http.StatusInternalServerError)
		return
	}
	// log.Printf("[debug] Pack zunka handler: %+v", p)

	var deadlinePlus int
	var includeMotoboy bool
	switch strings.ToLower(p.Dealer) {
	case "aldo":
		includeMotoboy = true
		deadlinePlus = 3
	default:
		includeMotoboy = true
		deadlinePlus = 0
	}

	// Correios
	cCorreios := make(chan *freightsOk)
	go getCorreiosFreightByPack(cCorreios, &p)

	// Motoboy.
	cMotoboy := make(chan *freightsOk)
	go getMotoboyFreightByCEP(cMotoboy, p.CEPDestiny)

	// Region.
	cRegion := make(chan *freightsOk)
	go getFreightRegionByCEPAndWeight(cRegion, p.CEPDestiny, p.Weight)

	frsOkMotoboy, frsOkCorreios, frsOkRegion := <-cMotoboy, <-cCorreios, <-cRegion

	frInfoS := []freightInfo{}
	// Correios result.
	if frsOkCorreios.Ok {
		for _, pfr := range frsOkCorreios.Freights {
			frInfoS = append(frInfoS, freightInfo{
				Carrier:     pfr.Carrier,
				ServiceCode: pfr.ServiceCode,
				ServiceDesc: pfr.ServiceDesc,
				Deadline:    pfr.Deadline + deadlinePlus,
				Price:       pfr.Price,
			})
			// log.Printf("Correio freight: %+v", *pfr)
		}
	}

	// Region result, if no Correios result.
	if len(frInfoS) == 0 && frsOkRegion.Ok {
		for _, pfr := range frsOkRegion.Freights {
			frInfoS = append(frInfoS, freightInfo{
				Carrier:  pfr.Carrier,
				Deadline: pfr.Deadline + deadlinePlus,
				Price:    pfr.Price,
			})
			// log.Printf("Region freight: %+v", *pfr)
		}
	}

	// Motoboy result.
	if frsOkMotoboy.Ok && includeMotoboy {
		// Motoboy return only one freight.
		frInfoS = append(frInfoS, freightInfo{
			Carrier:  frsOkMotoboy.Freights[0].Carrier,
			Deadline: frsOkMotoboy.Freights[0].Deadline + deadlinePlus,
			Price:    frsOkMotoboy.Freights[0].Price,
		})
		// log.Printf("Motoboy freight: %+v", *frsOkMotoboy.Freights[0])
	}

	frInfoSJson, err := json.Marshal(frInfoS)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(frInfoSJson)
}

// Freight for Zoom.
func freightsZoomHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// Get products ids and CEP.
	body, err := ioutil.ReadAll(req.Body)
	if checkError(err) {
		http.Error(w, "can't read body", http.StatusInternalServerError)
		return
	}
	// log.Printf("body: %s", string(body))
	// log.Printf("[debug] zoom freight request: %s", body)
	fRequest := zoomFregihtRequest{}
	err = json.Unmarshal(body, &fRequest)
	if checkError(err) {
		http.Error(w, "Error unmarshalling body", http.StatusInternalServerError)
		return
	}
	// Get products information.
	prodIds := struct {
		Ids []string `json:"productsId"`
	}{}
	for _, item := range fRequest.Items {
		prodIds.Ids = append(prodIds.Ids, item.ProductId)
	}
	// Products ids request.
	reqBody, err := json.Marshal(prodIds)
	if checkError(err) {
		http.Error(w, "Error marshalling products ids.", http.StatusInternalServerError)
		return
	}
	// log.Printf("reqBody: %s", reqBody)
	start := time.Now()
	client := &http.Client{}
	req, err = http.NewRequest("GET", zunkaSiteHost()+"/setup/product-info", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	if checkError(err) {
		http.Error(w, "Error creating request to zunkasite", http.StatusInternalServerError)
		return
	}
	req.SetBasicAuth(zunkaSiteUser(), zunkaSitePass())
	res, err := client.Do(req)
	if checkError(err) {
		http.Error(w, "Error requesting products information to zunkasite.", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()
	// Result.
	resBody, err := ioutil.ReadAll(res.Body)
	if checkError(err) {
		http.Error(w, "Error reading body from zunkasite response.", http.StatusInternalServerError)
		return
	}
	log.Printf("[debug] Requesting product information from zunkasite, response time: %.3fs", time.Since(start).Seconds())
	// log.Printf("resBody: %s", resBody)
	// Bad request.
	if res.StatusCode == 400 {
		log.Printf("[warn] Requesting product information from zunkasite, status: %v, body: %v", res.StatusCode, string(resBody))
		http.Error(w, string(resBody), http.StatusBadRequest)
		return
	}
	// No 200 status.
	if res.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("Error requesting product information from zunkasite.\n\nstatus: %v\n\nbody: %v", res.StatusCode, string(resBody)))
		if checkError(err) {
			http.Error(w, "Error requesting product information from zunkasite.", http.StatusInternalServerError)
			return
		}
	}

	// Products informartions returned by zoom site.
	zProducts := []zunkaProduct{}
	// log.Printf("[debug] resBody: %v", resBody)
	err = json.Unmarshal(resBody, &zProducts)
	if checkError(err) {
		http.Error(w, "Can't read body from zunka.", http.StatusInternalServerError)
		return
	}
	// log.Printf("zProducts: %v", zProducts)
	if len(zProducts) != len(prodIds.Ids) {
		http.Error(w, "Some of product(s) was not found.", http.StatusBadRequest)
		return
	}

	// Create pack.
	p, ok := createPack(zProducts, fRequest.Zipcode)
	if !ok {
		http.Error(w, "Invalid product dimensions.", http.StatusInternalServerError)
		return
	}
	// log.Printf("[debug] Pack zoom handler: %+v\n", p)

	// Correios
	cCorreios := make(chan *freightsOk)
	go getCorreiosFreightByPack(cCorreios, &p)

	// Region.
	cRegion := make(chan *freightsOk)
	go getFreightRegionByCEPAndWeight(cRegion, p.CEPDestiny, p.Weight)

	frsOkCorreios, frsOkRegion := <-cCorreios, <-cRegion

	zoomFrEst := []zoomFregihtEstimate{}
	// Correios result.
	if frsOkCorreios.Ok {
		for _, pfr := range frsOkCorreios.Freights {
			zoomFrEst = append(zoomFrEst, zoomFregihtEstimate{
				Deadline:    pfr.Deadline + p.ShipmentDelay,
				Price:       pfr.Price,
				CarrierName: pfr.Carrier,
				CarrierCode: pfr.ServiceDesc,
			})
			// log.Printf("Correio freight: %+v", *pfr)
		}
	}

	// Region result, if no Correios result.
	if len(zoomFrEst) == 0 && frsOkRegion.Ok {
		for _, pfr := range frsOkRegion.Freights {
			zoomFrEst = append(zoomFrEst, zoomFregihtEstimate{
				Deadline:    pfr.Deadline + p.ShipmentDelay,
				Price:       pfr.Price,
				CarrierName: pfr.Carrier,
				CarrierCode: pfr.ServiceDesc,
			})
			// log.Printf("Region freight: %+v", *pfr)
		}
	}

	zoomFrResponse := zoomFregihtResponse{
		ID:        strconv.FormatInt(time.Now().Unix(), 10),
		Estimates: zoomFrEst,
	}

	// log.Printf("zoomFrEst: %v", zoomFrEst)
	zoomFrResponseJSON, err := json.Marshal(zoomFrResponse)
	if err != nil {
		HandleError(w, err)
		return
	}
	log.Printf("[debug] zoom freight response: %v", string(zoomFrResponseJSON))
	w.Header().Set("Content-Type", "application/json")
	w.Write(zoomFrResponseJSON)

	// // log.Printf("zoomFrEst: %v", zoomFrEst)
	// zoomFrEstJSON, err := json.Marshal(zoomFrEst)
	// if err != nil {
	// HandleError(w, err)
	// return
	// }
	// log.Printf("[debug] zoom freight response: %v", string(zoomFrEstJSON))
	// w.Header().Set("Content-Type", "application/json")
	// w.Write(zoomFrEstJSON)
}

// Create pack.
func createPackV2(CEPOrigin string, CEPDestiny string, products []zunkaProduct) (p pack, err error) {
	if len(products) == 0 {
		return
	}
	p.CEPOrigin = CEPOrigin
	p.CEPDestiny = CEPDestiny
	// Products loop.
	for _, product := range products {
		// Invalid lenght.
		if product.Length == 0 {
			return p, fmt.Errorf("Invalid product [%v] length [%v].", product.ID, product.Length)
		}
		// Invalid width.
		if product.Width == 0 {
			return p, fmt.Errorf("Invalid product [%v] width [%v].", product.ID, product.Width)
		}
		// Invalid height.
		if product.Height == 0 {
			return p, fmt.Errorf("Invalid product [%v] height [%v].", product.ID, product.Height)
		}
		// Invalid weight.
		if product.Weight == 0 {
			return p, fmt.Errorf("Invalid product [%v] weight [%v].", product.ID, product.Weight)
		}
		// Invalid price.
		if product.Price < 1.0 || product.Price > 1000000.0 {
			return p, fmt.Errorf("Invalid product [%v] price [%v].", product.ID, product.Price)
		}

		// Price.
		p.Price += (product.Price * float64(product.Quantity))

		// Sort dimensions as lenght > width > height.
		dim := []int{product.Length, product.Width, product.Height}
		sort.Ints(dim)
		// Length.
		if dim[2] > p.Length {
			p.Length = dim[2]
		}
		// Width.
		if dim[1] > p.Width {
			p.Width = dim[1]
		}
		// Height.
		p.Height += dim[0] * product.Quantity
		p.Weight += product.Weight * product.Quantity
	}
	return
}

// Create pack.
func createPack(products []zunkaProduct, CEPDestiny string) (p pack, ok bool) {
	if len(products) == 0 {
		return
	}
	p.CEPDestiny = CEPDestiny
	// Products loop.
	for _, product := range products {
		// Invalid measurments.
		if product.Length == 0 || product.Width == 0 || product.Height == 0 || product.Weight == 0 || product.Price == 0 {
			checkError(fmt.Errorf("Invalid product dimensions: %v", product))
			return
		}
		// Invalid price.
		if product.Price < 1.0 || product.Price > 1000000.0 {
			checkError(fmt.Errorf("Invalid product price: %v", product))
			return
		}
		// Price.
		p.Price += product.Price
		// Check delay.
		switch strings.ToLower(product.Dealer) {
		case "aldo":
			if p.ShipmentDelay < 3 {
				p.ShipmentDelay = 3
			}
		}
		// Sort dimensions as lenght > width > height.
		dim := []int{product.Length, product.Width, product.Height}
		sort.Ints(dim)
		// Length.
		if dim[2] > p.Length {
			p.Length = dim[2]
		}
		// Width.
		if dim[1] > p.Width {
			p.Width = dim[1]
		}
		// Height.
		p.Height += dim[0]
		p.Weight += product.Weight
	}
	return p, true
}
