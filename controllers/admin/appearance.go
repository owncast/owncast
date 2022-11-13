package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core/data"
)

// SetCustomColorVariableValues sets the custom color variables.
func SetCustomColorVariableValues(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	type request struct {
		Value map[string]string `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var values request

	if err := decoder.Decode(&values); err != nil {
		controllers.WriteSimpleResponse(w, false, "unable to update appearance variable values")
		return
	}

	if err := data.SetCustomColorVariableValues(values.Value); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "custom appearance variables updated")
}
