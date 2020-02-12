package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

// Freight deadline and price.
func TestFreightAPI(t *testing.T) {
	p := pack{
		DestinyCEP: "5-76-25-000",
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

	req, _ := http.NewRequest(http.MethodGet, "/freightsrv/freights", bytes.NewBuffer(reqBody))
	req.SetBasicAuth("bypass", "123456")
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	got := res.Body.String()
	want := "Some value\n"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
