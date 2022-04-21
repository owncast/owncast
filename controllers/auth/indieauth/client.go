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
		log.Debugln(err)
		msg := fmt.Sprintf("Unable to complete authentication. <a href=\"/\">Go back.</a><hr/> %s", err.Error())
		_ = controllers.WriteString(w, msg, http.StatusBadRequest)
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

		loginMessage := fmt.Sprintf("**%s** is now authenticated as **%s**", request.DisplayName, u.DisplayName)
		if err := chat.SendSystemAction(loginMessage, true); err != nil {
			log.Errorln(err)
		}

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return
	}

	// Otherwise, save this as new auth.
	log.Debug("indieauth token does not already exist, saving it as a new one for the current user")
	if err := auth.AddAuth(request.UserID, response.Me, auth.IndieAuth); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update the current user's authenticated flag so we can show it in
	// the chat UI.
	if err := user.SetUserAsAuthenticated(request.UserID); err != nil {
		log.Errorln(err)
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
