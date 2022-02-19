package indieauth

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"github.com/andybalholm/cascadia"
	"github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"golang.org/x/net/html"
)

func createAuthRequest(authDestination, userID, displayName, accessToken, baseServer string) (*IndieAuthRequest, error) {
	authURL, err := url.Parse(authDestination)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse IndieAuth destination")
	}

	authEndpointURL, err := getAuthEndpointFromURL(authURL.String())
	if err != nil {
		return nil, errors.Wrap(err, "unable to get IndieAuth endpoint from destination URL")
	}

	baseServerURL, err := url.Parse(baseServer)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse local owncast base server URL")
	}

	callbackURL := *baseServerURL
	callbackURL.Path = "/api/auth/indieauth/callback"

	codeChallenge := createCodeChallenge()
	state := shortid.MustGenerate()
	responseType := "code"
	clientID := baseServerURL.String() // Our local URL
	codeChallengeMethod := "S256"

	redirect := *authEndpointURL

	q := authURL.Query()
	q.Add("response_type", responseType)
	q.Add("client_id", clientID)
	q.Add("state", state)
	q.Add("code_challenge_method", codeChallengeMethod)
	q.Add("code_challenge", codeChallenge)
	q.Add("me", authURL.String())
	q.Add("redirect_uri", callbackURL.String())
	redirect.RawQuery = q.Encode()

	return &IndieAuthRequest{
		Me:                 authURL,
		UserID:             userID,
		DisplayName:        displayName,
		CurrentAccessToken: accessToken,
		Endpoint:           authEndpointURL,
		ClientID:           baseServer,
		CodeVerifier:       shortid.MustGenerate(),
		CodeChallenge:      codeChallenge,
		State:              state,
		Redirect:           &redirect,
		Callback:           &callbackURL,
	}, nil
}

func getAuthEndpointFromURL(urlstring string) (*url.URL, error) {
	r, err := http.Get(urlstring) // nolint:gosec
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	doc, err := html.Parse(r.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse html at remote auth host")
	}
	authorizationEndpointTag := cascadia.MustCompile("link[rel=authorization_endpoint]").MatchAll(doc)
	if len(authorizationEndpointTag) == 0 {
		return nil, fmt.Errorf("url does not support indieauth")
	}

	for _, attr := range authorizationEndpointTag[len(authorizationEndpointTag)-1].Attr {
		if attr.Key == "href" {
			return url.Parse(attr.Val)
		}
	}

	return nil, fmt.Errorf("unable to find href value for authorization_endpoint")
}

func createCodeChallenge() string {
	codeVerifier := shortid.MustGenerate()
	hasher := sha1.New() // nolint:gosec
	hasher.Write([]byte(codeVerifier))
	encodedHashedCode := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	return encodedHashedCode
}
