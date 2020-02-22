package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// All region freights.
func getAllRegionFreightHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// Get data.
	freights, ok := getAllFreightRegion()
	if !ok {
		http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		return
	}
	// Convert to json.
	freightJSON, err := json.Marshal(freights)
	if checkError(err) {
		http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		return
	}
	// Send response.
	w.Header().Set("Content-Type", "application/json")
	w.Write(freightJSON)
}

// One region freight.
func getOneRegionFreightHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// log.Printf("*** GET *** %v\n", ps.ByName("id"))
	// Get id.
	id, err := strconv.Atoi(ps.ByName("id"))
	if checkError(err) {
		http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		return
	}

	// Get data.
	fr, ok := getFreightRegionById(id)
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

// Create region freight.
func createRegionFreightHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// Data.
	fr := freightRegion{}
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

	// Create.
	ok := createFreightRegion(&fr)
	if !ok {
		http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
}

// Delete region freight.
func deleteRegionFreightHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// log.Printf("*** DELETE *** \n")
	id, err := strconv.Atoi(ps.ByName("id"))
	if checkError(err) {
		http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		return
	}
	// Delete.
	ok := deleteMotoboyFreight(id)
	if !ok {
		http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
}

// Update region freight.
func updateRegionFreightHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// log.Printf("*** PUT *** %v\n", ps.ByName("id"))
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
