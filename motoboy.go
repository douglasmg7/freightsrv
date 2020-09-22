package main

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

// Normalize city.
func (mf *motoboyFreight) NormalizeCity() {
	s := strings.ToLower(strings.TrimSpace(mf.City))
	reg := regexp.MustCompile(`\s+`)
	mf.CityNorm = normalizeString(reg.ReplaceAllString(s, `-`))
}

// Get motoboy freight by CEP.
func getMotoboyFreightByCEP(c chan *freightsOk, cep string) {
	result := &freightsOk{
		Freights: []*freight{},
	}
	result.CEPOrigin = CEP_ZUNKA
	result.CEPDestiny = cep
	// log.Printf("motoboy result: %+v", result)

	address, err := getAddressByCEP(cep)
	if checkError(err) {
		c <- result
		return
	}
	pmf, ok := getMotoboyFreightByLocation(address.State, address.City)
	// log.Printf("address: %+v", address)
	// log.Printf("pmf: %+v", *pmf)
	if !ok {
		c <- result
		return
	}

	fr := freight{}
	fr.Deadline = pmf.Deadline
	fr.Price = float64(pmf.Price) / 100
	fr.Carrier = "Motoboy"

	result.Freights = append(result.Freights, &fr)
	result.Ok = true

	// log.Printf("motoboy result: %+v", result)
	c <- result
}

// Get all motoboey feights.
func getAllMotoboyFreight() (result []motoboyFreight, err error) {
	err = sql3DB.Select(&result, "SELECT * FROM motoboy_freight ORDER BY state, city")
	if err != nil {
		return result, fmt.Errorf("getAllMotoboyFreight(). %s", err.Error())
	}
	return result, nil
}

// Get motoboy freight by id.
func getMotoboyFreightByID(id int) (mf *motoboyFreight, ok bool) {
	mf = &motoboyFreight{}
	err = sql3DB.Get(mf, "SELECT * FROM motoboy_freight  WHERE id=?", id)
	if checkError(err) {
		return mf, false
	}
	// log.Printf("by id, mf: %+v", *mf)
	return mf, true
}

// Get motoboy freight by location.
func getMotoboyFreightByLocation(state, city string) (mf *motoboyFreight, ok bool) {
	// Only for Minas Gerais.
	if strings.ToLower(state) != "mg" {
		return mf, false
	}
	mf = &motoboyFreight{}
	mf.State = "mg"
	mf.City = city
	mf.NormalizeCity()
	err = sql3DB.Get(mf, "SELECT * FROM motoboy_freight WHERE state=? AND city_norm=?", mf.State, mf.CityNorm)
	if err == sql.ErrNoRows {
		return mf, false
	}
	if checkError(err) {
		return mf, false
	}
	// log.Printf("by state and city, mf: %+v", *mf)
	return mf, true
}

// Update motoboy freight by id.
func updateMotoboyFreightById(freight *motoboyFreight) bool {
	freight.NormalizeCity()
	// log.Printf("freight: %+v\n", *freight)
	stm := "UPDATE motoboy_freight SET city_norm=?, city=?, deadline=?, price=?  WHERE id=?"
	// log.Printf("UPDATE motoboy_freight SET city_norm=%v, city=%v, deadline=%v, price=%v WHERE id=%v", freight.CityNorm, freight.City, freight.Deadline, freight.Price, freight.ID)
	_, err := sql3DB.Exec(stm, freight.CityNorm, freight.City, freight.Deadline, freight.Price, freight.ID)
	if checkError(err) {
		return false
	}
	return true
}

// Create motoboy freight.
func createMotoboyFreight(freight *motoboyFreight) bool {
	freight.NormalizeCity()
	// log.Printf("freight: %+v\n", *freight)
	stm := "INSERT INTO motoboy_freight(state, city, city_norm, deadline, price) VALUES(?, ?, ?, ?, ?)"
	_, err := sql3DB.Exec(stm, "mg", freight.City, freight.CityNorm, freight.Deadline, freight.Price)
	if checkError(err) {
		return false
	}
	return true
}

// Delete motoboy freight.
func deleteMotoboyFreight(id int) bool {
	// log.Printf("DELETE FROM motoboy_freight WHERE id=%d", id)
	stm := "DELETE FROM motoboy_freight WHERE id=?"
	result, err := sql3DB.Exec(stm, id)
	if checkError(err) {
		return false
	}
	rowsAffected, err := result.RowsAffected()
	if checkError(err) {
		return false
	}
	if rowsAffected == 0 {
		return false
	}
	return true
}

// Save motoboy freight.
func saveMotoboyFreight(mf *motoboyFreight) error {
	mf.NormalizeCity()
	tx := sql3DB.MustBegin()
	// Update.
	uStatement := "UPDATE motoboy_freight SET deadline=?, price=?, city=? WHERE state=? AND city_norm=?"
	uResult := tx.MustExec(uStatement, mf.Deadline, mf.Price, mf.City, "mg", mf.CityNorm)
	uRowsAffected, err := uResult.RowsAffected()
	if err != nil {
		return err
	}

	// Insert.
	if uRowsAffected == 0 {
		iStatement := "INSERT INTO motoboy_freight(state, city, city_norm, deadline, price) VALUES(?, ?, ?, ?, ?)"
		iResult, err := tx.Exec(iStatement, "mg", mf.City, mf.CityNorm, mf.Deadline, mf.Price)
		if err != nil {
			return err
		}
		iRowsAffected, err := iResult.RowsAffected()
		// log.Printf("iRowsAffected: %+v", iRowsAffected)
		if err != nil {
			return err
		}
		if iRowsAffected == 0 {
			return fmt.Errorf("Inserting into motoboy_freight table not affected any row.")
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Commiting insert/update into motoboy_freight table. %s", err)
	}
	return nil
}
