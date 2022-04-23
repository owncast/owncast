package apmodels

import (
	"fmt"
)

// WebfingerResponse represents a Webfinger response.
type WebfingerResponse struct {
	Aliases []string `json:"aliases"`
	Subject string   `json:"subject"`
	Links   []Link   `json:"links"`
}

// WebfingerProfileRequestResponse represents a Webfinger profile request response.
type WebfingerProfileRequestResponse struct {
	Self string
}

// Link represents a Webfinger response Link entity.
type Link struct {
	Rel  string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
}

// MakeWebfingerResponse will create a new Webfinger response.
func MakeWebfingerResponse(account string, inbox string, host string) WebfingerResponse {
	accountIRI := MakeLocalIRIForAccount(account)

	return WebfingerResponse{
		Subject: fmt.Sprintf("acct:%s@%s", account, host),
		Aliases: []string{
			accountIRI.String(),
		},
		Links: []Link{
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: accountIRI.String(),
			},
			{
				Rel:  "http://webfinger.net/rel/profile-page",
				Type: "text/html",
				Href: accountIRI.String(),
			},
		},
	}
}

// MakeWebFingerRequestResponseFromData converts WebFinger data to an easier
// to use model.
func MakeWebFingerRequestResponseFromData(data []map[string]interface{}) WebfingerProfileRequestResponse {
	response := WebfingerProfileRequestResponse{}
	for _, link := range data {
		if link["rel"] == "self" {
			return WebfingerProfileRequestResponse{
				Self: link["href"].(string),
			}
		}
	}

	return response
}
