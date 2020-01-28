package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
	// "github.com/go-redis/redis/v7"
)

type viaCEPAddress struct {
	Cep      string `json:"cep"`
	Street   string `json:"logradouro"`
	District string `json:"bairro"`
	City     string `json:"localidade"`
	State    string `json:"uf"`
}

func addressFromCEP(cep string) (viaCEPAddress, error) {
	// log.Printf("addressFromCEP init: %v", cep)
	address := viaCEPAddress{}

	// Change to "00000000"
	cep = strings.ReplaceAll(cep, "-", "")

	// Check if CEP is valid "00000000".
	cepRE := regexp.MustCompile(`^\d{8}$`)
	if !cepRE.MatchString(cep) {
		return address, fmt.Errorf("CEP \"%s\" inválid", cep)
	}

	// Get address from.
	start := time.Now()
	res, err := http.Get(`https://viacep.com.br/ws/` + cep + `/json/`)
	if checkError(err) {
		return address, err
	}
	log.Printf("[debug] Viacep response time: %.1fs", time.Since(start).Seconds())

	// Read response.
	resBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if checkError(err) {
		return address, err
	}
	// log.Printf("address: %s", resBody)

	err = json.Unmarshal(resBody, &address)
	// log.Printf("address: %+v", address)

	return address, nil
}

// Get brazilian region from state.
func regionFromState(state string) string {
	state = strings.TrimSpace(state)
	state = strings.ToLower(state)

	switch state {
	case "ro", "ac", "am", "rr", "pa", "ap", "to":
		return "north"
	case "ma", "pi", "ce", "rn", "pb", "pe", "al", "se", "ba":
		return "northeast"
	case "ms", "mt", "go", "df":
		return "midwest"
	case "mg", "es", "rj", "sp":
		return "southeast"
	case "pr", "sc", "rs":
		return "south"
	}
	return ""
}

// Get region from cep.
func regionFromCEP(cep string) (region string, err error) {
	// Try cache.
	if region = getCEPRegion(cep); region != "" {
		return region, nil
	}

	// Retrive address.
	address, err := addressFromCEP(cep)
	if err != nil {
		return "", err
	}

	// Retrive region.
	if region = regionFromState(address.State); region != "" {
		setCEPRegion(cep, region)
	}

	return region, nil
}
