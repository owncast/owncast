package controllers

import (
	"net/http"
)

// ServeCustomJavascript will serve optional custom Javascript.
func (s *Service) ServeCustomJavascript(w http.ResponseWriter, r *http.Request) {
	js := s.Data.GetCustomJavascript()
	_, _ = w.Write([]byte(js))
}
