package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/activitypub/crypto"
	"github.com/owncast/owncast/activitypub/persistence"
	"github.com/owncast/owncast/activitypub/requests"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	log "github.com/sirupsen/logrus"
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

	if !data.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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
	if err := writeResponse(res, w); err != nil {
		log.Errorln(err)
	}
}

// NodeInfoV2Controller returns the V2 node info response.
func NodeInfoV2Controller(w http.ResponseWriter, r *http.Request) {
	type metadata struct {
		ChatEnabled bool `json:"chat_enabled"`
	}
	type services struct {
		Outbound []string `json:"outbound"`
		Inbound  []string `json:"inbound"`
	}
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
		Services          services `json:"services"`
		Software          software `json:"software"`
		Protocols         []string `json:"protocols"`
		Usage             usage    `json:"usage"`
		OpenRegistrations bool     `json:"openRegistrations"`
		Metadata          metadata `json:"metadata"`
	}

	if !data.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	localPostCount, _ := persistence.GetLocalPostCount()

	res := response{
		Version: "2.0",
		Services: services{
			Inbound:  []string{},
			Outbound: []string{},
		},
		Software: software{
			Name:    "owncast",
			Version: config.VersionNumber,
		},
		Usage: usage{
			Users: users{
				Total:          1,
				ActiveMonth:    1,
				ActiveHalfyear: 1,
			},
			LocalPosts: int(localPostCount),
		},
		OpenRegistrations: false,
		Protocols:         []string{"activitypub"},
		Metadata: metadata{
			ChatEnabled: !data.GetChatDisabled(),
		},
	}

	if err := writeResponse(res, w); err != nil {
		log.Errorln(err)
	}
}

// XNodeInfo2Controller returns the x-nodeinfo2.
func XNodeInfo2Controller(w http.ResponseWriter, r *http.Request) {
	type Organization struct {
		Name    string `json:"name"`
		Contact string `json:"contact"`
	}
	type Server struct {
		BaseURL  string `json:"baseUrl"`
		Version  string `json:"version"`
		Name     string `json:"name"`
		Software string `json:"software"`
	}
	type Services struct {
		Outbound []string `json:"outbound"`
		Inbound  []string `json:"inbound"`
	}
	type Users struct {
		ActiveWeek     int `json:"activeWeek"`
		Total          int `json:"total"`
		ActiveMonth    int `json:"activeMonth"`
		ActiveHalfyear int `json:"activeHalfyear"`
	}
	type Usage struct {
		Users         Users `json:"users"`
		LocalPosts    int   `json:"localPosts"`
		LocalComments int   `json:"localComments"`
	}
	type response struct {
		Server            Server       `json:"server"`
		Organization      Organization `json:"organization"`
		Version           string       `json:"version"`
		Services          Services     `json:"services"`
		Protocols         []string     `json:"protocols"`
		Usage             Usage        `json:"usage"`
		OpenRegistrations bool         `json:"openRegistrations"`
	}

	if !data.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	serverURL := data.GetServerURL()
	if serverURL == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	localPostCount, _ := persistence.GetLocalPostCount()

	res := &response{
		Organization: Organization{
			Name:    data.GetServerName(),
			Contact: serverURL,
		},
		Server: Server{
			BaseURL:  serverURL,
			Version:  config.VersionNumber,
			Name:     "owncast",
			Software: "owncast",
		},
		Services: Services{
			Inbound:  []string{"activitypub"},
			Outbound: []string{"activitypub"},
		},
		Protocols: []string{"activitypub"},
		Version:   config.VersionNumber,
		Usage: Usage{
			Users: Users{
				ActiveWeek:     1,
				Total:          1,
				ActiveMonth:    1,
				ActiveHalfyear: 1,
			},

			LocalPosts:    int(localPostCount),
			LocalComments: 0,
		},
	}

	if err := writeResponse(res, w); err != nil {
		log.Errorln(err)
	}
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
		Thumbnail        string   `json:"thumbnail"`
		Languages        []string `json:"languages"`
		Stats            Stats    `json:"stats"`
		Registrations    bool     `json:"registrations"`
		ApprovalRequired bool     `json:"approval_required"`
		InvitesEnabled   bool     `json:"invites_enabled"`
	}

	if !data.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
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

	thumbnail.Path = "/logo/external"
	localPostCount, _ := persistence.GetLocalPostCount()

	res := response{
		URI:              serverURL,
		Title:            data.GetServerName(),
		ShortDescription: data.GetServerSummary(),
		Description:      data.GetServerSummary(),
		Version:          config.GetReleaseString(),
		Stats: Stats{
			UserCount:   1,
			StatusCount: int(localPostCount),
			DomainCount: 0,
		},
		Thumbnail:        thumbnail.String(),
		Registrations:    false,
		ApprovalRequired: false,
		InvitesEnabled:   false,
	}

	if err := writeResponse(res, w); err != nil {
		log.Errorln(err)
	}
}

func writeResponse(payload interface{}, w http.ResponseWriter) error {
	accountName := data.GetDefaultFederationUsername()
	actorIRI := apmodels.MakeLocalIRIForAccount(accountName)
	publicKey := crypto.GetPublicKey(actorIRI)

	return requests.WritePayloadResponse(payload, w, publicKey)
}

// HostMetaController points to webfinger.
func HostMetaController(w http.ResponseWriter, r *http.Request) {
	if !data.GetFederationEnabled() {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Debugln("host meta request rejected! Federation is not enabled")
		return
	}

	serverURL := data.GetServerURL()
	if serverURL == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
	<XRD xmlns="http://docs.oasis-open.org/ns/xri/xrd-1.0">
		<Link rel="lrdd" type="application/json" template="%s/.well-known/webfinger?resource={uri}"/>
	</XRD>`, serverURL)

	if _, err := w.Write([]byte(res)); err != nil {
		log.Errorln(err)
	}
}
