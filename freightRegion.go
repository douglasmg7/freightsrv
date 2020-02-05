package main

import "time"

type freightRegion struct {
	id        int       `db:"id"`
	region    string    `db:"region"`
	weight    int       `db:"weight"`   // g
	deadline  int       `db:"deadline"` // days
	price     int       `db:"price"`    // R$ X 100
	createdAt time.Time `db:"created_at"`
	updatedAt time.Time `db:"updated_at"`
}
