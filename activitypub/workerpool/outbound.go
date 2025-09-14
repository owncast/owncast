package workerpool

import (
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Job struct bundling the ActivityPub and the payload in one struct.
type Job struct {
	request *http.Request
}

var queue chan Job

// Circuit breaker backoff durations for exponential backoff.
var circuitBreakerBackoffDurations = []time.Duration{
	1 * time.Minute,  // 1st failure: 1 minute
	5 * time.Minute,  // 2nd failure: 5 minutes
	15 * time.Minute, // 3rd failure: 15 minutes
	30 * time.Minute, // 4th failure: 30 minutes
	60 * time.Minute, // 5+ failures: 1 hour (max)
}

// httpClient is a configured HTTP client with timeouts and connection limits.
var httpClient *http.Client

// failedDomains tracks domains that are consistently failing with their failure count and backoff time.
var (
	failedDomains      = make(map[string]*domainFailure)
	failedDomainsMutex sync.RWMutex
)

type domainFailure struct {
	count        int
	lastFailed   time.Time
	backoffUntil time.Time
}

// InitOutboundWorkerPool starts n go routines that await ActivityPub jobs.
func InitOutboundWorkerPool(workerPoolSize int) {
	queue = make(chan Job, workerPoolSize)

	// Initialize HTTP client with conservative timeouts and connection limits
	// to prevent resource exhaustion and hanging requests
	httpClient = &http.Client{
		Timeout: 15 * time.Second, // Reduced from 30s for faster failure detection
		Transport: &http.Transport{
			MaxIdleConns:        20,               // Reduced from 100 to limit resource usage
			MaxIdleConnsPerHost: 2,                // Reduced from 10 to be more conservative
			IdleConnTimeout:     10 * time.Second, // Reduced from 30s for faster cleanup
			DisableKeepAlives:   false,
		},
	}

	// start workers
	for i := 1; i <= workerPoolSize; i++ {
		go worker(i, queue)
	}
}

// AddToOutboundQueue will queue up an outbound http request.
func AddToOutboundQueue(req *http.Request) {
	// Check if domain should be skipped due to circuit breaker
	if shouldSkipDomain(req.URL.Host) {
		log.Debugf("Skipping request to %s due to circuit breaker", req.URL.Host)
		return
	}

	select {
	case queue <- Job{req}:
	default:
		log.Debugln("Outbound ActivityPub job queue is full")
		queue <- Job{req} // will block until received by a worker at this point
	}
	log.Tracef("Queued request for ActivityPub destination %s", req.RequestURI)
}

func worker(workerID int, queue <-chan Job) {
	log.Debugf("Started ActivityPub worker %d", workerID)

	for job := range queue {
		if err := sendActivityPubMessageToInbox(job); err != nil {
			log.Errorf("ActivityPub destination %s failed to send Error: %s", job.request.RequestURI, err)
			recordDomainFailure(job.request.URL.Host)
		} else {
			// Reset domain failure count on success
			resetDomainFailure(job.request.URL.Host)
		}
		log.Tracef("Done with ActivityPub destination %s using worker %d", job.request.RequestURI, workerID)
	}
}

func sendActivityPubMessageToInbox(job Job) error {
	resp, err := httpClient.Do(job.request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Consider HTTP 4xx and 5xx as failures for circuit breaker purposes
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

func (e *httpError) Error() string {
	return e.message
}

// shouldSkipDomain checks if a domain should be skipped due to circuit breaker.
func shouldSkipDomain(domain string) bool {
	failedDomainsMutex.RLock()
	defer failedDomainsMutex.RUnlock()

	failure, exists := failedDomains[domain]
	if !exists {
		return false
	}

	// If we're still in backoff period, skip this domain
	return time.Now().Before(failure.backoffUntil)
}

// recordDomainFailure records a failure for a domain and implements exponential backoff.
func recordDomainFailure(domain string) {
	failedDomainsMutex.Lock()
	defer failedDomainsMutex.Unlock()

	failure, exists := failedDomains[domain]
	if !exists {
		failure = &domainFailure{}
		failedDomains[domain] = failure
	}

	failure.count++
	failure.lastFailed = time.Now()

	// Use exponential backoff with pre-defined durations
	backoffIndex := failure.count - 1
	if backoffIndex >= len(circuitBreakerBackoffDurations) {
		backoffIndex = len(circuitBreakerBackoffDurations) - 1
	}

	backoffDuration := circuitBreakerBackoffDurations[backoffIndex]
	failure.backoffUntil = time.Now().Add(backoffDuration)

	log.Warnf("Domain %s failed %d times, backing off for %v", domain, failure.count, backoffDuration)
}

// resetDomainFailure resets the failure count for a domain on successful delivery.
func resetDomainFailure(domain string) {
	failedDomainsMutex.Lock()
	defer failedDomainsMutex.Unlock()

	if failure, exists := failedDomains[domain]; exists && failure.count > 0 {
		log.Debugf("Resetting failure count for domain %s after successful delivery", domain)
		delete(failedDomains, domain)
	}
}
