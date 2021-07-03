package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
)

type deleteTokenRequest struct {
	Token string `json:"token"`
}

type createTokenRequest struct {
	Name   string   `json:"name"`
	Scopes []string `json:"scopes"`
}

// CreateAccessToken will generate a 3rd party access token.
func CreateAccessToken(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request createTokenRequest
	if err := decoder.Decode(&request); err != nil {
		controllers.BadRequestHandler(w, err)
		return
	}

	// Verify all the scopes provided are valid
	if !models.HasValidScopes(request.Scopes) {
		controllers.BadRequestHandler(w, errors.New("one or more invalid scopes provided"))
		return
	}

	token, err := utils.GenerateAccessToken()
	if err != nil {
		controllers.InternalErrorHandler(w, err)
		return
	}

	if err := data.InsertToken(token, request.Name, request.Scopes); err != nil {
		controllers.InternalErrorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	controllers.WriteResponse(w, models.AccessToken{
		Token:     token,
		Name:      request.Name,
		Scopes:    request.Scopes,
		Timestamp: time.Now(),
		LastUsed:  nil,
	})
}

// GetAccessTokens will return all 3rd party access tokens.
func GetAccessTokens(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokens, err := data.GetIntegrationAccessTokens()
	if err != nil {
		controllers.InternalErrorHandler(w, err)
		return
	}

	controllers.WriteResponse(w, tokens)
}

// DeleteAccessToken will return a single 3rd party access token.
func DeleteAccessToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != controllers.POST {
		controllers.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request deleteTokenRequest
	if err := decoder.Decode(&request); err != nil {
		controllers.BadRequestHandler(w, err)
		return
	}

	if request.Token == "" {
		controllers.BadRequestHandler(w, errors.New("must provide a token"))
		return
	}

	if err := data.DeleteToken(request.Token); err != nil {
		controllers.InternalErrorHandler(w, err)
		return
	}

	controllers.WriteSimpleResponse(w, true, "deleted token")
}
