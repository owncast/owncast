package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
)

type j map[string]interface{}

func errorResponseHandler(statusCode int, w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(j{"error": err.Error()}); err != nil {
		internalErrorHandler(w, err)
	}
}

func internalErrorHandler(w http.ResponseWriter, err error) {
	errorResponseHandler(http.StatusInternalServerError, w, err)
}

func badRequestHandler(w http.ResponseWriter, err error) {
	errorResponseHandler(http.StatusBadRequest, w, err)
}

func notFoundHandler(w http.ResponseWriter, err error) {
	errorResponseHandler(http.StatusNotFound, w, err)
}

func WriteSimpleResponse(w http.ResponseWriter, success bool, message string) {
	response := models.BaseAPIResponse{
		Success: success,
		Message: message,
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		internalErrorHandler(w, err)
	}
}

func urlFor(path string) string {
	return fmt.Sprintf("%s%s", config.Config.Federation.Site, path)
}
