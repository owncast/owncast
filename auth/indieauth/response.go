package indieauth

// Profile represents optional profile data that is returned
// when completing the IndieAuth flow.
type Profile struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Photo string `json:"photo"`
}

// Response the response returned when completing
// the IndieAuth flow.
type Response struct {
	Me      string  `json:"me"`
	Profile Profile `json:"profile"`
}
