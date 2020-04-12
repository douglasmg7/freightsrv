package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	CORREIOS_URL              = `http://ws.correios.com.br/calculador/CalcPrecoPrazo.asmx/CalcPrecoPrazo`
	CORREIOS_SERVICES_CODE    = "4596,4553" // 4596-PAC, 4553-SEDEX
	CORREIOS_PACKAGE_FORMAT   = "1"         // 1 - caixa/pacote, 2 - rolo/prisma, 3 - Envelope.
	CORREIOS_PACKAGE_DIAMETER = "0"         // Diâmetro em cm.
	CORREIOS_OWN_HAND         = "N"         // Se a encomenda será entregue com o serviço adicional mão própria.
	// CORREIOS_DECLARED_VALUE         = "0"          // Valor delcarado.
	CORREIOS_ACKNOWLEDGMENT_RECEIPT = "N" // Aviso de recebimento.
)

func (p *pack) ValidateCorreios() bool {
	// Basic validation.
	if !p.Validate() {
		return false
	}
	// Length in cm.
	minLength := 15
	maxLength := 105
	if p.Length < minLength {
		log.Printf("[warning] [correios] Pack length changed from %v cm to %v cm", p.Length, minLength)
		p.Length = minLength
	}
	if p.Length > maxLength {
		log.Printf("[warning] [correios] Correios shipping will not be estimated. Length of %v cm greater than %v cm.", p.Length, maxLength)
		return false
	}

	// Width in cm.
	minWidth := 10
	maxWidth := 105
	if p.Width < minWidth {
		log.Printf("[warning] [correios] Pack width changed from %v cm to %v cm", p.Width, minWidth)
		p.Width = minWidth
	}
	if p.Width > maxWidth {
		log.Printf("[warning] [correios] Correios shipping will not be estimated. Width of %v cm greater than %v cm.", p.Width, maxWidth)
		return false
	}

	// Height in cm.
	minHeight := 1
	maxHeight := 105
	if p.Height < minHeight {
		log.Printf("[warning] [correios] Pack height changed from %v cm to %v cm", p.Height, minHeight)
		p.Height = minHeight
	}
	if p.Height > maxHeight {
		log.Printf("[warning] [correios] Correios shipping will not be estimated. Height of %v cm greater than %v cm.", p.Height, maxHeight)
		return false
	}

	// Dimensions sum.
	sum := p.Length + p.Width + p.Height
	minSum := 26
	maxSum := 200
	if sum < minSum {
		log.Printf("[warning] [correios] Correios shipping will not be estimated. Sum dimensions of %v cm less than %v cm.", sum, minSum)
		return false
	}
	if sum > maxSum {
		log.Printf("[warning] [correios] Correios shipping will not be estimated. Sum dimensions of %v cm greater than %v cm.", sum, maxSum)
		return false
	}

	return true
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

// Get correios freight by pack.
func getCorreiosFreightByPack(c chan *freightsOk, p *pack) {
	result := &freightsOk{
		Freights: []*freight{},
	}

	if !p.ValidateCorreios() {
		c <- result
		return
	}

	// Get from cache.
	temp, ok := getCorreiosCache(p)
	if ok {
		// log.Printf("result: %+v", temp)
		result.Freights = temp
		result.Ok = true
		c <- result
		return
	}
	// Not in the cache.
	reqBody := []byte(`nCdEmpresa=` + CORREIOS_COMPANY_ADMIN_CODE +
		`&sDsSenha=` + CORREIOS_COMPANY_PASSWORD +
		`&nCdServico=` + CORREIOS_SERVICES_CODE +
		`&sCepOrigem=` + p.CEPOrigin +
		`&sCepDestino=` + p.CEPDestiny +
		`&nCdFormato=` + CORREIOS_PACKAGE_FORMAT +
		`&nVlComprimento=` + strconv.Itoa(p.Length) +
		`&nVlAltura=` + strconv.Itoa(p.Height) +
		`&nVlLargura=` + strconv.Itoa(p.Width) +
		`&nVlPeso=` + fmt.Sprintf("%.3f", (float64(p.Weight)/1000)) + // Kg.
		`&nVlDiametro=` + CORREIOS_PACKAGE_DIAMETER +
		`&sCdMaoPropria=` + CORREIOS_OWN_HAND +
		`&nVlValorDeclarado=` + fmt.Sprintf("%.2f", p.Price) +
		// `&nVlValorDeclarado=` + "0" +
		`&sCdAvisoRecebimento=` + CORREIOS_ACKNOWLEDGMENT_RECEIPT)

	// Log request.
	// log.Println("[debug] Correios request body: " + string(reqBody))

	// Request product add.
	client := &http.Client{}
	req, err := http.NewRequest("POST", CORREIOS_URL, bytes.NewBuffer(reqBody))
	if checkError(err) {
		c <- result
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	start := time.Now()
	res, err := client.Do(req)
	if checkError(err) {
		c <- result
		return
	}
	log.Printf("[debug] Correios response time: %.1fs", time.Since(start).Seconds())

	defer res.Body.Close()
	if checkError(err) {
		c <- result
		return
	}

	// Result.
	resBody, err := ioutil.ReadAll(res.Body)
	if checkError(err) {
		c <- result
		return
	}
	// log.Println("resBody:", string(resBody))

	rCorreios := correiosXMLResult{}
	err = xml.Unmarshal(resBody, &rCorreios)
	if checkError(err) {
		c <- result
		return
	}
	// log.Printf("\n\nresult: %+v", result)

	for _, service := range rCorreios.Result.Services {
		// log.Printf("service: %+v", service)
		if service.Error != 0 {
			// log.Printf("[warning] [correios] pack: %+v, code: %d, error: %d, message: %v", p, service.Code, service.Error, service.MsgError)
			log.Printf("[warning] [correios] service code: %d, error: %d\n\tmessage: %v\n\tpack: %+v\n\treqBody: %s", service.Code, service.Error, service.MsgError, p, reqBody)
			continue
		}
		// Service description.
		var serviceDesc string
		switch service.Code {
		case 4596:
			serviceDesc = "PAC"
		case 4553:
			serviceDesc = "SEDEX"
		}

		// Convert price to float64.
		price := strings.ReplaceAll(service.Price, ".", "")
		price = strings.ReplaceAll(service.Price, ",", ".")
		priceF, err := strconv.ParseFloat(price, 64)
		if checkError(err) {
			continue
		}
		// log.Printf("Price: %v", priceF)
		result.Freights = append(result.Freights, &freight{Carrier: "Correios", ServiceCode: strconv.Itoa(service.Code), ServiceDesc: serviceDesc, Price: priceF, Deadline: service.DeadLine})
	}
	// log.Printf("result: %+v", result)
	// log.Printf("result.Freights: %+v", result.Freights[0])
	// Not cache empty values.
	if len(result.Freights) > 0 {
		setCorreiosCache(p, result.Freights)
	}
	result.Ok = true
	c <- result
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
