package gatewayb

import "encoding/xml"

type Request struct {
	XMLName     xml.Name `xml:"SOAP-ENV:Envelope"`
	Amount      float64  `xml:"SOAP-ENV:Body>amount"`
	CallbackURL string   `xml:"SOAP-ENV:Body>callback_url"`
}

type Response struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		Status      string `xml:"status"`
		Message     string `xml:"message"`
		ReferenceID string `xml:"id"`
	} `xml:"Body"`
}
