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

const (
	registrationTimeout = time.Minute * 10
	maxPendingRequests  = 1000
)

// Service bundles the per-instance state for the Fediverse OTP flow.
// Construct once in main() with New() and call Start() to launch the
// expired-request pruner; inject into the handlers that need to register
// and validate OTPs.
type Service struct {
	// Key by access token to limit one OTP request for a person
	// to be active at a time.
	pendingAuthRequests map[string]OTPRegistration
	lock                sync.Mutex
}

// New constructs a Service.
func New() *Service {
	return &Service{
		pendingAuthRequests: make(map[string]OTPRegistration),
	}
}

// Start launches the background goroutine that prunes expired OTP
// requests. Safe to call exactly once after New().
func (s *Service) Start() {
	go s.runExpiredRequestPruner()
}

// runExpiredRequestPruner clears out any pending requests that have been
// pending for greater than the specified timeout value.
func (s *Service) runExpiredRequestPruner() {
	pruneExpiredRequestsTimer := time.NewTicker(registrationTimeout)

	for range pruneExpiredRequestsTimer.C {
		s.lock.Lock()
		log.Debugln("Pruning expired OTP requests.")
		for k, v := range s.pendingAuthRequests {
			if time.Since(v.Timestamp) > registrationTimeout {
				delete(s.pendingAuthRequests, k)
			}
		}
		s.lock.Unlock()
	}
}

// RegisterFediverseOTP will start the OTP flow for a user, creating a new
// code and returning it to be sent to a destination.
func (s *Service) RegisterFediverseOTP(accessToken, userID, userDisplayName, account string) (OTPRegistration, bool, error) {
	account = strings.ToLower(account)

	s.lock.Lock()
	defer s.lock.Unlock()

	existing, exists := s.pendingAuthRequests[accessToken]

	// Block a duplicate request for the same account so we don't re-send a code,
	// but allow a request for a different account to replace a pending one so a
	// user who entered the wrong address can immediately retry.
	if exists && time.Since(existing.Timestamp) < registrationTimeout && existing.Account == account {
		return existing, false, nil
	}

	// Only the global cap applies when adding a new entry; replacing one does
	// not grow the map.
	if !exists && len(s.pendingAuthRequests)+1 > maxPendingRequests {
		return OTPRegistration{}, false, errors.New("too many pending requests, please try again later")
	}

	code, _ := createCode()
	r := OTPRegistration{
		Code:            code,
		UserID:          userID,
		UserDisplayName: userDisplayName,
		Account:         account,
		Timestamp:       time.Now(),
	}
	s.pendingAuthRequests[accessToken] = r

	return r, true, nil
}

// ValidateFediverseOTP will verify a OTP code for a auth request.
func (s *Service) ValidateFediverseOTP(accessToken, code string) (bool, *OTPRegistration) {
	s.lock.Lock()
	defer s.lock.Unlock()

	request, ok := s.pendingAuthRequests[accessToken]
	if !ok || request.Code != code || time.Since(request.Timestamp) > registrationTimeout {
		return false, nil
	}

	delete(s.pendingAuthRequests, accessToken)
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
