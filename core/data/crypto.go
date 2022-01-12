package data

// GetPublicKey will return the public key.
func GetPublicKey() string {
	value, _ := _datastore.GetString(publicKeyKey)
	return value
}

// SetPublicKey will save the public key.
func SetPublicKey(key string) error {
	return _datastore.SetString(publicKeyKey, key)
}

// GetPrivateKey will return the private key.
func GetPrivateKey() string {
	value, _ := _datastore.GetString(privateKeyKey)
	return value
}

// SetPrivateKey will save the private key.
func SetPrivateKey(key string) error {
	return _datastore.SetString(privateKeyKey, key)
}
