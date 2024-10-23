package web

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// Response represents the structure of an HTTP response.
type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// httpStatus is an interface for types that can provide an HTTP status code.
type httpStatus interface {
	HTTPStatus() int
}

// validationError is an interface for types that can provide validation error details.
type validationError interface {
	Fields() map[string]string
}

// RenderJSON renders a value as JSON in the HTTP response.
func RenderJSON(w http.ResponseWriter, code int, value interface{}, msg string, det interface{}) error {
	res := Response{
		Code:    code,
		Data:    value,
		Message: msg,
		Details: det,
	}

	buffer, err := json.Marshal(res)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(buffer)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		return err
	}

	return nil
}

// RenderOk renders a successful response with HTTP status 200.
func RenderOk(w http.ResponseWriter, value interface{}) error {
	return RenderJSON(w, http.StatusOK, value, "", nil)
}

// RenderNoContent renders a response with HTTP status 204 (No Content).
func RenderNoContent(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}

// RenderErr responds with an error as JSON in the HTTP response.
func RenderErr(w http.ResponseWriter, err error) error {
	var statusCode = http.StatusInternalServerError
	var det interface{} = nil
	switch v := err.(type) {
	case validationError:
		statusCode = http.StatusBadRequest
		det = v.Fields()
		err = errors.New("validation error")
	case httpStatus:
		statusCode = v.HTTPStatus()
	}
	return RenderJSON(w, statusCode, nil, err.Error(), det)
}
