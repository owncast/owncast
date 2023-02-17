package admin

import (
	"net/http"

	"github.com/owncast/owncast/internal/controllers"
	"github.com/owncast/owncast/internal/core"
	"github.com/owncast/owncast/internal/core/rtmp"
)

// DisconnectInboundConnection will force-disconnect an inbound stream.
func DisconnectInboundConnection(w http.ResponseWriter, r *http.Request) {
	if !core.GetStatus().Online {
		controllers.WriteSimpleResponse(w, false, "no inbound stream connected")
		return
	}

	rtmp.Disconnect()
	controllers.WriteSimpleResponse(w, true, "inbound stream disconnected")
}
