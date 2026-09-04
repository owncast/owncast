package main

import (
	"os"
	"testing"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/services/datastore"
)

func newTestConfigRepository(t *testing.T) configrepository.ConfigRepository {
	t.Helper()
	ds, err := datastore.SetupPersistence(":memory:", os.TempDir())
	if err != nil {
		t.Fatalf("failed to set up datastore: %v", err)
	}
	return configrepository.New(ds)
}

func TestHandleAutoHTTPSEnvironment(t *testing.T) {
	tests := []struct {
		name            string
		hostName        string
		enable          string
		wantEnabled     bool
		wantHost        string
		wantServerURL   string
		presetServerURL string
	}{
		{
			name: "neither set does nothing",
		},
		{
			name:          "hostname alone seeds the server url without enabling",
			hostName:      "live.example.com",
			wantHost:      "live.example.com",
			wantServerURL: "https://live.example.com",
		},
		{
			name:   "enable alone does not enable",
			enable: "true",
		},
		{
			name:          "hostname and enable turn on automatic https",
			hostName:      "live.example.com",
			enable:        "true",
			wantEnabled:   true,
			wantHost:      "live.example.com",
			wantServerURL: "https://live.example.com",
		},
		{
			name:          "hostname is trimmed and lowercased",
			hostName:      "  Live.Example.COM ",
			enable:        "TRUE",
			wantEnabled:   true,
			wantHost:      "live.example.com",
			wantServerURL: "https://live.example.com",
		},
		{
			name:          "a trailing dot is removed",
			hostName:      "Live.Example.COM.",
			enable:        "true",
			wantEnabled:   true,
			wantHost:      "live.example.com",
			wantServerURL: "https://live.example.com",
		},
		{
			name:     "a url instead of a hostname is rejected",
			hostName: "https://live.example.com",
			enable:   "true",
		},
		{
			name:     "a hostname with a port is rejected",
			hostName: "live.example.com:8080",
			enable:   "true",
		},
		{
			name:            "a hostname with a query is rejected",
			hostName:        "live.example.com?unexpected=true",
			enable:          "true",
			presetServerURL: "https://old.example.com",
		},
		{
			name:     "a hostname with a fragment is rejected",
			hostName: "live.example.com#unexpected",
			enable:   "true",
		},
		{
			name:     "a hostname with internal whitespace is rejected",
			hostName: "live.example.com\n.evil.example",
			enable:   "true",
		},
		{
			name:            "a differing stored server url is overwritten",
			hostName:        "live.example.com",
			enable:          "true",
			presetServerURL: "https://old.example.com",
			wantEnabled:     true,
			wantHost:        "live.example.com",
			wantServerURL:   "https://live.example.com",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("OWNCAST_HOST_NAME", tc.hostName)
			t.Setenv("OWNCAST_ENABLE_AUTO_HTTPS", tc.enable)

			repo := newTestConfigRepository(t)
			if tc.presetServerURL != "" {
				if err := repo.SetServerURL(tc.presetServerURL); err != nil {
					t.Fatal(err)
				}
			}

			cfg := config.NewDefault()
			handleAutoHTTPSEnvironment(cfg, repo)

			if cfg.AutoHTTPSEnabled != tc.wantEnabled {
				t.Errorf("AutoHTTPSEnabled = %v, want %v", cfg.AutoHTTPSEnabled, tc.wantEnabled)
			}
			if cfg.AutoHTTPSHost != tc.wantHost {
				t.Errorf("AutoHTTPSHost = %q, want %q", cfg.AutoHTTPSHost, tc.wantHost)
			}

			wantStored := tc.wantServerURL
			if wantStored == "" {
				wantStored = tc.presetServerURL
			}
			if got := repo.GetServerURL(); got != wantStored {
				t.Errorf("stored server URL = %q, want %q", got, wantStored)
			}
		})
	}
}
