package config

import (
	"fmt"
	"path/filepath"
	"time"
)

// Config holds runtime-mutable configuration for a single execution of the
// Owncast server. It is populated by the composition root (main.go) from
// CLI flags + persisted settings, then threaded by value-of-pointer into
// every service that reads any of these fields.
//
// All fields are runtime state — values that change between executions or
// are overridden at startup. Compile-time constants (see constants.go),
// build-time-set vars (VersionNumber, GitCommit, BuildPlatform), and the
// build-tag-driven EnableAutoUpdate remain package-level on purpose: they
// are not runtime state.
type Config struct {
	// DatabaseFilePath is the path to the file used as the global database
	// for this run of the application.
	DatabaseFilePath string

	// LogDirectory is the path to various log files.
	LogDirectory string

	// TempDir is where we store temporary files.
	TempDir string

	// BackupDirectory is the directory we write backup files to.
	BackupDirectory string

	// WebServerPort is the port the public webserver listens on.
	WebServerPort int

	// WebServerIP is the IP address to bind the web server to. All
	// interfaces by default.
	WebServerIP string

	// InternalHLSListenerPort is the port the in-process HLS receiver
	// binds to. Discovered at runtime when the receiver opens its
	// listener and read by the transcoder when it spins up.
	InternalHLSListenerPort string

	// InternalHLSListenerHost is the interface the in-process HLS receiver
	// binds to, and the host the transcoder's ffmpeg PUTs finished segments
	// to. Defaults to 127.0.0.1, where receiver and transcoder share one
	// box. A remote-engine deployment that PUTs HLS across the network sets
	// this to a reachable interface.
	InternalHLSListenerHost string

	// EnableDebugFeatures prints additional data to help in debugging.
	EnableDebugFeatures bool

	// TemporaryStreamKey is a stream key set via the command line that
	// overrides the persisted keys for the duration of this process.
	TemporaryStreamKey string

	// FollowerValidationInterval is how often the follower validation
	// job runs. A zero value means use the package default.
	FollowerValidationInterval time.Duration
}

// NewDefault returns a *Config populated with the default startup values.
// main.go overlays CLI flag values on top of this before injecting it
// into services.
func NewDefault() *Config {
	return &Config{
		DatabaseFilePath:        "data/owncast.db",
		LogDirectory:            "./data/logs",
		TempDir:                 "./data/tmp",
		BackupDirectory:         filepath.Join(DataDirectory, "backup"),
		WebServerPort:           8080,
		WebServerIP:             "0.0.0.0",
		InternalHLSListenerPort: "8927",
		InternalHLSListenerHost: "127.0.0.1",
		EnableDebugFeatures:     false,
		TemporaryStreamKey:      "",
		// FollowerValidationInterval defaults to 0; consumers use their
		// package-level default when zero.
	}
}

// VersionNumber is the current version string. Set at build time via
// linker flags; defaults to StaticVersionNumber for local builds.
var VersionNumber = StaticVersionNumber

// GitCommit is an optional commit this build was made from. Set at
// build time via linker flags.
var GitCommit = ""

// BuildPlatform is the optional platform this release was built for.
// Set at build time via linker flags; "dev" for local builds.
var BuildPlatform = "dev"

// EnableAutoUpdate will explicitly enable in-place auto-updates via the
// admin. Toggled by the enable_updates build tag (see
// updaterConfig_enabled.go); package-level because the build tag, not
// the runtime, owns its value.
var EnableAutoUpdate = false

// GetCommit will return an identifier used for identifying the point in time this build took place.
func GetCommit() string {
	if GitCommit == "" {
		GitCommit = time.Now().Format("20060102")
	}

	return GitCommit
}

// DefaultForbiddenUsernames are a list of usernames forbidden from being used in chat.
var DefaultForbiddenUsernames = []string{
	"owncast", "operator", "admin", "system",
}

// MaxSocketPayloadSize is the maximum payload we will allow to to be received via the chat socket.
const MaxSocketPayloadSize = 2048

// GetReleaseString gets the version string.
func GetReleaseString() string {
	versionNumber := VersionNumber
	buildPlatform := BuildPlatform
	gitCommit := GetCommit()

	return fmt.Sprintf("Owncast v%s-%s (%s)", versionNumber, buildPlatform, gitCommit)
}
