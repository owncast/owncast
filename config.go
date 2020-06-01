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
}

type VideoSettings struct {
	ResolutionWidth      int `yaml:"resolutionWidth"`
	ChunkLengthInSeconds int `yaml:"chunkLengthInSeconds"`
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

func getConfig() Config {
	filePath := "config/config.yaml"

	if !fileExists(filePath) {
		log.Fatal("ERROR: valid config/config.yaml is required")
	}

	yamlFile, err := ioutil.ReadFile(filePath)

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func checkConfig(config Config) {
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
