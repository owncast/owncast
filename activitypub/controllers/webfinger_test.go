package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/persistence/configrepository"
)

func TestMain(m *testing.M) {
	dbFile, err := os.CreateTemp(os.TempDir(), "owncast-test-db.db")
	if err != nil {
		panic(err)
	}
	defer os.Remove(dbFile.Name())

	if err := data.SetupPersistence(dbFile.Name()); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestWebfingerHandlerWithIDNHost(t *testing.T) {
	tests := []struct {
		name             string
		serverURL        string
		resource         string
		expectedHost     string
		expectedActorURL string
	}{
		{
			name:             "unicode account host",
			serverURL:        "https://live.retrospection.みんな",
			resource:         "acct:retrots3m@live.retrospection.みんな",
			expectedHost:     "live.retrospection.xn--q9jyb4c",
			expectedActorURL: "https://live.retrospection.xn--q9jyb4c/federation/user/retrots3m",
		},
		{
			name:             "punycode account host",
			serverURL:        "https://live.retrospection.みんな",
			resource:         "acct:retrots3m@live.retrospection.xn--q9jyb4c",
			expectedHost:     "live.retrospection.xn--q9jyb4c",
			expectedActorURL: "https://live.retrospection.xn--q9jyb4c/federation/user/retrots3m",
		},
		{
			name:             "unicode account host with port",
			serverURL:        "https://live.retrospection.みんな:8443",
			resource:         "acct:retrots3m@live.retrospection.みんな:8443",
			expectedHost:     "live.retrospection.xn--q9jyb4c:8443",
			expectedActorURL: "https://live.retrospection.xn--q9jyb4c:8443/federation/user/retrots3m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configRepository := configrepository.Get()
			configRepository.SetFederationEnabled(true)
			configRepository.SetFederationUsername("retrots3m")
			configRepository.SetServerURL(tt.serverURL)

			req := httptest.NewRequest(http.MethodGet, "/.well-known/webfinger?resource="+url.QueryEscape(tt.resource), nil)
			w := httptest.NewRecorder()

			WebfingerHandler(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
			}

			var response struct {
				Subject string   `json:"subject"`
				Aliases []string `json:"aliases"`
				Links   []struct {
					Rel  string `json:"rel"`
					Href string `json:"href"`
				} `json:"links"`
			}
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatal(err)
			}

			if response.Subject != "acct:retrots3m@"+tt.expectedHost {
				t.Errorf("subject = %v, want acct:retrots3m@%s", response.Subject, tt.expectedHost)
			}
			if len(response.Aliases) != 1 || response.Aliases[0] != tt.expectedActorURL {
				t.Errorf("aliases = %v, want [%s]", response.Aliases, tt.expectedActorURL)
			}

			var self string
			var avatar string
			var alternate string
			for _, link := range response.Links {
				switch link.Rel {
				case "self":
					self = link.Href
				case "http://webfinger.net/rel/avatar":
					avatar = link.Href
				case "alternate":
					alternate = link.Href
				}
			}

			if self != tt.expectedActorURL {
				t.Errorf("self href = %v, want %v", self, tt.expectedActorURL)
			}
			if avatar != "https://"+tt.expectedHost+"/logo/external" {
				t.Errorf("avatar href = %v, want punycode logo URL", avatar)
			}
			if alternate != "https://"+tt.expectedHost+"/hls/stream.m3u8" {
				t.Errorf("alternate href = %v, want punycode stream URL", alternate)
			}
		})
	}
}
