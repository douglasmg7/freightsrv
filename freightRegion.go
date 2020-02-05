package main

import "time"

type freightRegion struct {
	Id        int       `db:"id"`
	Region    string    `db:"region"`
	Weight    int       `db:"weight"`   // g
	Deadline  int       `db:"deadline"` // days
	Price     int       `db:"price"`    // R$ X 100
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
