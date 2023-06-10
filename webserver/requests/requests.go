package requests

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/owncast/owncast/webserver/responses"
	log "github.com/sirupsen/logrus"
)

func RequirePOST(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		responses.WriteSimpleResponse(w, false, r.Method+" not supported")
		return false
	}

	return true
}

func GetValueFromRequest(w http.ResponseWriter, r *http.Request) (ConfigValue, bool) {
	decoder := json.NewDecoder(r.Body)
	var configValue ConfigValue
	if err := decoder.Decode(&configValue); err != nil {
		log.Warnln(err)
		responses.WriteSimpleResponse(w, false, "unable to parse new value")
		return configValue, false
	}

	return configValue, true
}

func GetValuesFromRequest(w http.ResponseWriter, r *http.Request) ([]ConfigValue, bool) {
	var values []ConfigValue

	decoder := json.NewDecoder(r.Body)
	var configValue ConfigValue
	if err := decoder.Decode(&configValue); err != nil {
		responses.WriteSimpleResponse(w, false, "unable to parse array of values")
		return values, false
	}

	object := reflect.ValueOf(configValue.Value)

	for i := 0; i < object.Len(); i++ {
		values = append(values, ConfigValue{Value: object.Index(i).Interface()})
	}

	return values, true
}
