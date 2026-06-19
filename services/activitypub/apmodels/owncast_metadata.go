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
	// IsDirectory is true when the activity carries the ns#directory marker,
	// meaning the sender is a directory (or an Owncast server featuring another)
	// asking to receive stream pings. It is set only by the explicit marker.
	IsDirectory bool
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

	// A peer is a directory only when it explicitly says so with the ns#directory
	// marker on its Follow. Presence of stream-metadata fields does not imply a
	// directory: a normal go-live notification carries those too, and a directory
	// has no stream of its own to describe.
	metadata.IsDirectory = extractBoolProp(unknownProps, config.APOwncastNamespaceDirectory)

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

// extractBoolProp extracts a boolean value from the unknown properties map for a
// given key. It accepts a real boolean or the strings "true"/"false", since a
// JSON-LD value may arrive either way depending on the sender.
func extractBoolProp(unknownProps map[string]interface{}, key string) bool {
	val, exists := unknownProps[key]
	if !exists {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	case string:
		return v == "true"
	default:
		return false
	}
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

	// This helper is used to build the Follow we send when featuring another
	// server, so mark ourselves as a directory. That is what makes the remote
	// server hold the follow for its operator to approve and, once approved,
	// deliver its stream pings to us.
	unknownProps[config.APOwncastNamespaceDirectory] = true

	// Always include current stream status
	if isStreamConnected {
		unknownProps[config.APOwncastNamespaceStreamStatus] = config.APStreamStatusLive
	} else {
		unknownProps[config.APOwncastNamespaceStreamStatus] = config.APStreamStatusOffline
	}
}
