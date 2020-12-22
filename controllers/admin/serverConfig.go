package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/config"
	log "github.com/sirupsen/logrus"
)

// GetServerConfig gets the config details of the server.
func GetServerConfig(w http.ResponseWriter, r *http.Request) {
	var videoQualityVariants = make([]config.StreamQuality, 0)
	for _, variant := range config.Config.GetVideoStreamQualities() {
		videoQualityVariants = append(videoQualityVariants, config.StreamQuality{
			IsAudioPassthrough: variant.GetIsAudioPassthrough(),
			IsVideoPassthrough: variant.IsVideoPassthrough,
			Framerate:          variant.GetFramerate(),
			EncoderPreset:      variant.GetEncoderPreset(),
			VideoBitrate:       variant.VideoBitrate,
			AudioBitrate:       variant.AudioBitrate,
		})
	}
	response := serverConfigAdminResponse{
		InstanceDetails: config.Config.InstanceDetails,
		FFmpegPath:      config.Config.GetFFMpegPath(),
		StreamKey:       config.Config.VideoSettings.StreamingKey,
		WebServerPort:   config.Config.GetPublicWebServerPort(),
		RTMPServerPort:  config.Config.GetRTMPServerPort(),
		VideoSettings: videoSettings{
			VideoQualityVariants:  videoQualityVariants,
			SegmentLengthSeconds:  config.Config.GetVideoSegmentSecondsLength(),
			NumberOfPlaylistItems: config.Config.GetMaxNumberOfReferencedSegmentsInPlaylist(),
		},
		YP: config.Config.YP,
		S3: config.Config.S3,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorln(err)
	}
}

type serverConfigAdminResponse struct {
	InstanceDetails config.InstanceDetails `json:"instanceDetails"`
	FFmpegPath      string                 `json:"ffmpegPath"`
	StreamKey       string                 `json:"streamKey"`
	WebServerPort   int                    `json:"webServerPort"`
	RTMPServerPort  int                    `json:"rtmpServerPort"`
	S3              config.S3              `json:"s3"`
	VideoSettings   videoSettings          `json:"videoSettings"`
	YP              config.YP              `json:"yp"`
}

type videoSettings struct {
	VideoQualityVariants  []config.StreamQuality `json:"videoQualityVariants"`
	SegmentLengthSeconds  int                    `json:"segmentLengthSeconds"`
	NumberOfPlaylistItems int                    `json:"numberOfPlaylistItems"`
}
