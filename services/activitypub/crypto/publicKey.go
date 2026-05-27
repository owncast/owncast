package crypto

import "net/url"

// PublicKey represents a public key with associated ownership.
type PublicKey struct {
	ID           *url.URL `json:"id"`
	Owner        *url.URL `json:"owner"`
	PublicKeyPem string   `json:"publicKeyPem"`
}
