package models

type Actor struct {
	Context           []string  `json:"@context"`
	Id                string    `json:"id"`
	Type              string    `json:"type"`
	PreferredUsername string    `json:"preferredUsername"`
	Name              string    `json:"name"`
	Summary           string    `json:"summary"`
	URL               string    `json:"url"`
	Inbox             string    `json:"inbox"`
	Outbox            string    `json:"outbox"`
	Following         string    `json:"following"`
	Followers         string    `json:"followers"`
	PublicKey         PublicKey `json:"publicKey"`
}

type PublicKey struct {
	Id           string `json:"id"`
	Owner        string `json:"owner"`
	PublicKeyPem string `json:"publicKeyPem"`
}
