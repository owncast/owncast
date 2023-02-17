package utils

import (
	"testing"
)

func TestCreateAccessToken(t *testing.T) {
	if _, err := GenerateAccessToken(); err != nil {
		t.Error(err)
	}
}
