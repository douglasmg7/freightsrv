package main

import (
	"testing"
)

func TestRedis(t *testing.T) {
	want := viaCEPAddress{
		Cep:      "31170-210",
		Street:   "Rua Deputado Bernardino de Sena Figueiredo",
		District: "Cidade Nova",
		City:     "Belo Horizonte",
		State:    "MG",
	}

	address, err := addressFromCEP("3-1170210")
	if checkError(err) {
		t.Error(err)
	}

	if want.Cep != address.Cep || want.Street != address.Street || want.District != address.District || want.City != address.City || want.State != address.State {
		t.Errorf("addressFromCEP() = %q, want %q", address, want)
	}
}
