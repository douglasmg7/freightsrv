package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Freight for Zunka.
func addressHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	cep := string(body)
	// log.Printf("body: %s", cep)

	address, err := getAddressByCEP(cep)
	if err != nil {
		log.Printf("Getting address from CEP %s. %v", cep, err)
		http.Error(w, fmt.Sprintf("Can't get address from CEP %s", cep), http.StatusBadRequest)
		return
	}

	addressJSON, err := json.Marshal(address)
	if err != nil {
		log.Printf("Getting address from CEP %s. %v", cep, err)
		http.Error(w, "Internal error", http.StatusBadRequest)
		return
	}

	log.Printf("addressJSON: %s", addressJSON)
	w.Header().Set("Content-Type", "application/json")
	w.Write(addressJSON)
}
