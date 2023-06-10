package handlers

import (
	"net/http"

	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/webserver/responses"

	"github.com/owncast/owncast/core/rtmp"
)

// DisconnectInboundConnection will force-disconnect an inbound stream.
func (h *Handlers) DisconnectInboundConnection(w http.ResponseWriter, r *http.Request) {
	if !core.GetStatus().Online {
		responses.WriteSimpleResponse(w, false, "no inbound stream connected")
		return
	}

	rtmp.Disconnect()
	responses.WriteSimpleResponse(w, true, "inbound stream disconnected")
}
