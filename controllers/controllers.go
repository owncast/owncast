package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/models"
)

type j map[string]interface{}

func internalErrorHandler(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(j{"error": err.Error()}); err != nil {
		internalErrorHandler(w, err)
	}
}

func badRequestHandler(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(j{"error": err.Error()}); err != nil {
		internalErrorHandler(w, err)
	}
}

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
		internalErrorHandler(w, err)
	}
}
