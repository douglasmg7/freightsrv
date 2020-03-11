package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Freight for Zunka.
func freightsZunkaHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	// log.Printf("body: %s", string(body))

	p := pack{}

	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Printf("Error unmarshalling body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	// log.Printf("pack: %+v", p)

	var deadlinePlus int
	var includeMotoboy bool
	switch p.Dealer {
	case "aldo":
		includeMotoboy = false
		deadlinePlus = 4
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
	// Get product id and CEP.
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	// log.Printf("body: %s", string(body))
	pIdCEP := productIdCEP{}
	err = json.Unmarshal(body, &pIdCEP)
	if err != nil {
		log.Printf("Error unmarshalling body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	// log.Printf("pack: %+v", pIdCEP)

	// Get product information from zunkasite.
	start := time.Now()
	res, err := http.Get(`https://www.zunka.com.br/product-info/` + pIdCEP.ProductId)
	if checkError(err) {
		http.Error(w, "can't get product information", http.StatusBadRequest)
		return
	}
	log.Printf("[debug] zunkasite response time: %.1fs", time.Since(start).Seconds())
	// Read response.
	resBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if checkError(err) {
		http.Error(w, "can't read body from zunka", http.StatusBadRequest)
		return
	}
	// log.Printf("address: %s", resBody)
	zProduct := zunkaProduct{}
	err = json.Unmarshal(resBody, &zProduct)
	if err != nil {
		http.Error(w, "can't read body from zunka", http.StatusBadRequest)
		return
	}

	deadlinePlus := 0
	switch zProduct.Dealer {
	case "aldo":
		deadlinePlus = 4
	}

	p := pack{
		CEPDestiny: pIdCEP.CEPDestiny,
		// DestinyCEP: "31170210",
		Weight: zProduct.Weight, // g.
		Length: zProduct.Length, // cm.
		Height: zProduct.Height, // cm.
		Width:  zProduct.Width,  // cm.
	}

	// Correios
	cCorreios := make(chan *freightsOk)
	go getCorreiosFreightByPack(cCorreios, &p)

	// Region.
	cRegion := make(chan *freightsOk)
	go getFreightRegionByCEPAndWeight(cRegion, p.CEPDestiny, p.Weight)

	frsOkCorreios, frsOkRegion := <-cCorreios, <-cRegion

	frInfoBasicS := []freightInfoBasic{}
	// Correios result.
	if frsOkCorreios.Ok {
		for _, pfr := range frsOkCorreios.Freights {
			frInfoBasicS = append(frInfoBasicS, freightInfoBasic{
				Deadline: pfr.Deadline + deadlinePlus,
				Price:    pfr.Price,
			})
			// log.Printf("Correio freight: %+v", *pfr)
		}
	}

	// Region result, if no Correios result.
	if len(frInfoBasicS) == 0 && frsOkRegion.Ok {
		for _, pfr := range frsOkRegion.Freights {
			frInfoBasicS = append(frInfoBasicS, freightInfoBasic{
				Deadline: pfr.Deadline + deadlinePlus,
				Price:    pfr.Price,
			})
			// log.Printf("Region freight: %+v", *pfr)
		}
	}

	frInfoBasicSJson, err := json.Marshal(frInfoBasicS)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(frInfoBasicSJson)
}

// // Freight deadline and price.
// func freightsHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params, includeMotoboy bool) {
// body, err := ioutil.ReadAll(req.Body)
// if err != nil {
// log.Printf("Error reading body: %v", err)
// http.Error(w, "can't read body", http.StatusBadRequest)
// return
// }
// // log.Printf("body: %s", string(body))

// p := pack{}
// err = json.Unmarshal(body, &p)
// if err != nil {
// log.Printf("Error unmarshalling body: %v", err)
// http.Error(w, "can't read body", http.StatusBadRequest)
// return
// }
// // log.Printf("pack: %+v", p)

// // Correios
// cCorreios := make(chan *freightsOk)
// go getCorreiosFreightByPack(cCorreios, &p)

// // Motoboy.
// cMotoboy := make(chan *freightsOk)
// go getMotoboyFreightByCEP(cMotoboy, p.CEPDestiny)

// // Region.
// cRegion := make(chan *freightsOk)
// go getFreightRegionByCEPAndWeight(cRegion, p.CEPDestiny, p.Weight)

// frsOkMotoboy, frsOkCorreios, frsOkRegion := <-cMotoboy, <-cCorreios, <-cRegion

// frInfoS := []freightInfo{}
// // Correios result.
// if frsOkCorreios.Ok {
// for _, pfr := range frsOkCorreios.Freights {
// frInfoS = append(frInfoS, freightInfo{
// Carrier:  pfr.Carrier,
// Deadline: pfr.Deadline,
// Price:    pfr.Price,
// })
// // log.Printf("Correio freight: %+v", *pfr)
// }
// }

// // Region result, if no Correios result.
// if len(frInfoS) == 0 && frsOkRegion.Ok {
// for _, pfr := range frsOkRegion.Freights {
// frInfoS = append(frInfoS, freightInfo{
// Carrier:  pfr.Carrier,
// Deadline: pfr.Deadline,
// Price:    pfr.Price,
// })
// // log.Printf("Region freight: %+v", *pfr)
// }
// }

// // Motoboy result.
// if frsOkMotoboy.Ok && includeMotoboy {
// // Motoboy return only one freight.
// frInfoS = append(frInfoS, freightInfo{
// Carrier:  frsOkMotoboy.Freights[0].Carrier,
// Deadline: frsOkMotoboy.Freights[0].Deadline,
// Price:    frsOkMotoboy.Freights[0].Price,
// })
// // log.Printf("Motoboy freight: %+v", *frsOkMotoboy.Freights[0])
// }

// frInfoSJson, err := json.Marshal(frInfoS)
// if err != nil {
// HandleError(w, err)
// return
// }
// w.Header().Set("Content-Type", "application/json")
// w.Write(frInfoSJson)
// }