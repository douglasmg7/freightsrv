package main

import (
	"testing"
)

var cepDestiny = "5-76-25-000"

func TestCorreios(t *testing.T) {
	// testXML()
	p := pack{
		DestinyCEP: cepDestiny,
		Weight:     1500, // g.
		Length:     20,   // cm.
		Height:     30,   // cm.
		Width:      40,   // cm.
	}
	freights, err := correiosFreight(p)
	if checkError(err) {
		t.Errorf("correiosFreight() returned error. %v", err)
	}
	// log.Printf("Estimate freights: %+v", freights)

	for _, freight := range freights {
		result := freight.Carrier
		// Carrier.
		want := "Correios"
		if freight.Carrier != want {
			t.Errorf("Coerreios freight carrier name, result = %q, want %q", result, want)
		}
		if freight.Service == "" {
			t.Errorf("Correios freight service code must be != \"\"")
		}
		if freight.Price <= 0 {
			t.Errorf("Correios freight price must be more than 0")
		}
		if freight.DeadLine <= 0 {
			t.Errorf("Correios freight dead line must be more than 0")
		}
	}
}
