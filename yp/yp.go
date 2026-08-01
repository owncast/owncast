// Package yp is the Owncast directory listing service. Construct via
// New(Deps) and call Start when a stream goes live to begin pinging the
// directory.
package yp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/models"
	"github.com/owncast/owncast/persistence/configrepository"

	log "github.com/sirupsen/logrus"
)

const pingInterval = 4 * time.Minute

// YP is a service for handling listing in the Owncast directory.
type YP struct {
	timer *time.Ticker

	// getStatus returns the current stream status; consulted on each
	// ping cycle to skip pings while offline.
	getStatus func() models.Status

	// configRepository provides directory-listing settings (enabled flag,
	// server URL/name, registration key, etc.) read on each ping and on
	// every YP API response.
	configRepository configrepository.ConfigRepository
	// directoryAvailable is controlled by the active auth gate. A blocked
	// stream-status policy disables both pings and the public YP response.
	directoryAvailable func() bool
	inErrorState       bool
}

// Deps lists the explicit dependencies for the YP service.
type Deps struct {
	GetStatus        func() models.Status
	ConfigRepository configrepository.ConfigRepository
}

type ypPingResponse struct {
	Key       string `json:"key"`
	Error     string `json:"error"`
	ErrorCode int    `json:"errorCode"`
	Success   bool   `json:"success"`
}

type ypPingRequest struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

// New constructs a new instance of the YP service handler.
func New(deps Deps) *YP {
	return &YP{
		getStatus:        deps.GetStatus,
		configRepository: deps.ConfigRepository,
	}
}

// SetGetStatus wires (or rewires) the stream-status callback. Exists
// because yp must be constructed before stream (stream takes yp via
// Deps), but yp needs stream.GetStatus to skip pings while offline.
// main.go constructs yp first with a nil callback, then fills it in once
// streamSvc exists. Must be called before Start.
func (yp *YP) SetGetStatus(fn func() models.Status) {
	yp.getStatus = fn
}

// SetDirectoryAvailable wires the auth-gate policy after plugin host startup.
func (yp *YP) SetDirectoryAvailable(fn func() bool) {
	yp.directoryAvailable = fn
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

// shouldPing reports whether this instance is currently in a state where
// announcing itself to the directory makes sense: the gate's access policy
// allows directory listing, the operator has listing switched on, the stream
// is live, and the configured public URL is usable.
func (yp *YP) shouldPing() (instanceURL string, ok bool) {
	if yp.directoryAvailable != nil && !yp.directoryAvailable() {
		return "", false
	}
	if !yp.configRepository.GetDirectoryEnabled() {
		return "", false
	}

	// Hack: Don't allow ping'ing when offline.
	// It shouldn't even be trying to, but on some instances the ping timer isn't stopping.
	if !yp.getStatus().Online {
		return "", false
	}

	myInstanceURL := yp.configRepository.GetServerURL()
	if myInstanceURL == "" {
		log.Warnln("Server URL not set in the configuration. Directory access is disabled until this is set.")
		return "", false
	}
	if !isURL(myInstanceURL) {
		if !yp.inErrorState {
			log.Warnln("YP Error: unable to use", myInstanceURL, "as a public instance URL. Fix this value in your configuration.")
		}
		yp.inErrorState = true
		return "", false
	}
	return myInstanceURL, true
}

func (yp *YP) ping() {
	myInstanceURL, ok := yp.shouldPing()
	if !ok {
		return
	}

	key := yp.configRepository.GetDirectoryRegistrationKey()

	log.Traceln("Pinging YP as: ", yp.configRepository.GetServerName(), "with key", key)

	request := ypPingRequest{
		Key: key,
		URL: myInstanceURL,
	}

	req, err := json.Marshal(request)
	if err != nil {
		log.Errorln(err)
		return
	}

	pingURL := config.GetDefaults().YPServer + "/api/ping"
	resp, err := http.Post(pingURL, "application/json", bytes.NewBuffer(req)) //nolint
	if err != nil {
		log.Errorln(err)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorln(err)
	}

	pingResponse := ypPingResponse{}
	if err := json.Unmarshal(body, &pingResponse); err != nil {
		log.Errorln(err)
	}

	if !pingResponse.Success {
		if !yp.inErrorState {
			log.Warnln("YP Ping error returned from service:", pingResponse.Error)
		}
		yp.inErrorState = true
		return
	}

	yp.inErrorState = false

	if pingResponse.Key != key {
		if err := yp.configRepository.SetDirectoryRegistrationKey(key); err != nil {
			log.Errorln("unable to save directory key:", err)
		}
	}
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
