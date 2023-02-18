package indieauth

import (
	"net/http"
	"net/url"

	"github.com/owncast/owncast/app/middleware"
	ia "github.com/owncast/owncast/auth/indieauth"
)

// HandleAuthEndpoint will handle the IndieAuth auth endpoint.
func (c *Controller) HandleAuthEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Require the GET request for IndieAuth to be behind admin login.
		f := middleware.RequireAdminAuth(c.handleAuthEndpointGet, c.Data)
		f(w, r)
		return
	} else if r.Method == http.MethodPost {
		c.handleAuthEndpointPost(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (c *Controller) handleAuthEndpointGet(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	codeChallenge := r.URL.Query().Get("code_challenge")
	state := r.URL.Query().Get("state")
	me := r.URL.Query().Get("me")

	request, err := ia.StartServerAuth(clientID, redirectURI, codeChallenge, state, me)
	if err != nil {
		_ = c.WriteString(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect the client browser with the values we generated to continue
	// the IndieAuth flow.
	// If the URL is invalid then return with specific "invalid_request" error.
	u, err := url.Parse(redirectURI)
	if err != nil {
		c.WriteResponse(w, ia.Response{
			Error:            "invalid_request",
			ErrorDescription: err.Error(),
		})
		return
	}

	redirectParams := u.Query()
	redirectParams.Set("code", request.Code)
	redirectParams.Set("state", request.State)
	u.RawQuery = redirectParams.Encode()

	http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
}

func (c *Controller) handleAuthEndpointPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	code := r.PostForm.Get("code")
	redirectURI := r.PostForm.Get("redirect_uri")
	clientID := r.PostForm.Get("client_id")
	codeVerifier := r.PostForm.Get("code_verifier")
	dataService := c.ActivityPub.Persistence.Data

	// If the server auth flow cannot be completed then return with specific
	// "invalid_client" error.
	response, err := ia.CompleteServerAuth(code, redirectURI, clientID, codeVerifier, dataService)
	if err != nil {
		c.WriteResponse(w, ia.Response{
			Error:            "invalid_client",
			ErrorDescription: err.Error(),
		})
		return
	}

	c.WriteResponse(w, response)
}
