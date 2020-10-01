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
	req, _ := http.NewRequest(http.MethodGet, "/freightsrv/", nil)

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
	req, _ := http.NewRequest(http.MethodGet, "/freightsrv/", nil)

	// Correct password.
	req.SetBasicAuth("bypass", "123456")
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
	req, _ := http.NewRequest(http.MethodGet, "/freightsrv/", nil)

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
	req, _ := http.NewRequest(http.MethodGet, "/freightsrv/", nil)

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
* Address.
*******************************************************************************/
// Address by CEP.
func TestGetAddressByCEPAPI(t *testing.T) {
	reqBody := strings.NewReader("31170-210")
	url := "/freightsrv/address"
	req, _ := http.NewRequest(http.MethodGet, url, reqBody)

	req.SetBasicAuth("bypass", "123456")
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Errorf("Returned code: %d", res.Code)
		return
	}

	address := viaCEPAddress{}

	err = json.Unmarshal(res.Body.Bytes(), &address)
	if err != nil {
		t.Errorf("Err: %s", err)
		return
	}

	want := viaCEPAddress{
		Cep:      "31170-210",
		Street:   "Rua Deputado Bernardino de Sena Figueiredo",
		District: "Cidade Nova",
		City:     "Belo Horizonte",
		State:    "MG",
	}

	if address.Cep != want.Cep || address.Street != want.Street || address.District != want.District || address.City != want.City || address.State != want.State {
		t.Errorf("got:  %+v\nwant %+v", address, want)
	}
}

/******************************************************************************
* Freights
*******************************************************************************/
// Zunka freight correios and motoboy
func Test_FreightZunkaAPIV2CorreiosAndMotoboyForOneProduct(t *testing.T) {
	productsIn := zunkaProducts{
		CepDestiny: "31170210",
		Products: []zunkaProduct{
			{
				ID:            "1234",
				Dealer:        "Dell",
				StockLocation: "",
				Length:        20,
				Width:         90,
				Height:        39,
				Weight:        1250,
				Quantity:      1,
				Price:         2512.22,
			},
		},
	}

	// log.Printf("productsIn: %+v\n", productsIn)
	reqBody, err := json.Marshal(productsIn)
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

	frs := []freight{}
	json.Unmarshal(res.Body.Bytes(), &frs)
	// log.Printf("frs: %+v", frs)

	got := res.Body.String()
	haveMotoboy := false
	haveCorreiosOrTransporter := true

	for _, fr := range frs {
		if fr.Carrier == "Correios" || fr.Carrier == "Transportadora" {
			haveCorreiosOrTransporter = true
		}
		if fr.Carrier == "Correios" {
			if fr.ServiceCode == "" {
				t.Errorf("Correios service code empty")
				return
			}
			if fr.ServiceDesc == "" {
				t.Errorf("Correios service description empty")
				return
			}
		}
		if fr.Carrier == "Motoboy" {
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

// Zunka freight only correios.
func Test_FreightZunkaAPIV2OnlyCorreiosOneProductNotForBigBH(t *testing.T) {
	productsIn := zunkaProducts{
		CepDestiny: "88512530",
		Products: []zunkaProduct{
			{
				ID:            "1234",
				Dealer:        "Dell",
				StockLocation: "",
				Length:        20,
				Width:         90,
				Height:        39,
				Weight:        250,
				Quantity:      2,
				Price:         512.22,
			},
		},
	}

	// log.Printf("productsIn: %+v\n", productsIn)
	reqBody, err := json.Marshal(productsIn)
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

	frs := []freight{}
	json.Unmarshal(res.Body.Bytes(), &frs)
	// log.Printf("frs: %+v", frs)

	got := res.Body.String()
	haveMotoboy := false
	haveCorreios := false
	haveTransporter := false

	for _, fr := range frs {
		if strings.HasPrefix(fr.Carrier, "Transportadora") {
			haveTransporter = true
		} else if fr.Carrier == "Correios" {
			haveCorreios = true
			if fr.ServiceCode == "" {
				t.Errorf("Correios service code empty")
				return
			}
			if fr.ServiceDesc == "" {
				t.Errorf("Correios service description empty")
				return
			}
		} else if fr.Carrier == "Motoboy" {
			haveMotoboy = true
		}
	}
	if !haveCorreios {
		t.Errorf("got:  %q, no Correios carrier", got)
	}
	if haveTransporter {
		t.Errorf("got:  %q, hava a Transportadora carrier", got)
	}
	if haveMotoboy {
		t.Errorf("got:  %q, hava a Motoboy carrier", got)
	}
}

// Zunka freight only transportadora.
func Test_FreightZunkaAPIV2OnlyTransportadoraOneLongProduct(t *testing.T) {
	productsIn := zunkaProducts{
		CepDestiny: "88512530",
		// Use a lenght not supported by Correios carrier.
		Products: []zunkaProduct{
			{
				ID:            "1234",
				Dealer:        "Dell",
				StockLocation: "",
				Length:        200,
				Width:         90,
				Height:        39,
				Weight:        250,
				Quantity:      2,
				Price:         512.22,
			},
		},
	}

	// log.Printf("productsIn: %+v\n", productsIn)
	reqBody, err := json.Marshal(productsIn)
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

	frs := []freight{}
	json.Unmarshal(res.Body.Bytes(), &frs)
	// log.Printf("frs: %+v", frs)

	got := res.Body.String()
	haveMotoboy := false
	haveCorreios := false
	haveTransporter := false

	for _, fr := range frs {
		if strings.HasPrefix(fr.Carrier, "Transportadora") {
			haveTransporter = true
		} else if fr.Carrier == "Correios" {
			haveCorreios = true
		} else if fr.Carrier == "Motoboy" {
			haveMotoboy = true
		}
	}
	if haveCorreios {
		t.Errorf("got:  %q, have Correios carrier", got)
	}
	if !haveTransporter {
		t.Errorf("got:  %q, not hava a Transportadora carrier", got)
	}
	if haveMotoboy {
		t.Errorf("got:  %q, hava a Motoboy carrier", got)
	}
}

// Zunka freight product from dealer to client only correios.
func Test_FreightZunkaAPIV2ProductFromDealerOnlyCorreios(t *testing.T) {
	productsIn := zunkaProducts{
		CepDestiny: "88512530",
		Products: []zunkaProduct{
			{
				ID:            "1234",
				Dealer:        "Aldo",
				StockLocation: "",
				Length:        20,
				Width:         90,
				Height:        39,
				Weight:        1250,
				Quantity:      2,
				Price:         2512.22,
			},
		},
	}

	// log.Printf("productsIn: %+v\n", productsIn)
	reqBody, err := json.Marshal(productsIn)
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

	frs := []freight{}
	json.Unmarshal(res.Body.Bytes(), &frs)
	// log.Printf("frs: %+v", frs)

	got := res.Body.String()
	haveMotoboy := false
	haveCorreios := false
	haveTransporter := false

	for _, fr := range frs {
		if strings.HasPrefix(fr.Carrier, "Transportadora") {
			haveTransporter = true
		} else if fr.Carrier == "Correios" {
			haveCorreios = true
			if fr.ServiceCode == "" {
				t.Errorf("Correios service code empty")
				return
			}
			if fr.ServiceDesc == "" {
				t.Errorf("Correios service description empty")
				return
			}
		}
		if fr.Carrier == "Motoboy" {
			haveMotoboy = true
		}
	}
	if !haveCorreios {
		t.Errorf("got:  %q, not have Correios carrier", got)
	}
	if haveTransporter {
		t.Errorf("got:  %q, have transportadora carrier", got)
	}
	if haveMotoboy {
		t.Errorf("got:  %q, have motoboy carrier", got)
	}
}

// Zunka freight product from dealer to BH client only correios
func TestFreightZunkaAPIV2ProductFromDealerToBHOnlyCorreios(t *testing.T) {
	productsIn := zunkaProducts{
		CepDestiny: "31170210",
		Products: []zunkaProduct{
			{
				ID:            "1234",
				Dealer:        "Allnations",
				StockLocation: "rj",
				Length:        20,
				Width:         90,
				Height:        39,
				Weight:        1250,
				Quantity:      2,
				Price:         2512.22,
			},
		},
	}

	// log.Printf("productsIn: %+v\n", productsIn)
	reqBody, err := json.Marshal(productsIn)
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

	frs := []freight{}
	json.Unmarshal(res.Body.Bytes(), &frs)
	// log.Printf("frs: %+v", frs)

	got := res.Body.String()
	haveMotoboy := false
	haveCorreios := false
	haveTransporter := false

	for _, fr := range frs {
		if strings.HasPrefix(fr.Carrier, "Transportadora") {
			haveTransporter = true
		} else if fr.Carrier == "Correios" {
			haveCorreios = true
			if fr.ServiceCode == "" {
				t.Errorf("Correios service code empty")
				return
			}
			if fr.ServiceDesc == "" {
				t.Errorf("Correios service description empty")
				return
			}
		}
		if fr.Carrier == "Motoboy" {
			haveMotoboy = true
		}
	}
	if !haveCorreios {
		t.Errorf("got:  %q, not have Correios carrier", got)
	}
	if haveTransporter {
		t.Errorf("got:  %q, have transportadora carrier", got)
	}
	if haveMotoboy {
		t.Errorf("got:  %q, have motoboy carrier", got)
	}
}

// Zunka freight product from dealer to BH client only trasportadora
// Using a lenght not suprted by correios
func TestFreightZunkaAPIV2ProductFromDealerToBHOnlyTransportadora(t *testing.T) {
	productsIn := zunkaProducts{
		CepDestiny: "31170210",
		Products: []zunkaProduct{
			{
				ID:            "1234",
				Dealer:        "Allnations",
				StockLocation: "rj",
				Length:        200,
				Width:         90,
				Height:        39,
				Weight:        1250,
				Quantity:      1,
				Price:         2512.22,
			},
		},
	}

	// log.Printf("productsIn: %+v\n", productsIn)
	reqBody, err := json.Marshal(productsIn)
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

	frs := []freight{}
	json.Unmarshal(res.Body.Bytes(), &frs)
	// log.Printf("frs: %+v", frs)

	got := res.Body.String()
	haveMotoboy := false
	haveCorreios := false
	haveTransporter := false

	for _, fr := range frs {
		if strings.HasPrefix(fr.Carrier, "Transportadora") {
			haveTransporter = true
		} else if fr.Carrier == "Correios" {
			haveCorreios = true
			if fr.ServiceCode == "" {
				t.Errorf("Correios service code empty")
				return
			}
			if fr.ServiceDesc == "" {
				t.Errorf("Correios service description empty")
				return
			}
		}
		if fr.Carrier == "Motoboy" {
			haveMotoboy = true
		}
	}
	if haveCorreios {
		t.Errorf("got:  %q, have Correios carrier", got)
	}
	if !haveTransporter {
		t.Errorf("got:  %q, not have transportadora carrier", got)
	}
	if haveMotoboy {
		t.Errorf("got:  %q, have motoboy carrier", got)
	}
}

// Zunka freight dealer to client.
func TestFreightZunkaAPIV2DealerAndZunkaSeveral(t *testing.T) {
	productsIn := zunkaProducts{
		CepDestiny: "88512530",
		Products: []zunkaProduct{
			{
				ID:            "1234",
				Dealer:        "",
				StockLocation: "",
				Length:        20,
				Width:         3,
				Height:        10,
				Weight:        120,
				Quantity:      1,
				Price:         200.0,
			},
			{
				ID:            "1234",
				Dealer:        "Aldo",
				StockLocation: "",
				Length:        20,
				Width:         3,
				Height:        10,
				Weight:        121,
				Quantity:      1,
				Price:         200.1,
			},
			{
				ID:            "1234",
				Dealer:        "Allnations",
				StockLocation: "ES",
				Length:        20,
				Width:         3,
				Height:        10,
				Weight:        122,
				Quantity:      1,
				Price:         200.2,
			},
			{
				ID:            "1234",
				Dealer:        "Allnations",
				StockLocation: "RJ",
				Length:        20,
				Width:         3,
				Height:        10,
				Weight:        123,
				Quantity:      1,
				Price:         220.3,
			},
		},
	}

	// log.Printf("productsIn: %+v\n", productsIn)
	reqBody, err := json.Marshal(productsIn)
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

	frs := []freight{}
	json.Unmarshal(res.Body.Bytes(), &frs)
	// log.Printf("frs: %+v", frs)

	got := res.Body.String()
	haveMotoboy := false
	haveCorreios := false
	haveTransporter := false

	for _, fr := range frs {
		if strings.HasPrefix(fr.Carrier, "Transportadora") {
			haveTransporter = true
		} else if fr.Carrier == "Correios" {
			haveCorreios = true
			if fr.ServiceCode == "" {
				t.Errorf("Correios service code empty")
				return
			}
			if fr.ServiceDesc == "" {
				t.Errorf("Correios service description empty")
				return
			}
		}
		if fr.Carrier == "Motoboy" {
			haveMotoboy = true
		}
	}
	if !haveCorreios {
		t.Errorf("got:  %q, Not have Correios carrier", got)
	}
	if haveTransporter {
		t.Errorf("got:  %q, Have transportadora carrier", got)
	}
	if haveMotoboy {
		t.Errorf("got:  %q, have motoboy carrier", got)
	}
}

// Zunka freight local stock to BH.
func Test_FreightZunkaAPI(t *testing.T) {
	p := pack{
		CEPDestiny: "31170210",
		Weight:     1500, // g.
		Length:     20,   // cm.
		Height:     30,   // cm.
		Width:      40,   // cm.
		Price:      2512.22,
	}
	if !p.ValidateCorreios() {
		t.Errorf("Not a valid pack to estimate correios shipping. Pack: %+v", p)
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
		if frInfo.Carrier == "Correios" {
			if frInfo.ServiceCode == "" {
				t.Errorf("Correios service code empty")
				return
			}
			if frInfo.ServiceDesc == "" {
				t.Errorf("Correios service description empty")
				return
			}
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

// Freight deadline and price by product.
func TestProductFreightZoomAPI(t *testing.T) {
	fRequest := zoomFregihtRequest{
		Zipcode: "31170210",
		Items:   []zoomFregihtRequestItem{},
	}
	// First product.
	fRequest.Items = append(fRequest.Items, zoomFregihtRequestItem{
		ProductId: "5e60eed63d13910785412eba",
		Quantity:  1,
		Price:     34.4,
		Weight:    2,
		Height:    .3,
		Width:     .2,
		Length:    .4,
	})
	// Second product.
	fRequest.Items = append(fRequest.Items, zoomFregihtRequestItem{
		ProductId: "5bcb336a4253f81781faca09",
		Quantity:  2,
		Price:     34.4,
		Weight:    2,
		Height:    .3,
		Width:     .2,
		Length:    .4,
	})

	reqBody, err := json.Marshal(fRequest)
	if err != nil {
		t.Error(err)
	}
	// log.Println("request body: " + string(reqBody))

	url := "/freightsrv/freights/zoom"
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))

	// req.SetBasicAuth("bypass", "123456")
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Errorf("Status: %d. body: %s", res.Code, res.Body.String())
		return
	}
	// log.Printf("res.Body: %s", res.Body.String())

	// fResponse := zoomFregihtResponse{
	// // No sense, because just the first product.
	// ID: fRequest.Items[0].ProductId,
	// }
	// json.Unmarshal(res.Body.Bytes(), &fResponse.Estimates)
	// // log.Printf("fResponse: %+v", fResponse)

	fResponse := zoomFregihtResponse{}
	json.Unmarshal(res.Body.Bytes(), &fResponse)
	// log.Printf("fResponse: %+v", fResponse)

	// Some estimate.
	if len(fResponse.Estimates) == 0 {
		t.Errorf("No freight info.")
		return
	}
	// Price
	if fResponse.Estimates[0].Price == 0 {
		t.Errorf("No valid Price.")
		return
	}
	// Deadline.
	if fResponse.Estimates[0].Deadline == 0 {
		t.Errorf("No valid deadline.")
		return
	}
	// Carrier name.
	if fResponse.Estimates[0].CarrierName == "" {
		t.Errorf("No valid carrier name.")
		return
	}
	// Carrier code.
	if fResponse.Estimates[0].CarrierCode == "" {
		t.Errorf("No valid carrier code.")
		return
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
*	DEALER FREIGHTS
*******************************************************************************/

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
