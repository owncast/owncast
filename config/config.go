package config

import (
	"io/ioutil"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Config contains a reference to the configuration.
var _default config

var DatabaseFilePath = "data/owncast.db"
var EnableDebugFeatures = false
var VersionInfo = "v" + CurrentBuildString
var VersionNumber = CurrentBuildString
var WebServerPort = 8080
var RTMPServerPort = 1935
var HighestQualityStreamIndex = 0

type config struct {
	DatabaseFilePath     string `yaml:"databaseFile"`
	EnableDebugFeatures  bool   `yaml:"-"`
	FFMpegPath           string
	Files                files           `yaml:"files"`
	InstanceDetails      InstanceDetails `yaml:"instanceDetails"`
	VersionInfo          string          `yaml:"-"` // For storing the version/build number
	VersionNumber        string          `yaml:"-"`
	VideoSettings        videoSettings   `yaml:"videoSettings"`
	WebServerPort        int
	RTMPServerPort       int
	YP                   YP `yaml:"yp"`
	DisableUpgradeChecks bool
}

// InstanceDetails defines the user-visible information about this particular instance.
type InstanceDetails struct {
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

type videoSettings struct {
	ChunkLengthInSeconds int                          `yaml:"chunkLengthInSeconds"`
	StreamingKey         string                       `yaml:"streamingKey"`
	StreamQualities      []models.StreamOutputVariant `yaml:"streamQualities"`
}

// YP allows registration to the central Owncast YP (Yellow pages) service operating as a directory.
type YP struct {
	Enabled      bool   `json:"enabled"`
	InstanceURL  string `json:"instanceUrl"` // The public URL the directory should link to
	YPServiceURL string `json:"-"`           // The base URL to the YP API to register with (optional)
}

// StreamQuality defines the specifics of a single HLS stream variant.
type StreamQuality struct {
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

// S3 is for configuring the S3 integration.
type S3 struct {
	Enabled         bool   `yaml:"enabled" json:"enabled"`
	Endpoint        string `yaml:"endpoint" json:"endpoint,omitempty"`
	ServingEndpoint string `yaml:"servingEndpoint" json:"servingEndpoint,omitempty"`
	AccessKey       string `yaml:"accessKey" json:"accessKey,omitempty"`
	Secret          string `yaml:"secret" json:"secret,omitempty"`
	Bucket          string `yaml:"bucket" json:"bucket,omitempty"`
	Region          string `yaml:"region" json:"region,omitempty"`
	ACL             string `yaml:"acl" json:"acl,omitempty"`
}

func (c *config) load(filePath string) error {
	if !utils.DoesFileExists(filePath) {
		return nil
	}

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return err
	}

	if err := yaml.Unmarshal(yamlFile, c); err != nil {
		log.Fatalf("Error reading the config file.\nHave you recently updated your version of Owncast?\nIf so there may be changes to the config.\nPlease read the change log for your version at https://owncast.online/posts/\n%v", err)
		return err
	}

	// c.VideoSettings.HighestQualityStreamIndex = findHighestQuality(c.VideoSettings.StreamQualities)

	return nil
}

func (c *config) GetVideoSegmentSecondsLength() int {
	if c.VideoSettings.ChunkLengthInSeconds != 0 {
		return c.VideoSettings.ChunkLengthInSeconds
	}

	return _default.GetVideoSegmentSecondsLength()
}

func (c *config) GetMaxNumberOfReferencedSegmentsInPlaylist() int {
	if c.Files.MaxNumberInPlaylist > 0 {
		return c.Files.MaxNumberInPlaylist
	}

	return _default.GetMaxNumberOfReferencedSegmentsInPlaylist()
}

func (c *config) GetYPServiceHost() string {
	if c.YP.YPServiceURL != "" {
		return c.YP.YPServiceURL
	}

	return _default.YP.YPServiceURL
}

func (c *config) GetDataFilePath() string {
	if c.DatabaseFilePath != "" {
		return c.DatabaseFilePath
	}

	return _default.DatabaseFilePath
}

// Load tries to load the configuration file.
// func Load(filePath string, versionInfo string, versionNumber string) error {
// 	Config = new(config)
// 	_default = GetDefaults()

// 	if err := Config.load(filePath); err != nil {
// 		return err
// 	}

// 	Config.VersionInfo = versionInfo
// 	Config.VersionNumber = versionNumber
// 	return Config.verifySettings()
// }
