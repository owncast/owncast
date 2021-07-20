package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core/user"
	"github.com/owncast/owncast/utils"
)

type deleteExternalAPIUserRequest struct {
	Token string `json:"token"`
}

type createExternalAPIUserRequest struct {
	Name   string   `json:"name"`
	Scopes []string `json:"scopes"`
}

// CreateExternalAPIUser will generate a 3rd party access token.
func CreateExternalAPIUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request createExternalAPIUserRequest
	if err := decoder.Decode(&request); err != nil {
		controllers.BadRequestHandler(w, err)
		return
	}

	// Verify all the scopes provided are valid
	if !user.HasValidScopes(request.Scopes) {
		controllers.BadRequestHandler(w, errors.New("one or more invalid scopes provided"))
		return
	}

	token, err := utils.GenerateAccessToken()
	if err != nil {
		controllers.InternalErrorHandler(w, err)
		return
	}

	color := utils.GenerateRandomDisplayColor()

	if err := user.InsertExternalAPIUser(token, request.Name, color, request.Scopes); err != nil {
		controllers.InternalErrorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	controllers.WriteResponse(w, user.ExternalAPIUser{
		AccessToken:  token,
		DisplayName:  request.Name,
		DisplayColor: color,
		Scopes:       request.Scopes,
		CreatedAt:    time.Now(),
		LastUsedAt:   nil,
	})
}

// GetExternalAPIUsers will return all 3rd party access tokens.
func GetExternalAPIUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokens, err := user.GetExternalAPIUser()
	if err != nil {
		controllers.InternalErrorHandler(w, err)
		return
	}

	controllers.WriteResponse(w, tokens)
}

// DeleteExternalAPIUser will return a single 3rd party access token.
func DeleteExternalAPIUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != controllers.POST {
		controllers.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request deleteExternalAPIUserRequest
	if err := decoder.Decode(&request); err != nil {
		controllers.BadRequestHandler(w, err)
		return
	}

	if request.Token == "" {
		controllers.BadRequestHandler(w, errors.New("must provide a token"))
		return
	}

	if err := user.DeleteExternalAPIUser(request.Token); err != nil {
		controllers.InternalErrorHandler(w, err)
		return
	}

	controllers.WriteSimpleResponse(w, true, "deleted token")
}
