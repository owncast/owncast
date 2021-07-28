package models

import (
	"net/url"
	"path"
)

type ResourceURL = string

func MakeURLForResource(resource string, host string) (ResourceURL, error) {
	generatedURL := "https://" + host
	u, err := url.Parse(generatedURL)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, "federation", resource)

	return u.String(), nil
}
