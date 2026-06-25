package fediverse

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	fediverseauth "github.com/owncast/owncast/auth/fediverse"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"
	"github.com/owncast/owncast/persistence/userrepository"
	"github.com/owncast/owncast/services/activitypub"
	"github.com/owncast/owncast/services/chat"
	"github.com/owncast/owncast/webserver/handlers/generated"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// Handler bundles the dependencies the fediverse auth handlers need:
// the activitypub service (to send OTP codes via direct federated
// message), the chat service (to broadcast a system action when an
// authenticated user re-logs in under a different display name), and
// the fediverse auth service that owns the in-memory OTP state.
type Handler struct {
	activitypub      *activitypub.Service
	chat             *chat.Service
	fediverseAuth    *fediverseauth.Service
	configRepository configrepository.ConfigRepository
	userRepository   userrepository.UserRepository
}

// Deps lists the dependencies of the fediverse auth Handler.
type Deps struct {
	Activitypub      *activitypub.Service
	Chat             *chat.Service
	FediverseAuth    *fediverseauth.Service
	ConfigRepository configrepository.ConfigRepository
	UserRepository   userrepository.UserRepository
}

// New constructs the Handler.
func New(deps Deps) *Handler {
	return &Handler{
		activitypub:      deps.Activitypub,
		chat:             deps.Chat,
		fediverseAuth:    deps.FediverseAuth,
		configRepository: deps.ConfigRepository,
		userRepository:   deps.UserRepository,
	}
}

// RegisterFediverseOTPRequest registers a new OTP request for the given access token.
func (h *Handler) RegisterFediverseOTPRequest(u models.User, w http.ResponseWriter, r *http.Request) {
	type request struct {
		FediverseAccount string `json:"account"`
	}
	var req request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		webutils.WriteSimpleResponse(w, false, "Could not decode request: "+err.Error())
		return
	}

	accessToken := r.URL.Query().Get("accessToken")
	reg, success, err := h.fediverseAuth.RegisterFediverseOTP(accessToken, u.ID, u.DisplayName, req.FediverseAccount)
	if err != nil {
		webutils.WriteSimpleResponse(w, false, "Could not register auth request: "+err.Error())
		return
	}

	if !success {
		webutils.WriteSimpleResponse(w, false, "Could not register auth request. One may already be pending. Try again later.")
		return
	}

	msg := fmt.Sprintf("<p>One-time code from %s: %s. If you did not request this message please ignore or block.</p>", h.configRepository.GetServerName(), reg.Code)
	if err := h.activitypub.SendDirectFederatedMessage(msg, reg.Account); err != nil {
		webutils.WriteSimpleResponse(w, false, "Could not send code to fediverse: "+err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "")
}

// VerifyFediverseOTPRequest verifies the given OTP code for the given access token.
func (h *Handler) VerifyFediverseOTPRequest(w http.ResponseWriter, r *http.Request) {
	var req generated.VerifyFediverseOTPRequestJSONBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		webutils.WriteSimpleResponse(w, false, "Could not decode request: "+err.Error())
		return
	}

	if req.Code == nil {
		webutils.WriteSimpleResponse(w, false, "Could not decode request: code is required")
		return
	}

	accessToken := r.URL.Query().Get("accessToken")
	valid, authRegistration := h.fediverseAuth.ValidateFediverseOTP(accessToken, *req.Code)
	if !valid {
		webutils.WriteSimpleResponse(w, false, "Incorrect or expired code. Please request a new one and try again.")
		return
	}

	// Check if a user with this auth already exists, if so, log them in.
	if u := h.userRepository.GetUserByAuth(authRegistration.Account, models.Fediverse); u != nil {
		// Handle existing auth.
		log.Debugln("user with provided fedvierse identity already exists, logging them in")

		// Update the current user's access token to point to the existing user id.
		userID := u.ID
		if err := h.userRepository.SetAccessTokenToOwner(accessToken, userID); err != nil {
			webutils.WriteSimpleResponse(w, false, err.Error())
			return
		}

		if authRegistration.UserDisplayName != u.DisplayName {
			loginMessage := fmt.Sprintf("**%s** is now authenticated as **%s**", authRegistration.UserDisplayName, u.DisplayName)
			if err := h.chat.SendSystemAction(loginMessage, true); err != nil {
				log.Errorln(err)
			}
		}

		webutils.WriteSimpleResponse(w, true, "")

		return
	}

	// Otherwise, save this as new auth. The @me@host account is the auth key and
	// the human-readable handle; the public profile URL is left to be resolved
	// from the handle (webfinger) when the public badge UI needs it.
	log.Debug("fediverse account does not already exist, saving it as a new one for the current user")
	if err := h.userRepository.AddAuth(authRegistration.UserID, authRegistration.Account, models.Fediverse, &models.LinkedIdentityFields{Handle: authRegistration.Account}); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update the current user's authenticated flag so we can show it in
	// the chat UI.
	if err := h.userRepository.SetUserAsAuthenticated(authRegistration.UserID); err != nil {
		log.Errorln(err)
	}

	webutils.WriteSimpleResponse(w, true, "")
}
