package data

// This is a mapping between account names and their outbox
func GetFederatedInboxMap() map[string]string {
	return map[string]string{
		GetDefaultFederationUsername(): GetDefaultFederationUsername(),
	}
}

// GetDefaultFederationUsername will return the username used for sending federation activities.
func GetDefaultFederationUsername() string {
	return "owncast-ap-test"
}
