package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Get freight region by id.
func Test_IndexHandler(t *testing.T) {
	log.Println("production:", production)

	req, _ := http.NewRequest(http.MethodGet, "/productsrv", nil)
	req.SetBasicAuth(zoomUser(), zoomPass())
	res := httptest.NewRecorder()

	// indexHandler(res, req, nil)
	router.ServeHTTP(res, req)

	got := res.Body.String()
	want := "Hello!\n"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
