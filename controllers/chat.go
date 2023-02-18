package controllers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/app/middleware"
	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/utils"
)

// ExternalGetChatMessages gets all of the chat messages.
func (s *Service) ExternalGetChatMessages(integration user.ExternalAPIUser, w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)
	s.getChatMessages(w, r)
}

// GetChatMessages gets all of the chat messages.
func (s *Service) GetChatMessages(u user.User, w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)
	s.getChatMessages(w, r)
}

func (s *Service) getChatMessages(w http.ResponseWriter, r *http.Request) {
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
			s.InternalErrorHandler(w, err)
		}
	}
}

// RegisterAnonymousChatUser will register a new user.
func (s *Service) RegisterAnonymousChatUser(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(w)

	if r.Method == "OPTIONS" {
		// All OPTIONS requests should have a wildcard CORS header.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		s.WriteSimpleResponse(w, false, r.Method+" not supported")
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
	newUser, accessToken, err := user.CreateAnonymousUser(proposedNewDisplayName, s.Data)
	if err != nil {
		s.WriteSimpleResponse(w, false, err.Error())
		return
	}

	response := registerAnonymousUserResponse{
		ID:          newUser.ID,
		AccessToken: accessToken,
		DisplayName: newUser.DisplayName,
	}

	w.Header().Set("Content-Type", "application/json")
	middleware.DisableCache(w)

	s.WriteResponse(w, response)
}
