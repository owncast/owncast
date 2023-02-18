package fediverse

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/auth"
	fediverseauth "github.com/owncast/owncast/auth/fediverse"
	"github.com/owncast/owncast/core/user"
)

// RegisterFediverseOTPRequest registers a new OTP request for the given access token.
func (c *Controller) RegisterFediverseOTPRequest(u user.User, w http.ResponseWriter, r *http.Request) {
	type request struct {
		FediverseAccount string `json:"account"`
	}
	var req request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		c.WriteSimpleResponse(w, false, "Could not decode request: "+err.Error())
		return
	}

	accessToken := r.URL.Query().Get("accessToken")
	reg, success, err := fediverseauth.RegisterFediverseOTP(accessToken, u.ID, u.DisplayName, req.FediverseAccount)
	if err != nil {
		c.WriteSimpleResponse(w, false, "Could not register auth request: "+err.Error())
		return
	}

	if !success {
		c.WriteSimpleResponse(w, false, "Could not register auth request. One may already be pending. Try again later.")
		return
	}

	const fmtMsg = "<p>This is an automated message from %s. If you did not request this message please ignore or block. Your requested one-time code is:</p><p>%s</p>"
	msg := fmt.Sprintf(fmtMsg, c.Data.GetServerName(), reg.Code)
	if err := c.ActivityPub.SendDirectFederatedMessage(msg, reg.Account); err != nil {
		c.WriteSimpleResponse(w, false, "Could not send code to fediverse: "+err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "")
}

// VerifyFediverseOTPRequest verifies the given OTP code for the given access token.
func (c *Controller) VerifyFediverseOTPRequest(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Code string `json:"code"`
	}

	var req request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		c.WriteSimpleResponse(w, false, "Could not decode request: "+err.Error())
		return
	}
	accessToken := r.URL.Query().Get("accessToken")
	valid, authRegistration := fediverseauth.ValidateFediverseOTP(accessToken, req.Code)
	if !valid {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Check if a user with this auth already exists, if so, log them in.
	if u := auth.GetUserByAuth(authRegistration.Account, auth.Fediverse); u != nil {
		// Handle existing auth.
		log.Debugln("user with provided fedvierse identity already exists, logging them in")

		// Update the current user's access token to point to the existing user id.
		userID := u.ID
		if err := user.SetAccessTokenToOwner(accessToken, userID, c.Data.Store); err != nil {
			c.WriteSimpleResponse(w, false, err.Error())
			return
		}

		if authRegistration.UserDisplayName != u.DisplayName {
			loginMessage := fmt.Sprintf("**%s** is now authenticated as **%s**", authRegistration.UserDisplayName, u.DisplayName)
			if err := c.Core.Chat.SendSystemAction(loginMessage, true); err != nil {
				log.Errorln(err)
			}
		}

		c.WriteSimpleResponse(w, true, "")

		return
	}

	// Otherwise, save this as new auth.
	log.Debug("fediverse account does not already exist, saving it as a new one for the current user")
	if err := auth.AddAuth(authRegistration.UserID, authRegistration.Account, auth.Fediverse); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	// Update the current user's authenticated flag so we can show it in
	// the chat UI.
	if err := user.SetUserAsAuthenticated(authRegistration.UserID, c.Data.Store); err != nil {
		log.Errorln(err)
	}

	c.WriteSimpleResponse(w, true, "")
}
