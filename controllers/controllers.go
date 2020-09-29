package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gabek/owncast/models"
)

type j map[string]interface{}

func internalErrorHandler(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(j{"error": err.Error()})
}

func badRequestHandler(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(j{"error": err.Error()})
}

func writeSimpleResponse(w http.ResponseWriter, success bool, message string) {
	response := models.BaseAPIResponse{
		Success: success,
		Message: message,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
