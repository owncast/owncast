package configrepository

// GetPublicKey will return the public key.
func (r *SqlConfigRepository) GetPublicKey() string {
	value, _ := r.datastore.GetString(publicKeyKey)
	return value
}

// SetPublicKey will save the public key.
func (r *SqlConfigRepository) SetPublicKey(key string) error {
	return r.datastore.SetString(publicKeyKey, key)
}

// GetPrivateKey will return the private key.
func (r *SqlConfigRepository) GetPrivateKey() string {
	value, _ := r.datastore.GetString(privateKeyKey)
	return value
}

// SetPrivateKey will save the private key.
func (r *SqlConfigRepository) SetPrivateKey(key string) error {
	return r.datastore.SetString(privateKeyKey, key)
}
