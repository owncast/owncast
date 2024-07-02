package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/userrepository"
	"github.com/owncast/owncast/router/middleware"
	"github.com/owncast/owncast/utils"
	log "github.com/sirupsen/logrus"
)

// ExternalGetChatMessages gets all of the chat messages.
func ExternalGetChatMessages(integration models.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)
	getChatMessages(w, r)
}

// GetChatMessages gets all of the chat messages.
func GetChatMessages(u models.User, w http.ResponseWriter, r *http.Request) {
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
		if err := json.NewEncoder(w).Encode(j{"error": "method not implemented (PRs are accepted)"}); err != nil {
			InternalErrorHandler(w, err)
		}
	}
}

// RegisterAnonymousChatUser will register a new user.
func RegisterAnonymousChatUser(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)

	userRepository := userrepository.Get()

	if r.Method == http.MethodOptions {
		// All OPTIONS requests should have a wildcard CORS header.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		// nolint:goconst
		WriteSimpleResponse(w, false, r.Method+" not supported")
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

	proposedNewDisplayName := r.Header.Get("X-Forwarded-User")
	if proposedNewDisplayName == "" {
		proposedNewDisplayName = request.DisplayName
	}
	if proposedNewDisplayName == "" {
		proposedNewDisplayName = generateDisplayName()
	}

	proposedNewDisplayName = utils.MakeSafeStringOfLength(proposedNewDisplayName, config.MaxChatDisplayNameLength)
	newUser, accessToken, err := userRepository.CreateAnonymousUser(proposedNewDisplayName)
	if err != nil {
		WriteSimpleResponse(w, false, err.Error())
		return
	}

	response := registerAnonymousUserResponse{
		ID:          newUser.ID,
		AccessToken: accessToken,
		DisplayName: newUser.DisplayName,
	}

	w.Header().Set("Content-Type", "application/json")
	middleware.DisableCache(w)

	WriteResponse(w, response)
}

func generateDisplayName() string {
	suggestedUsernamesList := data.GetSuggestedUsernamesList()
	minSuggestedUsernamePoolLength := 10

	if len(suggestedUsernamesList) >= minSuggestedUsernamePoolLength {
		index := utils.RandomIndex(len(suggestedUsernamesList))
		return suggestedUsernamesList[index]
	} else {
		return utils.GeneratePhrase()
	}
}
