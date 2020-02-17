package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Handler error.
func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		// http.Error(w, "Some thing wrong", 404)
		if production {
			http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Println(err.Error())
	}
}

// Index handler.
func indexHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("Hello!\n"))
}

// Freight for Zunka.
func freightsZunkaHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	freightsHandler(w, req, ps, true)
}

// Freight for Zoom.
func freightsZoomHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	freightsHandler(w, req, ps, false)
}

// Freight deadline and price.
func freightsHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params, includeMotoboy bool) {
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

	// Correios
	cCorreios := make(chan *freightsOk)
	go getCorreiosFreightByPack(cCorreios, &p)

	// Motoboy.
	cMotoboy := make(chan *freightsOk)
	go getMotoboyFreightByCEP(cMotoboy, p.DestinyCEP)

	// Region.
	cRegion := make(chan *freightsOk)
	go getFreightRegionByCEPAndWeight(cRegion, p.DestinyCEP, p.Weight)

	frsOkMotoboy, frsOkCorreios, frsOkRegion := <-cMotoboy, <-cCorreios, <-cRegion

	frInfoS := []freightInfo{}
	// Correios result.
	if frsOkCorreios.Ok {
		for _, pfr := range frsOkCorreios.Freights {
			frInfoS = append(frInfoS, freightInfo{
				Carrier:  pfr.Carrier,
				Deadline: pfr.Deadline,
				Price:    pfr.Price,
			})
			// log.Printf("Correio freight: %+v", *pfr)
		}
	}

	// Region result, if no Correios result.
	if len(frInfoS) == 0 && frsOkRegion.Ok {
		for _, pfr := range frsOkRegion.Freights {
			frInfoS = append(frInfoS, freightInfo{
				Carrier:  pfr.Carrier,
				Deadline: pfr.Deadline,
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
			Deadline: frsOkMotoboy.Freights[0].Deadline,
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
