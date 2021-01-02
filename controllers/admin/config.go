package admin

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/owncast/owncast/controllers"
	"github.com/owncast/owncast/core/data"
)

type ConfigValue struct {
	Value interface{} `json:"value"`
}

func ChangeTags(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValues, success := getValuesFromRequest(w, r)
	if !success {
		return
	}

	var tagStrings []string
	for _, tag := range configValues {
		tagStrings = append(tagStrings, tag.Value.(string))
	}

	if err := data.SetServerMetadataTags(tagStrings); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "changed")
}

func ChangeStreamTitle(w http.ResponseWriter, r *http.Request) {
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

func ChangeExtraPageContent(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetExtraPageBodyContent(configValue.Value.(string)); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "changed")
}

func ChangeStreamKey(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetStreamKey(configValue.Value.(string)); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "changed")
}

func ChangeLogoPath(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetLogoPath(configValue.Value.(string)); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "changed")
}

func ChangeNSFW(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetNSFW(configValue.Value.(bool)); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "changed")
}

func ChangeFfmpegPath(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetFfmpegPath(configValue.Value.(string)); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "changed")
}

func ChangeWebServerPort(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetHTTPPortNumber(configValue.Value.(int)); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "http port set")
}

func ChangeRTMPServerPort(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetRTMPPortNumber(configValue.Value.(int)); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "rtmp port set")
}

func ChangeServerURL(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetServerURL(configValue.Value.(string)); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}

	controllers.WriteSimpleResponse(w, true, "server url set")
}

func ChangeDirectoryEnabled(w http.ResponseWriter, r *http.Request) {
	if !requirePOST(w, r) {
		return
	}

	configValue, success := getValueFromRequest(w, r)
	if !success {
		return
	}

	if err := data.SetDirectoryEnabled(configValue.Value.(bool)); err != nil {
		controllers.WriteSimpleResponse(w, false, err.Error())
		return
	}
	controllers.WriteSimpleResponse(w, true, "directory state changed")
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

func getValuesFromRequest(w http.ResponseWriter, r *http.Request) ([]ConfigValue, bool) {
	var values []ConfigValue

	decoder := json.NewDecoder(r.Body)
	var configValue ConfigValue
	if err := decoder.Decode(&configValue); err != nil {
		controllers.WriteSimpleResponse(w, false, "unable to parse array of values")
		return values, false
	}

	object := reflect.ValueOf(configValue.Value)

	for i := 0; i < object.Len(); i++ {
		values = append(values, ConfigValue{Value: object.Index(i).Interface()})
	}

	return values, true
}
