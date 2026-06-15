package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExtractFederationUsername(t *testing.T) {
	tests := []struct {
		name        string
		nodeInfo    NodeInfoV2
		expected    string
		expectError bool
	}{
		{
			name: "Valid federation username",
			nodeInfo: NodeInfoV2{
				Software: struct {
					Name    string `json:"name"`
					Version string `json:"version"`
				}{
					Name:    "owncast",
					Version: "0.1.0",
				},
				Metadata: struct {
					Federation struct {
						Username        string `json:"username"`
						FeaturedStreams int    `json:"featured_streams"`
					} `json:"federation"`
					ChatEnabled bool `json:"chat_enabled"`
				}{
					Federation: struct {
						Username        string `json:"username"`
						FeaturedStreams int    `json:"featured_streams"`
					}{
						Username:        "testuser",
						FeaturedStreams: 1,
					},
				},
			},
			expected:    "testuser",
			expectError: false,
		},
		{
			name: "Missing username in federation",
			nodeInfo: NodeInfoV2{
				Software: struct {
					Name    string `json:"name"`
					Version string `json:"version"`
				}{
					Name:    "owncast",
					Version: "0.1.0",
				},
				Metadata: struct {
					Federation struct {
						Username        string `json:"username"`
						FeaturedStreams int    `json:"featured_streams"`
					} `json:"federation"`
					ChatEnabled bool `json:"chat_enabled"`
				}{
					Federation: struct {
						Username        string `json:"username"`
						FeaturedStreams int    `json:"featured_streams"`
					}{
						Username:        "",
						FeaturedStreams: 0,
					},
				},
			},
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractFederationUsername(&tt.nodeInfo)
			if tt.expectError && err == nil {
				t.Errorf("ExtractFederationUsername() expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("ExtractFederationUsername() unexpected error = %v", err)
			} else if result != tt.expected {
				t.Errorf("ExtractFederationUsername() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidateOwncastServer(t *testing.T) {
	tests := []struct {
		name        string
		nodeInfo    NodeInfoV2
		expectError bool
	}{
		{
			name: "Valid Owncast server",
			nodeInfo: NodeInfoV2{
				Software: struct {
					Name    string `json:"name"`
					Version string `json:"version"`
				}{
					Name:    "owncast",
					Version: "0.1.0",
				},
				Protocols: []string{"activitypub"},
			},
			expectError: false,
		},
		{
			name: "Non-Owncast server (Mastodon)",
			nodeInfo: NodeInfoV2{
				Software: struct {
					Name    string `json:"name"`
					Version string `json:"version"`
				}{
					Name:    "mastodon",
					Version: "4.0.0",
				},
				Protocols: []string{"activitypub"},
			},
			expectError: true,
		},
		{
			name: "Owncast without ActivityPub",
			nodeInfo: NodeInfoV2{
				Software: struct {
					Name    string `json:"name"`
					Version string `json:"version"`
				}{
					Name:    "owncast",
					Version: "0.1.0",
				},
				Protocols: []string{},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOwncastServer(&tt.nodeInfo)
			if tt.expectError && err == nil {
				t.Errorf("ValidateOwncastServer() expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("ValidateOwncastServer() unexpected error = %v", err)
			}
		})
	}
}

func TestValidateFeaturedStreamsSupport(t *testing.T) {
	tests := []struct {
		name        string
		nodeInfo    NodeInfoV2
		expectError bool
	}{
		{
			name: "Featured streams supported",
			nodeInfo: NodeInfoV2{
				Metadata: struct {
					Federation struct {
						Username        string `json:"username"`
						FeaturedStreams int    `json:"featured_streams"`
					} `json:"federation"`
					ChatEnabled bool `json:"chat_enabled"`
				}{
					Federation: struct {
						Username        string `json:"username"`
						FeaturedStreams int    `json:"featured_streams"`
					}{
						Username:        "testuser",
						FeaturedStreams: 1,
					},
				},
			},
			expectError: false,
		},
		{
			name: "Featured streams unsupported",
			nodeInfo: NodeInfoV2{
				Metadata: struct {
					Federation struct {
						Username        string `json:"username"`
						FeaturedStreams int    `json:"featured_streams"`
					} `json:"federation"`
					ChatEnabled bool `json:"chat_enabled"`
				}{
					Federation: struct {
						Username        string `json:"username"`
						FeaturedStreams int    `json:"featured_streams"`
					}{
						Username:        "testuser",
						FeaturedStreams: 0,
					},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFeaturedStreamsSupport(&tt.nodeInfo)
			if tt.expectError && err == nil {
				t.Errorf("ValidateFeaturedStreamsSupport() expected error but got none")
			} else if !tt.expectError && err != nil {
				t.Errorf("ValidateFeaturedStreamsSupport() unexpected error = %v", err)
			}
		})
	}
}

func TestFetchNodeInfo(t *testing.T) {
	// httptest.NewServer binds to a loopback address; opt into the same
	// integration-test bypass that the AP test scripts use.
	t.Setenv("OWNCAST_ALLOW_INTERNAL_FEDERATION", "true")

	tests := []struct {
		name                 string
		setupServer          func() *httptest.Server
		expectedError        bool
		expectedSoftwareName string
	}{
		{
			name: "Valid nodeinfo response",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/.well-known/nodeinfo" {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{
							"links": []map[string]interface{}{
								{
									"rel":  "http://nodeinfo.diaspora.software/ns/schema/2.0",
									"href": "http://" + r.Host + "/nodeinfo/2.0",
								},
							},
						})
					} else if r.URL.Path == "/nodeinfo/2.0" {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(NodeInfoV2{
							Software: struct {
								Name    string `json:"name"`
								Version string `json:"version"`
							}{
								Name:    "owncast",
								Version: "0.1.0",
							},
							Metadata: struct {
								Federation struct {
									Username        string `json:"username"`
									FeaturedStreams int    `json:"featured_streams"`
								} `json:"federation"`
								ChatEnabled bool `json:"chat_enabled"`
							}{
								Federation: struct {
									Username        string `json:"username"`
									FeaturedStreams int    `json:"featured_streams"`
								}{
									Username:        "testuser",
									FeaturedStreams: 1,
								},
							},
							Protocols: []string{"activitypub"},
						})
					}
				}))
			},
			expectedError:        false,
			expectedSoftwareName: "owncast",
		},
		{
			name: "Server returns 404",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}))
			},
			expectedError: true,
		},
		{
			name: "Invalid JSON response",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/.well-known/nodeinfo" {
						w.Header().Set("Content-Type", "application/json")
						w.Write([]byte("invalid json"))
					}
				}))
			},
			expectedError: true,
		},
		{
			name: "No nodeinfo 2.0 link",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/.well-known/nodeinfo" {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(map[string]interface{}{
							"links": []map[string]interface{}{
								{
									"rel":  "http://nodeinfo.diaspora.software/ns/schema/1.0",
									"href": "http://" + r.Host + "/nodeinfo/1.0",
								},
							},
						})
					}
				}))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer()
			defer server.Close()

			nodeInfo, err := FetchNodeInfo(server.URL)

			if tt.expectedError {
				if err == nil {
					t.Errorf("FetchNodeInfo() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("FetchNodeInfo() unexpected error = %v", err)
				}
				if nodeInfo == nil {
					t.Errorf("FetchNodeInfo() returned nil nodeinfo")
				} else if nodeInfo.Software.Name != tt.expectedSoftwareName {
					t.Errorf("FetchNodeInfo() software name = %v, want %v", nodeInfo.Software.Name, tt.expectedSoftwareName)
				}
			}
		})
	}
}
