package controllers

import (
	"net/http"
)

// DisconnectInboundConnection will force-disconnect an inbound stream.
func (s *Service) DisconnectInboundConnection(w http.ResponseWriter, r *http.Request) {
	s.Rtmp.Disconnect()
	w.WriteHeader(http.StatusOK)
}
