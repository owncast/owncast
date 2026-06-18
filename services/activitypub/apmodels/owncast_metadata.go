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

	metadata.StreamStatus = extractStringProp(unknownProps, config.APOwncastNamespaceStreamStatus)
	metadata.StreamTitle = extractStringProp(unknownProps, config.APOwncastNamespaceStreamTitle)
	metadata.StreamDescription = extractStringProp(unknownProps, config.APOwncastNamespaceStreamDescription)
	metadata.ServerName = extractStringProp(unknownProps, config.APOwncastNamespaceServerName)
	metadata.LogoURL = extractStringProp(unknownProps, config.APOwncastNamespaceLogoURL)
	metadata.ThumbnailURL = extractStringProp(unknownProps, config.APOwncastNamespaceThumbnailURL)
	metadata.Tags = extractTagsProp(unknownProps)

	metadata.IsOwncastServer = metadata.StreamStatus != "" ||
		metadata.StreamTitle != "" ||
		metadata.StreamDescription != "" ||
		metadata.ServerName != "" ||
		metadata.LogoURL != "" ||
		metadata.ThumbnailURL != "" ||
		len(metadata.Tags) > 0

	return metadata
}

// extractStringProp extracts a string value from the unknown properties map for a given key.
func extractStringProp(unknownProps map[string]interface{}, key string) string {
	if val, exists := unknownProps[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// extractTagsProp extracts the tags list from the unknown properties map.
func extractTagsProp(unknownProps map[string]interface{}) []string {
	tags, exists := unknownProps[config.APOwncastNamespaceStreamTags]
	if !exists {
		return nil
	}

	tagList, ok := tags.([]interface{})
	if !ok || len(tagList) == 0 {
		return nil
	}

	tagStrings := make([]string, 0, len(tagList))
	for _, tag := range tagList {
		if tagStr, ok := tag.(string); ok {
			tagStrings = append(tagStrings, tagStr)
		}
	}

	if len(tagStrings) == 0 {
		return nil
	}

	return tagStrings
}

// SetOwncastMetadata sets Owncast metadata properties in unknownProps map from ConfigRepository.
// It always includes stream status.
func SetOwncastMetadata(unknownProps map[string]interface{}, repo configrepository.ConfigRepository, isStreamConnected bool) {
	// Always include server identification
	unknownProps[config.APOwncastNamespaceServerName] = repo.GetServerName()
	unknownProps[config.APOwncastNamespaceStreamDescription] = repo.GetServerSummary()

	// Always include current stream status
	if isStreamConnected {
		unknownProps[config.APOwncastNamespaceStreamStatus] = config.APStreamStatusLive
		// The live preview thumbnail only exists while streaming. Advertise it so
		// following servers can show a preview in their featured-streams directory.
		unknownProps[config.APOwncastNamespaceThumbnailURL] = fmt.Sprintf("%s/thumbnail.jpg", repo.GetServerURL())
	} else {
		unknownProps[config.APOwncastNamespaceStreamStatus] = config.APStreamStatusOffline
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
		unknownProps[config.APOwncastNamespaceStreamStatus] = config.APStreamStatusLive
	} else {
		unknownProps[config.APOwncastNamespaceStreamStatus] = config.APStreamStatusOffline
	}
}
