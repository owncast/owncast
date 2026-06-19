package apmodels

import (
	"testing"

	"github.com/owncast/owncast/config"
)

// mockConfigRepository is a minimal mock implementation for testing metadata functions
type mockConfigRepository struct {
	serverName         string
	serverSummary      string
	streamTitle        string
	logoPath           string
	serverURL          string
	serverMetadataTags []string
}

func (m *mockConfigRepository) GetServerName() string {
	return m.serverName
}

func (m *mockConfigRepository) GetServerSummary() string {
	return m.serverSummary
}

func (m *mockConfigRepository) GetStreamTitle() string {
	return m.streamTitle
}

func (m *mockConfigRepository) GetLogoPath() string {
	return m.logoPath
}

func (m *mockConfigRepository) GetServerURL() string {
	return m.serverURL
}

func (m *mockConfigRepository) GetServerMetadataTags() []string {
	return m.serverMetadataTags
}

func TestSetOwncastMetadata(t *testing.T) {
	tests := []struct {
		name                string
		repo                *mockConfigRepository
		includeStreamStatus bool
		streamStatus        string
		expectedKeys        []string
		unexpectedKeys      []string
	}{
		{
			name: "basic metadata without stream status",
			repo: &mockConfigRepository{
				serverName:    "Test Server",
				serverSummary: "A test server description",
			},
			includeStreamStatus: false,
			streamStatus:        "",
			expectedKeys:        []string{config.APOwncastNamespaceServerName, config.APOwncastNamespaceStreamDescription},
			unexpectedKeys:      []string{config.APOwncastNamespaceStreamStatus, config.APOwncastNamespaceStreamTitle, config.APOwncastNamespaceLogoURL, config.APOwncastNamespaceStreamTags},
		},
		{
			name: "full metadata with stream status",
			repo: &mockConfigRepository{
				serverName:         "Test Server",
				serverSummary:      "A test server description",
				streamTitle:        "Test Stream",
				logoPath:           "logo.png",
				serverURL:          "https://example.com",
				serverMetadataTags: []string{"gaming", "tech"},
			},
			includeStreamStatus: true,
			streamStatus:        "live",
			expectedKeys:        []string{config.APOwncastNamespaceServerName, config.APOwncastNamespaceStreamDescription, config.APOwncastNamespaceStreamStatus, config.APOwncastNamespaceStreamTitle, config.APOwncastNamespaceLogoURL, config.APOwncastNamespaceStreamTags},
			unexpectedKeys:      []string{},
		},
		{
			name: "metadata with stream title but no stream status",
			repo: &mockConfigRepository{
				serverName:    "Test Server",
				serverSummary: "A test server description",
				streamTitle:   "Test Stream",
			},
			includeStreamStatus: false,
			streamStatus:        "",
			expectedKeys:        []string{config.APOwncastNamespaceServerName, config.APOwncastNamespaceStreamDescription, config.APOwncastNamespaceStreamTitle},
			unexpectedKeys:      []string{config.APOwncastNamespaceStreamStatus, config.APOwncastNamespaceLogoURL, config.APOwncastNamespaceStreamTags},
		},
		{
			name: "metadata with empty logo path",
			repo: &mockConfigRepository{
				serverName:    "Test Server",
				serverSummary: "A test server description",
				logoPath:      "",
				serverURL:     "https://example.com",
			},
			includeStreamStatus: false,
			streamStatus:        "",
			expectedKeys:        []string{config.APOwncastNamespaceServerName, config.APOwncastNamespaceStreamDescription},
			unexpectedKeys:      []string{config.APOwncastNamespaceLogoURL},
		},
		{
			name: "metadata with empty tags",
			repo: &mockConfigRepository{
				serverName:         "Test Server",
				serverSummary:      "A test server description",
				serverMetadataTags: []string{},
			},
			includeStreamStatus: false,
			streamStatus:        "",
			expectedKeys:        []string{config.APOwncastNamespaceServerName, config.APOwncastNamespaceStreamDescription},
			unexpectedKeys:      []string{config.APOwncastNamespaceStreamTags},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unknownProps := make(map[string]interface{})
			// Call the function directly without the interface
			setOwncastMetadataImpl(unknownProps, tt.repo, tt.includeStreamStatus, tt.streamStatus)

			// Check expected keys are present
			for _, key := range tt.expectedKeys {
				if _, exists := unknownProps[key]; !exists {
					t.Errorf("Expected key %s not found in unknownProps", key)
				}
			}

			// Check unexpected keys are not present
			for _, key := range tt.unexpectedKeys {
				if _, exists := unknownProps[key]; exists {
					t.Errorf("Unexpected key %s found in unknownProps", key)
				}
			}

			// Validate specific values
			if val, exists := unknownProps[config.APOwncastNamespaceServerName]; exists {
				if val != tt.repo.serverName {
					t.Errorf("Expected server name %s, got %v", tt.repo.serverName, val)
				}
			}

			if val, exists := unknownProps[config.APOwncastNamespaceStreamDescription]; exists {
				if val != tt.repo.serverSummary {
					t.Errorf("Expected server summary %s, got %v", tt.repo.serverSummary, val)
				}
			}

			if tt.includeStreamStatus && tt.streamStatus != "" {
				if val, exists := unknownProps[config.APOwncastNamespaceStreamStatus]; exists {
					if val != tt.streamStatus {
						t.Errorf("Expected stream status %s, got %v", tt.streamStatus, val)
					}
				}
			}

			if tt.repo.logoPath != "" && tt.repo.serverURL != "" {
				if val, exists := unknownProps[config.APOwncastNamespaceLogoURL]; exists {
					expectedURL := tt.repo.serverURL + "/" + tt.repo.logoPath
					if val != expectedURL {
						t.Errorf("Expected logo URL %s, got %v", expectedURL, val)
					}
				}
			}

			if len(tt.repo.serverMetadataTags) > 0 {
				if val, exists := unknownProps[config.APOwncastNamespaceStreamTags]; exists {
					tags, ok := val.([]interface{})
					if !ok {
						t.Errorf("Expected tags to be []interface{}, got %T", val)
					} else if len(tags) != len(tt.repo.serverMetadataTags) {
						t.Errorf("Expected %d tags, got %d", len(tt.repo.serverMetadataTags), len(tags))
					}
				}
			}
		})
	}
}

func TestSetBasicOwncastMetadata(t *testing.T) {
	repo := &mockConfigRepository{
		serverName:         "Basic Test Server",
		serverSummary:      "Basic server description",
		streamTitle:        "Should not be included",
		logoPath:           "logo.png",
		serverURL:          "https://example.com",
		serverMetadataTags: []string{"tag1", "tag2"},
	}

	unknownProps := make(map[string]interface{})
	// Call the implementation function directly
	setOwncastMetadataImpl(unknownProps, repo, false, "")

	// Basic metadata should always include server name and description
	if val, exists := unknownProps[config.APOwncastNamespaceServerName]; !exists || val != repo.serverName {
		t.Errorf("Expected server name %s, got %v", repo.serverName, val)
	}

	if val, exists := unknownProps[config.APOwncastNamespaceStreamDescription]; !exists || val != repo.serverSummary {
		t.Errorf("Expected server summary %s, got %v", repo.serverSummary, val)
	}

	// Stream status should never be included in basic metadata
	if _, exists := unknownProps[config.APOwncastNamespaceStreamStatus]; exists {
		t.Error("Stream status should not be included in basic metadata")
	}

	// Stream title, logo, and tags should still be included if available
	if _, exists := unknownProps[config.APOwncastNamespaceStreamTitle]; !exists {
		t.Error("Stream title should be included when available")
	}

	if _, exists := unknownProps[config.APOwncastNamespaceLogoURL]; !exists {
		t.Error("Logo URL should be included when available")
	}

	if _, exists := unknownProps[config.APOwncastNamespaceStreamTags]; !exists {
		t.Error("Tags should be included when available")
	}
}

func TestParseOwncastMetadata(t *testing.T) {
	tests := []struct {
		name            string
		unknownProps    map[string]interface{}
		expected        *OwncastMetadata
		expectDirectory bool
	}{
		{
			name: "parse full stream metadata without directory marker",
			unknownProps: map[string]interface{}{
				config.APOwncastNamespaceStreamStatus:      "live",
				config.APOwncastNamespaceStreamTitle:       "Test Stream",
				config.APOwncastNamespaceStreamDescription: "Test Description",
				config.APOwncastNamespaceServerName:        "Test Server",
				config.APOwncastNamespaceLogoURL:           "https://example.com/logo.png",
				config.APOwncastNamespaceThumbnailURL:      "https://example.com/thumb.png",
				config.APOwncastNamespaceStreamTags:        []interface{}{"gaming", "tech"},
			},
			expected: &OwncastMetadata{
				StreamStatus:      "live",
				StreamTitle:       "Test Stream",
				StreamDescription: "Test Description",
				ServerName:        "Test Server",
				LogoURL:           "https://example.com/logo.png",
				ThumbnailURL:      "https://example.com/thumb.png",
				Tags:              []string{"gaming", "tech"},
				IsDirectory:       false,
			},
			// Stream fields alone no longer flag a peer as a directory.
			expectDirectory: false,
		},
		{
			name: "parse directory marker",
			unknownProps: map[string]interface{}{
				config.APOwncastNamespaceDirectory: true,
			},
			expected: &OwncastMetadata{
				IsDirectory: true,
			},
			expectDirectory: true,
		},
		{
			name: "parse directory marker with server metadata",
			unknownProps: map[string]interface{}{
				config.APOwncastNamespaceDirectory:  true,
				config.APOwncastNamespaceServerName: "Some Directory",
			},
			expected: &OwncastMetadata{
				ServerName:  "Some Directory",
				IsDirectory: true,
			},
			expectDirectory: true,
		},
		{
			name: "parse directory marker sent as string true",
			unknownProps: map[string]interface{}{
				config.APOwncastNamespaceDirectory: "true",
			},
			expected: &OwncastMetadata{
				IsDirectory: true,
			},
			expectDirectory: true,
		},
		{
			name: "parse partial metadata",
			unknownProps: map[string]interface{}{
				config.APOwncastNamespaceServerName:        "Test Server",
				config.APOwncastNamespaceStreamDescription: "Test Description",
			},
			expected: &OwncastMetadata{
				ServerName:        "Test Server",
				StreamDescription: "Test Description",
				IsDirectory:       false,
			},
			expectDirectory: false,
		},
		{
			name:         "parse empty metadata",
			unknownProps: map[string]interface{}{},
			expected: &OwncastMetadata{
				IsDirectory: false,
			},
			expectDirectory: false,
		},
		{
			name: "parse metadata with wrong types",
			unknownProps: map[string]interface{}{
				config.APOwncastNamespaceStreamStatus: 123, // Wrong type
				config.APOwncastNamespaceServerName:   "Test Server",
			},
			expected: &OwncastMetadata{
				ServerName:  "Test Server",
				IsDirectory: false,
			},
			expectDirectory: false,
		},
		{
			name: "parse metadata with invalid tags",
			unknownProps: map[string]interface{}{
				config.APOwncastNamespaceServerName: "Test Server",
				config.APOwncastNamespaceStreamTags: "not an array", // Wrong type
			},
			expected: &OwncastMetadata{
				ServerName:  "Test Server",
				IsDirectory: false,
			},
			expectDirectory: false,
		},
		{
			name: "parse metadata with mixed tag types",
			unknownProps: map[string]interface{}{
				config.APOwncastNamespaceServerName: "Test Server",
				config.APOwncastNamespaceStreamTags: []interface{}{"gaming", 123, "tech"}, // Mixed types
			},
			expected: &OwncastMetadata{
				ServerName:  "Test Server",
				Tags:        []string{"gaming", "tech"}, // Only string values
				IsDirectory: false,
			},
			expectDirectory: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata := ParseOwncastMetadata(tt.unknownProps)

			if metadata.IsDirectory != tt.expectDirectory {
				t.Errorf("Expected IsDirectory=%v, got %v", tt.expectDirectory, metadata.IsDirectory)
			}

			if tt.expected.StreamStatus != "" && metadata.StreamStatus != tt.expected.StreamStatus {
				t.Errorf("Expected StreamStatus=%s, got %s", tt.expected.StreamStatus, metadata.StreamStatus)
			}

			if tt.expected.StreamTitle != "" && metadata.StreamTitle != tt.expected.StreamTitle {
				t.Errorf("Expected StreamTitle=%s, got %s", tt.expected.StreamTitle, metadata.StreamTitle)
			}

			if tt.expected.StreamDescription != "" && metadata.StreamDescription != tt.expected.StreamDescription {
				t.Errorf("Expected StreamDescription=%s, got %s", tt.expected.StreamDescription, metadata.StreamDescription)
			}

			if tt.expected.ServerName != "" && metadata.ServerName != tt.expected.ServerName {
				t.Errorf("Expected ServerName=%s, got %s", tt.expected.ServerName, metadata.ServerName)
			}

			if tt.expected.LogoURL != "" && metadata.LogoURL != tt.expected.LogoURL {
				t.Errorf("Expected LogoURL=%s, got %s", tt.expected.LogoURL, metadata.LogoURL)
			}

			if tt.expected.ThumbnailURL != "" && metadata.ThumbnailURL != tt.expected.ThumbnailURL {
				t.Errorf("Expected ThumbnailURL=%s, got %s", tt.expected.ThumbnailURL, metadata.ThumbnailURL)
			}

			if tt.expected.Tags != nil {
				if len(metadata.Tags) != len(tt.expected.Tags) {
					t.Errorf("Expected %d tags, got %d", len(tt.expected.Tags), len(metadata.Tags))
				} else {
					for i, tag := range tt.expected.Tags {
						if metadata.Tags[i] != tag {
							t.Errorf("Expected tag[%d]=%s, got %s", i, tag, metadata.Tags[i])
						}
					}
				}
			}
		})
	}
}

func TestRoundTripMetadata(t *testing.T) {
	// Test that setting and then parsing metadata works correctly
	repo := &mockConfigRepository{
		serverName:         "Round Trip Server",
		serverSummary:      "Round trip description",
		streamTitle:        "Round Trip Stream",
		logoPath:           "logo.png",
		serverURL:          "https://example.com",
		serverMetadataTags: []string{"tag1", "tag2"},
	}

	// Set metadata using the implementation function
	unknownProps := make(map[string]interface{})
	setOwncastMetadataImpl(unknownProps, repo, true, "live")

	// Parse it back
	metadata := ParseOwncastMetadata(unknownProps)

	// Verify round trip
	if metadata.ServerName != repo.serverName {
		t.Errorf("Round trip failed for ServerName: expected %s, got %s", repo.serverName, metadata.ServerName)
	}

	if metadata.StreamDescription != repo.serverSummary {
		t.Errorf("Round trip failed for StreamDescription: expected %s, got %s", repo.serverSummary, metadata.StreamDescription)
	}

	if metadata.StreamTitle != repo.streamTitle {
		t.Errorf("Round trip failed for StreamTitle: expected %s, got %s", repo.streamTitle, metadata.StreamTitle)
	}

	if metadata.StreamStatus != "live" {
		t.Errorf("Round trip failed for StreamStatus: expected live, got %s", metadata.StreamStatus)
	}

	expectedLogoURL := repo.serverURL + "/" + repo.logoPath
	if metadata.LogoURL != expectedLogoURL {
		t.Errorf("Round trip failed for LogoURL: expected %s, got %s", expectedLogoURL, metadata.LogoURL)
	}

	if len(metadata.Tags) != len(repo.serverMetadataTags) {
		t.Errorf("Round trip failed for Tags: expected %d tags, got %d", len(repo.serverMetadataTags), len(metadata.Tags))
	}

	// Stream metadata alone does not make a peer a directory; only the explicit
	// ns#directory marker does, and the round-trip helper does not set it.
	if metadata.IsDirectory {
		t.Error("Round trip of stream metadata should not be flagged as a directory")
	}
}

// Helper function that replicates the logic of SetOwncastMetadata for testing
// This avoids the ConfigRepository interface dependency
func setOwncastMetadataImpl(unknownProps map[string]interface{}, repo *mockConfigRepository, includeStreamStatus bool, streamStatus string) {
	// Always include server identification
	unknownProps[config.APOwncastNamespaceServerName] = repo.GetServerName()
	unknownProps[config.APOwncastNamespaceStreamDescription] = repo.GetServerSummary()

	// Add stream status if requested
	if includeStreamStatus && streamStatus != "" {
		unknownProps[config.APOwncastNamespaceStreamStatus] = streamStatus
	}

	// Add stream title if available
	if streamTitle := repo.GetStreamTitle(); streamTitle != "" {
		unknownProps[config.APOwncastNamespaceStreamTitle] = streamTitle
	}

	// Add logo if available
	if logoPath := repo.GetLogoPath(); logoPath != "" {
		logoURL := repo.GetServerURL() + "/" + logoPath
		unknownProps[config.APOwncastNamespaceLogoURL] = logoURL
	}

	// Add tags if available
	if tags := repo.GetServerMetadataTags(); len(tags) > 0 {
		// Convert []string to []interface{} for compatibility with ParseOwncastMetadata
		tagInterfaces := make([]interface{}, len(tags))
		for i, tag := range tags {
			tagInterfaces[i] = tag
		}
		unknownProps[config.APOwncastNamespaceStreamTags] = tagInterfaces
	}
}
