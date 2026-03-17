package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// NodeInfoV2 represents the nodeinfo 2.0 response structure.
type NodeInfoV2 struct {
	Metadata struct {
		Federation struct {
			Username string `json:"username"`
		} `json:"federation"`
		ChatEnabled bool `json:"chat_enabled"`
	} `json:"metadata"`
	Software struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"software"`
	Protocols []string `json:"protocols"`
}

// FetchNodeInfo fetches the nodeinfo from a given server URL.
func FetchNodeInfo(serverURL string) (*NodeInfoV2, error) {
	// Parse and validate the URL
	parsedURL, err := url.Parse(serverURL)
	if err != nil {
		return nil, fmt.Errorf("invalid server URL: %w", err)
	}

	// Ensure we're using HTTPS
	if parsedURL.Scheme != "https" && parsedURL.Scheme != "http" {
		return nil, errors.New("server URL must use http or https protocol")
	}

	// First, fetch the well-known nodeinfo endpoint
	wellKnownURL := fmt.Sprintf("%s://%s/.well-known/nodeinfo", parsedURL.Scheme, parsedURL.Host)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(wellKnownURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch well-known nodeinfo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from well-known nodeinfo", resp.StatusCode)
	}

	// Parse well-known response to get nodeinfo 2.0 URL
	var wellKnown struct {
		Links []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"links"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read well-known response: %w", err)
	}

	if err := json.Unmarshal(body, &wellKnown); err != nil {
		return nil, fmt.Errorf("failed to parse well-known response: %w", err)
	}

	// Find the nodeinfo 2.0 URL
	var nodeinfoURL string
	for _, link := range wellKnown.Links {
		if link.Rel == "http://nodeinfo.diaspora.software/ns/schema/2.0" {
			nodeinfoURL = link.Href
			break
		}
	}

	if nodeinfoURL == "" {
		return nil, errors.New("nodeinfo 2.0 URL not found")
	}

	// Fetch the actual nodeinfo
	resp2, err := client.Get(nodeinfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch nodeinfo: %w", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from nodeinfo", resp2.StatusCode)
	}

	body2, err := io.ReadAll(resp2.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read nodeinfo response: %w", err)
	}

	var nodeinfo NodeInfoV2
	if err := json.Unmarshal(body2, &nodeinfo); err != nil {
		return nil, fmt.Errorf("failed to parse nodeinfo: %w", err)
	}

	return &nodeinfo, nil
}

// ExtractFederationUsername extracts the federation username from nodeinfo.
func ExtractFederationUsername(nodeinfo *NodeInfoV2) (string, error) {
	if nodeinfo == nil {
		return "", errors.New("nodeinfo is nil")
	}

	username := nodeinfo.Metadata.Federation.Username
	if username == "" {
		return "", errors.New("federation username not found in nodeinfo")
	}

	return username, nil
}

// ValidateOwncastServer validates if the server is an Owncast instance.
func ValidateOwncastServer(nodeinfo *NodeInfoV2) error {
	if nodeinfo == nil {
		return errors.New("nodeinfo is nil")
	}

	if nodeinfo.Software.Name != "owncast" {
		return fmt.Errorf("server is not an Owncast instance (software: %s)", nodeinfo.Software.Name)
	}

	// Check if ActivityPub is enabled
	hasActivityPub := false
	for _, protocol := range nodeinfo.Protocols {
		if protocol == "activitypub" {
			hasActivityPub = true
			break
		}
	}

	if !hasActivityPub {
		return errors.New("server does not support ActivityPub")
	}

	return nil
}
