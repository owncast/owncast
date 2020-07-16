package config

import (
	"errors"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/gabek/owncast/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

//Config contains a reference to the configuration
var Config *config

type config struct {
	ChatDatabaseFilePath string          `yaml:"chatDatabaseFile"`
	DisableWebFeatures   bool            `yaml:"disableWebFeatures"`
	EnableDebugFeatures  bool            `yaml:"-"`
	FFMpegPath           string          `yaml:"ffmpegPath"`
	Files                files           `yaml:"files"`
	IPFS                 ipfs            `yaml:"ipfs"`
	InstanceDetails      InstanceDetails `yaml:"instanceDetails"`
	PrivateHLSPath       string          `yaml:"privateHLSPath"`
	PublicHLSPath        string          `yaml:"publicHLSPath"`
	S3                   s3              `yaml:"s3"`
	VersionInfo          string          `yaml:"-"`
	VideoSettings        videoSettings   `yaml:"videoSettings"`
	WebServerPort        int             `yaml:"webServerPort"`
}

// InstanceDetails defines the user-visible information about this particular instance.
type InstanceDetails struct {
	Name          string            `yaml:"name" json:"name"`
	Title         string            `yaml:"title" json:"title"`
	Summary       string            `yaml:"summary" json:"summary"`
	Logo          map[string]string `yaml:"logo" json:"logo"`
	Tags          []string          `yaml:"tags" json:"tags"`
	SocialHandles []socialHandle    `yaml:"socialHandles" json:"socialHandles"`
	ExtraInfoFile string            `yaml:"extraUserInfoFileName" json:"extraUserInfoFileName"`
	Version       string            `json:"version"`
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

type ipfs struct {
	Enabled bool   `yaml:"enabled"`
	Gateway string `yaml:"gateway"`
}

//s3 is for configuring the s3 integration
type s3 struct {
	Enabled   bool   `yaml:"enabled"`
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"accessKey"`
	Secret    string `yaml:"secret"`
	Bucket    string `yaml:"bucket"`
	Region    string `yaml:"region"`
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
	if c.S3.Enabled && c.IPFS.Enabled {
		return errors.New("s3 and IPFS support cannot be enabled at the same time; choose one")
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

func (c *config) GetFFMpegPath() string {
	if c.FFMpegPath != "" {
		return c.FFMpegPath
	}

	cmd := exec.Command("which", "ffmpeg")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Panicln("Unable to determine path to ffmpeg. Please specify it in the config file.")
	}

	path := strings.TrimSpace(string(out))

	// Memoize it for future access
	c.FFMpegPath = path

	return path
}

func (c *config) GetVideoSegmentSecondsLength() int {
	if c.VideoSettings.ChunkLengthInSeconds != 0 {
		return c.VideoSettings.ChunkLengthInSeconds
	}

	// Default
	return 4
}

func (c *config) GetPublicHLSSavePath() string {
	if c.PublicHLSPath != "" {
		return c.PublicHLSPath
	}

	return "webroot/hls"
}

func (c *config) GetPrivateHLSSavePath() string {
	if c.PrivateHLSPath != "" {
		return c.PrivateHLSPath
	}

	return "hls"
}

func (c *config) GetPublicWebServerPort() int {
	if c.WebServerPort != 0 {
		return c.WebServerPort
	}

	// Default web server port
	return 8080
}

func (c *config) GetMaxNumberOfReferencedSegmentsInPlaylist() int {
	if c.Files.MaxNumberInPlaylist > 0 {
		return c.Files.MaxNumberInPlaylist
	}

	return 20
}

func (c *config) GetOfflineContentPath() string {
	if c.VideoSettings.OfflineContent != "" {
		return c.VideoSettings.OfflineContent
	}

	// This is relative to the webroot, not the project root.
	return "static/offline.m4v"
}

//Load tries to load the configuration file
func Load(filePath string, versionInfo string) error {
	Config = new(config)

	if err := Config.load(filePath); err != nil {
		return err
	}

	Config.VersionInfo = versionInfo

	// Defaults

	// This is relative to the webroot, not the project root.
	if Config.InstanceDetails.ExtraInfoFile == "" {
		Config.InstanceDetails.ExtraInfoFile = "/static/content.md"
	}

	return Config.verifySettings()
}
