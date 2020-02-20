package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Motoboy deliveries.
func getAllMotoboyFreightHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// Get data.
	deliveries, err := getAllMotoboyFreight()
	if checkError(err) {
		http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		return
	}
	// Convert to json.
	deliveriesJSON, err := json.Marshal(deliveries)
	if checkError(err) {
		http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		return
	}
	// Send response.
	w.Header().Set("Content-Type", "application/json")
	w.Write(deliveriesJSON)
}

// Get motoboy freight.
func getMotoboyFreightHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// Get id.
	id, err := strconv.Atoi(ps.ByName("id"))
	if checkError(err) {
		http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		return
	}

	// Get data.
	fr, ok := getMotoboyFreightByID(id)
	if !ok {
		http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		return
	}
	// Convert to json.
	frJSON, err := json.Marshal(fr)
	if checkError(err) {
		http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		return
	}
	// Send response.
	w.Header().Set("Content-Type", "application/json")
	w.Write(frJSON)
}

// Update motoboy freight.
func updateMotoboyFreightHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.Printf("*** PUT *** %v\n", ps.ByName("id"))
	// Data.
	fr := motoboyFreight{}
	body, err := ioutil.ReadAll(req.Body)
	if checkError(err) {
		http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		return
	}
	// log.Printf("body: %v\n", string(body))
	err = json.Unmarshal(body, &fr)
	if checkError(err) {
		http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		return
	}

	// Update.
	ok := updateMotoboyFreightById(&fr)
	if !ok {
		http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
}

// Create motoboy freight.
func createMotoboyFreightHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	log.Printf("*** POST *** \n")
	// Data.
	fr := motoboyFreight{}
	body, err := ioutil.ReadAll(req.Body)
	if checkError(err) {
		http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		return
	}
	// log.Printf("body: %v\n", string(body))
	err = json.Unmarshal(body, &fr)
	if checkError(err) {
		http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		return
	}

	// Update.
	ok := updateMotoboyFreightById(&fr)
	if !ok {
		http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
}
