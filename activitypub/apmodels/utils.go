package apmodels

import (
	"encoding/json"
	"net/url"
	"path"
	"path/filepath"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/core/data"
	log "github.com/sirupsen/logrus"
)

// MakeRemoteIRIForResource will create an IRI for a remote location.
func MakeRemoteIRIForResource(resourcePath string, host string) (*url.URL, error) {
	generatedURL := "https://" + host
	u, err := url.Parse(generatedURL)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "federation", resourcePath)

	return u, nil
}

// MakeLocalIRIForResource will create an IRI for the local server.
func MakeLocalIRIForResource(resourcePath string) *url.URL {
	host := data.GetServerURL()
	u, err := url.Parse(host)
	if err != nil {
		log.Errorln("unable to parse local IRI url", host, err)
		return nil
	}

	u.Path = path.Join(u.Path, "federation", resourcePath)

	return u
}

// MakeLocalIRIForAccount will return a full IRI for the local server account username.
func MakeLocalIRIForAccount(account string) *url.URL {
	host := data.GetServerURL()
	u, err := url.Parse(host)
	if err != nil {
		log.Errorln("unable to parse local IRI account server url", err)
		return nil
	}

	u.Path = path.Join(u.Path, "federation", "user", account)

	return u
}

// Serialize will serialize an ActivityPub object to a byte slice.
func Serialize(obj vocab.Type) ([]byte, error) {
	var jsonmap map[string]interface{}
	jsonmap, _ = streams.Serialize(obj)
	b, err := json.Marshal(jsonmap)

	return b, err
}

// MakeLocalIRIForStreamURL will return a full IRI for the local server stream url.
func MakeLocalIRIForStreamURL() *url.URL {
	host := data.GetServerURL()
	u, err := url.Parse(host)
	if err != nil {
		log.Errorln("unable to parse local IRI stream url", err)
		return nil
	}

	u.Path = path.Join(u.Path, "/hls/stream.m3u8")

	return u
}

// MakeLocalIRIforLogo will return a full IRI for the local server logo.
func MakeLocalIRIforLogo() *url.URL {
	host := data.GetServerURL()
	u, err := url.Parse(host)
	if err != nil {
		log.Errorln("unable to parse local IRI stream url", err)
		return nil
	}

	u.Path = path.Join(u.Path, "/logo/external")

	return u
}

// GetLogoType will return the rel value for the webfinger response and
// the default static image is of type png.
func GetLogoType() string {
	imageFilename := data.GetLogoPath()
	if imageFilename == "" {
		return "image/png"
	}

	logoType := "image/jpeg"
	if filepath.Ext(imageFilename) == ".svg" {
		logoType = "image/svg+xml"
	} else if filepath.Ext(imageFilename) == ".gif" {
		logoType = "image/gif"
	} else if filepath.Ext(imageFilename) == ".png" {
		logoType = "image/png"
	}
	return logoType
}
