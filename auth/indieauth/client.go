package indieauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/owncast/owncast/core/data"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var pendingAuthRequests = make(map[string]*Request)

// StartAuthFlow will begin the IndieAuth flow by generating an auth request.
func StartAuthFlow(authHost, userID, accessToken, displayName string) (*url.URL, error) {
	serverURL := data.GetServerURL()
	if serverURL == "" {
		return nil, errors.New("Owncast server URL must be set when using auth")
	}

	r, err := createAuthRequest(authHost, userID, displayName, accessToken, serverURL)
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate IndieAuth request")
	}

	pendingAuthRequests[r.State] = r

	return r.Redirect, nil
}

// HandleCallbackCode will handle the callback from the IndieAuth server
// to continue the next step of the auth flow.
func HandleCallbackCode(code, state string) (*Request, *Response, error) {
	request, exists := pendingAuthRequests[state]
	if !exists {
		return nil, nil, errors.New("no auth requests pending")
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", request.ClientID)
	data.Set("redirect_uri", request.Callback.String())
	data.Set("code_verifier", request.CodeVerifier)

	client := &http.Client{}
	r, err := http.NewRequest("POST", request.Endpoint.String(), strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		return nil, nil, err
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := client.Do(r)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, nil, errors.Wrap(err, "unable to parse IndieAuth response")
	}

	if response.Error != "" || response.ErrorDescription != "" {
		errorText := makeIndieAuthClientErrorText(response.Error)
		log.Debugln("IndieAuth error:", response.Error, response.ErrorDescription)
		return nil, nil, fmt.Errorf("IndieAuth error: %s - %s", errorText, response.ErrorDescription)
	}

	// In case this IndieAuth server does not use OAuth error keys or has internal
	// issues resulting in unstructured errors.
	if res.StatusCode < 200 || res.StatusCode > 299 {
		log.Debugln("IndieAuth error. status code:", res.StatusCode, "body:", string(body))
		return nil, nil, errors.New("there was an error authenticating against IndieAuth server")
	}

	// Trim any trailing slash so we can accurately compare the two "me" values
	meResponseVerifier := strings.TrimRight(response.Me, "/")
	meRequestVerifier := strings.TrimRight(request.Me.String(), "/")

	// What we sent and what we got back must match
	if meRequestVerifier != meResponseVerifier {
		return nil, nil, errors.New("indieauth response does not match the initial anticipated auth destination")
	}

	return request, &response, nil
}

// Error value should be from this list:
// https://datatracker.ietf.org/doc/html/rfc6749#section-5.2
func makeIndieAuthClientErrorText(err string) string {
	switch err {
	case "invalid_request", "invalid_client":
		return "The authentication request was invalid. Please report this to the Owncast project."
	case "invalid_grant", "unauthorized_client":
		return "This authorization request is unauthorized."
	case "unsupported_grant_type":
		return "The authorization grant type is not supported by the authorization server."
	default:
		return err
	}
}
