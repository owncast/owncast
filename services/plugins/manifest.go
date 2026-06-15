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

// CommandInfo is one chat command a plugin advertises for the unified !help.
// Reported by the SDK in register() output; purely informational (the plugin's
// own router does the actual matching).
type CommandInfo struct {
	Name        string   `json:"name"`
	Prefix      string   `json:"prefix,omitempty"`
	Description string   `json:"description,omitempty"`
	Usage       string   `json:"usage,omitempty"`
	Aliases     []string `json:"aliases,omitempty"`
	ModOnly     bool     `json:"modOnly,omitempty"`
}

type ConfigField struct {
	Type        string `json:"type"`
	Default     any    `json:"default,omitempty"`
	Description string `json:"description,omitempty"`
}

// Plugin runtimes, inferred at load time from the plugin's code file (not
// authored in the manifest): plugin.js → "javascript" and plugin.py → "python"
// run the author's source on the shared embedded interpreter engine, while
// plugin.wasm → "wasm" is a self-contained module the host loads directly.
const (
	RuntimeWasm       = "wasm"
	RuntimeJavaScript = "javascript"
	RuntimePython     = "python"
)

type Manifest struct {
	API string `json:"api"`
	// Type is the plugin's runtime, set by the host at load time from the code
	// artifact's filename (plugin.js/plugin.py/plugin.wasm) — NOT authored in
	// the manifest. The json tag is retained only so any stray legacy value
	// round-trips harmlessly; the loader always overrides it.
	Type string `json:"type,omitempty"`
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
	Slug          string        `json:"slug,omitempty"`
	Version       string        `json:"version"`
	Description   string        `json:"description,omitempty"`
	Bot           BotConfig     `json:"bot,omitempty"`
	Subscriptions Subscriptions `json:"subscriptions"`
	// Commands is the plugin's chat-command metadata, derived by the SDK from
	// its defineCommands/plugin.commands table and reported via register() (not
	// authored in the sidecar manifest). The host aggregates these across
	// plugins to answer a unified !help.
	Commands    []CommandInfo          `json:"commands,omitempty"`
	Permissions []string               `json:"permissions,omitempty"`
	Config      map[string]ConfigField `json:"config,omitempty"`
	Admin       AdminConfig            `json:"admin,omitempty"`
	Actions     []ActionButton         `json:"actions,omitempty"`
	Network     NetworkConfig          `json:"network,omitempty"`
	// Styles is a list of CSS files the plugin contributes to the
	// viewer page (the public streaming page). Each entry is a path
	// relative to the plugin's own URL namespace: bare paths like
	// "theme.css" and "/theme.css" auto-prefix to
	// /plugins/<slug>/theme.css, fully-qualified plugin paths like
	// "/plugins/<slug>/theme.css" pass through as-is, and paths in
	// any other plugin's namespace are rejected. CSS injects into
	// the global scope, so a plugin author can restyle anything the
	// viewer page renders. Requires the ui.modify permission and the
	// plugin must ship the referenced files (typically under
	// `assets/`) with http.serve declared.
	Styles []string `json:"styles,omitempty"`
	// Scripts is a list of JavaScript files the plugin contributes
	// to the viewer page. Each entry follows the same path rules as
	// Styles (bare/relative paths auto-prefix to /plugins/<slug>/,
	// cross-plugin and absolute http(s):// URLs are rejected) and
	// has to end in .js. Each entry renders as a <script src=...
	// async> tag on the viewer page, so the code runs in the same
	// global window context as the viewer chrome. Requires the
	// ui.modify permission (it's running inside Owncast's chrome)
	// and http.serve (the host serves the bytes).
	Scripts []string `json:"scripts,omitempty"`
	// ExtraPageContent declares the plugin's contribution to the viewer
	// page's extra-content block. When Content is set, the host reads
	// that HTML file's bytes and prepends them to the admin's
	// extraPageContent in the /api/config response (after the admin's
	// markdown is rendered) so plugin HTML never goes through the
	// markdown processor and can't be mangled by it. Paths follow the
	// same rules as styles/scripts (bare paths auto-prefix to
	// /plugins/<slug>/, cross-plugin and http(s):// URLs rejected,
	// .html extension required). When Content is absent, the host calls
	// on_page_content(slug, user) to get rendered HTML dynamically.
	// Requires ui.modify.
	ExtraPageContent *ExtraPageContent `json:"extraPageContent,omitempty"`
	// Tabs declares viewer-page tabs the plugin contributes to the
	// row of tabs Owncast renders next to chat (alongside built-ins
	// like Followers). Each entry's `content` is a relative path to
	// an HTML file under `assets/`; the host reads the bytes and
	// inlines them into the tab body on /api/config. Path and
	// extension rules match ExtraPageContent.
	Tabs []Tab `json:"tabs,omitempty"`
}

// Tab declares one viewer-page tab. Title is the label shown in the
// UI. Slug is the stable identifier passed to the plugin's
// onTabContent handler. Content is the optional path to a static HTML
// file under assets/; when omitted the host calls on_tab_content with
// the slug to get rendered HTML.
type Tab struct {
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	Content string `json:"content,omitempty"`
}

// ExtraPageContent declares the plugin's contribution to the viewer
// page's extra-content block. Slug names the target slot in the page;
// Content is an optional static HTML asset path. When Content is
// absent the host calls on_page_content(slug, user) to get rendered
// HTML dynamically.
type ExtraPageContent struct {
	Slug    string `json:"slug"`
	Content string `json:"content,omitempty"`
}

// UnmarshalJSON accepts both the current object form:
//
//	{"slug":"banner","content":"content.html"}
//
// and the legacy string form older SDK examples emitted:
//
//	"content.html"
//
// The host only consumes Content today, so older packages remain loadable
// even though they never carried an explicit slot slug.
func (e *ExtraPageContent) UnmarshalJSON(data []byte) error {
	var legacy string
	if err := json.Unmarshal(data, &legacy); err == nil {
		e.Content = legacy
		e.Slug = ""
		return nil
	}

	type extraPageContentAlias ExtraPageContent
	var decoded extraPageContentAlias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	*e = ExtraPageContent(decoded)
	return nil
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
	if err := m.validateConfigKeys(); err != nil {
		return err
	}
	if err := m.validateAdminPages(); err != nil {
		return err
	}
	if err := m.validateNetwork(); err != nil {
		return err
	}
	if err := m.validateActions(); err != nil {
		return err
	}
	if err := m.validateStyles(); err != nil {
		return err
	}
	if err := m.validateScripts(); err != nil {
		return err
	}
	if err := m.validateExtraPageContent(); err != nil {
		return err
	}
	return m.validateTabs()
}

// usesSharedEngine reports whether this plugin runs on the shared embedded
// interpreter engine (its code ships as source — plugin.js / plugin.py) rather
// than as a self-contained wasm module. Type is set at load time from the code
// artifact's filename, not authored in the manifest.
func (m *Manifest) usesSharedEngine() bool {
	return m.Type == RuntimeJavaScript || m.Type == RuntimePython
}

// validateConfigKeys rejects author config keys reserved by the host. Keys
// starting with "__" are used to inject per-instance state (e.g. "__slug") into
// shared-engine plugins via Extism config, and must not collide with declared
// plugin config.
func (m *Manifest) validateConfigKeys() error {
	for key := range m.Config {
		if strings.HasPrefix(key, "__") {
			return fmt.Errorf("manifest.config key %q is reserved (keys starting with \"__\" are used internally)", key)
		}
	}
	return nil
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

// validateStyles checks manifest.styles entries and rewrites them
// into absolute plugin-namespace URLs. Rules mirror manifest.actions
// (relative paths auto-prefix to /plugins/<slug>/, cross-plugin paths
// rejected, http.serve required since the host serves the bytes), but
// stricter on the external-URL front: http(s):// targets are rejected
// outright. Admins reviewing the manifest see every URL that will be
// injected into their viewer's global CSS scope, so the list should
// be exhaustive and self-contained. Plugins that need external assets
// (fonts, etc.) can @import or @font-face them from inside the
// bundled CSS, which is what the admin reviewed.
func (m *Manifest) validateStyles() error {
	if len(m.Styles) == 0 {
		return nil
	}
	if !m.hasPermission(PermUIModify) {
		return errors.New(
			"manifest.styles is set but the manifest does not declare " +
				"the \"ui.modify\" permission; plugins that inject CSS " +
				"into the viewer's global scope must opt in to ui.modify " +
				"so it's visible to anyone reviewing the manifest that the " +
				"plugin restyles Owncast's UI")
	}
	// http.serve is not required: the host reads each file from the plugin's
	// assets/ directory and inlines the bytes into customStyles on /api/config.
	for i, raw := range m.Styles {
		rewritten, err := rewritePluginAssetPath(m.Slug, raw, ".css")
		if err != nil {
			return fmt.Errorf("manifest.styles[%d]: %w", i, err)
		}
		m.Styles[i] = rewritten
	}
	return nil
}

// validateScripts checks manifest.scripts entries and rewrites them
// into absolute plugin-namespace URLs. Same rules as validateStyles
// applied to .js files: each is inlined into /customjavascript and
// runs in the viewer's window context. Only ui.modify is required;
// http.serve is not needed since the bytes come from assets/, not a URL.
func (m *Manifest) validateScripts() error {
	if len(m.Scripts) == 0 {
		return nil
	}
	if !m.hasPermission(PermUIModify) {
		return errors.New(
			"manifest.scripts is set but the manifest does not declare " +
				"the \"ui.modify\" permission; plugins that inject " +
				"JavaScript into the viewer page must opt in to ui.modify " +
				"so it's visible to anyone reviewing the manifest that the " +
				"plugin runs code inside Owncast's chrome")
	}
	// http.serve is not required: the host reads each file from the plugin's
	// assets/ directory and inlines the bytes into /customjavascript.
	for i, raw := range m.Scripts {
		rewritten, err := rewritePluginAssetPath(m.Slug, raw, ".js")
		if err != nil {
			return fmt.Errorf("manifest.scripts[%d]: %w", i, err)
		}
		m.Scripts[i] = rewritten
	}
	return nil
}

// validateExtraPageContent checks manifest.extraPageContent. Static content
// paths are rewritten into the plugin's namespace and inlined into the
// /api/config extraPageContent response (prepended to the admin's content),
// so http.serve is not required, but the same path-shape and extension rules
// apply for consistency with styles and scripts. Dynamic content omits
// Content and instead requires a valid slot slug for on_page_content.
func (m *Manifest) validateExtraPageContent() error {
	if m.ExtraPageContent == nil {
		return nil
	}
	if !m.hasPermission(PermUIModify) {
		return errors.New(
			"manifest.extraPageContent is set but the manifest does not " +
				"declare the \"ui.modify\" permission; plugins that inject " +
				"HTML into the viewer page must opt in to ui.modify so " +
				"it's visible to anyone reviewing the manifest that the " +
				"plugin paints inside Owncast's chrome")
	}
	if m.ExtraPageContent.Content != "" {
		if strings.TrimSpace(m.ExtraPageContent.Slug) != "" {
			if err := validatePluginSlug("manifest.extraPageContent.slug", m.ExtraPageContent.Slug); err != nil {
				return err
			}
		}
		rewritten, err := rewritePluginAssetPath(m.Slug, m.ExtraPageContent.Content, ".html")
		if err != nil {
			return fmt.Errorf("manifest.extraPageContent.content: %w", err)
		}
		m.ExtraPageContent.Content = rewritten
		return nil
	}
	if err := validatePluginSlug("manifest.extraPageContent.slug", m.ExtraPageContent.Slug); err != nil {
		return err
	}
	return nil
}

func (m *Manifest) validateTabs() error {
	if len(m.Tabs) == 0 {
		return nil
	}
	if !m.hasPermission(PermUIModify) {
		return errors.New(
			"manifest.tabs is set but the manifest does not declare " +
				"the \"ui.modify\" permission; plugins that add tabs " +
				"to the viewer page must opt in to ui.modify so it's " +
				"visible to anyone reviewing the manifest that the " +
				"plugin paints inside Owncast's chrome")
	}
	seenTitles := make(map[string]bool, len(m.Tabs))
	seenSlugs := make(map[string]bool, len(m.Tabs))
	for i := range m.Tabs {
		if strings.TrimSpace(m.Tabs[i].Title) == "" {
			return fmt.Errorf("manifest.tabs[%d].title is required", i)
		}
		if seenTitles[m.Tabs[i].Title] {
			return fmt.Errorf("manifest.tabs[%d].title %q is a duplicate; tab titles must be unique within a plugin", i, m.Tabs[i].Title)
		}
		seenTitles[m.Tabs[i].Title] = true
		if strings.TrimSpace(m.Tabs[i].Slug) == "" {
			if m.Tabs[i].Content == "" {
				return fmt.Errorf("manifest.tabs[%d].slug is required", i)
			}
			derived, err := slugify(m.Tabs[i].Title)
			if err != nil {
				return fmt.Errorf("manifest.tabs[%d].slug is required and could not be derived from title %q: %w", i, m.Tabs[i].Title, err)
			}
			m.Tabs[i].Slug = derived
		} else if err := validatePluginSlug(fmt.Sprintf("manifest.tabs[%d].slug", i), m.Tabs[i].Slug); err != nil {
			return err
		}
		if seenSlugs[m.Tabs[i].Slug] {
			return fmt.Errorf("manifest.tabs[%d].slug %q is a duplicate; tab slugs must be unique within a plugin", i, m.Tabs[i].Slug)
		}
		seenSlugs[m.Tabs[i].Slug] = true
		if m.Tabs[i].Content != "" {
			rewritten, err := rewritePluginAssetPath(m.Slug, m.Tabs[i].Content, ".html")
			if err != nil {
				return fmt.Errorf("manifest.tabs[%d].content: %w", i, err)
			}
			m.Tabs[i].Content = rewritten
		}
	}
	return nil
}

// validatePluginSlug checks that s is a non-empty, valid plugin-style
// slug (lowercase letters/digits/hyphens, starting with a letter, max
// 64 chars). field is used in error messages.
func validatePluginSlug(field, s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("%s is required", field)
	}
	if !slugPattern.MatchString(s) {
		return fmt.Errorf("%s %q must be lowercase letters, digits, and hyphens starting with a letter (max 64 chars)", field, s)
	}
	return nil
}

// rewritePluginAssetPath normalizes a single asset entry from a
// styles/scripts list into the plugin's /plugins/<slug>/<file>
// namespace and enforces the extension. Shared by validateStyles and
// validateScripts because the rules are identical apart from the file
// extension and the error context the caller wraps around it.
func rewritePluginAssetPath(slug, raw, requiredExt string) (string, error) {
	if strings.TrimSpace(raw) == "" {
		return "", errors.New("entry is empty")
	}
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return "", fmt.Errorf("cannot be an absolute URL (%q); bundle the file under assets/ and reference it by path", raw)
	}
	pluginPrefix := "/plugins/" + slug + "/"
	path := raw
	switch {
	case strings.HasPrefix(path, "/plugins/"):
		// Already-absolute plugin path: must be in our own namespace.
	case strings.HasPrefix(path, "/"):
		path = pluginPrefix + strings.TrimPrefix(path, "/")
	default:
		path = pluginPrefix + path
	}
	if !strings.HasPrefix(path, pluginPrefix) {
		return "", fmt.Errorf("points at another plugin's namespace: %s", path)
	}
	if !strings.HasSuffix(strings.ToLower(path), requiredExt) {
		return "", fmt.Errorf("must end in %s (got %q)", requiredExt, raw)
	}
	return path, nil
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
	// Version is intentionally not compared. It's informational metadata the
	// host gates nothing on, and the SDK bakes it into register() output from
	// the same manifest at build time, so a mismatch only ever indicates a
	// stale build rather than an identity or security problem. Identity (slug)
	// and security (permissions) are the checks that matter and are verified
	// directly below.
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
