package models

// ExternalIntegration represents a single 3rd party integration that uses an access token.
type ExternalIntegration struct {
	Name        string
	AccessToken string
}
