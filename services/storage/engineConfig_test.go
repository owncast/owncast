package storage

import (
	"testing"

	"github.com/owncast/owncast/models"
)

// fakeEngineConfig is a database-free models.EngineConfig. Its existence is the
// point of the EngineConfig refactor: the storage providers (and the rtmp /
// transcoder packages) now depend only on this narrow read-only surface, so
// they can be constructed without a config repository or datastore.
type fakeEngineConfig struct {
	servingEndpoint string
	s3              models.S3
}

func (f fakeEngineConfig) GetRTMPPortNumber() int     { return 1935 }
func (f fakeEngineConfig) GetRTMPBindAddress() string { return "0.0.0.0" }
func (f fakeEngineConfig) GetFfMpegPath() string      { return "ffmpeg" }
func (f fakeEngineConfig) GetVideoCodec() string      { return "libx264" }
func (f fakeEngineConfig) GetS3Config() models.S3     { return f.s3 }
func (f fakeEngineConfig) GetStreamLatencyLevel() models.LatencyLevel {
	return models.GetLatencyLevel(2)
}
func (f fakeEngineConfig) GetStreamOutputVariants() []models.StreamOutputVariant { return nil }
func (f fakeEngineConfig) GetVideoServingEndpoint() string                       { return f.servingEndpoint }
func (f fakeEngineConfig) GetStreamKeys() []models.StreamKey                     { return nil }

var _ models.EngineConfig = fakeEngineConfig{}

// TestLocalStorageBuildsFromEngineConfig proves a storage provider can be
// constructed and set up from a plain EngineConfig with no database behind it,
// and that it reads the configured serving endpoint through the interface.
func TestLocalStorageBuildsFromEngineConfig(t *testing.T) {
	s := NewLocalStorage(fakeEngineConfig{servingEndpoint: "https://cdn.example/"})
	if err := s.Setup(); err != nil {
		t.Fatalf("Setup: %v", err)
	}
	if s.host != "https://cdn.example/" {
		t.Errorf("host = %q, want %q", s.host, "https://cdn.example/")
	}
}
