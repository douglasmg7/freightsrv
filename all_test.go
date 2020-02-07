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

// CEP.
func Test_CEP(t *testing.T) {
	// Alredy check cep into redis_test.go.
	t.SkipNow()

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

// Correios.
func Test_Correios(t *testing.T) {
	// testXML()
	p := pack{
		DestinyCEP: cepNortheast,
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

func Test_Redis(t *testing.T) {
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

func Test_RegionFromCEP(t *testing.T) {
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

//*****************************************************************************
// FREIGHT REGION
//*****************************************************************************
var validFreightRegionId int
var fr freightRegion

// Save freight region.
func Test_SaveFreightRegion(t *testing.T) {
	fr = freightRegion{
		Region:   "south",
		Weight:   4000,
		Deadline: 2,
		Price:    7845,
	}

	err := saveFreightRegion(fr)
	if err != nil {
		t.Errorf("Saving freight region. %s", err)
	}
	// log.Printf("updatedAt        : %+v", updatedAt)
	// log.Printf("updatedAt        : %+s", updatedAt.Format(time.RFC3339))

	// Check saved data.
	var frResult freightRegion
	err = sql3DB.Get(&frResult, "SELECT * FROM freight_region WHERE region=? AND weight=? AND deadline=?", fr.Region, fr.Weight, fr.Deadline)
	if err != nil {
		t.Errorf("Getting freight_region row. %s", err)
	}
	// log.Printf("frResult createdAt val: %+v", frResult.CreatedAt)
	// log.Printf("frResult updatedAt val: %+v", frResult.UpdatedAt)

	// log.Printf("frResult createdAt str: %+s", frResult.CreatedAt.Format(time.RFC3339))
	// log.Printf("frResult updatedAt str: %+s", frResult.UpdatedAt.Format(time.RFC3339))

	if fr.Price != frResult.Price {
		t.Errorf("Getting updated freight_region, price: %v, want %v", frResult.Price, fr.Price)
	}

	// Not worked.
	// strQuery := fmt.Sprintf("INSERT INTO freight_region(region, weight, deadline, price, created_at, updated_at) VALUES(\"%s\", %v, %v, %v, \"%s\", \"%s\")", fr.region, fr.weight, fr.deadline, fr.price, nowFormated, nowFormated)
	// strQueryConflit := fmt.Sprintf("%s ON CONFLICT DO UPDATE SET price=%v", strQuery, fr.price)
}

// Get all freight regions.
func Test_GetAllFreightRegion(t *testing.T) {
	frS, err := getAllFreightRegion()
	if err != nil {
		t.Errorf("TestGetAllFreightRegion(). %s", err)
	}
	frSLen := len(frS)
	if frSLen == 0 {
		t.Errorf("freigh_region rows: %v, want > 0.", frSLen)
	}
	validFreightRegionId = frS[0].ID
	// log.Printf("frS: %+v", frS)
}

// Get freight region by id.
func Test_GetFreightRegionById(t *testing.T) {
	fr, err := getFreightRegionById(validFreightRegionId)
	if err != nil {
		t.Errorf(" TestGetFreightRegionById(). %s", err)
	}
	if fr.Region == "" {
		t.Errorf("region: %s, want not \"\".", fr.Region)
	}
	// log.Printf("fr: %+v", fr)
}

// Delete freight region.
func Test_DelFreightRegion(t *testing.T) {
	err := delFreightRegion(validFreightRegionId)
	if err != nil {
		t.Error(err)
	}
}

//*****************************************************************************
// MOTOBOY FREIGHT
//*****************************************************************************
var validMotoboyFreightCity string
var validMotoboyFreightID int
var mf motoboyFreight

// Save motoboy freight.
func Test_SaveMotoboyFreight(t *testing.T) {
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

	// Check saved data.
	mfResult := motoboyFreight{
		City: "Barao de cocais",
	}
	err = getMotoboyFreight(&mfResult)
	if err != nil {
		t.Error(err)
		return
	}
	if mf.Price != mfResult.Price {
		t.Errorf("Getting updated motoboy_freight, price: %v, want %v", mfResult.Price, mf.Price)
		return
	}
}

// Get all motoboy freight.
func Test_GetAllMotoboyFreight(t *testing.T) {
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
func Test_GetMotoboyFreightByID(t *testing.T) {
	mf = motoboyFreight{
		ID: validMotoboyFreightID,
	}
	err := getMotoboyFreight(&mf)
	if err != nil {
		t.Error(err)
		return
	}
	if mf.City == "" {
		t.Errorf("motoboy fregiht citiy = %s, want: %s.", mf.City, validMotoboyFreightCity)
		return
	}
	// log.Printf("mf by id: %+v", mf)
}

// Get motoboy freight by city.
func Test_GetMotoboyFreightByCity(t *testing.T) {
	mf = motoboyFreight{
		City: validMotoboyFreightCity,
	}
	err := getMotoboyFreight(&mf)
	if err != nil {
		t.Error(err)
		return
	}
	if mf.ID == 0 {
		t.Errorf("motoboy fregiht ID = %d, want: %d.", mf.ID, validMotoboyFreightID)
		return
	}
	// log.Printf("mf by city: %+v", mf)
}

// Delete motoboy freight.
func Test_DelMotoboyFreight(t *testing.T) {
	err := delMotoboyFreight(validMotoboyFreightID)
	if err != nil {
		t.Error(err)
	}
}

//*****************************************************************************
// TEST SERVER
//*****************************************************************************
