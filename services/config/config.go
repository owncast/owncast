package config

import (
	"fmt"
	"path/filepath"
	"time"
)

// These are runtime-set values used for configuration.

type Config struct {
	// DatabaseFilePath is the path to the file ot be used as the global database for this run of the application.
	DatabaseFilePath string
	// LogDirectory is the path to various log files.
	LogDirectory string
	// TempDir is where we store temporary files.
	TempDir string
	// EnableDebugFeatures will print additional data to help in debugging.
	EnableDebugFeatures bool
	// VersionNumber is the current version string.
	VersionNumber string
	// WebServerPort is the port for Owncast's webserver that is used for this execution of the service.
	WebServerPort int
	// WebServerIP is the IP address to bind the web server to. All interfaces by default.
	WebServerIP string
	// InternalHLSListenerPort is the port for HLS writes that is used for this execution of the service.
	InternalHLSListenerPort string
	// GitCommit is an optional commit this build was made from.
	GitCommit string
	// BuildPlatform is the optional platform this release was built for.
	BuildPlatform string
	// EnableAutoUpdate will explicitly enable in-place auto-updates via the admin.
	EnableAutoUpdate bool
	// A temporary stream key that can be set via the command line.
	TemporaryStreamKey string
	// BackupDirectory is the directory we write backup files to.
	BackupDirectory string
	// HLSStoragePath is the directory HLS video is written to.
	HLSStoragePath string
	// CustomEmojiPath is the emoji directory.
	CustomEmojiPath string
	// PublicFilesPath is the optional directory for hosting public files.
	PublicFilesPath string
}

// NewFediAuth creates a new FediAuth instance.
func New() *Config {
	// Default config values.
	c := &Config{
		DatabaseFilePath:        "data/owncast.db",
		LogDirectory:            "./data/logs",
		TempDir:                 "./data/tmp",
		EnableDebugFeatures:     false,
		VersionNumber:           StaticVersionNumber,
		WebServerPort:           8080,
		WebServerIP:             "0.0.0.0",
		InternalHLSListenerPort: "8927",
		GitCommit:               "",
		BuildPlatform:           "dev",
		EnableAutoUpdate:        false,
		TemporaryStreamKey:      "",
		BackupDirectory:         filepath.Join(DataDirectory, "backup"),
		HLSStoragePath:          filepath.Join(DataDirectory, "hls"),
		CustomEmojiPath:         filepath.Join(DataDirectory, "emoji"),
		PublicFilesPath:         filepath.Join(DataDirectory, "public"),
	}
	return c
}

var temporaryGlobalInstance *Config

// GetConfig returns the temporary global instance.
// Remove this after dependency injection is implemented.
func Get() *Config {
	if temporaryGlobalInstance == nil {
		temporaryGlobalInstance = New()
	}

	return temporaryGlobalInstance
}

// GetCommit will return an identifier used for identifying the point in time this build took place.
func (c *Config) GetCommit() string {
	if c.GitCommit == "" {
		c.GitCommit = time.Now().Format("20060102")
	}

	return c.GitCommit
}

// DefaultForbiddenUsernames are a list of usernames forbidden from being used in chat.
var DefaultForbiddenUsernames = []string{
	"owncast", "operator", "admin", "system",
}

// MaxSocketPayloadSize is the maximum payload we will allow to to be received via the chat socket.
const MaxSocketPayloadSize = 2048

// GetReleaseString gets the version string.
func (c *Config) GetReleaseString() string {
	versionNumber := c.VersionNumber
	buildPlatform := c.BuildPlatform
	gitCommit := c.GetCommit()

	return fmt.Sprintf("Owncast v%s-%s (%s)", versionNumber, buildPlatform, gitCommit)
}

// GetTranscoderLogFilePath returns the logging path for the transcoder log output.
func (c *Config) GetTranscoderLogFilePath() string {
	return filepath.Join(c.LogDirectory, "transcoder.log")
}
