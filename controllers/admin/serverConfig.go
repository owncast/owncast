package admin

import (
	"encoding/json"
	"net/http"

	"github.com/gabek/owncast/config"
)

// GetServerConfig gets the config details of the server
func GetServerConfig(w http.ResponseWriter, r *http.Request) {
	response := serverConfigAdminResponse{
		InstanceDetails: config.Config.InstanceDetails,
		FFmpegPath:      config.Config.GetFFMpegPath(),
		WebServerPort:   config.Config.GetPublicWebServerPort(),
		VideoSettings: videoSettings{
			VideoQualityVariants:  config.Config.GetVideoStreamQualities(),
			SegmentLengthSeconds:  config.Config.GetVideoSegmentSecondsLength(),
			NumberOfPlaylistItems: config.Config.GetMaxNumberOfReferencedSegmentsInPlaylist(),
		},
		S3: config.Config.S3,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type serverConfigAdminResponse struct {
	InstanceDetails config.InstanceDetails `json:"instanceDetails"`
	FFmpegPath      string                 `json:"ffmpegPath"`
	WebServerPort   int                    `json:"webServerPort"`
	S3              config.S3              `json:"s3"`
	VideoSettings   videoSettings          `json:"videoSettings"`
}

type videoSettings struct {
	VideoQualityVariants  []config.StreamQuality `json:"videoQualityVariants"`
	SegmentLengthSeconds  int                    `json:"segmentLengthSeconds"`
	NumberOfPlaylistItems int                    `json:"numberOfPlaylistItems"`
}
