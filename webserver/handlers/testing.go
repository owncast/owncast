package handlers

import "net/http"

func (s *Handlers) HandleTesting(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("testing"))
}
