package models

// PluginStyleInfo describes the viewer-page styling one enabled plugin
// contributes. The admin Appearance UI uses it to explain that plugin
// styles and the admin's own appearance settings are combined into the
// final look, and to flag the specific color swatches a plugin also
// sets. DeclaredVars holds the theme custom properties the plugin
// declares (without the leading `--`, e.g. "theme-color-action"); it
// can be empty when a plugin styles the page without touching a
// recognized appearance token. The host emits this list as
// `styleContributors` on the admin server config.
type PluginStyleInfo struct {
	Slug         string   `json:"slug"`
	Name         string   `json:"name"`
	DeclaredVars []string `json:"declaredVars"`
}
