package admin

// this is endpoint logic

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/router/middleware"
	log "github.com/sirupsen/logrus"
)

// UpdateMessageVisibility updates an array of message IDs to have the same visiblity.
func UpdateMessageVisibility(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		controllers.WriteSimpleResponse(w, false, r.Method+" not supported")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var request messageVisibilityUpdateRequest // creates an empty struc

	err := decoder.Decode(&request) // decode the json into `request`
	if err != nil {
		log.Errorln(err)
		controllers.WriteSimpleResponse(w, false, "")
		return
	}

	// make sql update call here.
	// := means create a new var
	_db := data.GetDatabase()
	updateMessageVisibility(_db, request)

	controllers.WriteSimpleResponse(w, true, "changed")
}

type messageVisibilityUpdateRequest struct {
	IDArray []string `json:"idArray"`
	Visible bool     `json:"visible"`
}

func updateMessageVisibility(_db *sql.DB, request messageVisibilityUpdateRequest) {
	tx, err := _db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("UPDATE messages SET visible=? WHERE id IN (?)")

	strIDList := strings.Join(request.IDArray, ",")

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(request.Visible, strIDList)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// GetChatMessages returns all of the chat messages, unfiltered.
func GetChatMessages(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)
	w.Header().Set("Content-Type", "application/json")

	messages := core.GetAllChatMessages(false)

	if err := json.NewEncoder(w).Encode(messages); err != nil {
		log.Errorln(err)
	}
}
