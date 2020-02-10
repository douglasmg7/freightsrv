package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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

// // Get freight region by id.
// func Test_freightValuesAPI(t *testing.T) {
// req, _ := http.NewRequest(http.MethodGet, "/freightsrv/freights/21170210", nil)
// req.SetBasicAuth(zoomUser(), zoomPass())
// res := httptest.NewRecorder()

// // indexHandler(res, req, nil)
// router.ServeHTTP(res, req)

// got := res.Body.String()
// want := "Hello!\n"

// if got != want {
// t.Errorf("got %q, want %q", got, want)
// }
// }
