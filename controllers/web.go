package controllers

import (
	"net/http"
)

// serveWeb will serve web assets.
func (s *Service) serveWeb(w http.ResponseWriter, r *http.Request) {
	s.static.ServeHTTP(w, r)
}
