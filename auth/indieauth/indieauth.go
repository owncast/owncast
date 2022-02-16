package indieauth

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/owncast/owncast/core/data"
	"github.com/pkg/errors"
)

var pendingAuthRequests = make(map[string]*IndieAuthRequest)

// StartAuthFlow will begin the IndieAuth flow by generating an auth request.
func StartAuthFlow(authHost, userID, accessToken, displayName string) (*url.URL, error) {
	serverURL := data.GetServerURL()
	if serverURL == "" {
		return nil, errors.New("server URL must be set to use auth")
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
func HandleCallbackCode(code, state string) (*IndieAuthRequest, *IndieAuthResponse, error) {
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

	var response IndieAuthResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, nil, err
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
