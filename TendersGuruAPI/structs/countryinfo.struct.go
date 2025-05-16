package structs

import "encoding/xml"

type RequestSoapService struct {
	XMLName     xml.Name    `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	BodyEnvelop BodyEnvelop `xml:"Body"`
}

type BodyEnvelop struct {
	XMLName         xml.Name        `xml:"Body"`
	SoapCountryName SoapCountryName `xml:"CountryName"`
}

type SoapCountryName struct {
	XMLName        xml.Name `xml:"http://www.oorsprong.org/websamples.countryinfo CountryName"`
	CountryISOCode string   `xml:"sCountryISOCode"`
}

type ResponseSoapService struct {
	XMLName        xml.Name       `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	BodyEnvelopRes BodyEnvelopRes `xml:"Body"`
}

type BodyEnvelopRes struct {
	XMLName            xml.Name              `xml:"Body"`
	SoapCountryNameRes SoapCountryNameResult `xml:"CountryNameResponse"`
}

type SoapCountryNameResult struct {
	XMLName     xml.Name `xml:"http://www.oorsprong.org/websamples.countryinfo CountryNameResponse"`
	CountryName string   `xml:"CountryNameResult"`
}
