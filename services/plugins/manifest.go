package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const SupportedAPIVersion = "1"

// slugPattern is the validation pattern for both explicit and
// auto-derived plugin slugs. The slug becomes the plugin's URL
// segment, on-disk filename, KV namespace, and registry identifier,
// so it has to stay narrow: lowercase letters/digits/hyphens, has to
// start with a letter, max 64 chars.
var slugPattern = regexp.MustCompile(`^[a-z][a-z0-9-]{0,63}$`)

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
	API string `json:"api"`
	// DisplayName is the user-facing plugin name shown in admin lists,
	// the Browse cards on the registry, and (by default) as the
	// in-chat bot identity. Mapped from the JSON `name` field for
	// author ergonomics (authors write `"name": "Awesome Echo Bot"`
	// in their manifest; the Go side treats it as a display string).
	DisplayName string `json:"name"`
	// Slug is the plugin's canonical identifier: URL segment, KV
	// namespace, on-disk filename, registry primary key. Lowercase,
	// hyphenated, matches slugPattern. Authors can pin it via the
	// JSON `slug` field; if empty, it's auto-derived from
	// DisplayName at parse time by slugify.
	Slug          string                 `json:"slug,omitempty"`
	Version       string                 `json:"version"`
	Description   string                 `json:"description,omitempty"`
	Bot           BotConfig              `json:"bot,omitempty"`
	Subscriptions Subscriptions          `json:"subscriptions"`
	Permissions   []string               `json:"permissions,omitempty"`
	Config        map[string]ConfigField `json:"config,omitempty"`
	Admin         AdminConfig            `json:"admin,omitempty"`
	Actions       []ActionButton         `json:"actions,omitempty"`
	Network       NetworkConfig          `json:"network,omitempty"`
}

// BotConfig is the chat-bot-specific configuration for plugins that
// post to chat. Optional; defaults to the plugin's DisplayName when
// unset. Wrapped in a struct so future fields (avatar URL, name
// color) can land here without flat-namespace pollution on the
// manifest.
type BotConfig struct {
	// DisplayName is what viewers see when the plugin posts to chat.
	// Falls back to Manifest.DisplayName at the call site
	// (ChatDisplayName()).
	DisplayName string `json:"displayName,omitempty"`
}

// ChatDisplayName resolves the name the chat bot should post as.
// Bot.DisplayName wins when set, otherwise the plugin's DisplayName.
// Always non-empty post-Validate because DisplayName is required.
func (m *Manifest) ChatDisplayName() string {
	if m.Bot.DisplayName != "" {
		return m.Bot.DisplayName
	}
	return m.DisplayName
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
	if m.DisplayName == "" {
		return errors.New("manifest.name is required")
	}
	if m.Version == "" {
		return errors.New("manifest.version is required")
	}
	if err := m.resolveSlug(); err != nil {
		return err
	}
	if err := m.validateAdminPages(); err != nil {
		return err
	}
	if err := m.validateNetwork(); err != nil {
		return err
	}
	return m.validateActions()
}

// resolveSlug fills in m.Slug when the author didn't pin one
// explicitly, and validates the resulting slug against slugPattern.
// Called from Validate, so every Manifest returned by ParseManifest
// has Slug populated and well-formed.
func (m *Manifest) resolveSlug() error {
	if m.Slug == "" {
		derived, err := slugify(m.DisplayName)
		if err != nil {
			return fmt.Errorf("could not auto-generate a slug from manifest.name %q: %w; set manifest.slug explicitly", m.DisplayName, err)
		}
		m.Slug = derived
	}
	if !slugPattern.MatchString(m.Slug) {
		return fmt.Errorf("manifest.slug %q is invalid; must match %s", m.Slug, slugPattern.String())
	}
	return nil
}

// slugify turns a free-form display name into a URL-safe slug
// matching slugPattern. The pipeline: ASCII-lowercase letters and
// digits pass through; every other rune (whitespace, punctuation,
// non-ASCII) collapses to a single hyphen; leading and trailing
// hyphens are trimmed. Non-ASCII names degrade noisily (e.g.
// "Café Helper" becomes "caf-helper"), so authors with diacritics
// or non-Latin scripts should pin `slug` explicitly in the
// manifest instead of relying on derivation.
//
// Returns an error when the result is empty (the input had no
// usable characters) or doesn't start with a letter (e.g.
// "123 Plugin" produces "123-plugin", which fails the
// start-with-letter rule); in either case the author has to pin
// slug explicitly.
func slugify(input string) (string, error) {
	var sb strings.Builder
	prevHyphen := false
	for _, r := range input {
		switch {
		case r >= 'a' && r <= 'z':
			sb.WriteRune(r)
			prevHyphen = false
		case r >= 'A' && r <= 'Z':
			sb.WriteRune(r + ('a' - 'A'))
			prevHyphen = false
		case r >= '0' && r <= '9':
			sb.WriteRune(r)
			prevHyphen = false
		default:
			if !prevHyphen && sb.Len() > 0 {
				sb.WriteRune('-')
				prevHyphen = true
			}
		}
	}
	out := strings.TrimRight(sb.String(), "-")
	if out == "" {
		return "", errors.New("slugified value is empty")
	}
	if !slugPattern.MatchString(out) {
		return "", fmt.Errorf("slugified value %q does not match the required pattern", out)
	}
	return out, nil
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
	// Action buttons are the only UI surface a plugin can place inside
	// Owncast's own chrome (the viewer action bar). Self-contained admin
	// pages and static content served under /plugins/<name>/ are baseline
	// plugin functionality and don't gate on this.
	if len(m.Actions) > 0 && !m.hasPermission(PermUIModify) {
		return errors.New(
			"manifest.actions is set but the manifest does not declare " +
				"the \"ui.modify\" permission; plugins that contribute viewer " +
				"action buttons must opt in to ui.modify so it's visible to " +
				"anyone reviewing the manifest that the plugin places UI " +
				"inside Owncast's chrome")
	}
	pluginPrefix := "/plugins/" + m.Slug + "/"
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
	for i := range m.Actions {
		rewritten, err := rewriteActionIcon(pluginPrefix, hasHTTPServe, m.Actions[i].Icon)
		if err != nil {
			return fmt.Errorf("manifest.actions[%d].icon: %w", i, err)
		}
		m.Actions[i].Icon = rewritten
	}
	return nil
}

// rewriteActionIcon applies the same path-handling rules to a button's
// icon URL as we do to the button's url: a same-origin relative path is
// rewritten into this plugin's namespace; an http(s) URL is left alone;
// a cross-plugin path is rejected. Empty input passes through (icons
// are optional).
func rewriteActionIcon(pluginPrefix string, hasHTTPServe bool, icon string) (string, error) {
	if icon == "" {
		return "", nil
	}
	if strings.HasPrefix(icon, "http://") || strings.HasPrefix(icon, "https://") {
		return icon, nil
	}
	if strings.HasPrefix(icon, "/") && !strings.HasPrefix(icon, "/plugins/") {
		icon = pluginPrefix + strings.TrimPrefix(icon, "/")
	}
	if strings.HasPrefix(icon, pluginPrefix) && !hasHTTPServe {
		return "", fmt.Errorf("targets this plugin (%s) but http.serve permission is not declared", icon)
	}
	if strings.HasPrefix(icon, "/plugins/") && !strings.HasPrefix(icon, pluginPrefix) {
		return "", fmt.Errorf("points at another plugin's namespace: %s", icon)
	}
	return icon, nil
}

// AgreesWith reports whether the runtime registration `other` is consistent
// with the sidecar manifest. The sidecar declares identity and permissions;
// the runtime must not exceed declared permissions. Subscriptions are derived
// by the SDK at runtime, so they aren't validated here.
//
// Identity is checked on Slug (the canonical identifier), not DisplayName:
// the runtime side may not echo DisplayName back identically (the SDK doesn't
// have to), but Slug is mechanically derived/declared the same way on both
// sides, so it's the reliable identity column. When the runtime side ships
// a Slug field, the two slugs are compared directly. When it doesn't (older
// SDK that only emits DisplayName), the slug is derived from the runtime
// DisplayName the same way ParseManifest derives it on the sidecar side
// and the comparison falls through to the same result.
func (m *Manifest) AgreesWith(other *Manifest) error {
	// Resolve the slug on both sides: a sidecar that didn't go through
	// ParseManifest (e.g. a test fixture) won't have Slug populated, but
	// it should still produce the same identity from its DisplayName.
	resolveSlug := func(x *Manifest) string {
		if x.Slug != "" {
			return x.Slug
		}
		if x.DisplayName == "" {
			return ""
		}
		derived, err := slugify(x.DisplayName)
		if err != nil {
			return ""
		}
		return derived
	}
	mySlug := resolveSlug(m)
	otherSlug := resolveSlug(other)
	if mySlug != otherSlug {
		return fmt.Errorf("slug mismatch: manifest=%q register=%q", mySlug, otherSlug)
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
