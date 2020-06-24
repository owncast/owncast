package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/gabek/owncast/utils"
)

//Config contains a reference to the configuration
var Config *config

type config struct {
	IPFS           ipfs          `yaml:"ipfs"`
	PublicHLSPath  string        `yaml:"publicHLSPath"`
	PrivateHLSPath string        `yaml:"privateHLSPath"`
	VideoSettings  videoSettings `yaml:"videoSettings"`
	Files          files         `yaml:"files"`
	FFMpegPath     string        `yaml:"ffmpegPath"`
	WebServerPort  int           `yaml:"webServerPort"`
	S3             s3            `yaml:"s3"`
}

type videoSettings struct {
	ChunkLengthInSeconds int             `yaml:"chunkLengthInSeconds"`
	StreamingKey         string          `yaml:"streamingKey"`
	StreamQualities      []StreamQuality `yaml:"streamQualities"`
	OfflineContent       string          `yaml:"offlineContent"`
	EnablePassthrough    bool            `yaml:"passthrough"`
}

type StreamQuality struct {
	// Enable passthrough to copy the video and/or audio directly from the
	// incoming stream and disable any transcoding.  It will ignore any of
	// the below settings.
	IsVideoPassthrough bool `yaml:"videoPassthrough"`
	IsAudioPassthrough bool `yaml:"audioPassthrough"`

	Bitrate      int `yaml:"bitrate"`
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
		log.Fatal("ERROR: valid config/config.yaml is required.  Copy config-example.yaml to config.yaml and edit")
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

	// if !fileExists(config.PrivateHLSPath) {
	// 	os.MkdirAll(path.Join(config.PrivateHLSPath, strconv.Itoa(0)), 0777)
	// }

	// if !fileExists(config.PublicHLSPath) {
	// 	os.MkdirAll(path.Join(config.PublicHLSPath, strconv.Itoa(0)), 0777)
	// }

	if !utils.DoesFileExists(c.FFMpegPath) {
		return fmt.Errorf("ffmpeg does not exist at: %s", c.FFMpegPath)
	}

	return nil
}

//Load tries to load the configuration file
func Load(filePath string) error {
	Config = new(config)

	if err := Config.load(filePath); err != nil {
		return err
	}

	return Config.verifySettings()
}
