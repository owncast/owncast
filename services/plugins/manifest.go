package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const SupportedAPIVersion = "1"

type Subscription struct {
	Event    string `json:"event"`
	Priority int    `json:"priority,omitempty"`
}

type Subscriptions struct {
	Notify []Subscription `json:"notify,omitempty"`
	Filter []Subscription `json:"filter,omitempty"`
}

type ConfigField struct {
	Type        string `json:"type"`
	Default     any    `json:"default,omitempty"`
	Description string `json:"description,omitempty"`
}

type Manifest struct {
	API           string                 `json:"api"`
	Name          string                 `json:"name"`
	Version       string                 `json:"version"`
	Description   string                 `json:"description,omitempty"`
	Subscriptions Subscriptions          `json:"subscriptions"`
	Permissions   []string               `json:"permissions,omitempty"`
	Config        map[string]ConfigField `json:"config,omitempty"`
	Admin         AdminConfig            `json:"admin,omitempty"`
	Actions       []ActionButton         `json:"actions,omitempty"`
	Network       NetworkConfig          `json:"network,omitempty"`
}

// NetworkConfig narrows what hosts a plugin with the `network.fetch`
// permission can reach. The host passes AllowedHosts straight through to
// extism's manifest AllowedHosts; each entry is a hostname glob (e.g.
// "api.discord.com", "*.weather.com", "*").
//
// Plugins that declare `network.fetch` MUST declare a non-empty
// AllowedHosts list — the wildcard "*" is allowed but has to be written
// out so admins reviewing the manifest see the scope they're granting.
type NetworkConfig struct {
	AllowedHosts []string `json:"allowedHosts,omitempty"`
}

// ActionButton declares an entry the Owncast UI surfaces as an external
// action — a clickable button that loads a URL (in a modal or new tab) or
// shows raw HTML when pressed. Buttons declared here are merged with the
// admin-configured external actions while the plugin is enabled; when the
// plugin is disabled they disappear.
//
// Shape matches Owncast's existing ExternalAction. Exactly one of Url or
// Html must be set.
//
// Url ergonomics: if Url starts with "/" but not "/plugins/", it's treated
// as a relative path inside this plugin's namespace and the host rewrites
// it to "/plugins/<name><url>" at validation time. Absolute http(s) URLs
// and explicit "/plugins/<name>/..." paths are accepted as-is.
type ActionButton struct {
	Title          string `json:"title"`
	Url            string `json:"url,omitempty"`
	Html           string `json:"html,omitempty"`
	Icon           string `json:"icon,omitempty"`
	Color          string `json:"color,omitempty"`
	Description    string `json:"description,omitempty"`
	OpenExternally bool   `json:"openExternally,omitempty"`
}

// AdminConfig declares admin-only surfaces a plugin exposes. The Owncast
// admin web UI lists these in the "Plugins" section; each page renders the
// plugin's content at /plugins/<name>/<path>. Paths declared here are
// auth-gated by the host — unauthenticated requests get 401 before
// reaching the plugin.
type AdminConfig struct {
	Pages []AdminPage `json:"pages,omitempty"`
}

type AdminPage struct {
	Title string `json:"title"`
	// Path is a glob (e.g. "/admin", "/admin/*"). Requests under
	// /plugins/<name>/<rest> are checked against each glob and require
	// admin authentication when any match.
	Path string `json:"path"`
	Icon string `json:"icon,omitempty"`
}

func ParseManifest(b []byte) (*Manifest, error) {
	var m Manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("parse manifest json: %w", err)
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *Manifest) Validate() error {
	if m.API != SupportedAPIVersion {
		return fmt.Errorf("unsupported api version %q (host supports %q)", m.API, SupportedAPIVersion)
	}
	if m.Name == "" {
		return errors.New("manifest.name is required")
	}
	if m.Version == "" {
		return errors.New("manifest.version is required")
	}
	if err := m.validateAdminPages(); err != nil {
		return err
	}
	if err := m.validateNetwork(); err != nil {
		return err
	}
	return m.validateActions()
}

// hasPermission reports whether the manifest declares the given permission.
func (m *Manifest) hasPermission(perm string) bool {
	for _, p := range m.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}

func (m *Manifest) validateAdminPages() error {
	for i, page := range m.Admin.Pages {
		if page.Title == "" {
			return fmt.Errorf("manifest.admin.pages[%d].title is required", i)
		}
		if page.Path == "" {
			return fmt.Errorf("manifest.admin.pages[%d].path is required", i)
		}
	}
	return nil
}

func (m *Manifest) validateNetwork() error {
	if m.hasPermission(PermNetworkFetch) && len(m.Network.AllowedHosts) == 0 {
		return errors.New(
			"manifest declares network.fetch but no network.allowedHosts; " +
				"list the hostnames you'll reach (globs OK, e.g. \"api.discord.com\", " +
				"\"*.weather.com\") or [\"*\"] for any host")
	}
	for i, host := range m.Network.AllowedHosts {
		if strings.TrimSpace(host) == "" {
			return fmt.Errorf("manifest.network.allowedHosts[%d] is empty", i)
		}
	}
	return nil
}

func (m *Manifest) validateActions() error {
	pluginPrefix := "/plugins/" + m.Name + "/"
	hasHTTPServe := m.hasPermission(PermHttpServe)
	for i := range m.Actions {
		a := &m.Actions[i]
		if a.Title == "" {
			return fmt.Errorf("manifest.actions[%d].title is required", i)
		}
		hasURL, hasHTML := a.Url != "", a.Html != ""
		if hasURL == hasHTML {
			return fmt.Errorf("manifest.actions[%d]: exactly one of url or html is required", i)
		}
		if !hasURL {
			continue
		}
		// Relative path inside the plugin's own namespace? Rewrite.
		if strings.HasPrefix(a.Url, "/") && !strings.HasPrefix(a.Url, "/plugins/") {
			a.Url = pluginPrefix + strings.TrimPrefix(a.Url, "/")
		}
		// http.serve required when the action points back into the plugin.
		if strings.HasPrefix(a.Url, pluginPrefix) && !hasHTTPServe {
			return fmt.Errorf("manifest.actions[%d].url targets this plugin (%s) but http.serve permission is not declared",
				i, a.Url)
		}
		// Paths in other plugins' namespaces aren't allowed — catches typos
		// and prevents one plugin from advertising another's UI.
		if strings.HasPrefix(a.Url, "/plugins/") && !strings.HasPrefix(a.Url, pluginPrefix) {
			return fmt.Errorf("manifest.actions[%d].url points at another plugin's namespace: %s", i, a.Url)
		}
	}
	return nil
}

// AgreesWith reports whether the runtime registration `other` is consistent
// with the sidecar manifest. The sidecar declares identity and permissions;
// the runtime must not exceed declared permissions. Subscriptions are derived
// by the SDK at runtime, so they aren't validated here.
func (m *Manifest) AgreesWith(other *Manifest) error {
	if m.Name != other.Name {
		return fmt.Errorf("name mismatch: manifest=%q register=%q", m.Name, other.Name)
	}
	if m.Version != other.Version {
		return fmt.Errorf("version mismatch: manifest=%q register=%q", m.Version, other.Version)
	}
	declared := stringSet(m.Permissions)
	for _, p := range other.Permissions {
		if !declared[p] {
			return fmt.Errorf("plugin requested permission %q at runtime not declared in manifest", p)
		}
	}
	return nil
}

func stringSet(items []string) map[string]bool {
	out := make(map[string]bool, len(items))
	for _, i := range items {
		out[i] = true
	}
	return out
}
