package fediverse

import (
	"crypto/rand"
	"io"
	"strings"
	"time"
)

// OTPRegistration represents a single OTP request.
type OTPRegistration struct {
	UserID          string
	UserDisplayName string
	Code            string
	Account         string
	Timestamp       time.Time
}

// Key by access token to limit one OTP request for a person
// to be active at a time.
var pendingAuthRequests = make(map[string]OTPRegistration)

const registrationTimeout = time.Minute * 10

// RegisterFediverseOTP will start the OTP flow for a user, creating a new
// code and returning it to be sent to a destination.
func RegisterFediverseOTP(accessToken, userID, userDisplayName, account string) (OTPRegistration, bool) {
	request, requestExists := pendingAuthRequests[accessToken]

	// If a request is already registered and has not expired then return that
	// existing request.
	if requestExists && time.Since(request.Timestamp) < registrationTimeout {
		return request, false
	}

	code, _ := createCode()
	r := OTPRegistration{
		Code:            code,
		UserID:          userID,
		UserDisplayName: userDisplayName,
		Account:         strings.ToLower(account),
		Timestamp:       time.Now(),
	}
	pendingAuthRequests[accessToken] = r

	return r, true
}

// ValidateFediverseOTP will verify a OTP code for a auth request.
func ValidateFediverseOTP(accessToken, code string) (bool, *OTPRegistration) {
	request, ok := pendingAuthRequests[accessToken]

	if !ok || request.Code != code || time.Since(request.Timestamp) > registrationTimeout {
		return false, nil
	}

	delete(pendingAuthRequests, accessToken)
	return true, &request
}

func createCode() (string, error) {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	digits := 6
	b := make([]byte, digits)
	n, err := io.ReadAtLeast(rand.Reader, b, digits)
	if n != digits {
		return "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b), nil
}
