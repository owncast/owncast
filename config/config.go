package config

import (
	"errors"
	"io/ioutil"

	"github.com/gabek/owncast/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

//Config contains a reference to the configuration
var Config *config
var _default config

type config struct {
	ChatDatabaseFilePath string          `yaml:"chatDatabaseFile"`
	DisableWebFeatures   bool            `yaml:"disableWebFeatures"`
	EnableDebugFeatures  bool            `yaml:"-"`
	FFMpegPath           string          `yaml:"ffmpegPath"`
	Files                files           `yaml:"files"`
	InstanceDetails      InstanceDetails `yaml:"instanceDetails"`
	PrivateHLSPath       string          `yaml:"privateHLSPath"`
	PublicHLSPath        string          `yaml:"publicHLSPath"`
	S3                   s3              `yaml:"s3"`
	VersionInfo          string          `yaml:"-"`
	VideoSettings        videoSettings   `yaml:"videoSettings"`
	WebServerPort        int             `yaml:"webServerPort"`
	UsePeer2PeerHLS      bool            `yaml:"usePeer2PeerHLS"`
}

// InstanceDetails defines the user-visible information about this particular instance.
type InstanceDetails struct {
	Name            string            `yaml:"name" json:"name"`
	Title           string            `yaml:"title" json:"title"`
	Summary         string            `yaml:"summary" json:"summary"`
	Logo            map[string]string `yaml:"logo" json:"logo"`
	Tags            []string          `yaml:"tags" json:"tags"`
	SocialHandles   []socialHandle    `yaml:"socialHandles" json:"socialHandles"`
	ExtraInfoFile   string            `yaml:"extraUserInfoFileName" json:"extraUserInfoFileName"`
	Version         string            `json:"version"`
	PrivateHLSPath  string            `json:"privateHLSPath"`
	UsePeer2PeerHLS bool              `json:"usePeer2PeerHLS"`
}

type socialHandle struct {
	Platform string `yaml:"platform" json:"platform"`
	URL      string `yaml:"url" json:"url"`
}

type videoSettings struct {
	ChunkLengthInSeconds      int             `yaml:"chunkLengthInSeconds"`
	StreamingKey              string          `yaml:"streamingKey"`
	StreamQualities           []StreamQuality `yaml:"streamQualities"`
	OfflineContent            string          `yaml:"offlineContent"`
	HighestQualityStreamIndex int             `yaml"-"`
}

// StreamQuality defines the specifics of a single HLS stream variant.
type StreamQuality struct {
	// Enable passthrough to copy the video and/or audio directly from the
	// incoming stream and disable any transcoding.  It will ignore any of
	// the below settings.
	IsVideoPassthrough bool `yaml:"videoPassthrough"`
	IsAudioPassthrough bool `yaml:"audioPassthrough"`

	VideoBitrate int `yaml:"videoBitrate"`
	AudioBitrate int `yaml:"audioBitrate"`

	// Set only one of these in order to keep your current aspect ratio.
	// Or set neither to not scale the video.
	ScaledWidth  int `yaml:"scaledWidth"`
	ScaledHeight int `yaml:"scaledHeight"`

	Framerate     int    `yaml:"framerate"`
	EncoderPreset string `yaml:"encoderPreset"`
}

type files struct {
	MaxNumberInPlaylist int `yaml:"maxNumberInPlaylist"`
}

//s3 is for configuring the s3 integration
type s3 struct {
	Enabled         bool   `yaml:"enabled"`
	Endpoint        string `yaml:"endpoint"`
	ServingEndpoint string `yaml:"servingEndpoint"`
	AccessKey       string `yaml:"accessKey"`
	Secret          string `yaml:"secret"`
	Bucket          string `yaml:"bucket"`
	Region          string `yaml:"region"`
	ACL             string `yaml:"acl"`
}

func (c *config) load(filePath string) error {
	if !utils.DoesFileExists(filePath) {
		log.Fatal("ERROR: valid config.yaml is required.  Copy config-example.yaml to config.yaml and edit")
	}

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return err
	}

	if err := yaml.Unmarshal(yamlFile, c); err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return err
	}

	c.VideoSettings.HighestQualityStreamIndex = findHighestQuality(c.VideoSettings.StreamQualities)

	return nil
}

func (c *config) verifySettings() error {
	if c.VideoSettings.StreamingKey == "" {
		return errors.New("No stream key set. Please set one in your config file.")
	}

	if c.S3.Enabled {
		if c.S3.AccessKey == "" || c.S3.Secret == "" {
			return errors.New("s3 support requires an access key and secret")
		}

		if c.S3.Region == "" || c.S3.Endpoint == "" {
			return errors.New("s3 support requires a region and endpoint")
		}

		if c.S3.Bucket == "" {
			return errors.New("s3 support requires a bucket created for storing public video segments")
		}
	}

	return nil
}

func (c *config) GetUsePeer2PeerHLS() bool {
	if c.UsePeer2PeerHLS != nil {
		return c.UsePeer2PeerHLS
	}

	return _default.GetUsePeer2PeerHLS()
}

func (c *config) GetVideoSegmentSecondsLength() int {
	if c.VideoSettings.ChunkLengthInSeconds != 0 {
		return c.VideoSettings.ChunkLengthInSeconds
	}

	return _default.GetVideoSegmentSecondsLength()
}

func (c *config) GetPublicHLSSavePath() string {
	if c.PublicHLSPath != "" {
		return c.PublicHLSPath
	}

	return _default.PublicHLSPath
}

func (c *config) GetPrivateHLSSavePath() string {
	if c.PrivateHLSPath != "" {
		return c.PrivateHLSPath
	}

	return _default.PrivateHLSPath
}

func (c *config) GetPublicWebServerPort() int {
	if c.WebServerPort != 0 {
		return c.WebServerPort
	}

	return _default.WebServerPort
}

func (c *config) GetMaxNumberOfReferencedSegmentsInPlaylist() int {
	if c.Files.MaxNumberInPlaylist > 0 {
		return c.Files.MaxNumberInPlaylist
	}

	return _default.GetMaxNumberOfReferencedSegmentsInPlaylist()
}

func (c *config) GetOfflineContentPath() string {
	if c.VideoSettings.OfflineContent != "" {
		return c.VideoSettings.OfflineContent
	}

	// This is relative to the webroot, not the project root.
	return _default.VideoSettings.OfflineContent
}

func (c *config) GetFFMpegPath() string {
	if c.FFMpegPath != "" {
		return c.FFMpegPath
	}

	return _default.FFMpegPath
}

func (c *config) GetVideoStreamQualities() []StreamQuality {
	if len(c.VideoSettings.StreamQualities) > 0 {
		return c.VideoSettings.StreamQualities
	}

	return _default.VideoSettings.StreamQualities
}

// GetFramerate returns the framerate or default
func (q *StreamQuality) GetFramerate() int {
	if q.Framerate > 0 {
		return q.Framerate
	}

	return _default.VideoSettings.StreamQualities[0].Framerate
}

//Load tries to load the configuration file
func Load(filePath string, versionInfo string) error {
	Config = new(config)
	_default = getDefaults()

	if err := Config.load(filePath); err != nil {
		return err
	}

	Config.VersionInfo = versionInfo

	// Defaults

	// This is relative to the webroot, not the project root.
	// Has to be set here instead of pulled from a getter
	// since it's serialized to JSON.
	if Config.InstanceDetails.ExtraInfoFile == "" {
		Config.InstanceDetails.ExtraInfoFile = _default.InstanceDetails.ExtraInfoFile
	}

	return Config.verifySettings()
}
