package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/3bd-dev/wallet-service/pkg/web"
	"github.com/google/uuid"
)

// Struct for SOAP requests and responses
type Request struct {
	XMLName     xml.Name `xml:"Envelope"`
	Amount      float64  `xml:"Body>amount"`
	CallbackURL string   `xml:"Body>callback_url"`
}

type Response struct {
	XMLName     xml.Name `xml:"SOAP-ENV:Envelope"`
	Status      string   `xml:"SOAP-ENV:Body>status"`
	Message     string   `xml:"SOAP-ENV:Body>message"`
	ReferenceID string   `xml:"SOAP-ENV:Body>id"`
}

// Helper function to generate a random reference ID
func generateReferenceID() string {
	return uuid.NewString()
}

// Async callback function to simulate delayed processing
func triggerAsyncCallback(referenceID, callbackURL string, amount float64) {
	time.Sleep(5 * time.Second) // Simulate delay in processing
	status := "failed"
	if validateTransactionAmount(amount) {
		status = "success"
	}
	log.Printf("Sending async callback for reference ID: %s with status: %s to %s\n", referenceID, status, callbackURL)
	callbackPayload := Response{
		Status:      status,
		Message:     "Transaction processed",
		ReferenceID: referenceID,
	}

	xmlData, _ := xml.Marshal(callbackPayload)
	resp, err := http.Post(callbackURL, "application/xml", io.NopCloser(bytes.NewReader(xmlData)))
	if err != nil {
		log.Printf("Failed to send callback: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Callback response: %s - %s", body, callbackPayload)
}

// Validate the transaction based on the amount
func validateTransactionAmount(amount float64) bool {
	return int(amount)%2 == 0
}

// deposit to simulate a deposit request
func deposit(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var request Request
	err := xml.Unmarshal(body, &request)
	if err != nil {
		web.RenderErr(w, err)
		return
	}
	referenceID := generateReferenceID()

	// Initial response with pending status
	response := Response{
		Status:      "pending",
		Message:     "Deposit is being processed in Gateway B",
		ReferenceID: referenceID,
	}

	renderResponse(w, response)

	// Trigger async callback after a delay
	go triggerAsyncCallback(referenceID, request.CallbackURL, request.Amount)
}

// withdrawal to simulate a withdrawal request
func withdrawal(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var request Request
	err := xml.Unmarshal(body, &request)
	if err != nil {
		web.RenderErr(w, err)
		return
	}

	referenceID := generateReferenceID()

	response := Response{
		Status:      "pending",
		Message:     "Withdrawal is being processed in Gateway B",
		ReferenceID: referenceID,
	}

	renderResponse(w, response)

	go triggerAsyncCallback(referenceID, request.CallbackURL, request.Amount)
}

func renderResponse(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(res)
}

// Main function to set up routes and start the server
func main() {
	http.HandleFunc("/deposit", deposit)
	http.HandleFunc("/withdraw", withdrawal)

	fmt.Println("Mock server Gateway B running on port 8091...")
	log.Fatal(http.ListenAndServe(":8091", nil))
}
