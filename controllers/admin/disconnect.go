package admin

import (
	"net/http"
)

// DisconnectInboundConnection will force-disconnect an inbound stream.
func (c *Controller) DisconnectInboundConnection(w http.ResponseWriter, r *http.Request) {
	if !c.Service.Core.GetStatus().Online {
		c.WriteSimpleResponse(w, false, "no inbound stream connected")
		return
	}

	c.Rtmp.Disconnect()
	c.WriteSimpleResponse(w, true, "inbound stream disconnected")
}
