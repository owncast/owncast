package admin

import (
	"net/http"

	"github.com/owncast/owncast/core/rtmp"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// DisconnectInboundConnection will force-disconnect an inbound stream.
func (a *Admin) DisconnectInboundConnection(w http.ResponseWriter, r *http.Request) {
	if !a.stream.GetStatus().Online {
		webutils.WriteSimpleResponse(w, false, "no inbound stream connected")
		return
	}

	rtmp.Disconnect()
	webutils.WriteSimpleResponse(w, true, "inbound stream disconnected")
}
