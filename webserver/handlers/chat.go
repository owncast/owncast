package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/config"
	"github.com/owncast/owncast/storage/configrepository"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/middleware"
	"github.com/owncast/owncast/webserver/responses"
	log "github.com/sirupsen/logrus"
)

// ExternalGetChatMessages gets all of the chat messages.
func (h *Handlers) ExternalGetChatMessages(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)
	getChatMessages(w, r)
}

// GetChatMessages gets all of the chat messages.
func (h *Handlers) GetChatMessages(u models.User, w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)
	getChatMessages(w, r)
}

func getChatMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		messages := chat.GetChatHistory()

		if err := json.NewEncoder(w).Encode(messages); err != nil {
			log.Debugln(err)
		}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		responses.BadRequestHandler(w, errors.New("method not implemented"))
	}
}

// RegisterAnonymousChatUser will register a new user.
func (h *Handlers) RegisterAnonymousChatUser(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)

	if r.Method == "OPTIONS" {
		// All OPTIONS requests should have a wildcard CORS header.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		responses.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	type registerAnonymousUserRequest struct {
		DisplayName string `json:"displayName"`
	}

	type registerAnonymousUserResponse struct {
		ID          string `json:"id"`
		AccessToken string `json:"accessToken"`
		DisplayName string `json:"displayName"`
	}

	decoder := json.NewDecoder(r.Body)
	var request registerAnonymousUserRequest
	if err := decoder.Decode(&request); err != nil { //nolint
		// this is fine. register a new user anyway.
	}

	if request.DisplayName == "" {
		request.DisplayName = r.Header.Get("X-Forwarded-User")
	}

	proposedNewDisplayName := utils.MakeSafeStringOfLength(request.DisplayName, config.MaxChatDisplayNameLength)
	newUser, accessToken, err := userRepository.CreateAnonymousUser(proposedNewDisplayName)
	if err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	response := registerAnonymousUserResponse{
		ID:          newUser.ID,
		AccessToken: accessToken,
		DisplayName: newUser.DisplayName,
	}

	w.Header().Set("Content-Type", "application/json")
	middleware.DisableCache(w)

	responses.WriteResponse(w, response)
}

func generateDisplayName() string {
	cr := configrepository.Get()
	suggestedUsernamesList := cr.GetSuggestedUsernamesList()

	const minSuggestedUsernamePoolLength = 10
	if len(suggestedUsernamesList) >= minSuggestedUsernamePoolLength {
		index := utils.RandomIndex(len(suggestedUsernamesList))
		return suggestedUsernamesList[index]
	} else {
		return utils.GeneratePhrase()
	}
}
