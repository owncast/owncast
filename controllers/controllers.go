package controllers

import (
	"encoding/json"
	"net/http"
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
