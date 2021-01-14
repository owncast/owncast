package yp

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"encoding/json"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/data"
	"github.com/owncast/owncast/models"

	log "github.com/sirupsen/logrus"
)

const pingInterval = 4 * time.Minute

var getStatus func() models.Status
var _inErrorState = false

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

// NewYP creates a new instance of the YP service handler.
func NewYP(getStatusFunc func() models.Status) *YP {
	getStatus = getStatusFunc
	return &YP{}
}

// Start is run when a live stream begins to start pinging YP.
func (yp *YP) Start() {
	yp.timer = time.NewTicker(pingInterval)
	for range yp.timer.C {
		yp.ping()
	}

	yp.ping()
}

// Stop stops the pinging of YP.
func (yp *YP) Stop() {
	yp.timer.Stop()
}

func (yp *YP) ping() {
	myInstanceURL := data.GetServerURL()
	if myInstanceURL == "" {
		log.Warnln("Server URL not set in the configuration. Directory access is disabled until this is set.")
		return
	}
	isValidInstanceURL := isUrl(myInstanceURL)
	if myInstanceURL == "" || !isValidInstanceURL {
		if !_inErrorState {
			log.Warnln("YP Error: unable to use", myInstanceURL, "as a public instance URL. Fix this value in your configuration.")
		}
		_inErrorState = true
		return
	}

	key := data.GetDirectoryRegistrationKey()

	log.Traceln("Pinging YP as: ", data.GetServerName(), "with key", key)

	request := ypPingRequest{
		Key: key,
		URL: myInstanceURL,
	}

	req, err := json.Marshal(request)
	if err != nil {
		log.Errorln(err)
		return
	}

	pingURL := config.GetDefaults().YPServer + "/ping"
	resp, err := http.Post(pingURL, "application/json", bytes.NewBuffer(req)) //nolint
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
	err = json.Unmarshal(body, &pingResponse)
	if err != nil {
		log.Errorln(err)
	}

	if !pingResponse.Success {
		if !_inErrorState {
			log.Warnln("YP Ping error returned from service:", pingResponse.Error)
		}
		_inErrorState = true
		return
	}

	_inErrorState = false

	if pingResponse.Key != key {
		data.SetDirectoryRegistrationKey(key)
	}
}

// DisplayInstructions will let the user know they are not in the directory by default and
// how they can enable the feature.
func DisplayInstructions() {
	text := "Your instance can be listed on the Owncast directory at http://directory.owncast.online by enabling YP in your config.  Learn more at https://directory.owncast.online/get-listed."
	log.Debugln(text)
}
func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
