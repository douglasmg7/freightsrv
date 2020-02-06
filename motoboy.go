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

// Get all motoboey feights.
func getAllMotoboyFreight() (result []motoboyFreight, err error) {
	err = sql3DB.Select(&result, "SELECT * FROM motoboy_freight ORDER BY state, city")
	if err != nil {
		return result, fmt.Errorf("getAllMotoboyFreight(). %s", err.Error())
	}
	return result, nil
}

// Get motoboy freight.
func getMotoboyFreight(mf *motoboyFreight) error {
	// Find by ID.
	if mf.ID != 0 {
		err = sql3DB.Get(&fr, "SELECT * FROM motoboy_freight  WHERE id=?", mf.ID)
		if err != nil {
			return fmt.Errorf("getMotoboyFreight(). %s", err.Error())
		}
		return nil
	}
	mf.NormalizeCity()
	// Find by normalized city.
	err = sql3DB.Get(&fr, "SELECT * FROM motoboy_freight  WHERE state=?, city_norm=?", "mg", mf.CityNorm)
	if err != nil {
		return fmt.Errorf("getMotoboyFreight(). %s", err.Error())
	}
	return nil
}

// Save motoboy freight.
func saveMotoboyFreight(mf *motoboyFreight) error {
	mf.NormalizeCity()
	tx := sql3DB.MustBegin()
	// Update.
	uStatement := "UPDATE motoboy_freight SET deadline=?, price=?, city=? WHERE state=mg AND city_norm=?"
	uResult := tx.MustExec(uStatement, mf.Deadline, mf.Price, mf.City, mf.CityNorm)
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
