package configrepository

// GetFederatedInboxMap is a mapping between account names and their outbox.
func (cr *SqlConfigRepository) GetFederatedInboxMap() map[string]string {
	return map[string]string{
		cr.GetDefaultFederationUsername(): cr.GetDefaultFederationUsername(),
	}
}

// GetDefaultFederationUsername will return the username used for sending federation activities.
func (cr *SqlConfigRepository) GetDefaultFederationUsername() string {
	return cr.GetFederationUsername()
}
