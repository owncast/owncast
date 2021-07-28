package models

import (
	"net/url"
	"path"
)

type ResourceURL = string

func MakeURLForResource(resource string, host string) (*url.URL, error) {
	generatedURL := "https://" + host
	u, err := url.Parse(generatedURL)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, "federation", resource)

	return u, nil
}
