package configrepository

// GetFederatedInboxMap is a mapping between account names and their outbox.
func (r SqlConfigRepository) GetFederatedInboxMap() map[string]string {
	return map[string]string{
		r.GetDefaultFederationUsername(): r.GetDefaultFederationUsername(),
	}
}

// GetDefaultFederationUsername will return the username used for sending federation activities.
func (r *SqlConfigRepository) GetDefaultFederationUsername() string {
	return r.GetFederationUsername()
}
