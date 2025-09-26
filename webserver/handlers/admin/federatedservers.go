package admin

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/owncast/owncast/persistence/federatedserversrepository"
	"github.com/owncast/owncast/webserver/handlers/generated"
	log "github.com/sirupsen/logrus"
)

// GetFederatedServers returns the list of federated servers.
func GetFederatedServers(w http.ResponseWriter, r *http.Request) {
	repo := federatedserversrepository.Get()

	servers, err := repo.GetFederatedServers()
	if err != nil {
		writeSimpleResponse(w, false, "Failed to get federated servers: "+err.Error())
		return
	}

	response := struct {
		Servers interface{} `json:"servers"`
	}{
		Servers: servers,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Errorln("Failed to encode federated servers response:", err)
		writeSimpleResponse(w, false, "Failed to encode response")
		return
	}
}

// AddFederatedServer adds a new federated server to follow.
func AddFederatedServer(w http.ResponseWriter, r *http.Request) {
	var request generated.AddFederatedServerJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeSimpleResponse(w, false, "Invalid request body: "+err.Error())
		return
	}

	// Validate the URL
	serverURL, err := url.Parse(request.Url)
	if err != nil {
		writeSimpleResponse(w, false, "Invalid server URL: "+err.Error())
		return
	}

	if serverURL.Scheme != "https" && serverURL.Scheme != "http" {
		writeSimpleResponse(w, false, "Server URL must use http or https protocol")
		return
	}

	repo := federatedserversrepository.Get()

	// Check if we're already following this server
	existingServer, _ := repo.GetFederatedServer(serverURL.String())
	if existingServer != nil {
		writeSimpleResponse(w, false, "Already following this federated server")
		return
	}

	// Add the server (name and logo will be populated when we receive activities)
	err = repo.AddFederatedServer(serverURL.String(), "", "")
	if err != nil {
		log.Errorf("Failed to add federated server %s: %v", serverURL.String(), err)
		writeSimpleResponse(w, false, "Failed to add federated server: "+err.Error())
		return
	}

	log.Infof("Added federated server: %s", serverURL.String())
	writeSimpleResponse(w, true, "Federated server added successfully")
}

// RemoveFederatedServer removes a federated server by ID.
func RemoveFederatedServer(w http.ResponseWriter, r *http.Request, id int) {
	repo := federatedserversrepository.Get()

	// Convert int to int32 as expected by the repository
	serverID := int32(id)

	err := repo.RemoveFederatedServer(serverID)
	if err != nil {
		log.Errorf("Failed to remove federated server with ID %d: %v", id, err)
		writeSimpleResponse(w, false, "Failed to remove federated server: "+err.Error())
		return
	}

	log.Infof("Removed federated server with ID: %d", id)
	writeSimpleResponse(w, true, "Federated server removed successfully")
}

// AddFederatedServerOptions handles CORS preflight requests.
func AddFederatedServerOptions(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// RemoveFederatedServerOptions handles CORS preflight requests.
func RemoveFederatedServerOptions(w http.ResponseWriter, r *http.Request, id int) {
	w.WriteHeader(http.StatusNoContent)
}

// Helper function to write simple JSON responses
func writeSimpleResponse(w http.ResponseWriter, success bool, message string) {
	w.Header().Set("Content-Type", "application/json")

	status := http.StatusOK
	if !success {
		status = http.StatusBadRequest
	}
	w.WriteHeader(status)

	response := generated.BaseAPIResponse{
		Message: &message,
		Success: &success,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Errorln("Failed to encode response:", err)
	}
}
