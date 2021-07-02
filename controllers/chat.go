package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/owncast/owncast/core/chat"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/router/middleware"
	log "github.com/sirupsen/logrus"
)

// GetChatMessages gets all of the chat messages.
func GetChatMessages(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		messages := chat.GetChatHistory()

		if err := json.NewEncoder(w).Encode(messages); err != nil {
			log.Errorln(err)
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
	w.Header().Set("Content-Type", "application/json")

	if r.Method != POST {
		WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	type registerAnonymousUserRequest struct {
		DisplayName string `json:"displayName"`
	}

	type registerAnonymousUserResponse struct {
		Id          string `json:"id"`
		AccessToken string `json:"accessToken"`
		DisplayName string `json:"displayName"`
	}

	decoder := json.NewDecoder(r.Body)
	var request registerAnonymousUserRequest
	if err := decoder.Decode(&request); err != nil {
		// this is fine. register a new user anyway.
	}

	err, newUser := user.CreateAnonymousUser(request.DisplayName)
	if err != nil {
		WriteSimpleResponse(w, false, err.Error())
		return
	}

	response := registerAnonymousUserResponse{
		Id:          newUser.Id,
		AccessToken: newUser.AccessToken,
		DisplayName: newUser.DisplayName,
	}

	fmt.Println("Registering user....", newUser.AccessToken)

	WriteResponse(w, response)
}
