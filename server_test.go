package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// // Get freight region by id.
// func Test_IndexHandler(t *testing.T) {
// req, _ := http.NewRequest(http.MethodGet, "/", nil)
// res := httptest.NewRecorder()

// indexHandler(res, req, nil)

// got := res.Body.String()
// want := "Hello!\n"

// if got != want {
// t.Errorf("got %q, want %q", got, want)
// }
// }

// Get freight region by id.
func Test_IndexHandler(t *testing.T) {
	// curl -u zoomteste_zunka:H2VA79Ug4fjFsJb localhost:8084/productsrv

	// client := &http.Client{}
	// req, err = http.NewRequest("POST", zunkaSiteHost()+"/setup/product/add", bytes.NewBuffer(reqBody))
	// req.Header.Set("Content-Type", "application/json")
	// HandleError(w, err)
	// req.SetBasicAuth(zunkaSiteUser(), zunkaSitePass())
	// res, err := client.Do(req)
	// HandleError(w, err)

	req, _ := http.NewRequest(http.MethodGet, "/productsrv", nil)
	req.SetBasicAuth(zoomHost(), zoomUser())
	res := httptest.NewRecorder()

	// indexHandler(res, req, nil)
	router.ServeHTTP(res, req)

	got := res.Body.String()
	want := "Hello!\n"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
