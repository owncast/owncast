package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const restURLPatternHeaderKey = "Owncast-Resturl-Pattern"

// takes the segment pattern of an Url string and returns the segment before the first dynamic REST parameter.
func getPatternForRestEndpoint(pattern string) string {
	firstIndex := strings.Index(pattern, "/{")
	if firstIndex == -1 {
		return pattern
	}

	return strings.TrimRight(pattern[:firstIndex], "/") + "/"
}

func zip2D(iterable1 *[]string, iterable2 *[]string) map[string]string {
	var dict = make(map[string]string)
	for index, key := range *iterable1 {
		dict[key] = (*iterable2)[index]
	}
	return dict
}

func mapPatternWithRequestURL(pattern string, requestURL string) (map[string]string, error) {
	patternSplit := strings.Split(pattern, "/")
	requestURLSplit := strings.Split(requestURL, "/")

	if len(patternSplit) == len(requestURLSplit) {
		return zip2D(&patternSplit, &requestURLSplit), nil
	}
	return nil, errors.New("the length of pattern and request Url does not match")
}

func readParameter(pattern string, requestURL string, paramName string) (string, error) {
	all, err := mapPatternWithRequestURL(pattern, requestURL)
	if err != nil {
		return "", err
	}

	if value, exists := all[fmt.Sprintf("{%s}", paramName)]; exists {
		return value, nil
	}
	return "", fmt.Errorf("parameter with name %s not found", paramName)
}

// ReadRestURLParameter will return the parameter from the request of the requested name.
func ReadRestURLParameter(r *http.Request, parameterName string) (string, error) {
	pattern, found := r.Header[restURLPatternHeaderKey]
	if !found {
		return "", fmt.Errorf("this HandlerFunc is not marked as REST-Endpoint. Cannot read Parameter '%s' from Request", parameterName)
	}

	return readParameter(pattern[0], r.URL.Path, parameterName)
}

// RestEndpoint wraps a handler to use the rest endpoint helper.
func RestEndpoint(pattern string, handler http.HandlerFunc) (string, http.HandlerFunc) {
	baseURL := getPatternForRestEndpoint(pattern)
	return baseURL, func(w http.ResponseWriter, r *http.Request) {
		r.Header[restURLPatternHeaderKey] = []string{pattern}
		handler(w, r)
	}
}
