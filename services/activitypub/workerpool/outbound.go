// Package workerpool is the outbound HTTP delivery pool for the
// ActivityPub subsystem. It maintains a bounded worker pool that posts
// signed activities to follower inboxes, with per-domain circuit
// breaking on repeated failure.
package workerpool

import (
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/utils"
)

// Job bundles a single outbound HTTP request for a worker.
type Job struct {
	request *http.Request
}

// circuitBreakerBackoffDurations is the exponential backoff schedule
// applied to a domain after consecutive delivery failures.
var circuitBreakerBackoffDurations = []time.Duration{
	1 * time.Minute,
	5 * time.Minute,
	15 * time.Minute,
	30 * time.Minute,
	60 * time.Minute,
}

// Service owns the worker pool, HTTP client, and per-domain circuit
// breaker state. Construct with New(workerPoolSize) then call Start to
// spin up workers.
type Service struct {
	workerPoolSize int

	queue      chan Job
	httpClient *http.Client

	failedDomainsMu sync.RWMutex
	failedDomains   map[string]*domainFailure
}

type domainFailure struct {
	count        int
	lastFailed   time.Time
	backoffUntil time.Time
}

// New constructs an idle Service sized for the given worker count. Call
// Start to bind the HTTP client and launch the worker goroutines.
func New(workerPoolSize int) *Service {
	return &Service{
		workerPoolSize: workerPoolSize,
		failedDomains:  make(map[string]*domainFailure),
	}
}

// Start initializes the HTTP client and worker goroutines. Safe to call
// once; subsequent calls reset the pool.
func (s *Service) Start() {
	// Use a larger buffer to decouple request creation from processing.
	// This prevents SendToFollowers from blocking when many followers
	// need updates.
	const minQueueBuffer = 500
	queueBuffer := s.workerPoolSize * 10
	if queueBuffer < minQueueBuffer {
		queueBuffer = minQueueBuffer
	}
	s.queue = make(chan Job, queueBuffer)

	// HTTP client with retry logic for transient failures (502/503/504).
	s.httpClient = utils.GetRetryableHTTPClient()

	for i := 1; i <= s.workerPoolSize; i++ {
		go s.worker(i)
	}
}

// AddToOutboundQueue queues an outbound HTTP request for delivery.
func (s *Service) AddToOutboundQueue(req *http.Request) {
	if s.ShouldSkipDomain(req.URL.Host) {
		log.Debugf("Skipping request to %s due to circuit breaker", req.URL.Host)
		return
	}

	select {
	case s.queue <- Job{req}:
	default:
		log.Debugln("Outbound ActivityPub job queue is full")
		s.queue <- Job{req} // blocks until a worker drains
	}
	log.Tracef("Queued request for ActivityPub destination %s", req.RequestURI)
}

func (s *Service) worker(workerID int) {
	log.Debugf("Started ActivityPub worker %d", workerID)

	for job := range s.queue {
		if err := s.sendActivityPubMessageToInbox(job); err != nil {
			log.Errorf("ActivityPub destination %s failed to send Error: %s", job.request.RequestURI, err)
			s.recordDomainFailure(job.request.URL.Host)
		} else {
			s.resetDomainFailure(job.request.URL.Host)
		}
		log.Tracef("Done with ActivityPub destination %s using worker %d", job.request.RequestURI, workerID)
	}
}

func (s *Service) sendActivityPubMessageToInbox(job Job) error {
	resp, err := s.httpClient.Do(job.request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// HTTP 4xx and 5xx count as failures for circuit-breaker purposes.
	if resp.StatusCode >= 400 {
		return &httpError{statusCode: resp.StatusCode, message: resp.Status}
	}
	return nil
}

// httpError represents an HTTP error response.
type httpError struct {
	statusCode int
	message    string
}

func (e *httpError) Error() string { return e.message }

// ShouldSkipDomain reports whether the given domain is currently inside
// its circuit-breaker backoff window.
func (s *Service) ShouldSkipDomain(domain string) bool {
	s.failedDomainsMu.RLock()
	defer s.failedDomainsMu.RUnlock()

	failure, exists := s.failedDomains[domain]
	if !exists {
		return false
	}
	return time.Now().Before(failure.backoffUntil)
}

func (s *Service) recordDomainFailure(domain string) {
	s.failedDomainsMu.Lock()
	defer s.failedDomainsMu.Unlock()

	failure, exists := s.failedDomains[domain]
	if !exists {
		failure = &domainFailure{}
		s.failedDomains[domain] = failure
	}

	failure.count++
	failure.lastFailed = time.Now()

	backoffIndex := failure.count - 1
	if backoffIndex >= len(circuitBreakerBackoffDurations) {
		backoffIndex = len(circuitBreakerBackoffDurations) - 1
	}
	backoffDuration := circuitBreakerBackoffDurations[backoffIndex]
	failure.backoffUntil = time.Now().Add(backoffDuration)

	log.Warnf("Domain %s failed %d times, backing off for %v", domain, failure.count, backoffDuration)
}

func (s *Service) resetDomainFailure(domain string) {
	s.failedDomainsMu.Lock()
	defer s.failedDomainsMu.Unlock()

	if failure, exists := s.failedDomains[domain]; exists && failure.count > 0 {
		log.Debugf("Resetting failure count for domain %s after successful delivery", domain)
		delete(s.failedDomains, domain)
	}
}
