package main

import (
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

var cepNortheast = "5-76-25-000"

func TestMain(m *testing.M) {
	setupT()
	code := m.Run()
	shutdownT()
	os.Exit(code)
}

func setupT() {
	cep := strings.ReplaceAll(cepNortheast, "-", "")
	err := redisDel("cep-region-" + cep)
	if err != nil {
		log.Printf("Deleting cep-region. %v", err)
	}
}

func shutdownT() {
}

// p := pack{
// DestinyCEP: "35460000",
// Weight:     1500,
// Length:     20,
// Height:     30,
// Width:      40,
// }
// freights, err := correiosFreight(p)
// if !checkError(err) {
// log.Printf("Estimate freights: %+v", freights)
// }
// // testXML()

func TestRedis(t *testing.T) {
	want := "Hello!"
	key := "freightsrv-test"

	err := redisSet(key, want, time.Second*10)
	if err != nil {
		t.Errorf("Saving on redis DB. %v", err)
	}

	result := redisGet(key)
	if result != want {
		t.Errorf("result = %q, want %q", result, want)
	}
}

func TestRegionFromCEP(t *testing.T) {
	// First time get from rest api.
	want := "northeast"
	result, err := regionFromCEP(cepNortheast)
	if err != nil {
		t.Errorf("Getting region from CEP. %v", err)
	}
	if result != want {
		t.Errorf("result = %q, want %q", result, want)
	}

	// Second time get from cache.
	result, err = regionFromCEP(cepNortheast)
	if err != nil {
		t.Errorf("Getting region from CEP. %v", err)
	}
	if result != want {
		t.Errorf("result = %q, want %q", result, want)
	}
}
