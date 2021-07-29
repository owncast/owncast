package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/owncast/owncast/activitypub/models"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/utils"
)

// This is a mapping between account names and their outbox
var validAccounts = map[string]string{
	data.GetDefaultFederationUsername(): data.GetDefaultFederationUsername(),
}

func WebfingerHandler(w http.ResponseWriter, r *http.Request) {
	resource := r.URL.Query().Get("resource")
	resourceComponents := strings.Split(resource, ":")
	account := resourceComponents[1]

	userComponents := strings.Split(account, "@")
	if len(userComponents) < 2 {
		return
	}
	host := userComponents[1]
	user := userComponents[0]

	if _, valid := validAccounts[user]; !valid {
		// User is not valid
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// If the webfinger request doesn't match our server then it
	// should be rejected.
	instanceHostString := data.GetServerURL()
	if instanceHostString == "" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	instanceHostString = utils.GetHostnameFromURLString(instanceHostString)
	if instanceHostString == "" || instanceHostString != host {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	// instanceHost, err := url.Parse(instanceHostString)
	// if err != nil {
	// 	return
	// }

	// if host != instanceHost.Host {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	// if inbox, ok := validAccounts[user]; ok {
	webfingerResponse := models.MakeWebfingerResponse(user, user, host)

	w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(webfingerResponse); err != nil {
		panic(err)
	}
}
