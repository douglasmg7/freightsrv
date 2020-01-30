package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	CORREIOS_URL                    = `http://ws.correios.com.br/calculador/CalcPrecoPrazo.asmx/CalcPrecoPrazo`
	CORREIOS_SERVICES_CODE          = "4596, 4553"
	CORREIOS_PACKAGE_FORMAT         = "1" // 1 - caixa/pacote, 2 - rolo/prisma, 3 - Envelope.
	CORREIOS_PACKAGE_DIAMETER       = "0" // Diâmetro em cm.
	CORREIOS_OWN_HAND               = "N" // Se a encomenda será entregue com o serviço adicional mão própria.
	CORREIOS_DECLARED_VALUE         = "0" // Valor delcarado.
	CORREIOS_ACKNOWLEDGMENT_RECEIPT = "N" // Aviso de recebimento.
)

type pack struct {
	OriginCEP  string
	DestinyCEP string
	Weight     int // g.
	Length     int // cm.
	Height     int // cm.
	Width      int // cm.
}

func (p *pack) Validate() error {

	regCep := regexp.MustCompile(`^[0-9]{8}$`)

	// Origin CEP.
	p.OriginCEP = strings.ReplaceAll(p.OriginCEP, "-", "")
	if p.OriginCEP == "" {
		p.OriginCEP = CEP_ORIGIN
	}
	if !regCep.MatchString(p.OriginCEP) {
		return fmt.Errorf("Origin CEP \"%v\" invalid.", p.OriginCEP)
	}

	// Destiny CEP.
	p.DestinyCEP = strings.ReplaceAll(p.DestinyCEP, "-", "")
	if !regCep.MatchString(p.DestinyCEP) {
		return fmt.Errorf("Destiny CEP \"%v\" invalid.", p.DestinyCEP)
	}

	// Weight in kg.
	minWeight := 1
	maxWeight := 50000
	if p.Weight < minWeight {
		return fmt.Errorf("Wight must be more then %v g", minWeight)
	}
	if p.Weight > maxWeight {
		return fmt.Errorf("Wight must be less then %v kg", maxWeight)
	}

	// Length in cm.
	minLength := 1
	maxLength := 100
	if p.Length < minLength {
		return fmt.Errorf("Length must be more then %v cm", minLength)
	}
	if p.Length > maxLength {
		return fmt.Errorf("Length must be less then %v cm", maxLength)
	}

	// Height in cm.
	minHeight := 1
	maxHeight := 100
	if p.Height < minHeight {
		return fmt.Errorf("Height must be more then %v cm", minHeight)
	}
	if p.Height > maxHeight {
		return fmt.Errorf("Height must be less then %v cm", maxHeight)
	}

	// Width in cm.
	minWidth := 1
	maxWidth := 100
	if p.Width < minWidth {
		return fmt.Errorf("Width must be more then %v cm", minWidth)
	}
	if p.Width > maxWidth {
		return fmt.Errorf("Width must be less then %v cm", maxWidth)
	}

	return nil
}

type correiosXMLService struct {
	Code     int    `xml:"Codigo"`
	Price    string `xml:"Valor"`
	DeadLine int    `xml:"PrazoEntrega"`
	Error    int    `xml:"Erro"`
	MsgError string `xml:"MsgErro"`
}

type correiosXMLServices struct {
	Services []correiosXMLService `xml:"cServico"`
}

type correiosXMLResult struct {
	XMLName xml.Name            `xml:"cResultado"`
	Result  correiosXMLServices `xml:"Servicos"`
}

func correiosFreight(cPack pack) (freights []freight, err error) {
	err = cPack.Validate()
	if err != nil {
		return freights, err
	}

	reqBody := []byte(`nCdEmpresa=` + CORREIOS_COMPANY_ADMIN_CODE +
		`&sDsSenha=` + CORREIOS_COMPANY_PASSWORD +
		`&nCdServico=` + CORREIOS_SERVICES_CODE +
		`&sCepOrigem=` + cPack.OriginCEP +
		`&sCepDestino=` + cPack.DestinyCEP +
		`&nVlPeso=` + strconv.Itoa(cPack.Weight/1000) +
		`&nCdFormato=` + CORREIOS_PACKAGE_FORMAT +
		`&nVlComprimento=` + strconv.Itoa(cPack.Length) +
		`&nVlAltura=` + strconv.Itoa(cPack.Height) +
		`&nVlLargura=` + strconv.Itoa(cPack.Width) +
		`&nVlDiametro=` + CORREIOS_PACKAGE_DIAMETER +
		`&sCdMaoPropria=` + CORREIOS_OWN_HAND +
		`&nVlValorDeclarado=` + CORREIOS_DECLARED_VALUE +
		`&sCdAvisoRecebimento=` + CORREIOS_ACKNOWLEDGMENT_RECEIPT)

	// Log request.
	// log.Println("request body: " + string(reqBody))

	// Request product add.
	client := &http.Client{}
	req, err := http.NewRequest("POST", CORREIOS_URL, bytes.NewBuffer(reqBody))
	if checkError(err) {
		return freights, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	start := time.Now()
	res, err := client.Do(req)
	if checkError(err) {
		return freights, err
	}
	log.Printf("[debug] Correios response time: %.1fs", time.Since(start).Seconds())

	defer res.Body.Close()
	if checkError(err) {
		return freights, err
	}

	// Result.
	resBody, err := ioutil.ReadAll(res.Body)
	if checkError(err) {
		return freights, err
	}
	// log.Println("resBody:", string(resBody))

	result := correiosXMLResult{}
	err = xml.Unmarshal(resBody, &result)
	if checkError(err) {
		return freights, err
	}
	// log.Printf("\n\nresult: %+v", result)

	for _, service := range result.Result.Services {
		// log.Printf("service: %+v", service)
		if service.Error != 0 {
			log.Printf("[warning] [correios] origin: %s, destiny: %s, code: %d, error: %d, message: %v", cPack.OriginCEP, cPack.DestinyCEP, service.Code, service.Error, service.MsgError)
			continue
		}
		// Convert to float64.
		price := strings.ReplaceAll(service.Price, ".", "")
		price = strings.ReplaceAll(service.Price, ",", ".")
		priceF, err := strconv.ParseFloat(price, 64)
		if checkError(err) {
			continue
		}
		// log.Printf("Price: %v", priceF)
		freights = append(freights, freight{Carrier: "Correios", Service: strconv.Itoa(service.Code), Price: priceF, DeadLine: service.DeadLine})
	}

	return freights, nil
}

func testXML() {
	testString := []byte(`<cResultado>
	  <Servicos>
		<cServico>
		  <Codigo>4596</Codigo>
		  <Up>232</Up>
		</cServico>
		<cServico>
		  <Codigo>4553</Codigo>
		  <Up>333</Up>
		</cServico>
	  </Servicos>
	</cResultado>`)
	data := correiosXMLResult{}

	// data.Services = []correiosXMLService{}
	// data.Services = append(data.Services, correiosXMLService{123, 345})
	// data.Services = append(data.Services, correiosXMLService{223, 545})
	// log.Printf("\n\ndata: %+v", data)

	// res, err := xml.MarshalIndent(data, "", "    ")
	// checkError(err)
	// log.Printf("\n\nres: %v", string(res))
	// return

	err = xml.Unmarshal(testString, &data)

	checkError(err)
	log.Printf("\n\ndata: %+v", data)
}
