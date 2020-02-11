package main

import (
	"fmt"
	"time"
)

type freightRegion struct {
	ID        int       `db:"id"`
	Region    string    `db:"region"`
	Weight    int       `db:"weight"`   // g
	Deadline  int       `db:"deadline"` // days
	Price     int       `db:"price"`    // R$ X 100
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func getAllFreightRegion() (frS []freightRegion, err error) {
	err = sql3DB.Select(&frS, "SELECT * FROM freight_region ORDER BY region, weight, deadline")
	if err != nil {
		return frS, fmt.Errorf("getAllFreightRegion(). %s", err.Error())
	}
	return frS, nil
}

// Get freight region by CEP.
func getFreightRegionByCEP(cep string) (frs []freight, ok bool) {
	region, err := getRegionByCEP(cep)
	if checkError(err) {
		return frs, false
	}

	// pmf, ok := getMotoboyFreightByLocation(address.State, address.City)
	// // log.Printf("address: %+v", address)
	// // log.Printf("pmf: %+v", *pmf)
	// if !ok {
	// return pfr, false
	// }
	// pfr = &freight{}
	// pfr.Deadline = pmf.Deadline
	// pfr.Price = float64(pmf.Price) / 100
	// pfr.Carrier = "motoboy"
	// return pfr, true
}

// Get region freight by region.
func getFreightRegionByRegion(region string) (frs []freightRegion, ok bool) {
	err = sql3DB.Select(frs, "SELECT * FROM freight_region WHERE region=?", region)
	if checkError(err) {
		return frs, false
	}
	return frs, true
}

// Get region freight by id.
func getFreightRegionById(id int) (fr freightRegion, err error) {
	err = sql3DB.Get(&fr, "SELECT * FROM freight_region WHERE id=?", id)
	if err != nil {
		return fr, fmt.Errorf(" getFreightRegionById(). %s", err.Error())
	}
	return fr, nil
}

func saveFreightRegion(fr freightRegion) error {
	tx := sql3DB.MustBegin()
	// Update.
	uStatement := "UPDATE freight_region SET price=? WHERE region=? AND weight=? AND deadline=?"
	uResult := tx.MustExec(uStatement, fr.Price, fr.Region, fr.Weight, fr.Deadline)
	uRowsAffected, err := uResult.RowsAffected()
	if err != nil {
		return err
	}

	// Insert.
	if uRowsAffected == 0 {
		iStatement := "INSERT INTO freight_region(region, weight, deadline, price) VALUES(?, ?, ?, ?)"
		iResult, err := tx.Exec(iStatement, fr.Region, fr.Weight, fr.Deadline, fr.Price)
		if err != nil {
			return err
		}
		iRowsAffected, err := iResult.RowsAffected()
		// log.Printf("iRowsAffected: %+v", iRowsAffected)
		if err != nil {
			return err
		}
		if iRowsAffected == 0 {
			return fmt.Errorf("Inserting into freight_region table not affected any row.")
		}
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Commiting insert/update into freight_region table. %s", err)
	}
	return nil
}

// Delete freight region.
func delFreightRegion(id int) error {
	stm := "DELETE FROM freight_region WHERE id=?"
	result, err := sql3DB.Exec(stm, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no one row was affected, deleting by id: %d from freight_region table", id)
	}
	return nil
}
