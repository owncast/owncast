package models

// PluginTab is one viewer-page tab a plugin contributes via
// manifest.tabs. Slug identifies the source plugin (used as the React
// key on the frontend so unmount/remount only fires when a plugin
// goes away); Title is the label shown in the tab bar; HTML is the
// inlined body of the tab's content file. The host emits this list
// as `pluginTabs` on /api/config.
type PluginTab struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
	HTML  string `json:"html"`
}
