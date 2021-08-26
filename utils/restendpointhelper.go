package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// takes the segment pattern of an Url string and returns the segment before the first dynamic REST parameter
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

func mapPatternWithRequestUrl(pattern string, requestUrl string) (map[string]string, error) {
	patternSplit := strings.Split(pattern, "/")
	requestUrlSplit := strings.Split(requestUrl, "/")

	if len(patternSplit) == len(requestUrlSplit) {
		return zip2D(&patternSplit, &requestUrlSplit), nil
	}
	return nil, errors.New("The lenght of pattern and request Url does not match")
}

func readParameter(pattern string, requestUrl string, paramName string) (string, error) {
	all, err := mapPatternWithRequestUrl(pattern, requestUrl)
	if err != nil {
		return "", err
	}

	if value, exists := all[fmt.Sprintf("{%s}", paramName)]; exists {
		return value, nil
	}
	return "", errors.New(fmt.Sprintf("Paramater with name %s not found", paramName))
}

func ReadRestUrlParameter(r *http.Request, parameterName string) (string, error) {
	pattern, has_header := r.Header["OWNCAST_RESTURL_PATTERN"]
	if !has_header {
		return "", errors.New(fmt.Sprintf("This HandlerFunc is not marked as REST-Endpoint. Cannot read Paramet %s from Request", parameterName))
	}

	return readParameter(pattern[0], r.URL.Path, parameterName)
}

func RestEndpoint(pattern string, handler http.HandlerFunc) (string, http.HandlerFunc) {
	baseUrl := getPatternForRestEndpoint(pattern)
	return baseUrl, func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("OWNCAST_RESTURL_PATTERN", pattern)
		handler(w, r)
	}
}
