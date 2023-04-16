package indieauth

import (
	"net/url"
	"time"
)

// Request represents a single in-flight IndieAuth request.
type Request struct {
	Timestamp          time.Time
	Endpoint           *url.URL
	Redirect           *url.URL // Outbound redirect URL to continue auth flow
	Callback           *url.URL // Inbound URL to get auth flow results
	Me                 *url.URL
	UserID             string
	DisplayName        string
	CurrentAccessToken string
	ClientID           string
	CodeVerifier       string
	CodeChallenge      string
	State              string
}
