package models

import (
	"net/url"
	"path"

	"github.com/owncast/owncast/core/data"
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
		panic(err)
	}

	u.Path = path.Join(u.Path, "federation", resourcePath)

	return u
}

func MakeLocalIRIForAccount(account string) *url.URL {
	host := data.GetServerURL()
	u, err := url.Parse(host)
	if err != nil {
		panic(err)
	}

	u.Path = path.Join(u.Path, "federation", "user", account)

	return u

}
