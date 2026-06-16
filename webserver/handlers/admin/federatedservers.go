package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/federatedserversrepository"
	activitypubutils "github.com/owncast/owncast/services/activitypub/utils"
	"github.com/owncast/owncast/webserver/handlers/generated"
	webutils "github.com/owncast/owncast/webserver/utils"
)

const errCodeUnsupportedFeaturedStreams = "UNSUPPORTED_FEATURED_STREAMS"

// GetFederatedServers returns the list of federated servers we are
// following for the featured-streams mini-directory.
func (a *Admin) GetFederatedServers(w http.ResponseWriter, r *http.Request) {
	repo := federatedserversrepository.Get()
	if repo == nil {
		webutils.WriteSimpleResponse(w, false, "Federated servers repository is not initialised")
		return
	}

	servers, err := repo.GetFederatedServers()
	if err != nil {
		webutils.WriteSimpleResponse(w, false, "Failed to get federated servers: "+err.Error())
		return
	}

	// Ensure we return an empty array instead of null.
	if servers == nil {
		servers = []models.FederatedServer{}
	}

	response := struct {
		Servers interface{} `json:"servers"`
	}{
		Servers: servers,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Errorln("Failed to encode federated servers response:", err)
		webutils.WriteSimpleResponse(w, false, "Failed to encode response")
		return
	}
}

// AddFederatedServer kicks off a Follow request to another Owncast
// instance and stores a pending follow record.
func (a *Admin) AddFederatedServer(w http.ResponseWriter, r *http.Request) {
	var request generated.AddFederatedServerJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		webutils.WriteSimpleResponse(w, false, "Invalid request body: "+err.Error())
		return
	}

	serverURL, err := url.Parse(request.Url)
	if err != nil {
		webutils.WriteSimpleResponse(w, false, "Invalid server URL: "+err.Error())
		return
	}

	if serverURL.Scheme == "" {
		serverURL.Scheme = "https"
	}

	if serverURL.Scheme != "https" {
		webutils.WriteSimpleResponse(w, false, "Server URL must use https protocol for federation")
		return
	}

	repo := federatedserversrepository.Get()
	if repo == nil {
		webutils.WriteSimpleResponse(w, false, "Federated servers repository is not initialised")
		return
	}

	existingServer, _ := repo.GetFederatedServer(serverURL.String())
	if existingServer != nil {
		webutils.WriteSimpleResponse(w, false, "Already following this federated server")
		return
	}

	// Persist a pending follow record before sending the Follow so the
	// server is visible in the admin list immediately and so the Accept
	// that comes back has a record to transition to "accepted". The
	// record is keyed by the base server URL to match the lookups done
	// by the Accept/Reject inbox handlers.
	if err := repo.AddFederatedServer(serverURL.String(), "", "", time.Now(), true, "", "pending"); err != nil {
		log.Errorf("Failed to store pending federated server %s: %v", serverURL.String(), err)
		webutils.WriteSimpleResponse(w, false, "Failed to store federated server: "+err.Error())
		return
	}

	isStreamConnected := a.stream.GetStatus().Online
	if err := a.activitypub.Outbox().SendFollowRequestToOwncastServerURL(serverURL.String(), isStreamConnected); err != nil {
		// The follow never went out, so drop the pending record we just
		// created to keep the list honest and allow a later retry.
		_ = repo.RemoveFederatedServerByIRI(serverURL.String())
		log.Errorf("Failed to send follow request to %s: %v", serverURL.String(), err)
		if errors.Is(err, activitypubutils.ErrFeaturedStreamsUnsupported) {
			webutils.WriteSimpleResponseWithCode(w, false, "Featured streams unsupported by remote server", errCodeUnsupportedFeaturedStreams)
			return
		}
		webutils.WriteSimpleResponse(w, false, "Failed to send follow request: "+err.Error())
		return
	}

	log.Infof("Sent follow request to federated server: %s", serverURL.String())
	webutils.WriteSimpleResponse(w, true, "Follow request sent successfully. The server will appear in your list once they accept the follow.")
}

// RemoveFederatedServer removes a federated server by ID.
func (a *Admin) RemoveFederatedServer(w http.ResponseWriter, r *http.Request, id int) {
	repo := federatedserversrepository.Get()
	if repo == nil {
		webutils.WriteSimpleResponse(w, false, "Federated servers repository is not initialised")
		return
	}

	if err := repo.RemoveFederatedServer(int64(id)); err != nil {
		log.Errorf("Failed to remove federated server with ID %d: %v", id, err)
		webutils.WriteSimpleResponse(w, false, "Failed to remove federated server: "+err.Error())
		return
	}

	log.Infof("Removed federated server with ID: %d", id)
	webutils.WriteSimpleResponse(w, true, "Federated server removed successfully")
}

// AddFederatedServerOptions handles CORS preflight requests.
func (a *Admin) AddFederatedServerOptions(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// RemoveFederatedServerOptions handles CORS preflight requests.
func (a *Admin) RemoveFederatedServerOptions(w http.ResponseWriter, r *http.Request, id int) {
	w.WriteHeader(http.StatusNoContent)
}
