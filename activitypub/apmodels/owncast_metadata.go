package apmodels

import (
	"fmt"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/persistence/configrepository"
)

// OwncastMetadata represents parsed Owncast custom properties from ActivityPub activities.
type OwncastMetadata struct {
	StreamStatus      string
	StreamTitle       string
	StreamDescription string
	ServerName        string
	LogoURL           string
	ThumbnailURL      string
	Tags              []string
	IsOwncastServer   bool
}

// ParseOwncastMetadata extracts Owncast custom properties from unknown properties map.
func ParseOwncastMetadata(unknownProps map[string]interface{}) *OwncastMetadata {
	metadata := &OwncastMetadata{}

	// Check if this is an Owncast server by looking for any Owncast namespace properties
	if _, hasStreamStatus := unknownProps[config.APOwncastNamespaceStreamStatus]; hasStreamStatus {
		metadata.IsOwncastServer = true
	}

	// Parse stream status
	if streamStatus, exists := unknownProps[config.APOwncastNamespaceStreamStatus]; exists {
		if status, ok := streamStatus.(string); ok {
			metadata.StreamStatus = status
			metadata.IsOwncastServer = true
		}
	}

	// Parse stream title
	if streamTitle, exists := unknownProps[config.APOwncastNamespaceStreamTitle]; exists {
		if title, ok := streamTitle.(string); ok {
			metadata.StreamTitle = title
			metadata.IsOwncastServer = true
		}
	}

	// Parse stream description
	if streamDescription, exists := unknownProps[config.APOwncastNamespaceStreamDescription]; exists {
		if desc, ok := streamDescription.(string); ok {
			metadata.StreamDescription = desc
			metadata.IsOwncastServer = true
		}
	}

	// Parse server name
	if serverName, exists := unknownProps[config.APOwncastNamespaceServerName]; exists {
		if name, ok := serverName.(string); ok {
			metadata.ServerName = name
			metadata.IsOwncastServer = true
		}
	}

	// Parse logo URL
	if logoURL, exists := unknownProps[config.APOwncastNamespaceLogoURL]; exists {
		if logo, ok := logoURL.(string); ok {
			metadata.LogoURL = logo
			metadata.IsOwncastServer = true
		}
	}

	// Parse thumbnail URL
	if thumbnailURL, exists := unknownProps[config.APOwncastNamespaceThumbnailURL]; exists {
		if thumb, ok := thumbnailURL.(string); ok {
			metadata.ThumbnailURL = thumb
			metadata.IsOwncastServer = true
		}
	}

	// Parse tags
	if tags, exists := unknownProps[config.APOwncastNamespaceStreamTags]; exists {
		if tagList, ok := tags.([]interface{}); ok && len(tagList) > 0 {
			tagStrings := make([]string, 0, len(tagList))
			for _, tag := range tagList {
				if tagStr, ok := tag.(string); ok {
					tagStrings = append(tagStrings, tagStr)
				}
			}
			if len(tagStrings) > 0 {
				metadata.Tags = tagStrings
				metadata.IsOwncastServer = true
			}
		}
	}

	return metadata
}

// SetOwncastMetadata sets Owncast metadata properties in unknownProps map from ConfigRepository.
// It always includes stream status.
func SetOwncastMetadata(unknownProps map[string]interface{}, repo configrepository.ConfigRepository, isStreamConnected bool) {
	// Always include server identification
	unknownProps[config.APOwncastNamespaceServerName] = repo.GetServerName()
	unknownProps[config.APOwncastNamespaceStreamDescription] = repo.GetServerSummary()

	// Always include current stream status
	if isStreamConnected {
		unknownProps[config.APOwncastNamespaceStreamStatus] = "live"
	} else {
		unknownProps[config.APOwncastNamespaceStreamStatus] = "offline"
	}

	// Add stream title if available
	if streamTitle := repo.GetStreamTitle(); streamTitle != "" {
		unknownProps[config.APOwncastNamespaceStreamTitle] = streamTitle
	}

	// Add logo if available
	if logoPath := repo.GetLogoPath(); logoPath != "" {
		logoURL := fmt.Sprintf("%s/%s", repo.GetServerURL(), logoPath)
		unknownProps[config.APOwncastNamespaceLogoURL] = logoURL
	}

	// Add tags if available
	if tags := repo.GetServerMetadataTags(); len(tags) > 0 {
		unknownProps[config.APOwncastNamespaceStreamTags] = tags
	}
}

// SetBasicOwncastMetadata sets only the basic server identification metadata.
// This is useful for responses that don't need full stream information.
func SetBasicOwncastMetadata(unknownProps map[string]interface{}, repo configrepository.ConfigRepository, isStreamConnected bool) {
	// Always include server identification
	unknownProps[config.APOwncastNamespaceServerName] = repo.GetServerName()
	unknownProps[config.APOwncastNamespaceStreamDescription] = repo.GetServerSummary()

	// Always include current stream status
	if isStreamConnected {
		unknownProps[config.APOwncastNamespaceStreamStatus] = "live"
	} else {
		unknownProps[config.APOwncastNamespaceStreamStatus] = "offline"
	}
}
