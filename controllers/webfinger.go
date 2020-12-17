package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
)

func GetWebfinger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/jrd+json")

	resource := r.URL.Query().Get("resource")

	if resource == "" {
		badRequestHandler(w, errors.New("Required parameter `resource` not supplied"))
		return
	}

	subject, err := config.Config.Federation.GetSubject()

	if err != nil {
		internalErrorHandler(w, err)
		return
	}

	if resource != subject {
		notFoundHandler(w, errors.New("Unrecognized resource requested"))
		return
	}

	links := []models.WebfingerLink{
		{
			Rel:  "self",
			Type: "application/activity+json",
			Href: urlFor("/actor"),
		},
	}

	webfinger := models.Webfinger{
		Subject: subject,
		Links:   links,
	}

	if err := json.NewEncoder(w).Encode(webfinger); err != nil {
		internalErrorHandler(w, err)
	}
}
