package models

import (
	"fmt"
)

type Webfinger struct {
	Subject string `json:"subject"`
	Links   []Link `json:"links"`
}

type Link struct {
	Rel  string `json:"rel"`
	Type string `json:"type"`
	Href string `json:"href"`
}

func MakeWebfingerResponse(account string, inbox string, host string) Webfinger {
	accountIRI := MakeLocalIRIForAccount(account)

	return Webfinger{
		Subject: fmt.Sprintf("acct:%s@%s", account, host),
		Links: []Link{
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: accountIRI.String(),
			},
		},
	}
}
