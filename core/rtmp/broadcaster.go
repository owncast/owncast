package rtmp

import (
	"time"

	"github.com/nareix/joy5/format/flv/flvio"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

func setCurrentBroadcasterInfo(t flvio.Tag, remoteAddr string) {
	data, err := getInboundDetailsFromMetadata(t.DebugFields())
	if err != nil {
		log.Traceln("RTMP meadata:", err)
	}

	broadcaster := models.Broadcaster{
		RemoteAddr: remoteAddr,
		Time:       time.Now(),
		StreamDetails: models.InboundStreamDetails{
			Width:          data.Width,
			Height:         data.Height,
			VideoBitrate:   int(data.VideoBitrate),
			VideoCodec:     getVideoCodec(data.VideoCodec),
			VideoFramerate: data.VideoFramerate,
			AudioBitrate:   int(data.AudioBitrate),
			AudioCodec:     getAudioCodec(data.AudioCodec),
			Encoder:        data.Encoder,
			VideoOnly:      data.AudioCodec == nil,
		},
	}

	_setBroadcaster(broadcaster)
}
