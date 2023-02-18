package data

// GetPublicKey will return the public key.
func (s *Service) GetPublicKey() string {
	value, _ := s.Store.GetString(publicKeyKey)
	return value
}

// SetPublicKey will save the public key.
func (s *Service) SetPublicKey(key string) error {
	return s.Store.SetString(publicKeyKey, key)
}

// GetPrivateKey will return the private key.
func (s *Service) GetPrivateKey() string {
	value, _ := s.Store.GetString(privateKeyKey)
	return value
}

// SetPrivateKey will save the private key.
func (s *Service) SetPrivateKey(key string) error {
	return s.Store.SetString(privateKeyKey, key)
}
