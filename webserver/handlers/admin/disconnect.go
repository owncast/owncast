package admin

import (
	"net/http"

	"github.com/owncast/owncast/core"
	webutils "github.com/owncast/owncast/webserver/utils"

	"github.com/owncast/owncast/core/rtmp"
)

// DisconnectInboundConnection will force-disconnect an inbound stream.
func DisconnectInboundConnection(w http.ResponseWriter, r *http.Request) {
	if !core.GetStatus().Online {
		webutils.WriteSimpleResponse(w, false, "no inbound stream connected")
		return
	}

	rtmp.Disconnect()
	webutils.WriteSimpleResponse(w, true, "inbound stream disconnected")
}
