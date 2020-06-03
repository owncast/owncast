package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config struct
type Config struct {
	IPFS           IPFS          `yaml:"ipfs"`
	PublicHLSPath  string        `yaml:"publicHLSPath"`
	PrivateHLSPath string        `yaml:"privateHLSPath"`
	VideoSettings  VideoSettings `yaml:"videoSettings"`
	Files          Files         `yaml:"files"`
	FFMpegPath     string        `yaml:"ffmpegPath"`
	WebServerPort  int           `yaml:"webServerPort"`
	S3             S3            `yaml:"s3"`
}

type VideoSettings struct {
	ResolutionWidth      int    `yaml:"resolutionWidth"`
	ChunkLengthInSeconds int    `yaml:"chunkLengthInSeconds"`
	StreamingKey         string `yaml:"streamingKey"`
}

// MaxNumberOnDisk must be at least as large as MaxNumberInPlaylist
type Files struct {
	MaxNumberInPlaylist int `yaml:"maxNumberInPlaylist"`
	MaxNumberOnDisk     int `yaml:"maxNumberOnDisk"`
}

type IPFS struct {
	Enabled bool   `yaml:"enabled"`
	Gateway string `yaml:"gateway"`
}

type S3 struct {
	Enabled   bool   `yaml:"enabled"`
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

	checkConfig(config)

	// fmt.Printf("%+v\n", config)

	return config
}

func checkConfig(config Config) {
	if config.S3.Enabled && config.IPFS.Enabled {
		panic("S3 and IPFS support cannot be enabled at the same time.  Choose one.")
	}

	if !fileExists(config.PrivateHLSPath) {
		panic(fmt.Sprintf("%s does not exist.", config.PrivateHLSPath))
	}

	if !fileExists(config.PublicHLSPath) {
		panic(fmt.Sprintf("%s does not exist.", config.PublicHLSPath))
	}

	if !fileExists(config.FFMpegPath) {
		panic(fmt.Sprintf("ffmpeg does not exist at %s.", config.FFMpegPath))
	}
}
