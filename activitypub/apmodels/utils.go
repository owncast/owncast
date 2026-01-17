package apmodels

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"path/filepath"

	"github.com/go-fed/activity/streams"
	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/persistence/configrepository"
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
	configRepository := configrepository.Get()

	host := configRepository.GetServerURL()
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
	configRepository := configrepository.Get()

	host := configRepository.GetServerURL()
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
	configRepository := configrepository.Get()

	host := configRepository.GetServerURL()
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
	configRepository := configrepository.Get()

	host := configRepository.GetServerURL()
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
	configRepository := configrepository.Get()

	imageFilename := configRepository.GetLogoPath()
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

// ErrMissingIRI is returned when an IRI cannot be extracted from an ActivityStreams property.
var ErrMissingIRI = fmt.Errorf("missing IRI")

// GetIRIFromActorProperty safely extracts the IRI from an ActivityStreamsActorProperty.
// Returns the IRI and nil error on success, or nil and an error if the IRI cannot be extracted.
func GetIRIFromActorProperty(actor vocab.ActivityStreamsActorProperty) (*url.URL, error) {
	if actor == nil || actor.Empty() || actor.Len() == 0 {
		return nil, fmt.Errorf("%w: actor property is empty or nil", ErrMissingIRI)
	}
	first := actor.At(0)
	if first == nil {
		return nil, fmt.Errorf("%w: actor property first element is nil", ErrMissingIRI)
	}
	iri := first.GetIRI()
	if iri == nil {
		return nil, fmt.Errorf("%w: actor IRI is nil", ErrMissingIRI)
	}
	return iri, nil
}

// GetIRIStringFromActorProperty safely extracts the IRI string from an ActivityStreamsActorProperty.
// Returns the IRI string and nil error on success, or empty string and an error if extraction fails.
func GetIRIStringFromActorProperty(actor vocab.ActivityStreamsActorProperty) (string, error) {
	iri, err := GetIRIFromActorProperty(actor)
	if err != nil {
		return "", err
	}
	return iri.String(), nil
}

// GetIRIFromObjectProperty safely extracts the IRI from an ActivityStreamsObjectProperty.
// Returns the IRI and nil error on success, or nil and an error if the IRI cannot be extracted.
func GetIRIFromObjectProperty(object vocab.ActivityStreamsObjectProperty) (*url.URL, error) {
	if object == nil || object.Len() == 0 {
		return nil, fmt.Errorf("%w: object property is empty or nil", ErrMissingIRI)
	}
	first := object.At(0)
	if first == nil {
		return nil, fmt.Errorf("%w: object property first element is nil", ErrMissingIRI)
	}
	iri := first.GetIRI()
	if iri == nil {
		return nil, fmt.Errorf("%w: object IRI is nil", ErrMissingIRI)
	}
	return iri, nil
}

// GetIRIStringFromObjectProperty safely extracts the IRI string from an ActivityStreamsObjectProperty.
// Returns the IRI string and nil error on success, or empty string and an error if extraction fails.
func GetIRIStringFromObjectProperty(object vocab.ActivityStreamsObjectProperty) (string, error) {
	iri, err := GetIRIFromObjectProperty(object)
	if err != nil {
		return "", err
	}
	return iri.String(), nil
}

// GetIRIFromJSONLDIdProperty safely extracts the IRI from a JSONLDIdProperty.
// Returns the IRI and nil error on success, or nil and an error if the IRI cannot be extracted.
func GetIRIFromJSONLDIdProperty(id vocab.JSONLDIdProperty) (*url.URL, error) {
	if id == nil {
		return nil, fmt.Errorf("%w: JSONLD id property is nil", ErrMissingIRI)
	}
	iri := id.GetIRI()
	if iri == nil {
		return nil, fmt.Errorf("%w: JSONLD id IRI is nil", ErrMissingIRI)
	}
	return iri, nil
}

// GetIRIStringFromJSONLDIdProperty safely extracts the IRI string from a JSONLDIdProperty.
// Returns the IRI string and nil error on success, or empty string and an error if extraction fails.
func GetIRIStringFromJSONLDIdProperty(id vocab.JSONLDIdProperty) (string, error) {
	iri, err := GetIRIFromJSONLDIdProperty(id)
	if err != nil {
		return "", err
	}
	return iri.String(), nil
}

// GetPublicKeyPem safely extracts the public key PEM from a W3IDSecurityV1PublicKey.
// Returns the PEM string and nil error on success, or empty string and an error if extraction fails.
func GetPublicKeyPem(publicKey vocab.W3IDSecurityV1PublicKey) (string, error) {
	if publicKey == nil {
		return "", fmt.Errorf("public key is nil")
	}
	pemProp := publicKey.GetW3IDSecurityV1PublicKeyPem()
	if pemProp == nil {
		return "", fmt.Errorf("public key PEM property is nil")
	}
	return pemProp.Get(), nil
}

// IsFirstObjectActivityStreamsPerson safely checks if the first element of an
// ActivityStreamsObjectProperty is an ActivityStreamsPerson.
// Returns false if the object is nil, empty, or the first element is not a Person.
func IsFirstObjectActivityStreamsPerson(object vocab.ActivityStreamsObjectProperty) bool {
	if object == nil || object.Len() == 0 {
		return false
	}
	first := object.At(0)
	if first == nil {
		return false
	}
	return first.IsActivityStreamsPerson()
}

// GetImageFromIcon safely extracts the image URL from an ActivityStreamsIconProperty.
// Returns the URL and nil error on success, or nil and nil if the icon is not present or invalid.
// This handles the common pattern of icon -> image -> url -> iri.
func GetImageFromIcon(icon vocab.ActivityStreamsIconProperty) *url.URL {
	if icon == nil || icon.Empty() {
		return nil
	}
	first := icon.At(0)
	if first == nil {
		return nil
	}
	image := first.GetActivityStreamsImage()
	if image == nil {
		return nil
	}
	urlProp := image.GetActivityStreamsUrl()
	if urlProp == nil {
		return nil
	}
	begin := urlProp.Begin()
	if begin == nil {
		return nil
	}
	return begin.GetIRI()
}

// GetHostnameFromJSONLDId safely extracts the hostname from a JSONLDIdProperty.
// Returns the hostname string, or empty string if extraction fails.
func GetHostnameFromJSONLDId(id vocab.JSONLDIdProperty) string {
	if id == nil {
		return ""
	}
	iri := id.GetIRI()
	if iri == nil {
		return ""
	}
	return iri.Hostname()
}
