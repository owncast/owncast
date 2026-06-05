package models

// PluginTab is one viewer-page tab a plugin contributes via
// manifest.tabs. Slug is the composite unique key for this tab
// (pluginSlug/tabSlug) used as the React key on the frontend so
// unmount/remount only fires when a tab's source changes.
// PluginSlug identifies the source plugin for filtering.
// Title is the label shown in the tab bar; HTML is the inlined body.
// The host emits this list as `pluginTabs` on /api/config.
type PluginTab struct {
	Slug       string `json:"slug"`
	PluginSlug string `json:"pluginSlug"`
	Title      string `json:"title"`
	HTML       string `json:"html"`
}
