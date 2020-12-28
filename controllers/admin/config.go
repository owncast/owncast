package admin

import (
	"encoding/json"
	"net/http"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/router/middleware"
)

type ConfigValue struct {
	Value interface{} `json:"value"`
}

func ChangeStreamTitle(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)

	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetStreamTitle(configValue.Value.(string)); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "changed")
}

func ChangeServerName(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)

	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetServerName(configValue.Value.(string)); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "changed")
}

func ChangeServerTitle(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)

	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetServerTitle(configValue.Value.(string)); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "changed")
}

func ChangeServerSummary(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCors(&w)

	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetServerSummary(configValue.Value.(string)); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "changed")
}

func requirePOST(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != "POST" {
		controllers.WriteSimpleResponse(w, false, r.Method+" not supported")
		return false
	}

	return true
}

func getValueFromRequest(w http.ResponseWriter, r *http.Request) (ConfigValue, bool) {
	decoder := json.NewDecoder(r.Body)
	var configValue ConfigValue
	if err := decoder.Decode(&configValue); err != nil {
		controllers.WriteSimpleResponse(w, false, "unable to parse new value")
		return configValue, false
	}

	return configValue, true
}
