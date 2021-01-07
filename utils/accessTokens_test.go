package utils

import (
	"testing"
)

func TestCreateAccessToken(t *testing.T) {
	_, err := GenerateAccessToken()
	if err != nil {
		t.Error(err)
	}
}
