package data

// GetFederatedInboxMap is a mapping between account names and their outbox.
func (s *Service) GetFederatedInboxMap() map[string]string {
	return map[string]string{
		s.GetDefaultFederationUsername(): s.GetDefaultFederationUsername(),
	}
}

// GetDefaultFederationUsername will return the username used for sending federation activities.
func (s *Service) GetDefaultFederationUsername() string {
	return s.GetFederationUsername()
}
