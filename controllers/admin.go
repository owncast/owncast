package controllers

import (
	"net/http"

	"github.com/gabek/owncast/core"
	"github.com/gabek/owncast/core/rtmp"
)

// DisconnectInboundConnection will force-disconnect an inbound stream
func DisconnectInboundConnection(w http.ResponseWriter, r *http.Request) {
	if !core.GetStatus().Online {
		writeSimpleResponse(w, false, "no inbound stream connected")
		return
	}

	rtmp.Disconnect()
	writeSimpleResponse(w, true, "inbound stream disconnected")
}
