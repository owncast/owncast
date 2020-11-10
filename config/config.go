package config

import (
	"errors"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

//Config contains a reference to the configuration
var Config *config
var _default config

type config struct {
	DatabaseFilePath     string          `yaml:"databaseFile"`
	EnableDebugFeatures  bool            `yaml:"-"`
	FFMpegPath           string          `yaml:"ffmpegPath"`
	Files                files           `yaml:"files"`
	InstanceDetails      InstanceDetails `yaml:"instanceDetails"`
	S3                   S3              `yaml:"s3"`
	VersionInfo          string          `yaml:"-"` // For storing the version/build number
	VersionNumber        string          `yaml:"-"`
	VideoSettings        videoSettings   `yaml:"videoSettings"`
	WebServerPort        int             `yaml:"webServerPort"`
	DisableUpgradeChecks bool            `yaml:"disableUpgradeChecks"`
	YP                   YP              `yaml:"yp"`
}

// InstanceDetails defines the user-visible information about this particular instance.
type InstanceDetails struct {
	Name             string         `yaml:"name" json:"name"`
	Title            string         `yaml:"title" json:"title"`
	Summary          string         `yaml:"summary" json:"summary"`
	Logo             logo           `yaml:"logo" json:"logo"`
	Tags             []string       `yaml:"tags" json:"tags"`
	SocialHandles    []socialHandle `yaml:"socialHandles" json:"socialHandles"`
	Version          string         `json:"version"`
	NSFW             bool           `yaml:"nsfw" json:"nsfw"`
	ExtraPageContent string         `json:"extraPageContent"`
}

type logo struct {
	Large string `yaml:"large" json:"large"`
	Small string `yaml:"small" json:"small"`
}

type socialHandle struct {
	Platform string `yaml:"platform" json:"platform"`
	URL      string `yaml:"url" json:"url"`
}

type videoSettings struct {
	ChunkLengthInSeconds      int             `yaml:"chunkLengthInSeconds"`
	StreamingKey              string          `yaml:"streamingKey"`
	StreamQualities           []StreamQuality `yaml:"streamQualities"`
	HighestQualityStreamIndex int             `yaml:"-"`
}

// YP allows registration to the central Owncast YP (Yellow pages) service operating as a directory.
type YP struct {
	Enabled      bool   `yaml:"enabled" json:"enabled"`
	InstanceURL  string `yaml:"instanceURL" json:"instanceUrl"` // The public URL the directory should link to
	YPServiceURL string `yaml:"ypServiceURL" json:"-"`          // The base URL to the YP API to register with (optional)
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

//S3 is for configuring the S3 integration
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

	// Add custom page content to the instance details.
	customContentMarkdownData, err := ioutil.ReadFile(ExtraInfoFile)
	if err == nil {
		customContentMarkdownString := string(customContentMarkdownData)
		c.InstanceDetails.ExtraPageContent = utils.RenderSimpleMarkdown(customContentMarkdownString)
	}

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

	if c.YP.Enabled && c.YP.InstanceURL == "" {
		return errors.New("YP is enabled but instance url is not set")
	}

	return nil
}

func (c *config) GetVideoSegmentSecondsLength() int {
	if c.VideoSettings.ChunkLengthInSeconds != 0 {
		return c.VideoSettings.ChunkLengthInSeconds
	}

	return _default.GetVideoSegmentSecondsLength()
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

func (c *config) GetFFMpegPath() string {
	if c.FFMpegPath != "" {
		return c.FFMpegPath
	}

	// First look to see if ffmpeg is in the current working directory
	localCopy := "./ffmpeg"
	hasLocalCopyError := verifyFFMpegPath(localCopy)
	if hasLocalCopyError == nil {
		// No error, so all is good.  Use the local copy.
		return localCopy
	}

	cmd := exec.Command("which", "ffmpeg")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Debugln("Unable to determine path to ffmpeg. Please specify it in the config file.")
	}

	path := strings.TrimSpace(string(out))

	return path
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

// GetEncoderPreset returns the preset or default
func (q *StreamQuality) GetEncoderPreset() string {
	if q.EncoderPreset != "" {
		return q.EncoderPreset
	}

	return _default.VideoSettings.StreamQualities[0].EncoderPreset
}

//Load tries to load the configuration file
func Load(filePath string, versionInfo string, versionNumber string) error {
	Config = new(config)
	_default = getDefaults()

	if err := Config.load(filePath); err != nil {
		return err
	}

	Config.VersionInfo = versionInfo
	Config.VersionNumber = versionNumber
	return Config.verifySettings()
}
