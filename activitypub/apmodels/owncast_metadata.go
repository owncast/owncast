package apmodels

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
	if _, hasStreamStatus := unknownProps["https://owncast.online/ns#streamStatus"]; hasStreamStatus {
		metadata.IsOwncastServer = true
	}

	// Parse stream status
	if streamStatus, exists := unknownProps["https://owncast.online/ns#streamStatus"]; exists {
		if status, ok := streamStatus.(string); ok {
			metadata.StreamStatus = status
			metadata.IsOwncastServer = true
		}
	}

	// Parse stream title
	if streamTitle, exists := unknownProps["https://owncast.online/ns#streamTitle"]; exists {
		if title, ok := streamTitle.(string); ok {
			metadata.StreamTitle = title
			metadata.IsOwncastServer = true
		}
	}

	// Parse stream description
	if streamDescription, exists := unknownProps["https://owncast.online/ns#streamDescription"]; exists {
		if desc, ok := streamDescription.(string); ok {
			metadata.StreamDescription = desc
			metadata.IsOwncastServer = true
		}
	}

	// Parse server name
	if serverName, exists := unknownProps["https://owncast.online/ns#serverName"]; exists {
		if name, ok := serverName.(string); ok {
			metadata.ServerName = name
			metadata.IsOwncastServer = true
		}
	}

	// Parse logo URL
	if logoURL, exists := unknownProps["https://owncast.online/ns#logoUrl"]; exists {
		if logo, ok := logoURL.(string); ok {
			metadata.LogoURL = logo
			metadata.IsOwncastServer = true
		}
	}

	// Parse thumbnail URL
	if thumbnailURL, exists := unknownProps["https://owncast.online/ns#thumbnailUrl"]; exists {
		if thumb, ok := thumbnailURL.(string); ok {
			metadata.ThumbnailURL = thumb
			metadata.IsOwncastServer = true
		}
	}

	// Parse tags
	if tags, exists := unknownProps["https://owncast.online/ns#streamTags"]; exists {
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
