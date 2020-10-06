package controllers

import (
	"net/http"

	"github.com/owncast/owncast/core/rtmp"
)

// DisconnectInboundConnection will force-disconnect an inbound stream
func DisconnectInboundConnection(w http.ResponseWriter, r *http.Request) {
	rtmp.Disconnect()
	w.WriteHeader(http.StatusOK)
}
