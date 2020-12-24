package admin

// this is endpoint logic

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core/data"
	log "github.com/sirupsen/logrus"
)

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
	ID      string `json:"id"`
	Visible bool   `json:"visible"`
}

func updateMessageVisibility(_db *sql.DB, request messageVisibilityUpdateRequest) {
	tx, err := _db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("UPDATE messages SET visible=? WHERE id=?")

	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(request.Visible, request.ID)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
