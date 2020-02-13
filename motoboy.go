package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type motoboyFreight struct {
	ID        int       `db:"id"`
	State     string    `db:"state"`
	City      string    `db:"city"`
	CityNorm  string    `db:"city_norm"` // Normalized city
	Deadline  int       `db:"deadline"`  // days
	Price     int       `db:"price"`     // R$ X 100
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

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
	if checkError(err) {
		return mf, false
	}
	// log.Printf("by state and city, mf: %+v", *mf)
	return mf, true
}

// Get motoboy freight.
func getMotoboyFreightOld(mf *motoboyFreight) error {
	// Find by ID.
	if mf.ID != 0 {
		err = sql3DB.Get(mf, "SELECT * FROM motoboy_freight  WHERE id=?", mf.ID)
		if err != nil {
			return fmt.Errorf("Getting motoboy_freight by ID. %s", err.Error())
		}
		return nil
	}

	// Find by city.
	mf.State = "mg"
	mf.NormalizeCity()
	err = sql3DB.Get(mf, "SELECT * FROM motoboy_freight WHERE state=? AND city_norm=?", mf.State, mf.CityNorm)
	if err != nil {
		return fmt.Errorf("getting motoboy freight by city. %s", err.Error())
	}
	// log.Printf("mf: %+v", mf)
	return nil
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

// Delete motoboy freight.
func delMotoboyFreight(id int) error {
	stm := "DELETE FROM motoboy_freight WHERE id=?"
	result, err := sql3DB.Exec(stm, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no one row was affected, deleting by id: %d from motoboy_freight table", id)
	}
	return nil
}
