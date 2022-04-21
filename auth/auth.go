package auth

// Type represents a form of authentication.
type Type string

// The different auth types we support.
// Currently only IndieAuth.
const (
	// IndieAuth https://indieauth.spec.indieweb.org/.
	IndieAuth Type = "indieauth"
)
