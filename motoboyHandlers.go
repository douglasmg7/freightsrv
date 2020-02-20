package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Motoboy deliveries.
func motoboyFreightsHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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

// Get motoboy deliverie.
func motoboyFreightHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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

// Save/Update motoboy deliverie.
func motoboyFreightUpdateHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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
