package main

import (
	"errors"
	"fmt"
	"strings"
)

// Get all dealer freights.
func getAllDealerFreight() (frS []dealerFreight, ok bool) {
	err = sql3DB.Select(&frS, "SELECT * FROM dealer_freight ORDER BY dealer, weight, deadline")
	if checkError(err) {
		return frS, false
	}
	return frS, true
}

// Get dealer freight by dealer, location  and weight.
func getDealerFreightByDealerAndLocationAndWeight(c chan *freightsOk, dealer string, stockLocation string, weight int) {
	result := &freightsOk{
		Freights: []*freight{},
	}

	result.CEPOrigin = getCEPByDealerLocation(dealer, stockLocation)
	result.CEPDestiny = CEP_ZUNKA

	// Inválid CEP origin.
	if result.CEPOrigin == "" {
		c <- result
		return
	}

	// Inválid weight.
	if weight == 0 {
		c <- result
		return
	}

	delaerLocation := strings.ToLower(dealer)
	if stockLocation != "" {
		delaerLocation = delaerLocation + "_" + strings.ToLower(stockLocation)
	}

	frs, ok := getFreightRegionByRegionAndWeight(delaerLocation, weight)
	// log.Printf("frs: %+v", frs)
	if !ok {
		c <- result
		return
	}

	for i, fr := range frs {
		frOut := freight{
			Carrier:  fmt.Sprintf("Transportadora %d", i+1),
			Deadline: fr.Deadline,
			Price:    float64(fr.Price) / 100,
		}
		result.Freights = append(result.Freights, &frOut)
	}
	result.Ok = true
	c <- result
}

// Get dealer freight by dealer and weight.
func getDealerFreightByDealerAndWeight(dealer string, weight int) (frs []regionFreight, ok bool) {
	// Inváid weight.
	if weight == 0 {
		return frs, false
	}
	// Select weight.
	var weightSel int
	err = sql3DB.Get(&weightSel, "SELECT CASE WHEN MIN(weight) IS NULL THEN 0 ELSE MIN(weight) END FROM dealer_freight WHERE dealer==? AND weight>=? ORDER BY deadline;", dealer, weight)
	if checkError(err) {
		return frs, false
	}
	// log.Printf("weightSel: %v", weightSel)
	// NULL from sqlite, no record for selected dealer and weight.
	if weightSel == 0 {
		return frs, false
	}

	err = sql3DB.Select(&frs, "SELECT * FROM dealer_freight WHERE dealer=? AND weight==? ORDER BY deadline", dealer, weightSel)
	if checkError(err) {
		return frs, false
	}
	// log.Printf("getFreightRegionByRegionAndWeight: %+v", frs)
	return frs, true
}

// Get dealer freight by id.
func getDealerFreightById(id int) (fr dealerFreight, ok bool) {
	err = sql3DB.Get(&fr, "SELECT * FROM dealer_freight WHERE id=?", id)
	if checkError(err) {
		return fr, false
	}
	return fr, true
}

// Create dealer freight.
func createDealerFreight(fr *dealerFreight) bool {
	stm := "INSERT INTO dealer_freight(region, weight, deadline, price) VALUES(?, ?, ?, ?)"
	result, err := sql3DB.Exec(stm, fr.dealer, fr.Weight, fr.Deadline, fr.Price)
	if checkError(err) {
		return false
	}
	rowsAffected, err := result.RowsAffected()
	if checkError(err) {
		return false
	}
	// log.Printf("iRowsAffected: %+v", iRowsAffected)
	if rowsAffected == 0 {
		checkError(errors.New("Inserting into dealer_freight table, no affected row."))
		return false
	}
	return true
}

// Update freight region.
func updateDealerFreight(fr *dealerFreight) bool {
	stm := "UPDATE dealer_freight SET dealer=?, weight=?, deadline=?, price=? WHERE id=?"
	result, err := sql3DB.Exec(stm, fr.dealer, fr.Weight, fr.Deadline, fr.Price, fr.ID)
	if checkError(err) {
		return false
	}
	rowsAffected, err := result.RowsAffected()
	if checkError(err) {
		return false
	}
	if rowsAffected == 0 {
		checkError(errors.New("Updateing dealer_freight table, no affected row."))
		return false
	}
	return true
}

// Delete freight region.
func deleteDealerFreight(id int) bool {
	// log.Printf("DELETE FROM dealer_freight WHERE id=%d", id)
	stm := "DELETE FROM dealer_freight WHERE id=?"
	result, err := sql3DB.Exec(stm, id)
	if checkError(err) {
		return false
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		checkError(errors.New(fmt.Sprintf("No rows was affected by deleting freight id: %d from dealer_freight table", id)))
		return false
	}
	return true
}
