package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Valid no user and no password.
func Test_NoUserNoPassAuth(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/freightsrv", nil)

	// Correct password.
	res := httptest.NewRecorder()

	// indexHandler(res, req, nil)
	router.ServeHTTP(res, req)

	got := res.Body.String()
	want := "Unauthorised\n"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// Valid user and password.
func Test_ValidUserAndPassAuth(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/freightsrv", nil)

	// Correct password.
	req.SetBasicAuth("test", "1234")
	res := httptest.NewRecorder()

	// indexHandler(res, req, nil)
	router.ServeHTTP(res, req)

	got := res.Body.String()
	want := "Hello!\n"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// Invalid user.
func Test_InvalidUserAuth(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/freightsrv", nil)

	// Correct password.
	req.SetBasicAuth("test-", "1234")
	res := httptest.NewRecorder()

	// indexHandler(res, req, nil)
	router.ServeHTTP(res, req)

	got := res.Body.String()
	want := "Unauthorised\n"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// Invalid password.
func Test_InvalidPassAuth(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/freightsrv", nil)

	// Correct password.
	req.SetBasicAuth("test", "12345")
	res := httptest.NewRecorder()

	// indexHandler(res, req, nil)
	router.ServeHTTP(res, req)

	got := res.Body.String()
	want := "Unauthorised\n"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// Freight for Zunka.
func TestFreightZunkaAPI(t *testing.T) {
	TFreightAPI(t, Zunka)
}

// Freight for Zoom.
func TestFreightZoomAPI(t *testing.T) {
	TFreightAPI(t, Zoom)
}

// Freight deadline and price.
func TFreightAPI(t *testing.T, client Client) {
	p := pack{
		DestinyCEP: "5-76-25-000",
		// DestinyCEP: "31170210",
		Weight: 1500, // g.
		Length: 20,   // cm.
		Height: 30,   // cm.
		Width:  40,   // cm.
	}
	err := p.Validate()
	if err != nil {
		t.Errorf("Invalid pack. %v", err)
	}

	reqBody, err := json.Marshal(p)
	if err != nil {
		t.Error(err)
	}
	// log.Println("request body: " + string(reqBody))

	var url string
	var want []string

	switch client {
	case Zunka:
		url = "/freightsrv/freights/zunka"
		want = []string{"Correios", "Transportadora", "Motoboy"}
	case Zoom:
		url = "/freightsrv/freights/zoom"
		want = []string{"Correios", "Transportadora"}
	}
	req, _ := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(reqBody))

	req.SetBasicAuth("bypass", "123456")
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	// log.Printf("res.Body: %s", res.Body.String())

	frInfoS := []freightInfo{}
	json.Unmarshal(res.Body.Bytes(), &frInfoS)
	// log.Printf("frInfoS: %+v", frInfoS)

	// got := res.Body.String()

	for _, frInfo := range frInfoS {
		valid := false
		for _, wantCarrier := range want {
			if strings.Contains(frInfo.Carrier, wantCarrier) {
				valid = true
				break
			}
		}
		if !valid {
			t.Errorf("got:  %q, want some of %q", frInfo.Carrier, want)
		}
	}
}

/******************************************************************************
*	MOTOBOY
*******************************************************************************/
// Get motoboy freights.
func TestMotoboyFreightsAPI(t *testing.T) {
	url := "/freightsrv/motoboy-freights"
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.SetBasicAuth("bypass", "123456")
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	// log.Printf("res.Body: %s", res.Body.String())

	deliveries := []motoboyFreight{}
	json.Unmarshal(res.Body.Bytes(), &deliveries)

	wantCity := "Guarupé"
	wantDeadline := 3
	wantPrice := 9545
	valid := false
	for _, deliverie := range deliveries {
		if deliverie.City == wantCity && deliverie.Deadline == wantDeadline && deliverie.Price == wantPrice {
			valid = true
		}
	}
	if !valid {
		t.Errorf("got:  %v\nwant  %v, %v, %v", deliveries, wantCity, wantDeadline, wantPrice)
	}
}

// Get one motoboy freights.
func TestMotoboyFreightAPI(t *testing.T) {
	url := "/freightsrv/motoboy-freight/1"
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.SetBasicAuth("bypass", "123456")
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	// log.Printf("res.Body: %s", res.Body.String())

	deliverie := motoboyFreight{}
	json.Unmarshal(res.Body.Bytes(), &deliverie)

	wantCity := "Belo Horizonte"
	wantDeadline := 1
	wantPrice := 7520
	valid := false
	if deliverie.City == wantCity && deliverie.Deadline == wantDeadline && deliverie.Price == wantPrice {
		valid = true
	}
	if !valid {
		t.Errorf("got:  %v\nwant  %v, %v, %v", deliverie, wantCity, wantDeadline, wantPrice)
	}
}

// New motoboy freight.
func TestNewMotoboyFreightAPI(t *testing.T) {
	// Url.
	url := "/freightsrv/motoboy-freight"

	// Data.
	fr := motoboyFreight{
		City:     "Nova Lima",
		Deadline: 1,
		Price:    5620,
	}
	frJSON, err := json.Marshal(fr)
	if err != nil {
		t.Error(err)
	}

	// Request.
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(frJSON))
	req.SetBasicAuth("bypass", "123456")
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	// log.Printf("res.Body: %s", res.Body.String())

	want := 200
	if res.Code != want {
		t.Errorf("got:  %v, want  %v\n", res.Code, want)
		t.Errorf("res.Body:  %s\n", res.Body.String())
	}
}

// Update motoboy freight.
func TestUpdateMotoboyFreightAPI(t *testing.T) {
	// Url.
	url := "/freightsrv/motoboy-freight/4"

	// Data.
	fr := motoboyFreight{
		ID:       4,
		City:     "Sabará",
		Deadline: 2,
		Price:    9999,
	}
	frJSON, err := json.Marshal(fr)
	if err != nil {
		t.Error(err)
	}

	// Request.
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(frJSON))
	req.SetBasicAuth("bypass", "123456")
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	// log.Printf("res.Body: %s", res.Body.String())

	want := 200
	if res.Code != want {
		t.Errorf("got:  %v, want  %v\n", res.Code, want)
		t.Errorf("res.Body:  %s\n", res.Body.String())
	}
}
