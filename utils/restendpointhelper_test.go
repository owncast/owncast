package utils

import (
	"strings"
	"testing"
)

func TestGetPatternForRestEndpoint(t *testing.T) {
	expected := "/hello/"
	endpoints := [...]string{"/hello/{param1}", "/hello/{param1}/{param2}", "/hello/{param1}/world/{param2}"}
	for _, endpoint := range endpoints {
		if ep := getPatternForRestEndpoint(endpoint); ep != expected {
			t.Errorf("%s p does not match expected %s", ep, expected)
		}
	}
}

func TestReadParameter(t *testing.T) {
	expected := "world"
	endpoints := [...]string{
		"/hello/{p1}",
		"/hello/cruel/{p1}",
		"/hello/{p1}/my/friend",
		"/hello/{p1}/{p2}/friend",
		"/hello/{p2}/{p3}/{p1}",
		"/{p1}/is/nice",
		"/{p1}/{p1}/{p1}",
	}

	for _, ep := range endpoints {
		v, err := readParameter(ep, strings.Replace(ep, "{p1}", expected, -1), "p1")
		if err != nil {
			t.Errorf("Unexpected error when reading parameter: %s", err.Error())
		}
		if v != expected {
			t.Errorf("'%s' should have returned %s", ep, expected)
		}
	}
}
