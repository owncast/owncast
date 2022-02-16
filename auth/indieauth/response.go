package indieauth

type IndieAuthProfile struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Photo string `json:"photo"`
}

type IndieAuthResponse struct {
	Me      string           `json:"me"`
	Profile IndieAuthProfile `json:"profile"`
}
