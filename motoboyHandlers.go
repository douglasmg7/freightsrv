package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Motoboy deliveries.
func motoboyDeliveriesHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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
