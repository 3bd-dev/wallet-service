package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Struct for JSON requests and responses
type Request struct {
	Amount      float64 `json:"amount"`
	CallbackURL string  `json:"callback_url"`
}

type Response struct {
	Status      string `json:"status"`
	Message     string `json:"message"`
	ReferenceID string `json:"id"`
}

// Helper function to generate a random reference ID
func generateReferenceID() string {
	return uuid.NewString()
}

// Async callback function to simulate delayed processing
func triggerAsyncCallback(referenceID, callbackURL string, amount float64) {
	time.Sleep(5 * time.Second)
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

	jsonData, _ := json.Marshal(callbackPayload)
	resp, err := http.Post(callbackURL, "application/json", io.NopCloser(bytes.NewReader(jsonData)))
	if err != nil {
		log.Printf("Failed to send callback: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Callback response: %s", body)
}

// Validate the transaction based on the amount
func validateTransactionAmount(amount float64) bool {
	return int(amount)%2 == 0
}

func deposit(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var request Request
	json.Unmarshal(body, &request)

	referenceID := generateReferenceID()

	response := Response{
		Status:      "pending",
		Message:     "Deposit is being processed in Gateway A",
		ReferenceID: referenceID,
	}

	renderResponse(w, response)

	go triggerAsyncCallback(referenceID, request.CallbackURL, request.Amount)
}

// withdrawal simulates a withdrawal request
func withdrawal(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var request Request
	json.Unmarshal(body, &request)

	referenceID := generateReferenceID()

	response := Response{
		Status:      "pending",
		Message:     "Withdrawal is being processed in Gateway A",
		ReferenceID: referenceID,
	}
	renderResponse(w, response)

	go triggerAsyncCallback(referenceID, request.CallbackURL, request.Amount)
}

func renderResponse(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// Main function to set up routes and start the server
func main() {
	http.HandleFunc("/deposit", deposit)
	http.HandleFunc("/withdrawal", withdrawal)

	fmt.Println("Mock server Gateway A running on port 8090...")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
