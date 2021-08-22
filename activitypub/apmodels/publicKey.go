package apmodels

import "net/url"

type PublicKey struct {
	Id           *url.URL `json:"id"`
	Owner        *url.URL `json:"owner"`
	PublicKeyPem string   `json:"publicKeyPem"`
}
