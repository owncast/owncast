package configrepository

// GetPublicKey will return the public key.
func (cr *SqlConfigRepository) GetPublicKey() string {
	value, _ := cr.datastore.GetString(publicKeyKey)
	return value
}

// SetPublicKey will save the public key.
func (cr *SqlConfigRepository) SetPublicKey(key string) error {
	return cr.datastore.SetString(publicKeyKey, key)
}

// GetPrivateKey will return the private key.
func (cr *SqlConfigRepository) GetPrivateKey() string {
	value, _ := cr.datastore.GetString(privateKeyKey)
	return value
}

// SetPrivateKey will save the private key.
func (cr *SqlConfigRepository) SetPrivateKey(key string) error {
	return cr.datastore.SetString(privateKeyKey, key)
}
