package indieauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/owncast/owncast/auth"
	ia "github.com/owncast/owncast/auth/indieauth"
	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/user"
	log "github.com/sirupsen/logrus"
)

// StartAuthFlow will begin the IndieAuth flow for the current user.
func StartAuthFlow(u user.User, w http.ResponseWriter, r *http.Request) {
	type request struct {
		AuthHost string `json:"authHost"`
	}

	type response struct {
		Redirect string `json:"redirect"`
	}

	var authRequest request
	p, err := io.ReadAll(r.Body)
	if err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	if err := json.Unmarshal(p, &authRequest); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	accessToken := r.URL.Query().Get("accessToken")

	redirectURL, err := ia.StartAuthFlow(authRequest.AuthHost, u.ID, accessToken, u.DisplayName)
	if err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	redirectResponse := response{
		Redirect: redirectURL.String(),
	}
	controllers.WriteResponse(w, redirectResponse)
}

// HandleRedirect will handle the redirect from an IndieAuth server to
// continue the auth flow.
func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	request, response, err := ia.HandleCallbackCode(code, state)
	if err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Check if a user with this auth already exists, if so, log them in.
	if u := auth.GetUserByAuth(response.Me, auth.IndieAuth); u != nil {
		// Handle existing auth.
		log.Debugln("user with provided indieauth already exists, logging them in")

		// Update the current user's access token to point to the existing user id.
		accessToken := request.CurrentAccessToken
		userID := u.ID
		if err := user.SetAccessTokenToOwner(accessToken, userID); err != nil {
			controllers.WriteSimpleResponse(w, false, err.Error())
			return
		}

		loginMessage := fmt.Sprintf("%s logged in as %s", request.DisplayName, u.DisplayName)
		chat.SendSystemAction(loginMessage, true)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return
	}

	// Otherwise, save this as new auth.
	log.Debug("indieauth token does not already exist, saving it as a new one for the current user")
	if err := auth.AddAuth(request.UserID, response.Me, auth.IndieAuth); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
