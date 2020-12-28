package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

type j map[string]interface{}

func InternalErrorHandler(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	log.Errorln(err)

	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(j{"error": err.Error()}); err != nil {
		InternalErrorHandler(w, err)
	}
}

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

func WriteResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		InternalErrorHandler(w, err)
	}
}
