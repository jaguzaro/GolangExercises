package controllers

import (
	stc "TendersGuruAPI/structs"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
)

func GetCountry(ctx context.Context) (string, error) {
	var soapRequest stc.RequestSoapService
	soapRequest.BodyEnvelop.SoapCountryName.CountryISOCode = "CO"
	endpoint := `http://webservices.oorsprong.org/websamples.countryinfo/CountryInfoService.wso`

	var byteValue bytes.Buffer
	err := xml.NewEncoder(&byteValue).Encode(soapRequest)
	if err != nil {
		fmt.Println("Error to parse request")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, &byteValue)
	if err != nil {
		fmt.Println("Error")
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://www.oorsprong.org/websamples.countryinfo/CountryName")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error")
	}
	defer resp.Body.Close()

	var responseSOAP stc.ResponseSoapService
	if err := xml.NewDecoder(resp.Body).Decode(&responseSOAP); err != nil {
		fmt.Println("Error")
	}

	countryName := responseSOAP.BodyEnvelopRes.SoapCountryNameRes.CountryName
	return countryName, nil
}
