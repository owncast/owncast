package indieauth

import (
	"net/http"
	"net/url"

	ia "github.com/owncast/owncast/auth/indieauth"
	"github.com/owncast/owncast/webserver/router/middleware"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// HandleAuthEndpoint will handle the IndieAuth auth endpoint.
func (h *Handler) HandleAuthEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Require the GET request for IndieAuth to be behind admin login.
		f := middleware.RequireAdminAuth(h.HandleAuthEndpointGet)
		f(w, r)
		return
	case http.MethodPost:
		h.HandleAuthEndpointPost(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) HandleAuthEndpointGet(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	codeChallenge := r.URL.Query().Get("code_challenge")
	state := r.URL.Query().Get("state")
	me := r.URL.Query().Get("me")

	request, err := ia.StartServerAuth(clientID, redirectURI, codeChallenge, state, me)
	if err != nil {
		_ = webutils.WriteString(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect the client browser with the values we generated to continue
	// the IndieAuth flow.
	// If the URL is invalid then return with specific "invalid_request" error.
	u, err := url.Parse(redirectURI)
	if err != nil {
		webutils.WriteResponse(w, ia.Response{
			Error:            "invalid_request",
			ErrorDescription: err.Error(),
		})
		return
	}

	redirectParams := u.Query()
	redirectParams.Set("code", request.Code)
	redirectParams.Set("state", request.State)
	u.RawQuery = redirectParams.Encode()

	// The redirect target is the IndieAuth client's registered redirect_uri,
	// already parsed and validated by ia.StartServerAuth above. This handler
	// is also gated by middleware.RequireAdminAuth, so only authenticated
	// admins reach this code path. The "open redirect" warning is the
	// protocol working as designed.
	http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect) //nolint:gosec // G710: redirect target is the validated IndieAuth client redirect_uri
}

func (h *Handler) HandleAuthEndpointPost(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := r.ParseForm(); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	code := r.PostForm.Get("code")
	redirectURI := r.PostForm.Get("redirect_uri")
	clientID := r.PostForm.Get("client_id")
	codeVerifier := r.PostForm.Get("code_verifier")

	// If the server auth flow cannot be completed then return with specific
	// "invalid_client" error.
	response, err := ia.CompleteServerAuth(code, redirectURI, clientID, codeVerifier)
	if err != nil {
		webutils.WriteResponse(w, ia.Response{
			Error:            "invalid_client",
			ErrorDescription: err.Error(),
		})
		return
	}

	webutils.WriteResponse(w, response)
}
