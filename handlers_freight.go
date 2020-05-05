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
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

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
	err = json.Unmarshal(resBody, &zProducts)
	if checkError(err) {
		http.Error(w, "Can't read body from zunka.", http.StatusInternalServerError)
		return
	}
	// log.Printf("zProducts: %v", zProducts)

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

	// log.Printf("zoomFrEst: %v", zoomFrEst)
	zoomFrEstJSON, err := json.Marshal(zoomFrEst)
	if err != nil {
		HandleError(w, err)
		return
	}
	log.Printf("[debug] zoom freight response: %v", string(zoomFrEstJSON))
	w.Header().Set("Content-Type", "application/json")
	w.Write(zoomFrEstJSON)
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
