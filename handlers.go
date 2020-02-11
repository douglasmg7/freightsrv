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
		return
	}
}

// Index handler.
func indexHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("Hello!\n"))
}

// Freight deadline and price.
func freightsHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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

	// Correios.
	cFreights, err := correiosFreight(p)
	log.Printf("correios freights: %+v", cFreights)

	// Motoboy.
	mFreight, ok := getMotoboyFreightByCEP(p.OriginCEP)
	if ok {
		log.Printf("motoboy freight: %+v", mFreight)
	}

	// Default.

	w.WriteHeader(200)
	w.Write([]byte("Some value\n"))
}
