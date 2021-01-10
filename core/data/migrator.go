package data

import (
	"io/ioutil"
	"os"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func MigrateConfigFile() {
	filePath := "config.yaml"

	if !utils.DoesFileExists(filePath) {
		return
	}

	log.Infoln("Migrating", filePath, "to new datastore")

	var oldConfig configFile

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorln("config file err   #%v ", err)
		return
	}

	if err := yaml.Unmarshal(yamlFile, &oldConfig); err != nil {
		log.Errorln("Error reading the config file.", err)
		return
	}

	SetServerTitle(oldConfig.InstanceDetails.Title)
	SetServerName(oldConfig.InstanceDetails.Name)
	SetServerSummary(oldConfig.InstanceDetails.Summary)
	SetServerMetadataTags(oldConfig.InstanceDetails.Tags)
	SetStreamKey(oldConfig.VideoSettings.StreamingKey)
	SetLogoPath(oldConfig.InstanceDetails.Logo)
	SetNSFW(oldConfig.InstanceDetails.NSFW)
	SetServerURL(oldConfig.YP.InstanceURL)
	SetDirectoryEnabled(oldConfig.YP.Enabled)
	SetSocialHandles(oldConfig.InstanceDetails.SocialHandles)
	SetFfmpegPath(oldConfig.FFMpegPath)
	SetHTTPPortNumber(oldConfig.WebServerPort)
	SetRTMPPortNumber(oldConfig.RTMPServerPort)
	SetDisableUpgradeChecks(oldConfig.DisableUpgradeChecks)

	// Migrate the old content.md file
	content, err := ioutil.ReadFile(config.ExtraInfoFile)
	if err == nil && len(content) > 0 {
		SetExtraPageBodyContent(string(content))
	}

	if !utils.DoesFileExists(config.BackupDirectory) {
		if err := os.Mkdir(config.BackupDirectory, 0700); err != nil {
			log.Errorln("Unable to create backup directory", err)
			return
		}
	}

	if err := utils.Move(filePath, "backup/config.old"); err != nil {
		log.Warnln(err)
	}
}

type configFile struct {
	DatabaseFilePath     string `yaml:"databaseFile"`
	EnableDebugFeatures  bool   `yaml:"-"`
	FFMpegPath           string
	Files                files           `yaml:"files"`
	InstanceDetails      instanceDetails `yaml:"instanceDetails"`
	VersionInfo          string          `yaml:"-"` // For storing the version/build number
	VersionNumber        string          `yaml:"-"`
	VideoSettings        videoSettings   `yaml:"videoSettings"`
	WebServerPort        int
	RTMPServerPort       int
	YP                   yp `yaml:"yp"`
	DisableUpgradeChecks bool
}

// instanceDetails defines the user-visible information about this particular instance.
type instanceDetails struct {
	Name             string                `yaml:"name"`
	Title            string                `yaml:"title"`
	Summary          string                `yaml:"summary"`
	Logo             string                `yaml:"logo"`
	Tags             []string              `yaml:"tags"`
	Version          string                `yaml:"version"`
	NSFW             bool                  `yaml:"nsfw"`
	ExtraPageContent string                `yaml:"extraPageContent"`
	StreamTitle      string                `yaml:"streamTitle"`
	SocialHandles    []models.SocialHandle `yaml:"socialHandles"`
}

type videoSettings struct {
	ChunkLengthInSeconds      int             `yaml:"chunkLengthInSeconds"`
	StreamingKey              string          `yaml:"streamingKey"`
	StreamQualities           []streamQuality `yaml:"streamQualities"`
	HighestQualityStreamIndex int             `yaml:"-"`
}

// yp allows registration to the central Owncast yp (Yellow pages) service operating as a directory.
type yp struct {
	Enabled     bool   `yaml:"enabled"`
	InstanceURL string `yaml:"instanceUrl"` // The public URL the directory should link to
}

// streamQuality defines the specifics of a single HLS stream variant.
type streamQuality struct {
	// Enable passthrough to copy the video and/or audio directly from the
	// incoming stream and disable any transcoding.  It will ignore any of
	// the below settings.
	IsVideoPassthrough bool `yaml:"videoPassthrough" json:"videoPassthrough"`
	IsAudioPassthrough bool `yaml:"audioPassthrough" json:"audioPassthrough"`

	VideoBitrate int `yaml:"videoBitrate" json:"videoBitrate"`
	AudioBitrate int `yaml:"audioBitrate" json:"audioBitrate"`

	// Set only one of these in order to keep your current aspect ratio.
	// Or set neither to not scale the video.
	ScaledWidth  int `yaml:"scaledWidth" json:"scaledWidth,omitempty"`
	ScaledHeight int `yaml:"scaledHeight" json:"scaledHeight,omitempty"`

	Framerate     int    `yaml:"framerate" json:"framerate"`
	EncoderPreset string `yaml:"encoderPreset" json:"encoderPreset"`
}

type files struct {
	MaxNumberInPlaylist int `yaml:"maxNumberInPlaylist"`
}

// s3 is for configuring the s3 integration.
type s3 struct {
	Enabled         bool   `yaml:"enabled" json:"enabled"`
	Endpoint        string `yaml:"endpoint" json:"endpoint,omitempty"`
	ServingEndpoint string `yaml:"servingEndpoint" json:"servingEndpoint,omitempty"`
	AccessKey       string `yaml:"accessKey" json:"accessKey,omitempty"`
	Secret          string `yaml:"secret" json:"secret,omitempty"`
	Bucket          string `yaml:"bucket" json:"bucket,omitempty"`
	Region          string `yaml:"region" json:"region,omitempty"`
	ACL             string `yaml:"acl" json:"acl,omitempty"`
}
