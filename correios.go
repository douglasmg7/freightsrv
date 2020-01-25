package main

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

// var correiosUrl = "http://ws.correios.com.br/calculador/CalcPrecoPrazo.asmx?wsdl"
var correiosUrl = `http://ws.correios.com.br/calculador/CalcPrecoPrazo.asmx/CalcPrecoPrazo`

var nCdEmpresa = "18037020" // Código administrativo junto à ECT (para clientes com contrato) .
var sDsSenha = "15178404"
var nCdServico = "4596, 4553"
var sCepOrigem = "35460000"
var sCepDestino = "31170210"
var nVlPeso = "1"         // Weight in Kg.
var nCdFormato = "1"      // 1 - caixa/pacote, 2 - rolo/prisma, 3 - Envelope.
var nVlComprimento = "20" // Lenght in cm.
var nVlAltura = "30"      // Height in cm.
var nVlLargura = "40"     // Width in cm.
var nVlDiametro = "0"     // Diâmetro em cm.
var sCdMaoPropria = "N"   // Se a encomenda será entregue com o serviço adicional mão própria.
var nVlValorDeclarado = "0"
var sCdAvisoRecebimento = "N"

type correiosPackage struct {
	cepOrigin  string
	CepDestiny string
	Weight     int
	Length     int
	Height     int
	Width      int
}

type xmlService struct {
	Code     int    `xml:"Codigo"`
	Price    string `xml:"Valor"`
	DeadLine int    `xml:"PrazoEntrega"`
	Error    int    `xml:"Erro"`
	MsgError string `xml:"MsgErro"`
}

type xmlServices struct {
	Services []xmlService `xml:"cServico"`
}

type xmlResult struct {
	XMLName xml.Name    `xml:"cResultado"`
	Result  xmlServices `xml:"Servicos"`
}

func correiosFreight() {

	reqBody := []byte(`nCdEmpresa=` + nCdEmpresa +
		`&sDsSenha=` + sDsSenha +
		`&nCdServico=` + nCdServico +
		`&sCepOrigem=` + sCepOrigem +
		`&sCepDestino=` + sCepDestino +
		`&nVlPeso=` + nVlPeso +
		`&nCdFormato=` + nCdFormato +
		`&nVlComprimento=` + nVlComprimento +
		`&nVlAltura=` + nVlAltura +
		`&nVlLargura=` + nVlLargura +
		`&nVlDiametro=` + nVlDiametro +
		`&sCdMaoPropria=` + sCdMaoPropria +
		`&nVlValorDeclarado=` + nVlValorDeclarado +
		`&sCdAvisoRecebimento=` + sCdAvisoRecebimento)

	// Log request.
	log.Println("request body: " + string(reqBody))

	// Request product add.
	client := &http.Client{}
	req, err := http.NewRequest("POST", correiosUrl, bytes.NewBuffer(reqBody))
	checkError(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	checkError(err)

	defer res.Body.Close()
	checkError(err)

	// Result.
	resBody, err := ioutil.ReadAll(res.Body)
	checkError(err)
	log.Println("resBody:", string(resBody))

	result := xmlResult{}
	err = xml.Unmarshal(resBody, &result)
	checkError(err)
	log.Printf("\n\nresult: %+v", result)
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
	data := xmlResult{}

	// data.Services = []xmlService{}
	// data.Services = append(data.Services, xmlService{123, 345})
	// data.Services = append(data.Services, xmlService{223, 545})
	// log.Printf("\n\ndata: %+v", data)

	// res, err := xml.MarshalIndent(data, "", "    ")
	// checkError(err)
	// log.Printf("\n\nres: %v", string(res))
	// return

	err = xml.Unmarshal(testString, &data)

	checkError(err)
	log.Printf("\n\ndata: %+v", data)
}
