package indieauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	ia "github.com/owncast/owncast/auth/indieauth"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/userrepository"
	webutils "github.com/owncast/owncast/webserver/utils"
	log "github.com/sirupsen/logrus"
)

// StartAuthFlow will begin the IndieAuth flow for the current user.
func StartAuthFlow(u models.User, w http.ResponseWriter, r *http.Request) {
	type request struct {
		AuthHost string `json:"authHost"`
	}

	type response struct {
		Redirect string `json:"redirect"`
	}

	var authRequest request
	p, err := io.ReadAll(r.Body)
	if err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	if err := json.Unmarshal(p, &authRequest); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	accessToken := r.URL.Query().Get("accessToken")

	redirectURL, err := ia.StartAuthFlow(authRequest.AuthHost, u.ID, accessToken, u.DisplayName)
	if err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	redirectResponse := response{
		Redirect: redirectURL.String(),
	}
	webutils.WriteResponse(w, redirectResponse)
}

// HandleRedirect will handle the redirect from an IndieAuth server to
// continue the auth flow.
func HandleRedirect(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	request, response, err := ia.HandleCallbackCode(code, state)
	if err != nil {
		log.Debugln(err)
		msg := `Unable to complete authentication. <a href="/">Go back.</a><hr/>`
		_ = webutils.WriteString(w, msg, http.StatusBadRequest)
		return
	}

	userRepository := userrepository.Get()

	// Check if a user with this auth already exists, if so, log them in.
	if u := userRepository.GetUserByAuth(response.Me, models.IndieAuth); u != nil {
		// Handle existing auth.
		log.Debugln("user with provided indieauth already exists, logging them in")

		// Update the current user's access token to point to the existing user id.
		accessToken := request.CurrentAccessToken
		userID := u.ID
		if err := userRepository.SetAccessTokenToOwner(accessToken, userID); err != nil {
			webutils.WriteSimpleResponse(w, false, err.Error())
			return
		}

		if request.DisplayName != u.DisplayName {
			loginMessage := fmt.Sprintf("**%s** is now authenticated as **%s**", request.DisplayName, u.DisplayName)
			if err := chat.SendSystemAction(loginMessage, true); err != nil {
				log.Errorln(err)
			}
		}

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return
	}

	// Otherwise, save this as new auth.
	log.Debug("indieauth token does not already exist, saving it as a new one for the current user")
	if err := userRepository.AddAuth(request.UserID, response.Me, models.IndieAuth); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update the current user's authenticated flag so we can show it in
	// the chat UI.
	if err := userRepository.SetUserAsAuthenticated(request.UserID); err != nil {
		log.Errorln(err)
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
