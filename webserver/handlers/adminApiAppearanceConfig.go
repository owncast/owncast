package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/webserver/requests"
	"github.com/owncast/owncast/webserver/responses"
)

// SetCustomColorVariableValues sets the custom color variables.
func (h *Handlers) SetCustomColorVariableValues(w http.ResponseWriter, r *http.Request) {
	if !requests.RequirePOST(w, r) {
		return
	}

	type request struct {
		Value map[string]string `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var values request

	if err := decoder.Decode(&values); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to update appearance variable values")
		return
	}

	if err := data.SetCustomColorVariableValues(values.Value); err != nil {
		responses.WriteSimpleResponse(w, false, err.Error())
		return
	}

	responses.WriteSimpleResponse(w, true, "custom appearance variables updated")
}
