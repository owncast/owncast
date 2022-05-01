package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

type j map[string]interface{}

// InternalErrorHandler will return an error message as an HTTP response.
func InternalErrorHandler(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(j{"error": err.Error()}); err != nil {
		InternalErrorHandler(w, err)
	}
}

// BadRequestHandler will return an HTTP 500 as an HTTP response.
func BadRequestHandler(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	log.Debugln(err)

	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(j{"error": err.Error()}); err != nil {
		InternalErrorHandler(w, err)
	}
}

// WriteSimpleResponse will return a message as a response.
func WriteSimpleResponse(w http.ResponseWriter, success bool, message string) {
	response := models.BaseAPIResponse{
		Success: success,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")

	if success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		InternalErrorHandler(w, err)
	}
}

// WriteResponse will return an object as a JSON encoded response.
func WriteResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		InternalErrorHandler(w, err)
	}
}

// WriteString will return a basic string and a status code to the client.
func WriteString(w http.ResponseWriter, text string, status int) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	_, err := w.Write([]byte(text))
	return err
}
