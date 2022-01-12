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
