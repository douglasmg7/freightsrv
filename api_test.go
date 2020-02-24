package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

/******************************************************************************
* Freights
*******************************************************************************/
// Zunka freight no local stock, equipment from Aldo.
func Test_FreightZunkaNoLocalStockAldoProductAPI(t *testing.T) {
	p := pack{
		Dealer: "aldo",
		// CEPDestiny: "5-76-25-000",
		CEPDestiny: "31170210",
		Weight:     1500, // g.
		Length:     20,   // cm.
		Height:     30,   // cm.
		Width:      40,   // cm.
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

	url := "/freightsrv/freights/zunka"
	req, _ := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(reqBody))

	req.SetBasicAuth("bypass", "123456")
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	// log.Printf("res.Body: %s", res.Body.String())

	frInfoS := []freightInfo{}
	json.Unmarshal(res.Body.Bytes(), &frInfoS)
	// log.Printf("frInfoS: %+v", frInfoS)

	got := res.Body.String()
	haveMotoboy := false
	haveCorreiosOrTransporter := true

	for _, frInfo := range frInfoS {
		if frInfo.Carrier == "Correios" || frInfo.Carrier == "Transportadora" {
			haveCorreiosOrTransporter = true
		}
		if frInfo.Carrier == "Motoboy" {
			haveMotoboy = true
		}
	}
	if !haveCorreiosOrTransporter {
		t.Errorf("got:  %q, no Correios neither Transportadora carrier", got)
	}
	if haveMotoboy {
		t.Errorf("Can have Motoboy carrier, got:  %q", got)
	}
}

// Zunka freight local stock to BH.
func Test_FreightZunkaBHLocalStockAPI(t *testing.T) {
	p := pack{
		CEPDestiny: "31170210",
		Weight:     1500, // g.
		Length:     20,   // cm.
		Height:     30,   // cm.
		Width:      40,   // cm.
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

	url := "/freightsrv/freights/zunka"
	req, _ := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(reqBody))

	req.SetBasicAuth("bypass", "123456")
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	// log.Printf("res.Body: %s", res.Body.String())

	frInfoS := []freightInfo{}
	json.Unmarshal(res.Body.Bytes(), &frInfoS)
	// log.Printf("frInfoS: %+v", frInfoS)

	got := res.Body.String()
	haveMotoboy := false
	haveCorreiosOrTransporter := true

	for _, frInfo := range frInfoS {
		if frInfo.Carrier == "Correios" || frInfo.Carrier == "Transportadora" {
			haveCorreiosOrTransporter = true
		}
		if frInfo.Carrier == "Motoboy" {
			haveMotoboy = true
		}
	}
	if !haveCorreiosOrTransporter {
		t.Errorf("got:  %q, no Correios neither Transportadora carrier", got)
	}
	if !haveMotoboy {
		t.Errorf("got:  %q, no Motoboy carrier", got)
	}
}

// Freight deadline and price.
func TestFreightZoomAPI(t *testing.T) {
	t.SkipNow()
	p := pack{
		CEPDestiny: "5-76-25-000",
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

	url := "/freightsrv/freights/zoom"
	want := []string{"Correios", "Transportadora"}

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
*	Region freights
*******************************************************************************/
var regionFreightTemp = regionFreight{
	Region:   "south",
	Weight:   4000,
	Deadline: 8,
	Price:    12345,
}

// Create region freight.
func TestCreateRegionFreightAPI(t *testing.T) {
	// Url.
	url := "/freightsrv/region-freight"

	frJSON, err := json.Marshal(regionFreightTemp)
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

// All region freights.
func TestGetAllRegionFreightsAPI(t *testing.T) {
	url := "/freightsrv/region-freights"
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.SetBasicAuth("bypass", "123456")
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Errorf("Returned code: %d", res.Code)
		return
	}

	freights := []regionFreight{}
	err = json.Unmarshal(res.Body.Bytes(), &freights)
	if err != nil {
		t.Errorf("Err: %s", err)
		return
	}

	valid := false
	want := regionFreightTemp
	for _, freight := range freights {
		if freight.Region == want.Region && freight.Weight == want.Weight && freight.Deadline == want.Deadline && freight.Price == want.Price {
			valid = true
			regionFreightTemp.ID = freight.ID
		}
	}
	if !valid {
		t.Errorf("got:  %v\nwant %v, %v, %v, %v", freights, want.Region, want.Weight, want.Deadline, want.Price)
	}
}

// Get one region freight.
func TestGetOneRegionFreightAPI(t *testing.T) {
	url := fmt.Sprintf("/freightsrv/region-freight/%d", regionFreightTemp.ID)
	// log.Printf("url: %v", url)

	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.SetBasicAuth("bypass", "123456")
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Errorf("Returned code: %d", res.Code)
		return
	}

	freight := regionFreight{}
	err = json.Unmarshal(res.Body.Bytes(), &freight)
	if err != nil {
		t.Error(err)
		return
	}
	// log.Printf("Freight: %+v", freight)

	want := regionFreightTemp
	if freight.Region != want.Region || freight.Weight != want.Weight || freight.Deadline != want.Deadline || freight.Price != want.Price {
		t.Errorf("got:  %v\nwant %v, %v, %v, %v", freight, want.Region, want.Weight, want.Deadline, want.Price)
	}
}

// Update region freight.
func TestUpdateRegionFreightAPI(t *testing.T) {
	// Url.
	url := "/freightsrv/region-freight"

	regionFreightTemp.Price = 54321
	frJSON, err := json.Marshal(regionFreightTemp)
	if err != nil {
		t.Error(err)
	}

	// Request.
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(frJSON))
	req.SetBasicAuth("bypass", "123456")
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Errorf("Returned code: %d", res.Code)
		t.Errorf("res.Body:  %s\n", res.Body.String())
		return
	}
}

// Delete region freight.
func TestDeleteRegionFreightAPI(t *testing.T) {
	// Url.
	url := fmt.Sprintf("/freightsrv/region-freight/%d", regionFreightTemp.ID)

	// Request.
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	req.SetBasicAuth("bypass", "123456")
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Errorf("Returned code: %d", res.Code)
		t.Errorf("res.Body:  %s\n", res.Body.String())
		return
	}
}

/******************************************************************************
*	MOTOBOY
*******************************************************************************/
// Get all motoboy freights.
func TestGetAllMotoboyFreightsAPI(t *testing.T) {
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

// Get one motoboy freight.
func TestGetOneMotoboyFreightAPI(t *testing.T) {
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

// Create motoboy freight.
func TestCreateMotoboyFreightAPI(t *testing.T) {
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
	url := "/freightsrv/motoboy-freight"

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

// Delete motoboy freight.
func TestDeleteMotoboyFreightAPI(t *testing.T) {
	// Get city id to delete.
	city := "Nova Lima"
	freight, ok := getMotoboyFreightByLocation("mg", city)
	if !ok {
		t.Errorf("No city %s to test delete.", city)
	}
	// fmt.Printf("freight to delete: %+v", freight)

	// Url.
	url := fmt.Sprintf("/freightsrv/motoboy-freight/%d", freight.ID)

	// Request.
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
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
