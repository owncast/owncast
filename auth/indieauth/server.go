package indieauth

import (
	"github.com/pkg/errors"
	"github.com/teris-io/shortid"
)

// ServerAuthRequest is n inbound request to authenticate against
// this Owncast instance.
type ServerAuthRequest struct {
	ClientID    string
	RedirectURI string
	State       string
	Me          string
	Code        string
}

// ServerProfile represents basic user-provided data about this Owncast instance.
type ServerProfile struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Photo string `json:"photo"`
}

// ServerProfileResponse is returned when an auth flow requests the final
// confirmation of the IndieAuth flow.
type ServerProfileResponse struct {
	Me      string        `json:"me"`
	Profile ServerProfile `json:"profile"`
}

var pendingServerAuthRequests = map[string]ServerAuthRequest{}

// StartServerAuth will handle the authentication for the admin user of this
// Owncast server. Initiated via a GET of the auth endpoint.
// https://indieweb.org/authorization-endpoint
func StartServerAuth(clientID, redirectURI, state, me string) (*ServerAuthRequest, error) {
	code := shortid.MustGenerate()

	r := ServerAuthRequest{
		ClientID:    clientID,
		RedirectURI: redirectURI,
		State:       state,
		Me:          me,
		Code:        code,
	}

	pendingServerAuthRequests[code] = r

	return &r, nil
}

// CompleteServerAuth will verify that the values provided in the final step
// of the IndieAuth flow are correct, and return some basic profile info.
func CompleteServerAuth(code, redirectURI, clientID string) (*ServerProfileResponse, error) {
	request, pending := pendingServerAuthRequests[code]
	if !pending {
		return nil, errors.New("no pending authentication request")
	}

	if request.RedirectURI != redirectURI {
		return nil, errors.New("redirect URI does not match")
	}

	if request.ClientID != clientID {
		return nil, errors.New("client ID does not match")
	}

	response := ServerProfileResponse{
		Me: request.Me,
		Profile: ServerProfile{
			Name:  "name",
			URL:   "URL",
			Photo: "photo",
		},
	}

	return &response, nil
}
