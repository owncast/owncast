package indieauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	ia "github.com/owncast/owncast/auth/indieauth"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/userrepository"
	"github.com/owncast/owncast/services/chat"
	"github.com/owncast/owncast/webserver/router/middleware"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// Handler bundles the dependencies the IndieAuth handlers need.
type Handler struct {
	chat           *chat.Service
	userRepository userrepository.UserRepository
	indieAuth      *ia.Service
	middleware     *middleware.Middleware
}

// Deps lists the dependencies of the IndieAuth Handler.
type Deps struct {
	Chat           *chat.Service
	UserRepository userrepository.UserRepository
	IndieAuth      *ia.Service
	Middleware     *middleware.Middleware
}

// New constructs the Handler.
func New(deps Deps) *Handler {
	return &Handler{
		chat:           deps.Chat,
		userRepository: deps.UserRepository,
		indieAuth:      deps.IndieAuth,
		middleware:     deps.Middleware,
	}
}

// StartAuthFlow will begin the IndieAuth flow for the current user.
func (h *Handler) StartAuthFlow(u models.User, w http.ResponseWriter, r *http.Request) {
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

	redirectURL, err := h.indieAuth.StartAuthFlow(authRequest.AuthHost, u.ID, accessToken, u.DisplayName)
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
func (h *Handler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	request, response, err := h.indieAuth.HandleCallbackCode(code, state)
	if err != nil {
		log.Debugln(err)
		msg := `Unable to complete authentication. <a href="/">Go back.</a><hr/>`
		_ = webutils.WriteString(w, msg, http.StatusBadRequest)
		return
	}

	// Check if a user with this auth already exists, if so, log them in.
	if u := h.userRepository.GetUserByAuth(response.Me, models.IndieAuth); u != nil {
		// Handle existing auth.
		log.Debugln("user with provided indieauth already exists, logging them in")

		// Update the current user's access token to point to the existing user id.
		accessToken := request.CurrentAccessToken
		userID := u.ID
		if err := h.userRepository.SetAccessTokenToOwner(accessToken, userID); err != nil {
			webutils.WriteSimpleResponse(w, false, err.Error())
			return
		}

		if request.DisplayName != u.DisplayName {
			loginMessage := fmt.Sprintf("**%s** is now authenticated as **%s**", request.DisplayName, u.DisplayName)
			if err := h.chat.SendSystemAction(loginMessage, true); err != nil {
				log.Errorln(err)
			}
		}

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return
	}

	// Otherwise, save this as new auth. The IndieAuth "me" URL is both the auth
	// key and the public profile link, so record it as the profile URL too.
	log.Debug("indieauth token does not already exist, saving it as a new one for the current user")
	if err := h.userRepository.AddAuth(request.UserID, response.Me, models.IndieAuth, &models.LinkedIdentityFields{ProfileURL: response.Me}); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update the current user's authenticated flag so we can show it in
	// the chat UI.
	if err := h.userRepository.SetUserAsAuthenticated(request.UserID); err != nil {
		log.Errorln(err)
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
