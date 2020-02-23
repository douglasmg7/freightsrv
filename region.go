package main

import (
	"errors"
	"fmt"
	"time"
)

type freightRegion struct {
	ID        int       `db:"id" json:"id"`
	Region    string    `db:"region" json:"region"`
	Weight    int       `db:"weight" json:"weight"`     // g
	Deadline  int       `db:"deadline" json:"deadline"` // days
	Price     int       `db:"price" json:"price"`       // R$ X 100
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

func getAllFreightRegion() (frS []freightRegion, ok bool) {
	err = sql3DB.Select(&frS, "SELECT * FROM freight_region ORDER BY region, weight, deadline")
	if checkError(err) {
		return frS, false
	}
	return frS, true
}

// Get freight region by CEP and weight.
func getFreightRegionByCEPAndWeight(c chan *freightsOk, cep string, weight int) {
	result := &freightsOk{
		Freights: []*freight{},
	}

	region, err := getRegionByCEP(cep)
	if checkError(err) {
		c <- result
		return
	}

	frrs, ok := getFreightRegionByRegionAndWeight(region, weight)
	// log.Printf("frrs: %+v", frrs)
	if !ok {
		c <- result
		return
	}

	for i, frr := range frrs {
		fr := freight{
			Carrier:  fmt.Sprintf("Transportadora %d", i+1),
			Deadline: frr.Deadline,
			Price:    float64(frr.Price) / 100,
		}
		result.Freights = append(result.Freights, &fr)
	}
	result.Ok = true
	c <- result
}

// Get region freight by region.
func getFreightRegionByRegionAndWeight(region string, weight int) (frs []freightRegion, ok bool) {
	var weightSel int
	// Get min weight freight for current weight.
	// log.Printf("SELECT MIN(weight) FROM freight_region WHERE region=%s AND weight>=%d ORDER BY deadline", region, weight)
	// err = sql3DB.Get(&weightSel, "SELECT MIN(weight) FROM freight_region WHERE region=? AND weight>=? ORDER BY deadline", region, weight)
	err = sql3DB.Get(&weightSel, "SELECT CASE WHEN MIN(weight) IS NULL THEN 0 ELSE MIN(weight) END FROM freight_region WHERE region==? AND weight>=? ORDER BY deadline;", region, weight)
	if checkError(err) {
		return frs, false
	}
	// log.Printf("weightSel: %v", weightSel)
	// NULL from sqlite, no record for selected region and weight.
	if weightSel == 0 {
		return frs, false
	}

	err = sql3DB.Select(&frs, "SELECT * FROM freight_region WHERE region=? AND weight==? ORDER BY deadline", region, weightSel)
	if checkError(err) {
		return frs, false
	}
	// log.Printf("getFreightRegionByRegionAndWeight: %+v", frs)
	return frs, true
}

// Get region freight by id.
func getFreightRegionById(id int) (fr freightRegion, ok bool) {
	err = sql3DB.Get(&fr, "SELECT * FROM freight_region WHERE id=?", id)
	if checkError(err) {
		return fr, false
	}
	return fr, true
}

// Create freight region.
func createFreightRegion(fr *freightRegion) bool {
	stm := "INSERT INTO freight_region(region, weight, deadline, price) VALUES(?, ?, ?, ?)"
	result, err := sql3DB.Exec(stm, fr.Region, fr.Weight, fr.Deadline, fr.Price)
	if checkError(err) {
		return false
	}
	rowsAffected, err := result.RowsAffected()
	if checkError(err) {
		return false
	}
	// log.Printf("iRowsAffected: %+v", iRowsAffected)
	if rowsAffected == 0 {
		checkError(errors.New("Inserting into freight_region table, no affected row."))
		return false
	}
	return true
}

// Update freight region.
func updateFreightRegion(fr *freightRegion) bool {
	// log.Printf("UPDATE freight_region SET price=%d WHERE region=%v AND weight=%d AND deadline=%d", fr.Price, fr.Region, fr.Weight, fr.Deadline)
	// stm := "UPDATE freight_region SET price=? WHERE region=? AND weight=? AND deadline=?"
	stm := "UPDATE freight_region SET region=?, weight=?, deadline=?, price=? WHERE id=?"
	result, err := sql3DB.Exec(stm, fr.Region, fr.Weight, fr.Deadline, fr.Price, fr.ID)
	if checkError(err) {
		return false
	}
	rowsAffected, err := result.RowsAffected()
	if checkError(err) {
		return false
	}
	if rowsAffected == 0 {
		checkError(errors.New("Updateing freight_region table, no affected row."))
		return false
	}
	return true
}

// Delete freight region.
func deleteFreightRegion(id int) bool {
	// log.Printf("DELETE FROM freight_region WHERE id=%d", id)
	stm := "DELETE FROM freight_region WHERE id=?"
	result, err := sql3DB.Exec(stm, id)
	if checkError(err) {
		return false
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		checkError(errors.New(fmt.Sprintf("No rows was affected by deleting freight id: %d from freight_region table", id)))
		return false
	}
	return true
}

// func saveFreightRegion(fr freightRegion) error {
// tx := sql3DB.MustBegin()
// // Update.
// uStatement := "UPDATE freight_region SET price=? WHERE region=? AND weight=? AND deadline=?"
// uResult := tx.MustExec(uStatement, fr.Price, fr.Region, fr.Weight, fr.Deadline)
// uRowsAffected, err := uResult.RowsAffected()
// if err != nil {
// return err
// }

// // Insert.
// if uRowsAffected == 0 {
// iStatement := "INSERT INTO freight_region(region, weight, deadline, price) VALUES(?, ?, ?, ?)"
// iResult, err := tx.Exec(iStatement, fr.Region, fr.Weight, fr.Deadline, fr.Price)
// if err != nil {
// return err
// }
// iRowsAffected, err := iResult.RowsAffected()
// // log.Printf("iRowsAffected: %+v", iRowsAffected)
// if err != nil {
// return err
// }
// if iRowsAffected == 0 {
// return fmt.Errorf("Inserting into freight_region table not affected any row.")
// }
// }
// err = tx.Commit()
// if err != nil {
// return fmt.Errorf("Commiting insert/update into freight_region table. %s", err)
// }
// return nil
// }
