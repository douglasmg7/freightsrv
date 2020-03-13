package main

import "time"

type freight struct {
	Carrier     string  `json:"carrier"`
	ServiceCode string  `json:"serviceCode"`
	ServiceDesc string  `json:"serviceDesc"`
	Price       float64 `json:"price"`
	Deadline    int     `json:"deadline"` // Days.
}

type freightInfo struct {
	Carrier     string  `json:"carrier"`
	ServiceCode string  `json:"serviceCode"`
	ServiceDesc string  `json:"serviceDesc"`
	Price       float64 `json:"price"`
	Deadline    int     `json:"deadline"` // Days.
}

type freightInfoBasic struct {
	Price    float64 `json:"price"`
	Deadline int     `json:"deadline"` // Days.
}

type freightsOk struct {
	Freights []*freight
	Ok       bool
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

type pack struct {
	Dealer     string `json:"dealer"` // Aldo, Allnations, etc...
	CEPOrigin  string `json:"cepOrigin"`
	CEPDestiny string `json:"cepDestiny"`
	Weight     int    `json:"weight"` // g.
	Length     int    `json:"length"` // cm.
	Height     int    `json:"height"` // cm.
	Width      int    `json:"width"`  // cm.
}

type productIdCEP struct {
	ProductId  string `json:"productId"`
	CEPDestiny string `json:"cepDestiny"`
}

// Zunka product.
type zunkaProduct struct {
	Dealer string `json:"dealer"` // Dealer.
	Length int    `json:"length"` // cm.
	Width  int    `json:"width"`  // cm.
	Height int    `json:"height"` // cm.
	Weight int    `json:"weight"` // grams.
}

// Zoom freight request.
type zoomFregihtRequest struct {
	Zipcode string                   `json:"zipcode"` // Dealer.
	Items   []zoomFregihtRequestItem `json:"items"`   // cm.
}

// Zoom freight request.
type zoomFregihtRequestItem struct {
	Quantity  int     `json:"amount"` // Quantity.
	ProductId string  `json:"sku"`    // Product id.
	Price     float64 `json:"price"`
	Weight    float64 `json:"weight"` // Kg.
	Height    float64 `json:"height"` // Meter.
	Width     float64 `json:"width"`  // Meter.
	Length    float64 `json:"length"` // Meter.
}

// Zoom freight response.
type zoomFregihtResponse struct {
	ProductId string                `json:"id"`        // Dealer.
	Estimates []zoomFregihtEstimate `json:"estimates"` // cm.
}

// Zoom freight request.
type zoomFregihtEstimate struct {
	Price       int    `json:"shippingPrice"`
	Deadline    string `json:"daysToDelivery"`
	CarrierCode string `json:"methodId"`
	CarrierName string `json:"methodName"`
}
