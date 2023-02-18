package admin

import (
	"encoding/json"
	"net/http"
)

// SetCustomColorVariableValues sets the custom color variables.
func (c *Controller) SetCustomColorVariableValues(w http.ResponseWriter, r *http.Request) {
	if !c.requirePOST(w, r) {
		return
	}

	type request struct {
		Value map[string]string `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	var values request

	if err := decoder.Decode(&values); err != nil {
		c.WriteSimpleResponse(w, false, "unable to update appearance variable values")
		return
	}

	if err := c.Data.SetCustomColorVariableValues(values.Value); err != nil {
		c.WriteSimpleResponse(w, false, err.Error())
		return
	}

	c.WriteSimpleResponse(w, true, "custom appearance variables updated")
}
