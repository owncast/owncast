package models

// Type represents a form of authentication.
type AuthType string

// The different auth types we support.
const (
	// IndieAuth https://indieauth.spec.indieweb.org/.
	IndieAuth AuthType = "indieauth"
	Fediverse AuthType = "fediverse"
)
