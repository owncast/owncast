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

const pingInterval = 4 * time.Minute

var getStatus func() models.Status

//YP is a service for handling listing in the Owncast directory.
type YP struct {
	timer *time.Ticker
}

type ypPingResponse struct {
	Key       string `json:"key"`
	Success   bool   `json:"success"`
	Error     string `json:"error"`
	ErrorCode int    `json:"errorCode"`
}

type ypPingRequest struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

// NewYP creates a new instance of the YP service handler
func NewYP(getStatusFunc func() models.Status) *YP {
	getStatus = getStatusFunc
	return &YP{}
}

// Start is run when a live stream begins to start pinging YP
func (yp *YP) Start() {
	yp.timer = time.NewTicker(pingInterval)

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
	myInstanceURL := config.Config.YP.InstanceURL
	key := yp.getSavedKey()

	log.Traceln("Pinging YP as: ", config.Config.InstanceDetails.Name)

	request := ypPingRequest{
		Key: key,
		URL: myInstanceURL,
	}

	req, err := json.Marshal(request)
	if err != nil {
		log.Errorln(err)
		return
	}

	pingURL := config.Config.YP.YPServiceURL + "/ping"
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
		log.Errorln("YP Ping error:", pingResponse.Error)
		return
	}

	if pingResponse.Key != key {
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

// DisplayInstructions will let the user know they are not in the directory by default and
// how they can enable the feature.
func DisplayInstructions() {
	text := "Your instance can be listed on the Owncast directory at http://something.something by enabling YP in your config.  Learn more at http://something.something."
	log.Infoln(text)
}
