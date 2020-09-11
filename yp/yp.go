package yp

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"encoding/json"

	"github.com/gabek/owncast/config"
	"github.com/gabek/owncast/models"

	log "github.com/sirupsen/logrus"
)

const pingURL = "http://localhost:8088/ping"
const pingInterval = 4 * time.Minute

var getStatus func() models.Status

//YP is a service for handling listing in the Owncast directory.
type YP struct {
	timer *time.Ticker
	key   string
}

type ypPingResponse struct {
	Key     string `json:"key"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type ypPingRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Logo        string   `json:"logo"`
	NSFW        bool     `json:"nsfw"`
	Key         string   `json:"key"`
	Tags        []string `json:"tags"`
}

// NewYP creates a new instance of the YP service handler
func NewYP(getStatusFunc func() models.Status) *YP {
	getStatus = getStatusFunc
	return &YP{}
}

// Start is run when a live stream begins to start pinging YP
func (yp *YP) Start() {
	yp.timer = time.NewTicker(pingInterval)
	yp.key = yp.getSavedKey()

	go func() {
		for {
			select {
			case <-yp.timer.C:
				yp.ping()
			}
		}
	}()

	yp.ping()
}

// Stop stops the pinging of YP
func (yp *YP) Stop() {
	yp.timer.Stop()
}

func (yp *YP) ping() {
	url := config.Config.YP.YPServiceURL

	log.Traceln("Pinging YP as: ", config.Config.InstanceDetails.Name)

	request := ypPingRequest{
		config.Config.InstanceDetails.Name,
		config.Config.InstanceDetails.Summary,
		url,
		url + "/" + config.Config.InstanceDetails.Logo["large"],
		false,
		yp.key,
		config.Config.InstanceDetails.Tags,
	}

	req, err := json.Marshal(request)
	if err != nil {
		log.Errorln(err)
		return
	}

	resp, err := http.Post(pingURL, "application/json", bytes.NewBuffer(req))
	if err != nil {
		log.Errorln(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorln(err)
	}

	pingResponse := ypPingResponse{}
	json.Unmarshal(body, &pingResponse)

	if !pingResponse.Success {
		log.Errorln("YP Ping:", pingResponse.Error)
		return
	}

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
