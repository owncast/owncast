package admin

// this is endpoint logic

import (
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/webserver/handlers/generated"
	webutils "github.com/owncast/owncast/webserver/utils"
)

// defaultUserPageSize is used when a caller doesn't provide a page limit.
const defaultUserPageSize = 50

// GetUsers returns a paginated, optionally name-filtered list of every user
// (chat viewers, authenticated/plugin users, and API integrations) for the
// admin user-management page.
func (a *Admin) GetUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	offset, _ := strconv.Atoi(query.Get("offset"))
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = defaultUserPageSize
	}
	search := query.Get("search")
	status := query.Get("status")

	users, total, err := a.userRepository.GetUsersPaginated(offset, limit, search, status)
	if err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	response := struct {
		Total   int            `json:"total"`
		Results []*models.User `json:"results"`
	}{
		Total:   total,
		Results: users,
	}

	webutils.WriteResponse(w, response)
}

// DeleteUser permanently removes a user and all of their data. This is
// irreversible, unlike disabling a user (a reversible ban). Any live chat
// connections for the user are dropped first.
func (a *Admin) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request generated.DeleteUserJSONBody
	if err := decoder.Decode(&request); err != nil {
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	if request.UserId == nil || *request.UserId == "" {
		webutils.WriteSimpleResponse(w, false, "must provide userId")
		return
	}

	// Drop any live connections before removing the identity.
	if clients, err := a.chat.GetClientsForUser(*request.UserId); err == nil && len(clients) > 0 {
		a.chat.DisconnectClients(clients)
	}

	if err := a.userRepository.DeleteUser(*request.UserId); err != nil {
		log.Errorln("error deleting user", err)
		webutils.WriteSimpleResponse(w, false, err.Error())
		return
	}

	webutils.WriteSimpleResponse(w, true, "user deleted")
}
