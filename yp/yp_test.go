package yp

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
)

type pingConfig struct {
	configrepository.ConfigRepository
	storedKey string
}

func (*pingConfig) GetDirectoryEnabled() bool             { return true }
func (*pingConfig) GetServerURL() string                  { return "https://example.com" }
func (*pingConfig) GetServerName() string                 { return "Test Server" }
func (c *pingConfig) GetDirectoryRegistrationKey() string { return c.storedKey }
func (c *pingConfig) SetDirectoryRegistrationKey(key string) error {
	c.storedKey = key
	return nil
}

func TestPingPersistsReturnedRegistrationKey(t *testing.T) {
	cfg := &pingConfig{storedKey: "old-key"}
	y := New(Deps{
		GetStatus:        func() models.Status { return models.Status{Online: true} },
		ConfigRepository: cfg,
	})
	y.post = func(string, string, io.Reader) (*http.Response, error) {
		return &http.Response{
			Body: io.NopCloser(strings.NewReader(`{"key":"new-key","success":true}`)),
		}, nil
	}

	y.ping()

	if cfg.storedKey != "new-key" {
		t.Fatalf("stored key = %q, want returned key", cfg.storedKey)
	}
}
