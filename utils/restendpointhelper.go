package utils

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetURLParam retrieves the specified URL param from a given request.
func GetURLParam(r *http.Request, key string) (value string, err error) {
	value = chi.URLParam(r, key)
	if value == "" {
		err = errors.New("Request does not contain requested URL param")
	}
	return
}
