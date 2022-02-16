package indieauth

import "net/url"

type IndieAuthRequest struct {
	UserID             string
	DisplayName        string
	CurrentAccessToken string
	Endpoint           *url.URL
	Redirect           *url.URL // Outbound redirect URL to continue auth flow
	Callback           *url.URL // Inbound URL to get auth flow results
	ClientID           string
	CodeVerifier       string
	CodeChallenge      string
	State              string
	Me                 *url.URL
}
