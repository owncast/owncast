package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

// Config struct
type Config struct {
	IPFS               IPFS          `yaml:"ipfs"`
	PublicHLSPath      string        `yaml:"publicHLSPath"`
	PrivateHLSPath     string        `yaml:"privateHLSPath"`
	VideoSettings      VideoSettings `yaml:"videoSettings"`
	Files              Files         `yaml:"files"`
	FFMpegPath         string        `yaml:"ffmpegPath"`
	WebServerPort      int           `yaml:"webServerPort"`
	S3                 S3            `yaml:"s3"`
	EnableOfflineImage bool          `yaml:"enableOfflineImage"`
}

type VideoSettings struct {
	ChunkLengthInSeconds int             `yaml:"chunkLengthInSeconds"`
	StreamingKey         string          `yaml:"streamingKey"`
	EncoderPreset        string          `yaml:"encoderPreset"`
	StreamQualities      []StreamQuality `yaml:"streamQualities"`
	EnablePassthrough    bool            `yaml:"passthrough"`
	OfflineImage         string          `yaml:"offlineImage"`
}

type StreamQuality struct {
	Bitrate int `yaml:"bitrate"`
}

// MaxNumberOnDisk must be at least as large as MaxNumberInPlaylist
type Files struct {
	MaxNumberInPlaylist int `yaml:"maxNumberInPlaylist"`
}

type IPFS struct {
	Enabled bool   `yaml:"enabled"`
	Gateway string `yaml:"gateway"`
}

type S3 struct {
	Enabled   bool   `yaml:"enabled"`
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"accessKey"`
	Secret    string `yaml:"secret"`
	Bucket    string `yaml:"bucket"`
	Region    string `yaml:"region"`
}

func getConfig() Config {
	filePath := "config/config.yaml"

	if !fileExists(filePath) {
		log.Fatal("ERROR: valid config/config.yaml is required.  Copy config/config-example.yaml to config/config.yaml and edit.")
	}

	yamlFile, err := ioutil.ReadFile(filePath)

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%+v\n", config)

	return config
}

func checkConfig(config Config) {
	if config.S3.Enabled && config.IPFS.Enabled {
		panic("S3 and IPFS support cannot be enabled at the same time.  Choose one.")
	}

	if config.S3.Enabled {
		if config.S3.AccessKey == "" || config.S3.Secret == "" {
			panic("S3 support requires an access key and secret.")
		}

		if config.S3.Region == "" || config.S3.Endpoint == "" {
			panic("S3 support requires a region and endpoint.")
		}

		if config.S3.Bucket == "" {
			panic("S3 support requires a bucket created for storing public video segments.")
		}
	}

	if !fileExists(config.PrivateHLSPath) {
		os.MkdirAll(path.Join(config.PrivateHLSPath, strconv.Itoa(0)), 0777)
	}

	if !fileExists(config.PublicHLSPath) {
		os.MkdirAll(path.Join(config.PublicHLSPath, strconv.Itoa(0)), 0777)
	}

	if !fileExists(config.FFMpegPath) {
		panic(fmt.Sprintf("ffmpeg does not exist at %s.", config.FFMpegPath))
	}

	if config.VideoSettings.EncoderPreset == "" {
		panic("A video encoder preset is required to be set in the config file.")
	}

}
