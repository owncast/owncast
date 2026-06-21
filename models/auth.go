package models

// Type represents a form of authentication.
type AuthType string

// The different auth types we support.
const (
	// IndieAuth https://indieauth.spec.indieweb.org/.
	IndieAuth AuthType = "indieauth"
	Fediverse AuthType = "fediverse"
	// PluginAuth is an identity established by a viewer-auth plugin
	// (owncast.users.register). The stored auth token is the plugin's external
	// identity, namespaced by the plugin slug (e.g. "github-auth:github:583231").
	PluginAuth AuthType = "plugin.auth"
)
