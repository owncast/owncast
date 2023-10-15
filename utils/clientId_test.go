package utils

import (
	"net/http/httptest"
	"testing"
)

func TestGetIPAddressFromRequest(t *testing.T) {

	expectedResult := "203.0.113.195"

	request := httptest.NewRequest("GET", "/", nil)
	request.RemoteAddr = "203.0.113.195:41237"
	//First Test without X-FORWARDED-FOR header
	result := GetIPAddressFromRequest(request)
	if result != expectedResult {
		t.Errorf("Remote address only test failed: Expected %s, got %s", expectedResult, result)
	}

	//Test with X-FORWARDED-FOR header
	request.Header.Set("X-FORWARDED-FOR header", "203.0.113.195:41237, 198.51.100.100:38523, 140.248.67.176:12345")
	result = GetIPAddressFromRequest(request)
	if result != expectedResult {
		t.Errorf("X-FORWARD-FOR header test failed: Expected %s, got %s", expectedResult, result)
	}

}
