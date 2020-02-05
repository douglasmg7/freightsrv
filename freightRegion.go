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

func getFreightRegionById(id int) (fr freightRegion, err error) {
	err = sql3DB.Get(&fr, "SELECT * FROM freight_region WHERE id=?", id)
	if err != nil {
		return fr, fmt.Errorf(" getFreightRegionById(). %s", err.Error())
	}
	return fr, nil
}
