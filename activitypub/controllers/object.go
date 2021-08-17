package controllers

import (
	"net/http"
	"strings"

	"github.com/owncast/owncast/activitypub/persistence"
)

func ObjectHandler(w http.ResponseWriter, r *http.Request) {
	pathComponents := strings.Split(r.URL.Path, "/")
	objectId := pathComponents[2]

	object, err := persistence.GetObject(objectId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		// controllers.WriteSimpleResponse(w, false, err.Error())
	}

	w.Write([]byte(object))

	// controllers.WriteSimpleResponse(w, true, object)
}
