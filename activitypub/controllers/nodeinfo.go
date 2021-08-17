package controllers

import (
	"net/http"
	"net/url"

	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/models"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
)

// NodeInfoController returns the V1 node info response.
func NodeInfoController(w http.ResponseWriter, r *http.Request) {
	type links struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	}

	type response struct {
		Links []links `json:"links"`
	}

	serverURL := data.GetServerURL()
	if serverURL == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	v2, err := url.Parse(serverURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	v2.Path = "nodeinfo/2.0"

	res := response{
		Links: []links{
			{
				Rel:  "http://nodeinfo.diaspora.software/ns/schema/2.0",
				Href: v2.String(),
			},
		},
	}
	writeResponse(res, w)
}

// NodeInfoV2Controller returns the V2 node info response.
func NodeInfoV2Controller(w http.ResponseWriter, r *http.Request) {
	type software struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	type users struct {
		Total          int `json:"total"`
		ActiveMonth    int `json:"activeMonth"`
		ActiveHalfyear int `json:"activeHalfyear"`
	}
	type usage struct {
		Users      users `json:"users"`
		LocalPosts int   `json:"localPosts"`
	}
	type response struct {
		Version           string   `json:"version"`
		Software          software `json:"software"`
		Protocols         []string `json:"protocols"`
		Usage             usage    `json:"usage"`
		OpenRegistrations bool     `json:"openRegistrations"`
	}

	res := response{
		Version: "2.0",
		Software: software{
			Name:    "Owncast",
			Version: config.VersionNumber,
		},
		Usage: usage{
			Users: users{
				Total:          1,
				ActiveMonth:    1,
				ActiveHalfyear: 1,
			},
			LocalPosts: persistence.GetLocalPostCount(),
		},
		OpenRegistrations: false,
		Protocols:         []string{"activitypub"},
	}

	writeResponse(res, w)
}

// InstanceV1Controller returns the v1 instance details.
func InstanceV1Controller(w http.ResponseWriter, r *http.Request) {
	type Stats struct {
		UserCount   int `json:"user_count"`
		StatusCount int `json:"status_count"`
		DomainCount int `json:"domain_count"`
	}
	type response struct {
		URI              string   `json:"uri"`
		Title            string   `json:"title"`
		ShortDescription string   `json:"short_description"`
		Description      string   `json:"description"`
		Version          string   `json:"version"`
		Stats            Stats    `json:"stats"`
		Thumbnail        string   `json:"thumbnail"`
		Languages        []string `json:"languages"`
		Registrations    bool     `json:"registrations"`
		ApprovalRequired bool     `json:"approval_required"`
		InvitesEnabled   bool     `json:"invites_enabled"`
	}

	serverURL := data.GetServerURL()
	if serverURL == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	thumbnail, err := url.Parse(serverURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	thumbnail.Path = "/logo"

	res := response{
		URI:              serverURL,
		Title:            data.GetServerName(),
		ShortDescription: data.GetServerSummary(),
		Description:      data.GetServerSummary(),
		Version:          config.GetReleaseString(),
		Stats: Stats{
			UserCount:   1,
			StatusCount: persistence.GetLocalPostCount(),
			DomainCount: 0,
		},
		Thumbnail:        thumbnail.String(),
		Registrations:    false,
		ApprovalRequired: false,
		InvitesEnabled:   false,
	}

	writeResponse(res, w)
}

func writeResponse(payload interface{}, w http.ResponseWriter) error {
	accountName := data.GetDefaultFederationUsername()
	actorIRI := models.MakeLocalIRIForAccount(accountName)
	publicKey := crypto.GetPublicKey(actorIRI)

	return requests.WritePayloadResponse(payload, w, publicKey)
}
