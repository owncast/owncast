package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	log "github.com/sirupsen/logrus"
)

// GetServerConfig gets the config details of the server.
func GetServerConfig(w http.ResponseWriter, r *http.Request) {
	var videoQualityVariants = make([]models.StreamOutputVariant, 0)
	for _, variant := range data.GetStreamOutputVariants() {
		videoQualityVariants = append(videoQualityVariants, models.StreamOutputVariant{
			IsAudioPassthrough: variant.GetIsAudioPassthrough(),
			IsVideoPassthrough: variant.IsVideoPassthrough,
			Framerate:          variant.GetFramerate(),
			EncoderPreset:      variant.GetEncoderPreset(),
			VideoBitrate:       variant.VideoBitrate,
			AudioBitrate:       variant.AudioBitrate,
		})
	}
	response := serverConfigAdminResponse{
		InstanceDetails: webConfigResponse{
			Title:            data.GetServerTitle(),
			Name:             data.GetServerName(),
			Summary:          data.GetServerSummary(),
			Tags:             data.GetServerMetadataTags(),
			ExtraPageContent: data.GetExtraPageBodyContent(),
			StreamTitle:      data.GetStreamTitle(),
			Logo:             data.GetLogoPath(),
			SocialHandles:    data.GetSocialHandles(),
			NSFW:             data.GetNSFW(),
		},
		FFmpegPath:     data.GetFfMpegPath(),
		StreamKey:      data.GetStreamKey(),
		WebServerPort:  config.WebServerPort,
		RTMPServerPort: data.GetRTMPPortNumber(),
		VideoSettings: videoSettings{
			VideoQualityVariants: videoQualityVariants,
			LatencyLevel:         data.GetStreamLatencyLevel().Level,
		},
		YP: yp{
			Enabled:     data.GetDirectoryEnabled(),
			InstanceURL: data.GetServerURL(),
		},
		S3: data.GetS3Config(),
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Errorln(err)
	}
}

type serverConfigAdminResponse struct {
	InstanceDetails webConfigResponse `json:"instanceDetails"`
	FFmpegPath      string            `json:"ffmpegPath"`
	StreamKey       string            `json:"streamKey"`
	WebServerPort   int               `json:"webServerPort"`
	RTMPServerPort  int               `json:"rtmpServerPort"`
	S3              models.S3         `json:"s3"`
	VideoSettings   videoSettings     `json:"videoSettings"`
	LatencyLevel    int               `json:"latencyLevel"`
	YP              yp                `json:"yp"`
}

type videoSettings struct {
	VideoQualityVariants []models.StreamOutputVariant `json:"videoQualityVariants"`
	LatencyLevel         int                          `json:"latencyLevel"`
}

type webConfigResponse struct {
	Name             string                `json:"name"`
	Title            string                `json:"title"`
	Summary          string                `json:"summary"`
	Logo             string                `json:"logo"`
	Tags             []string              `json:"tags"`
	Version          string                `json:"version"`
	NSFW             bool                  `json:"nsfw"`
	ExtraPageContent string                `json:"extraPageContent"`
	StreamTitle      string                `json:"streamTitle"` // What's going on with the current stream
	SocialHandles    []models.SocialHandle `json:"socialHandles"`
}

type yp struct {
	Enabled      bool   `json:"enabled"`
	InstanceURL  string `json:"instanceUrl"` // The public URL the directory should link to
	YPServiceURL string `json:"-"`           // The base URL to the YP API to register with (optional)
}
