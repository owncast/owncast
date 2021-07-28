package activitypub

import (
	"net/http"

	"github.com/owncast/owncast/activitypub/controllers"
)

func Start() {
	// WebFinger
	http.HandleFunc("/.well-known/webfinger", controllers.WebfingerHandler)
}
