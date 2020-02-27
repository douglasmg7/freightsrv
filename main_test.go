package main

import (
	"errors"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

var cepNortheast = "5-76-25-000"

func TestMain(m *testing.M) {

	setupTest()
	code := m.Run()
	shutdownTest()

	os.Exit(code)
}

func setupTest() {
	initRedis()
	initSql3DB()

	// Clean cep region.
	cep := strings.ReplaceAll(cepNortheast, "-", "")
	err := redisDel("cep-region-" + cep)
	if err != nil {
		log.Printf("Deleting cep-region. %v\n", err)
	}
}

func shutdownTest() {
	closeRedis()
	closeSql3DB()
}

func Test_TextNormalization(t *testing.T) {
	want := "aaacu a"
	s := normalizeString("áàãçú â")
	// log.Println("string: ", s)
	if s != want {
		t.Errorf("Received %s, want %s", s, want)
	}
}

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

//*****************************************************************************
// CEP
//*****************************************************************************
// Get address by CEP.
func TestGetAddressByCEP(t *testing.T) {
	// Alredy check cep into redis_test.go.
	t.SkipNow()

	want := viaCEPAddress{
		Cep:      "31170-210",
		Street:   "Rua Deputado Bernardino de Sena Figueiredo",
		District: "Cidade Nova",
		City:     "Belo Horizonte",
		State:    "MG",
	}

	address, err := getAddressByCEP("3-1170210")
	if checkError(err) {
		t.Error(err)
	}

	if want.Cep != address.Cep || want.Street != address.Street || want.District != address.District || want.City != address.City || want.State != address.State {
		t.Errorf("addressFromCEP() = %q, want %q", address, want)
	}
}

// Get region by CEP.
func TestGetRegionByCEP(t *testing.T) {
	// First time get from rest api.
	want := "northeast"
	result, err := getRegionByCEP(cepNortheast)
	if err != nil {
		t.Errorf("Getting region from CEP. %v", err)
	}
	if result != want {
		t.Errorf("result = %q, want %q", result, want)
	}

	// Second time get from cache.
	result, err = getRegionByCEP(cepNortheast)
	if err != nil {
		t.Errorf("Getting region from CEP. %v", err)
	}
	if result != want {
		t.Errorf("result = %q, want %q", result, want)
	}
}

//*****************************************************************************
// CORREIOS
//*****************************************************************************
// Get Correios freight by pack.
func TestGetCorreiosFreightByPack(t *testing.T) {
	// testXML()
	p := &pack{
		CEPDestiny: cepNortheast,
		Weight:     1500, // g.
		Length:     20,   // cm.
		Height:     30,   // cm.
		Width:      40,   // cm.
	}

	c := make(chan *freightsOk)
	go getCorreiosFreightByPack(c, p)
	frsOk := <-c

	if !frsOk.Ok {
		t.Errorf("getCorreiosFreightByPack() not returned ok.")
	}

	for _, pFreight := range frsOk.Freights {
		// log.Printf("Correios freight: %+v", *pFreight)
		want := "Correios"
		if pFreight.Carrier != want {
			t.Errorf("Correios freight carrier name, carrier = %q, want %q", pFreight.Carrier, want)
		}
		if pFreight.ServiceCode == "" {
			t.Errorf("Correios freight service code must be != \"\"")
		}
		if pFreight.ServiceDesc == "" {
			t.Errorf("Correios freight service description must be != \"\"")
		}
		if pFreight.Price <= 0 {
			t.Errorf("Correios freight price must be more than 0")
		}
		if pFreight.Deadline <= 0 {
			t.Errorf("Correios freight dead line must be more than 0")
		}
	}
}

//*****************************************************************************
// Freight region
//*****************************************************************************
var createdFreightRegionID int
var fr regionFreight

// Get all freight regions.
func TestGetAllFreightRegion(t *testing.T) {
	frS, ok := getAllFreightRegion()
	if !ok {
		t.Errorf("Get all freight region returned not ok.")
	}
	frSLen := len(frS)
	if frSLen == 0 {
		t.Errorf("freigh_region rows: %v, want > 0.", frSLen)
	}
	// log.Printf("frS: %+v", frS)
}

// Create freight region.
func TestCreateFreightRegion(t *testing.T) {
	fr = regionFreight{
		Region:   "south",
		Weight:   4000,
		Deadline: 2,
		Price:    34567,
	}

	// Create freight.
	ok := createFreightRegion(&fr)
	if !ok {
		t.Errorf("Create freight region not returned ok.")
	}

	// Check saved data.
	var frResult regionFreight
	err = sql3DB.Get(&frResult, "SELECT * FROM freight_region WHERE region=? AND weight=? AND deadline=?", fr.Region, fr.Weight, fr.Deadline)
	if err != nil {
		t.Errorf("Getting created freight. %s", err)
	}

	if fr.Price != frResult.Price {
		t.Errorf("Getting updated freight_region, price: %v, want %v", frResult.Price, fr.Price)
	}
	fr.ID = frResult.ID

	// Not worked.
	// strQuery := fmt.Sprintf("INSERT INTO freight_region(region, weight, deadline, price, created_at, updated_at) VALUES(\"%s\", %v, %v, %v, \"%s\", \"%s\")", fr.region, fr.weight, fr.deadline, fr.price, nowFormated, nowFormated)
	// strQueryConflit := fmt.Sprintf("%s ON CONFLICT DO UPDATE SET price=%v", strQuery, fr.price)
}

// Get freight region by id.
func TestGetFreightRegionById(t *testing.T) {
	fr, ok := getFreightRegionById(fr.ID)
	if !ok {
		t.Errorf("Get freight region by id returned not ok.")
	}
	if fr.Region == "" {
		t.Errorf("region: %s, want not \"\".", fr.Region)
	}
}

// Update freight region.
func TestUpdateFreightRegion(t *testing.T) {
	fr.Price = 76543
	ok := updateFreightRegion(&fr)
	if !ok {
		t.Error("Update freight region retuned not ok.")
	}
}

// Delete freight region.
func TestDeleteFreightRegion(t *testing.T) {
	ok := deleteFreightRegion(fr.ID)
	if !ok {
		t.Error("Delete freight region retuned not ok.")
	}
}

// Get freight region by region and weight.
func TestGetFreightRegionByRegionAndWeight(t *testing.T) {
	frs, ok := getFreightRegionByRegionAndWeight("south", 3000)
	if !ok {
		t.Errorf("getFreightRegionByRegionAndWeight() returned not ok.")
	}
	if len(frs) == 0 {
		t.Errorf("getFreightRegionByRegionAndWeight() returned no one freight region.")
	}
	// All freight region must return the same weight.
	wantWeight := frs[0].Weight
	for _, fr := range frs {
		if fr.Weight != wantWeight {
			t.Errorf("Getting freight region by region and weight, weight: %v, want %v", fr.Weight, wantWeight)
		}
	}
	// log.Printf("frs: %+v", frs)
}

// Get freight region by CEP and wight.
func TestGetFreightRegionByCEPAndWeight(t *testing.T) {
	c := make(chan *freightsOk)
	go getFreightRegionByCEPAndWeight(c, "31-170210", 3000)
	frsOk := <-c
	if !frsOk.Ok {
		t.Errorf("getFreightRegionByCEPAndWeight() returned not ok.")
	}
	if len(frsOk.Freights) == 0 {
		t.Errorf("getFreightRegionByCEPAndWeight() returned no one freight.")
	}
	// Must have a valid price.
	for _, pFr := range frsOk.Freights {
		// log.Printf("region freight: %+v", *pFr)
		if pFr.Price <= 0 {
			t.Errorf("Getting freight region by CEP and weight, Price must be > 0")
		}
	}
}

//*****************************************************************************
// MOTOBOY FREIGHT
//*****************************************************************************
var validMotoboyFreightCity string
var validMotoboyFreightID int
var mf motoboyFreight

// Save motoboy freight.
func TestSaveMotoboyFreight(t *testing.T) {
	mf = motoboyFreight{
		City:     "Barão de Cocais",
		Deadline: 2,
		Price:    12570,
	}

	err := saveMotoboyFreight(&mf)
	if err != nil {
		t.Errorf("Saving freight region. %s", err)
		return
	}

	pmfr, ok := getMotoboyFreightByLocation("mg", "Barao de cocais")
	if !ok {
		t.Error("Motoboy freight returned not ok.")
	}
	if mf.Price != pmfr.Price {
		t.Errorf("Getting updated motoboy_freight, price: %v, want %v", pmfr.Price, mf.Price)
		return
	}
}

// Get all motoboy freight.
func TestGetAllMotoboyFreight(t *testing.T) {
	sl, err := getAllMotoboyFreight()
	if err != nil {
		t.Error(err)
		return
	}
	if len(sl) == 0 {
		t.Error("No returned row for motoboy_freight")
		return
	}
	// log.Printf("sl: %+v", sl)
	validMotoboyFreightID = sl[0].ID
	validMotoboyFreightCity = sl[0].City
}

// Get motoboy freight by id.
func TestGetMotoboyFreightByID(t *testing.T) {
	mf = motoboyFreight{
		ID: validMotoboyFreightID,
	}
	pmf, ok := getMotoboyFreightByID(validMotoboyFreightID)
	if !ok {
		t.Error("Motoboy freight returned not ok.")
		return
	}
	if pmf.City != validMotoboyFreightCity {
		t.Errorf("Motoboy fregiht citiy = %s, want: %s.", pmf.City, validMotoboyFreightCity)
		return
	}
	// log.Printf("mf by id: %+v", mf)
}

// Get motoboy freight by location.
func TestGetMotoboyFreightByLocation(t *testing.T) {
	mf = motoboyFreight{
		City: validMotoboyFreightCity,
	}
	pmf, ok := getMotoboyFreightByLocation("MG", validMotoboyFreightCity)
	if !ok {
		t.Error("Motoboy freight returned not ok.")
	}

	if pmf.ID == 0 {
		t.Errorf("motoboy fregiht ID = %d, want: %d.", pmf.ID, validMotoboyFreightID)
		return
	}
	// log.Printf("mf by city: %+v", mf)
}

// Get motoboy freight by location.
func TestGetMotoboyFreightByCEP(t *testing.T) {
	c := make(chan *freightsOk)
	go getMotoboyFreightByCEP(c, "31130210")

	frsOk := <-c
	if !frsOk.Ok {
		t.Error("Motoboy freight returned not ok.")
	}

	if frsOk.Freights[0].Deadline == 0 {
		t.Errorf("motoboy fregiht deadline = 0, want > 0.")
		return
	}
	// log.Printf("*frsOk.Freights[0]: %+v", *frsOk.Freights[0])
}

// Delete motoboy freight.
func TestDelMotoboyFreight(t *testing.T) {
	ok := deleteMotoboyFreight(validMotoboyFreightID)
	if !ok {
		t.Error(errors.New("Delete motoboy freight returned no ok."))
	}
}
