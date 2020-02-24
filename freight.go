package main

import "time"

type pack struct {
	OriginCEP  string `json:"cepOrigin"`
	DestinyCEP string `json:"cepDestiny"`
	Weight     int    `json:"weight"` // g.
	Length     int    `json:"length"` // cm.
	Height     int    `json:"height"` // cm.
	Width      int    `json:"width"`  // cm.
}

type regionFreight struct {
	ID        int       `db:"id" json:"id"`
	Region    string    `db:"region" json:"region"`
	Weight    int       `db:"weight" json:"weight"`     // g
	Deadline  int       `db:"deadline" json:"deadline"` // days
	Price     int       `db:"price" json:"price"`       // R$ X 100
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

type motoboyFreight struct {
	ID        int       `db:"id" json:"id"`
	State     string    `db:"state" json:"-"`
	City      string    `db:"city" json:"city"`
	CityNorm  string    `db:"city_norm" json:"-"`       // Normalized city
	Deadline  int       `db:"deadline" json:"deadline"` // days
	Price     int       `db:"price" json:"price"`       // R$ X 100
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}
