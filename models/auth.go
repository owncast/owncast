package models

// Type represents a form of authentication.
type AuthType string

// The different auth types we support.
const (
	// IndieAuth https://indieauth.spec.indieweb.org/.
	IndieAuth AuthType = "indieauth"
	Fediverse AuthType = "fediverse"
	// PluginAuth is an identity established by a viewer-auth plugin
	// (owncast.users.register). type stays the coarse namespace; the auth row's
	// provider column carries the plugin slug and auth_key the raw external id.
	PluginAuth AuthType = "plugin.auth"
)

// LinkedIdentityFields carries the optional public-profile metadata for a
// linked identity (see the auth table's provider/profile_url/handle/is_public
// columns). All fields are optional. Provider overrides the default provider
// (the auth type) — set it to the plugin slug for plugin auth.
type LinkedIdentityFields struct {
	// Provider is the first-class provider id. Empty defaults to the auth type
	// (correct for built-in IndieAuth/Fediverse); set to the plugin slug for
	// plugin auth.
	Provider string
	// ProfileURL is the public, clickable link to the external profile.
	ProfileURL string
	// Handle is the human label for the identity, e.g. @me@host.
	Handle string
	// Public is the user's consent to surface this identity publicly. Defaults
	// to false; nothing shows publicly until the user opts in.
	Public bool
}
