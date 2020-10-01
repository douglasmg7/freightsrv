package main

import (
	"log"
	"regexp"
	"strings"
	"time"
)

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

type freightsOk struct {
	Freights   []*freight
	Ok         bool
	CEPOrigin  string
	CEPDestiny string
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

type dealerFreight struct {
	ID        int       `db:"id" json:"id"`
	Dealer    string    `db:"dealer" json:"dealer"`
	Weight    int       `db:"weight" json:"weight"`     // g
	Deadline  int       `db:"deadline" json:"deadline"` // days
	Price     int       `db:"price" json:"price"`       // R$ X 100
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

type pack struct {
	Dealer        string  `json:"dealer"` // Aldo, Allnations, etc...
	ShipmentDelay int     `json:"-"`      // Some product not in store yet.
	CEPOrigin     string  `json:"cepOrigin"`
	CEPDestiny    string  `json:"cepDestiny"`
	Length        int     `json:"length"` // cm.
	Width         int     `json:"width"`  // cm.
	Height        int     `json:"height"` // cm.
	Weight        int     `json:"weight"` // g.
	Price         float64 `json:"price"`  // R$.
}

func (p *pack) Validate() bool {

	regCep := regexp.MustCompile(`^[0-9]{8}$`)

	// Origin CEP.
	p.CEPOrigin = strings.ReplaceAll(p.CEPOrigin, "-", "")
	if p.CEPOrigin == "" {
		p.CEPOrigin = CEP_ORIGIN
	}
	if !regCep.MatchString(p.CEPOrigin) {
		log.Printf("[warning] Invalid CEP origin: %v", p.CEPOrigin)
		return false
	}

	// Destiny CEP.
	p.CEPDestiny = strings.ReplaceAll(p.CEPDestiny, "-", "")
	if !regCep.MatchString(p.CEPDestiny) {
		log.Printf("[warning] Invalid CEP destiny: %v", p.CEPDestiny)
		return false
	}

	// Weight in kg.
	minWeight := 1
	maxWeight := 50000
	if p.Weight < minWeight {
		log.Printf("[warning] Invalid weight of %v grams. Must be more than %v grams", p.Weight, minWeight)
		return false
	}
	if p.Weight > maxWeight {
		log.Printf("[warning] Invalid weight of %v grams. Must be less than %v grams", p.Weight, maxWeight)
		return false
	}

	// Price in R$.
	minPrice := 1.0
	maxPrice := 1000000.0
	if p.Price < minPrice {
		log.Printf("[warning] Invalid price of R$ %v. Must be more than R$ %v", p.Price, minPrice)
		return false
	}
	if p.Price > maxPrice {
		log.Printf("[warning] Invalid price of R$ %v. Must be less than R$ %v", p.Price, maxPrice)
		return false
	}

	return true
}

type zunkaProducts struct {
	CepDestiny string         `json:"cepDestiny"`
	Products   []zunkaProduct `json:"products"`
}

// Zunka product.
type zunkaProduct struct {
	ID            string  `json:"id"`
	Dealer        string  `json:"dealer"`        // Dealer.
	StockLocation string  `json:"stockLocation"` // SC, SP, RJ...
	Length        int     `json:"length"`        // cm.
	Width         int     `json:"width"`         // cm.
	Height        int     `json:"height"`        // cm.
	Weight        int     `json:"weight"`        // grams.
	Quantity      int     `json:"quantity"`
	Price         float64 `json:"price"` // R$.
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
	ID        string                `json:"id"`        // Dealer.
	Estimates []zoomFregihtEstimate `json:"estimates"` // cm.
}

// Zoom freight request.
type zoomFregihtEstimate struct {
	Price       float64 `json:"shippingPrice"`
	Deadline    int     `json:"daysToDelivery"`
	CarrierName string  `json:"methodName"`
	CarrierCode string  `json:"methodId"`
}
