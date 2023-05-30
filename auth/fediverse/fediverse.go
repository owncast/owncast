package fediverse

import (
	"crypto/rand"
	"errors"
	"io"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// OTPRegistration represents a single OTP request.
type OTPRegistration struct {
	Timestamp       time.Time
	UserID          string
	UserDisplayName string
	Code            string
	Account         string
}

// Key by access token to limit one OTP request for a person
// to be active at a time.
var (
	pendingAuthRequests = make(map[string]OTPRegistration)
	lock                = sync.Mutex{}
)

const (
	registrationTimeout = time.Minute * 10
	maxPendingRequests  = 1000
)

func init() {
	go setupExpiredRequestPruner()
}

// Clear out any pending requests that have been pending for greater than
// the specified timeout value.
func setupExpiredRequestPruner() {
	pruneExpiredRequestsTimer := time.NewTicker(registrationTimeout)

	for range pruneExpiredRequestsTimer.C {
		lock.Lock()
		log.Debugln("Pruning expired OTP requests.")
		for k, v := range pendingAuthRequests {
			if time.Since(v.Timestamp) > registrationTimeout {
				delete(pendingAuthRequests, k)
			}
		}
		lock.Unlock()
	}
}

// RegisterFediverseOTP will start the OTP flow for a user, creating a new
// code and returning it to be sent to a destination.
func RegisterFediverseOTP(accessToken, userID, userDisplayName, account string) (OTPRegistration, bool, error) {
	request, requestExists := pendingAuthRequests[accessToken]

	// If a request is already registered and has not expired then return that
	// existing request.
	if requestExists && time.Since(request.Timestamp) < registrationTimeout {
		return request, false, nil
	}

	lock.Lock()
	defer lock.Unlock()

	if len(pendingAuthRequests)+1 > maxPendingRequests {
		return request, false, errors.New("Please try again later. Too many pending requests.")
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

	return r, true, nil
}

// ValidateFediverseOTP will verify a OTP code for a auth request.
func ValidateFediverseOTP(accessToken, code string) (bool, *OTPRegistration) {
	request, ok := pendingAuthRequests[accessToken]

	if !ok || request.Code != code || time.Since(request.Timestamp) > registrationTimeout {
		return false, nil
	}

	lock.Lock()
	defer lock.Unlock()

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
