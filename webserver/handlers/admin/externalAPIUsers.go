package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/utils"
	"github.com/owncast/owncast/webserver/handlers/generated"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// CreateExternalAPIUser will generate a 3rd party access token.
func (a *Admin) CreateExternalAPIUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request generated.CreateExternalAPIUserJSONBody
	if err := decoder.Decode(&request); err != nil {
		webutils.BadRequestHandler(w, err)
		return
	}

	if request.Scopes == nil {
		webutils.BadRequestHandler(w, errors.New("scopes field is required"))
		return
	}

	if request.Name == nil {
		webutils.BadRequestHandler(w, errors.New("name field is required"))
		return
	}

	// Verify all the scopes provided are valid
	if !a.userRepository.HasValidScopes(*request.Scopes) {
		webutils.BadRequestHandler(w, errors.New("one or more invalid scopes provided"))
		return
	}

	token, err := utils.GenerateAccessToken()
	if err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	color := utils.GenerateRandomDisplayColor(config.MaxUserColor)

	if err := a.userRepository.InsertExternalAPIUser(token, *request.Name, color, *request.Scopes); err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	webutils.WriteResponse(w, models.ExternalAPIUser{
		AccessToken:  token,
		DisplayName:  *request.Name,
		DisplayColor: color,
		Scopes:       *request.Scopes,
		CreatedAt:    time.Now(),
		LastUsedAt:   nil,
	})
}

// GetExternalAPIUsers will return all 3rd party access tokens.
func (a *Admin) GetExternalAPIUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokens, err := a.userRepository.GetExternalAPIUser()
	if err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}
	webutils.WriteResponse(w, tokens)
}

// DeleteExternalAPIUser will return a single 3rd party access token.
func (a *Admin) DeleteExternalAPIUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		webutils.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request generated.DeleteExternalAPIUserJSONBody
	if err := decoder.Decode(&request); err != nil {
		webutils.BadRequestHandler(w, err)
		return
	}

	if request.Token == nil {
		webutils.BadRequestHandler(w, errors.New("token field is required"))
		return
	}

	if *request.Token == "" {
		webutils.BadRequestHandler(w, errors.New("must provide a token"))
		return
	}

	if err := a.userRepository.DeleteExternalAPIUser(*request.Token); err != nil {
		webutils.InternalErrorHandler(w, err)
		return
	}

	webutils.WriteSimpleResponse(w, true, "deleted token")
}
