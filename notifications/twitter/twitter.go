package twitter

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dghubble/oauth1"
	"github.com/g8rswimmer/go-twitter/v2"
)

/*
1. developer.twitter.com. Apply to be a developer if needed.
2. Projects and apps -> Your project name
3. Settings.
4. Scroll down to"User authentication settings" Edit
5. Enable OAuth 1.0a with Read/Write permissions.
6. Fill out the form with your information. Callback can be anything.
7. Go to your project "Keys and tokens"
8. Generate API key and secret.
9. Generate access token and secret. Verify it says "Read and write permissions."
10. Generate bearer token.
*/

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

// Twitter is an instance of the Twitter notifier.
type Twitter struct {
	apiKey            string
	apiSecret         string
	accessToken       string
	accessTokenSecret string
	bearerToken       string
}

// New returns a new instance of the Twitter notifier.
func New(apiKey, apiSecret, accessToken, accessTokenSecret, bearerToken string) (*Twitter, error) {
	if apiKey == "" || apiSecret == "" || accessToken == "" || accessTokenSecret == "" || bearerToken == "" {
		return nil, errors.New("missing some or all of the required twitter configuration values")
	}

	return &Twitter{
		apiKey:            apiKey,
		apiSecret:         apiSecret,
		accessToken:       accessToken,
		accessTokenSecret: accessTokenSecret,
		bearerToken:       bearerToken,
	}, nil
}

// Notify will send a notification to Twitter with the supplied text.
func (t *Twitter) Notify(text string) error {
	config := oauth1.NewConfig(t.apiKey, t.apiSecret)
	token := oauth1.NewToken(t.accessToken, t.accessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	client := &twitter.Client{
		Authorizer: authorize{
			Token: t.bearerToken,
		},
		Client: httpClient,
		Host:   "https://api.twitter.com",
	}

	req := twitter.CreateTweetRequest{
		Text: text,
	}

	_, err := client.CreateTweet(context.Background(), req)
	return err
}
