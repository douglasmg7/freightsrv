package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
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

// Motoboy deliverie.
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
