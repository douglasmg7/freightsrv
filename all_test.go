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

// CEP.
func TestCEP(t *testing.T) {
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
func TestCorreios(t *testing.T) {
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

// Freight region.
// go test -run freightRegion, to run only this.
func TestFreightRegionDB(t *testing.T) {
	now := time.Now()
	nowFormated := now.Format(time.RFC3339)
	// log.Printf("datetime: %v", nowFormated)

	fr := freightRegion{
		Region:    "south",
		Weight:    4000,
		Deadline:  2,
		Price:     7845,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tx := sql3DB.MustBegin()
	// Update.
	uStatement := "UPDATE freight_region SET price=?, updated_at=? WHERE region=? AND weight=? AND deadline=?"
	uResult := tx.MustExec(uStatement, fr.Price, nowFormated, fr.Region, fr.Weight, fr.Deadline)
	uRowsAffected, err := uResult.RowsAffected()
	if err != nil {
		t.Errorf("Updating freight_region table. %s", err)
	}

	// Insert.
	if uRowsAffected == 0 {
		iStatement := "INSERT INTO freight_region(region, weight, deadline, price, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?)"
		iResult, err := tx.Exec(iStatement, fr.Region, fr.Weight, fr.Deadline, fr.Price, nowFormated, nowFormated)
		if err != nil {
			t.Errorf("Insert into freight_region table. %s", err)
		}
		iRowsAffected, err := iResult.RowsAffected()
		// log.Printf("iRowsAffected: %+v", iRowsAffected)
		if err != nil {
			t.Errorf("Insert into freight_region table. %s", err)
		}
		if iRowsAffected == 0 {
			t.Errorf("Insert into freight_region table not affected any row.")
		}
	}
	err = tx.Commit()
	if err != nil {
		t.Errorf("Commiting insert/update into freight_region table. %s", err)
	}

	// Select.
	var frResult freightRegion
	// err = sql3DB.Get(&frResult, "SELECT * FROM freight_region WHERE region=? AND weight=? AND deadline=?", fr.region, fr.weight, fr.deadline)
	err = sql3DB.Get(&frResult, "SELECT * FROM freight_region")
	if err != nil {
		t.Errorf("Getting freight_region row. %s", err)
	}
	// log.Printf("freightRegion: %+v", fr)
	// log.Printf("freightRegion: %+v", frResult)

	frUpdateAt := fr.UpdatedAt.Format(time.RFC3339)
	frResultUpdatedAt := frResult.UpdatedAt.Format(time.RFC3339)

	if frUpdateAt != frResultUpdatedAt {
		t.Errorf("Getting updated freight_region, UpdatedAt: %s, want %s", frUpdateAt, frResultUpdatedAt)
	}

	// Not worked.
	// strQuery := fmt.Sprintf("INSERT INTO freight_region(region, weight, deadline, price, created_at, updated_at) VALUES(\"%s\", %v, %v, %v, \"%s\", \"%s\")", fr.region, fr.weight, fr.deadline, fr.price, nowFormated, nowFormated)
	// strQueryConflit := fmt.Sprintf("%s ON CONFLICT DO UPDATE SET price=%v", strQuery, fr.price)
}
