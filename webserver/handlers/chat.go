package handlers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/chat"
	plugins "github.com/owncast/owncast/services/plugins"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/handlers/generated"
	"github.com/owncast/owncast/webserver/router/middleware"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// ExternalGetChatMessages gets all of the chat messages.
func (h *Handlers) ExternalGetChatMessages(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)
	h.getChatMessages(w, r)
}

// GetChatMessages gets all of the chat messages.
func (h *Handlers) GetChatMessages(u models.User, w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)
	h.getChatMessages(w, r)
}

func (h *Handlers) getChatMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		messages := h.chatMessageRepository.GetChatHistory()

		if err := json.NewEncoder(w).Encode(messages); err != nil {
			log.Debugln(err)
		}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		if err := json.NewEncoder(w).Encode(webutils.J{"error": "method not implemented (PRs are accepted)"}); err != nil {
			webutils.InternalErrorHandler(w, err)
		}
	}
}

// RegisterAnonymousChatUser will register a new user.
func (h *Handlers) RegisterAnonymousChatUser(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)

	if r.Method == http.MethodOptions {
		// All OPTIONS requests should have a wildcard CORS header.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		// nolint:goconst
		webutils.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	type registerAnonymousUserResponse struct {
		ID          string `json:"id"`
		AccessToken string `json:"accessToken"`
		DisplayName string `json:"displayName"`
	}
	// If the request carries a verified gate session (set by the auth-gate
	// middleware; empty on servers with no auth.gate plugin), return that
	// user's existing token rather than creating a new anonymous one.
	// Fail closed on any problem: a gate session with no resolvable user
	// must never fall through to CreateAnonymousUser — that would create an
	// anon row and trigger an infinite re-registration loop via the /ws
	// mismatch guard.
	if tok := plugins.SessionTokenFromContext(r.Context()); tok != "" {
		u := h.userRepository.GetUserByToken(tok)
		if u == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(models.BaseAPIResponse{Success: false, Message: "gate session invalid"})
			return
		}
		if u.DisabledAt != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(models.BaseAPIResponse{Success: false, Message: "user is disabled"})
			return
		}
		utils.SetChatAccessTokenCookie(w, r, tok)
		webutils.WriteResponse(w, registerAnonymousUserResponse{
			ID:          u.ID,
			AccessToken: tok,
			DisplayName: u.DisplayName,
		})
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request generated.RegisterAnonymousChatUserJSONBody // registerAnonymousUserRequest
	if err := decoder.Decode(&request); err != nil {        //nolint
		// this is fine. register a new user anyway.
	}

	proposedNewDisplayName := r.Header.Get("X-Forwarded-User")
	if proposedNewDisplayName == "" && request.DisplayName != nil {
		proposedNewDisplayName = *request.DisplayName
	}

	// Sanitize before the empty check so a proposed name that sanitizes to
	// nothing (e.g. only HTML tags or whitespace) still falls back to a
	proposedNewDisplayName = utils.MakeSafeStringOfLength(proposedNewDisplayName, config.MaxChatDisplayNameLength)
	if _, forbidden := chat.ForbiddenUsername(proposedNewDisplayName, h.configRepository.GetForbiddenUsernameList()); forbidden {
		proposedNewDisplayName = ""
	}
	if proposedNewDisplayName == "" {
		proposedNewDisplayName = utils.MakeSafeStringOfLength(h.generateDisplayName(), config.MaxChatDisplayNameLength)
	}
	if proposedNewDisplayName == "" {
		// An admin-configured suggested username can itself sanitize to
		// empty; GeneratePhrase always yields a usable word-word name.
		proposedNewDisplayName = utils.MakeSafeStringOfLength(utils.GeneratePhrase(), config.MaxChatDisplayNameLength)
	}
	newUser, accessToken, err := h.userRepository.CreateAnonymousUser(proposedNewDisplayName)
	if err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	response := registerAnonymousUserResponse{
		ID:          newUser.ID,
		AccessToken: accessToken,
		DisplayName: newUser.DisplayName,
	}

	w.Header().Set("Content-Type", "application/json")
	middleware.DisableCache(w)

	// Set the chat identity cookie so subsequent same-origin web requests
	// (e.g. to plugin HTTP handlers) carry this user's identity.
	utils.SetChatAccessTokenCookie(w, r, accessToken)

	webutils.WriteResponse(w, response)
}

func (h *Handlers) generateDisplayName() string {
	suggestedUsernamesList := h.configRepository.GetSuggestedUsernamesList()
	minSuggestedUsernamePoolLength := 10

	if len(suggestedUsernamesList) >= minSuggestedUsernamePoolLength {
		index := utils.RandomIndex(len(suggestedUsernamesList))
		return suggestedUsernamesList[index]
	} else {
		return utils.GeneratePhrase()
	}
}
