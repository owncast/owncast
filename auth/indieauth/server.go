package indieauth

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/teris-io/shortid"
)

// ServerAuthRequest is n inbound request to authenticate against
// this Owncast instance.
type ServerAuthRequest struct {
	Timestamp     time.Time
	ClientID      string
	RedirectURI   string
	CodeChallenge string
	State         string
	Me            string
	Code          string
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
	Me      string        `json:"me,omitempty"`
	Profile ServerProfile `json:"profile,omitempty"`
	// Error keys need to match the OAuth spec.
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

const maxPendingRequests = 100

const (
	schemeHTTP  = "http"
	schemeHTTPS = "https"
)

// StartServerAuth will handle the authentication for the admin user of this
// Owncast server. Initiated via a GET of the auth endpoint.
// https://indieweb.org/authorization-endpoint
func (s *Service) StartServerAuth(clientID, redirectURI, codeChallenge, state, me string) (*ServerAuthRequest, error) {
	// Validate the redirect URI against the client ID before the auth
	// endpoint hands it to the browser. Without this the endpoint is an open
	// redirect: a crafted redirect_uri would send the generated auth code to
	// an arbitrary origin.
	if err := validateClientRedirect(clientID, redirectURI); err != nil {
		return nil, err
	}

	s.pendingServerAuthRequestsLock.Lock()
	defer s.pendingServerAuthRequestsLock.Unlock()

	if len(s.pendingServerAuthRequests)+1 >= maxPendingRequests {
		return nil, errors.New("Please try again later. Too many pending requests.")
	}

	code := shortid.MustGenerate()

	r := ServerAuthRequest{
		ClientID:      clientID,
		RedirectURI:   redirectURI,
		CodeChallenge: codeChallenge,
		State:         state,
		Me:            me,
		Code:          code,
		Timestamp:     time.Now(),
	}

	s.pendingServerAuthRequests[code] = r

	return &r, nil
}

// validateClientRedirect enforces the IndieAuth requirement that the
// redirect_uri belong to the requesting client. Both values must be
// http(s) URLs with a host, and — because Owncast does not fetch the
// client's published metadata to discover alternate redirect URIs — the
// redirect host must match the client_id host. This prevents the
// authorization endpoint from being used as an open redirect.
func validateClientRedirect(clientID, redirectURI string) error {
	client, err := url.Parse(clientID)
	if err != nil || !isWebURL(client) {
		return errors.New("invalid client_id")
	}

	redirect, err := url.Parse(redirectURI)
	if err != nil || !isWebURL(redirect) {
		return errors.New("invalid redirect_uri")
	}

	if !strings.EqualFold(client.Hostname(), redirect.Hostname()) {
		return errors.New("redirect_uri host does not match client_id host")
	}

	return nil
}

// isWebURL reports whether u is an absolute http(s) URL with a non-empty
// host. Requiring the scheme and host rejects opaque/non-hierarchical URIs
// such as "javascript:..." or "data:..." — those have an empty hostname, so
// two of them would otherwise pass the host-match check and slip a hostile
// target into the redirect Location header.
func isWebURL(u *url.URL) bool {
	if u == nil {
		return false
	}

	scheme := strings.ToLower(u.Scheme)
	if scheme != schemeHTTP && scheme != schemeHTTPS {
		return false
	}

	return u.Hostname() != ""
}

// CompleteServerAuth will verify that the values provided in the final step
// of the IndieAuth flow are correct, and return some basic profile info.
func (s *Service) CompleteServerAuth(code, redirectURI, clientID string, codeVerifier string) (*ServerProfileResponse, error) {
	s.pendingServerAuthRequestsLock.Lock()
	request, pending := s.pendingServerAuthRequests[code]
	s.pendingServerAuthRequestsLock.Unlock()
	if !pending {
		return nil, errors.New("no pending authentication request")
	}

	if request.RedirectURI != redirectURI {
		return nil, errors.New("redirect URI does not match")
	}

	if request.ClientID != clientID {
		return nil, errors.New("client ID does not match")
	}

	codeChallengeFromRequest := createCodeChallenge(codeVerifier)
	if request.CodeChallenge != codeChallengeFromRequest {
		return nil, errors.New("code verifier is incorrect")
	}

	response := ServerProfileResponse{
		Me: s.configRepository.GetServerURL(),
		Profile: ServerProfile{
			Name:  s.configRepository.GetServerName(),
			URL:   s.configRepository.GetServerURL(),
			Photo: fmt.Sprintf("%s/%s", s.configRepository.GetServerURL(), s.configRepository.GetLogoPath()),
		},
	}

	return &response, nil
}
