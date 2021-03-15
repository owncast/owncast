package models

// ExternalAction is a link that will open as a 3rd party action.
type ExternalAction struct {
	// URL is the URL to load.
	URL string `json:"url"`
	// Title is the name of this action, displayed in the modal.
	Title string `json:"title"`
	// Description is the description of this action.
	Description string `json:"description"`
	// Icon is the optional icon for the button associated with this action.
	Icon string `json:"icon"`
	// Color is the optional color for the button associated with this action.
	Color string `json:"color"`
	// OpenExternally states if the action should open a new tab/window instead of an internal modal.
	OpenExternally bool `json:"openExternally"`
}
