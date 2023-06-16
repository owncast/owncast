package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/services/config"
	"github.com/owncast/owncast/storage"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/responses"
)

type deleteExternalAPIUserRequest struct {
	Token string `json:"token"`
}

type createExternalAPIUserRequest struct {
	Name   string   `json:"name"`
	Scopes []string `json:"scopes"`
}

// CreateExternalAPIUser will generate a 3rd party access token.
func (h *Handlers) CreateExternalAPIUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request createExternalAPIUserRequest
	if err := decoder.Decode(&request); err != nil {
		responses.BadRequestHandler(w, err)
		return
	}

	// Verify all the scopes provided are valid
	if !user.HasValidScopes(request.Scopes) {
		responses.BadRequestHandler(w, errors.New("one or more invalid scopes provided"))
		return
	}

	token, err := utils.GenerateAccessToken()
	if err != nil {
		responses.InternalErrorHandler(w, err)
		return
	}

	color := utils.GenerateRandomDisplayColor(config.MaxUserColor)

	if err := user.InsertExternalAPIUser(token, request.Name, color, request.Scopes); err != nil {
		responses.InternalErrorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	responses.WriteResponse(w, user.ExternalAPIUser{
		AccessToken:  token,
		DisplayName:  request.Name,
		DisplayColor: color,
		Scopes:       request.Scopes,
		CreatedAt:    time.Now(),
		LastUsedAt:   nil,
	})
}

// GetExternalAPIUsers will return all 3rd party access tokens.
func (h *Handlers) GetExternalAPIUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokens, err := user.GetExternalAPIUser()
	if err != nil {
		responses.InternalErrorHandler(w, err)
		return
	}
	responses.WriteResponse(w, tokens)
}

// DeleteExternalAPIUser will return a single 3rd party access token.
func (h *Handlers) DeleteExternalAPIUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		responses.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request deleteExternalAPIUserRequest
	if err := decoder.Decode(&request); err != nil {
		responses.BadRequestHandler(w, err)
		return
	}

	if request.Token == "" {
		responses.BadRequestHandler(w, errors.New("must provide a token"))
		return
	}

	if err := user.DeleteExternalAPIUser(request.Token); err != nil {
		responses.InternalErrorHandler(w, err)
		return
	}

	responses.WriteSimpleResponse(w, true, "deleted token")
}
