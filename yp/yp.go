package yp

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"encoding/json"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/utils"

	log "github.com/sirupsen/logrus"
)

const pingURL = "http://localhost:8088/ping"
const pingInterval = 1 * time.Minute

//YP
type YP struct {
	ticker            time.Ticker
	key               string
	getStreamingSince func() utils.NullTime
}

type ypPingResponse struct {
	Key string
}

type ypPingRequest struct {
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	URL            string          `json:"url"`
	Logo           string          `json:"logo"`
	NSFW           bool            `json:"nsfw"`
	Key            string          `json:"key"`
	Tags           []string        `json:"tags"`
	StreamingSince *utils.NullTime `json:"streamingSince"`
}

func (yp *YP) Register(getStreamingSince func() utils.NullTime) {
	yp.key = yp.getSavedKey()
	yp.getStreamingSince = getStreamingSince
	yp.ping()

	for range time.Tick(pingInterval) {
		yp.ping()
	}
}

func (yp *YP) ping() {
	url := "http://localhost:8080/" // TODO: Get the real location

	// Don't allow pinging if a live stream is not active.
	streamingSince := yp.getStreamingSince()
	if !streamingSince.Valid {
		return
	}

	log.Println("Pinging YP: ", config.Config.InstanceDetails.Name)

	request := ypPingRequest{
		config.Config.InstanceDetails.Name,
		config.Config.InstanceDetails.Summary,
		url,
		url + "/" + config.Config.InstanceDetails.Logo["large"],
		false,
		yp.key,
		config.Config.InstanceDetails.Tags,
		&streamingSince,
	}

	req, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(pingURL, "application/json", bytes.NewBuffer(req))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	pingResponse := ypPingResponse{}
	json.Unmarshal(body, &pingResponse)

	if pingResponse.Key != yp.key {
		yp.key = pingResponse.Key
		yp.writeSavedKey(pingResponse.Key)
	}
}

func (yp *YP) writeSavedKey(key string) {
	f, err := os.Create(".yp.key")
	defer f.Close()

	if err != nil {
		log.Errorln(err)
		return
	}

	_, err = f.WriteString(key)
	if err != nil {
		log.Errorln(err)
		return
	}
}

func (yp *YP) getSavedKey() string {
	fileBytes, err := ioutil.ReadFile(".yp.key")
	if err != nil {
		return ""
	}

	return string(fileBytes)
}
