// Package workerpool is the durable outbound delivery queue for ActivityPub.
// It persists unsigned activities in SQLite, signs each delivery attempt with a
// fresh Date header, and retries transient failures across process restarts.
package workerpool

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/db"
	apcrypto "github.com/owncast/owncast/services/activitypub/crypto"
	"github.com/owncast/owncast/services/datastore"
	"github.com/owncast/owncast/utils"
)

const (
	deliveryClaimDuration = 30 * time.Second
	deliveryPollInterval  = time.Second
	maxDeliveryAttempts   = int64(16)
	maxStoredErrorLength  = 2000
)

// Delivery is an unsigned ActivityPub activity waiting to be delivered.
// CoalesceKey replaces an older pending delivery to the same inbox when only
// the latest state matters, such as featured-stream live status.
type Delivery struct {
	Inbox        *url.URL
	Payload      []byte
	ActorIRI     *url.URL
	ActivityType string
	CoalesceKey  string
}

// Deps lists the dependencies for the outbound delivery queue.
type Deps struct {
	WorkerPoolSize int
	Datastore      *datastore.Datastore
	Signer         *apcrypto.Signer
}

// Service owns the durable queue dispatcher, delivery workers, HTTP client,
// and per-domain circuit breaker state.
type Service struct {
	workerPoolSize int
	datastore      *datastore.Datastore
	signer         *apcrypto.Signer
	httpClient     *http.Client

	work chan db.ClaimActivityPubDeliveryRow
	wake chan struct{}
	stop chan struct{}
	wg   sync.WaitGroup

	startOnce sync.Once
	stopOnce  sync.Once

	failedDomainsMu sync.RWMutex
	failedDomains   map[string]*domainFailure

	retryDelay func(int64) time.Duration
}

const circuitBreakerFailureThreshold = 3

type domainFailure struct {
	count        int
	backoffUntil time.Time
}

var circuitBreakerBackoffDurations = []time.Duration{
	1 * time.Minute,
	5 * time.Minute,
	15 * time.Minute,
	30 * time.Minute,
	60 * time.Minute,
}

// New constructs an idle durable delivery service. Call Start after the
// signing key has been generated.
func New(deps Deps) *Service {
	return &Service{
		workerPoolSize: deps.WorkerPoolSize,
		datastore:      deps.Datastore,
		signer:         deps.Signer,
		stop:           make(chan struct{}),
		failedDomains:  make(map[string]*domainFailure),
		retryDelay:     defaultRetryDelay,
	}
}

// Start launches the queue dispatcher and delivery workers. It is safe to call
// more than once.
func (s *Service) Start() {
	s.startOnce.Do(func() {
		s.httpClient = utils.GetFederationHTTPClient()
		s.work = make(chan db.ClaimActivityPubDeliveryRow, s.workerPoolSize)
		s.wake = make(chan struct{}, 1)

		s.wg.Add(s.workerPoolSize + 1)
		for i := 1; i <= s.workerPoolSize; i++ {
			go s.worker(i)
		}
		go s.dispatch()
		s.signalWake()
	})
}

// Stop waits for the dispatcher and workers to release their resources.
func (s *Service) Stop() {
	s.stopOnce.Do(func() { close(s.stop) })
	s.wg.Wait()
}

// Enqueue persists a delivery before returning.
func (s *Service) Enqueue(delivery Delivery) error {
	if err := validateDelivery(delivery); err != nil {
		return err
	}
	if err := s.queue(context.Background(), s.datastore.GetQueries(), delivery); err != nil {
		return err
	}
	s.signalWake()
	return nil
}

// EnqueueBatch persists a set of deliveries in one short transaction.
func (s *Service) EnqueueBatch(deliveries []Delivery) error {
	if len(deliveries) == 0 {
		return nil
	}
	tx, err := s.datastore.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("beginning ActivityPub delivery batch: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	queries := db.New(tx)
	for _, delivery := range deliveries {
		if err := validateDelivery(delivery); err != nil {
			return err
		}
		if err := s.queue(context.Background(), queries, delivery); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing ActivityPub delivery batch: %w", err)
	}
	s.signalWake()
	return nil
}

// EnqueueTx persists a delivery inside an existing state-change transaction.
func (s *Service) EnqueueTx(ctx context.Context, tx *sql.Tx, delivery Delivery) error {
	if err := validateDelivery(delivery); err != nil {
		return err
	}
	if err := s.queue(ctx, db.New(tx), delivery); err != nil {
		return err
	}
	s.signalWake()
	return nil
}

func validateDelivery(delivery Delivery) error {
	if delivery.Inbox == nil || delivery.ActorIRI == nil {
		return errors.New("ActivityPub delivery requires inbox and actor IRI")
	}
	if delivery.Inbox.Scheme != "https" || delivery.Inbox.Hostname() == "" {
		return fmt.Errorf("rejecting invalid inbox URL: %s", delivery.Inbox)
	}
	if utils.IsHostnameInternal(delivery.Inbox.Hostname()) {
		return fmt.Errorf("rejecting internal or loopback inbox URL: %s", delivery.Inbox)
	}
	if len(delivery.Payload) == 0 {
		return errors.New("ActivityPub delivery requires a payload")
	}
	if delivery.ActivityType == "" {
		return errors.New("ActivityPub delivery requires an activity type")
	}
	return nil
}

func (s *Service) queue(ctx context.Context, queries *db.Queries, delivery Delivery) error {
	_, err := queries.QueueActivityPubDelivery(ctx, db.QueueActivityPubDeliveryParams{
		Inbox:         delivery.Inbox.String(),
		Payload:       delivery.Payload,
		ActorIri:      delivery.ActorIRI.String(),
		ActivityType:  delivery.ActivityType,
		CoalesceKey:   sql.NullString{String: delivery.CoalesceKey, Valid: delivery.CoalesceKey != ""},
		NextAttemptAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("queueing ActivityPub %s delivery: %w", delivery.ActivityType, err)
	}
	return nil
}

func (s *Service) signalWake() {
	if s.wake == nil {
		return
	}
	select {
	case s.wake <- struct{}{}:
	default:
	}
}

func (s *Service) dispatch() {
	defer s.wg.Done()
	ticker := time.NewTicker(deliveryPollInterval)
	defer ticker.Stop()

	for {
		s.fillWorkers()
		select {
		case <-s.stop:
			return
		case <-s.wake:
		case <-ticker.C:
		}
	}
}

func (s *Service) fillWorkers() {
	for len(s.work) < cap(s.work) {
		job, err := s.datastore.GetQueries().ClaimActivityPubDelivery(context.Background(), db.ClaimActivityPubDeliveryParams{
			ClaimedUntil: sql.NullTime{Time: time.Now().Add(deliveryClaimDuration), Valid: true},
			Now:          time.Now(),
		})
		if errors.Is(err, sql.ErrNoRows) {
			return
		}
		if err != nil {
			log.Errorf("Unable to claim ActivityPub delivery: %v", err)
			return
		}
		select {
		case <-s.stop:
			return
		case s.work <- job:
		}
	}
}

func (s *Service) worker(workerID int) {
	defer s.wg.Done()
	log.Debugf("Started ActivityPub delivery worker %d", workerID)

	for {
		select {
		case <-s.stop:
			return
		case job := <-s.work:
			s.deliver(job)
			s.signalWake()
		}
	}
}

func (s *Service) deliver(job db.ClaimActivityPubDeliveryRow) {
	inbox, err := url.Parse(job.Inbox)
	if err != nil {
		s.fail(job, fmt.Errorf("invalid inbox URL: %w", err))
		return
	}
	domain := inbox.Hostname()

	if until, blocked := s.domainBackoff(domain); blocked {
		s.deferDelivery(job, until)
		return
	}

	actorIRI, err := url.Parse(job.ActorIri)
	if err != nil {
		s.fail(job, fmt.Errorf("invalid actor IRI: %w", err))
		return
	}

	req, err := s.signer.CreateSignedRequest(job.Payload, inbox, actorIRI)
	if err != nil {
		s.fail(job, fmt.Errorf("signing delivery: %w", err))
		return
	}

	attemptErr := s.send(req)
	if attemptErr == nil {
		s.resetDomainFailure(domain)
		s.complete(job)
		return
	}

	if attemptErr.permanent {
		s.fail(job, attemptErr)
		return
	}

	domainUntil := s.recordDomainFailure(domain)
	if job.Attempts >= maxDeliveryAttempts {
		s.fail(job, fmt.Errorf("retry limit reached: %w", attemptErr))
		return
	}

	delay := s.retryDelay(job.Attempts)
	if attemptErr.retryAfter > delay {
		delay = attemptErr.retryAfter
	}
	nextAttempt := time.Now().Add(delay)
	if domainUntil.After(nextAttempt) {
		nextAttempt = domainUntil
	}

	rows, err := s.datastore.GetQueries().RetryActivityPubDelivery(context.Background(), db.RetryActivityPubDeliveryParams{
		NextAttemptAt: nextAttempt,
		LastError:     nullableError(attemptErr),
		ID:            job.ID,
		Revision:      job.Revision,
	})
	if err != nil {
		log.Errorf("Unable to reschedule ActivityPub %s delivery to %s: %v", job.ActivityType, job.Inbox, err)
		return
	}
	if rows == 0 {
		s.releaseSuperseded(job)
		return
	}
	log.Warnf("ActivityPub %s delivery to %s failed on attempt %d; retrying at %s: %v", job.ActivityType, job.Inbox, job.Attempts, nextAttempt.Format(time.RFC3339), attemptErr)
}

type attemptError struct {
	statusCode int
	permanent  bool
	retryAfter time.Duration
	err        error
}

func (e *attemptError) Error() string { return e.err.Error() }
func (e *attemptError) Unwrap() error { return e.err }

func (s *Service) send(req *http.Request) *attemptError {
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return &attemptError{err: err}
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	retryable := resp.StatusCode == http.StatusUnauthorized ||
		resp.StatusCode == http.StatusRequestTimeout ||
		resp.StatusCode == http.StatusTooEarly ||
		resp.StatusCode == http.StatusTooManyRequests ||
		(resp.StatusCode >= http.StatusInternalServerError && resp.StatusCode != http.StatusNotImplemented)

	return &attemptError{
		statusCode: resp.StatusCode,
		permanent:  !retryable,
		retryAfter: parseRetryAfter(resp.Header.Get("Retry-After")),
		err:        fmt.Errorf("HTTP %s", resp.Status),
	}
}

func parseRetryAfter(value string) time.Duration {
	if value == "" {
		return 0
	}
	if seconds, err := strconv.Atoi(value); err == nil {
		if seconds > 0 {
			return time.Duration(seconds) * time.Second
		}
		return 0
	}
	when, err := http.ParseTime(value)
	if err != nil {
		return 0
	}
	if delay := time.Until(when); delay > 0 {
		return delay
	}
	return 0
}

func defaultRetryDelay(attempt int64) time.Duration {
	exponent := attempt - 1
	if exponent < 0 {
		exponent = 0
	}
	if exponent > 8 {
		exponent = 8
	}
	delay := time.Minute * time.Duration(1<<exponent)
	if delay > 6*time.Hour {
		delay = 6 * time.Hour
	}
	// #nosec G404: retry jitter does not protect a secret or security decision.
	return delay + time.Duration(rand.Int64N(int64(delay/2)+1))
}

func (s *Service) complete(job db.ClaimActivityPubDeliveryRow) {
	rows, err := s.datastore.GetQueries().CompleteActivityPubDelivery(context.Background(), db.CompleteActivityPubDeliveryParams{
		ID:       job.ID,
		Revision: job.Revision,
	})
	if err != nil {
		log.Errorf("Unable to complete ActivityPub %s delivery to %s: %v", job.ActivityType, job.Inbox, err)
		return
	}
	if rows == 0 {
		s.releaseSuperseded(job)
	}
}

func (s *Service) deferDelivery(job db.ClaimActivityPubDeliveryRow, until time.Time) {
	rows, err := s.datastore.GetQueries().DeferActivityPubDelivery(context.Background(), db.DeferActivityPubDeliveryParams{
		NextAttemptAt: until,
		ID:            job.ID,
		Revision:      job.Revision,
	})
	if err != nil {
		log.Errorf("Unable to defer ActivityPub %s delivery to %s: %v", job.ActivityType, job.Inbox, err)
		return
	}
	if rows == 0 {
		s.releaseSuperseded(job)
	}
}

func (s *Service) fail(job db.ClaimActivityPubDeliveryRow, deliveryErr error) {
	rows, err := s.datastore.GetQueries().FailActivityPubDelivery(context.Background(), db.FailActivityPubDeliveryParams{
		LastError: nullableError(deliveryErr),
		FailedAt:  sql.NullTime{Time: time.Now(), Valid: true},
		ID:        job.ID,
		Revision:  job.Revision,
	})
	if err != nil {
		log.Errorf("Unable to record failed ActivityPub %s delivery to %s: %v", job.ActivityType, job.Inbox, err)
		return
	}
	if rows == 0 {
		s.releaseSuperseded(job)
		return
	}
	log.Errorf("ActivityPub %s delivery to %s permanently failed after %d attempt(s): %v", job.ActivityType, job.Inbox, job.Attempts, deliveryErr)
}

func (s *Service) releaseSuperseded(job db.ClaimActivityPubDeliveryRow) {
	if err := s.datastore.GetQueries().ReleaseSupersededActivityPubDelivery(context.Background(), db.ReleaseSupersededActivityPubDeliveryParams{
		ID:       job.ID,
		Revision: job.Revision,
	}); err != nil {
		log.Errorf("Unable to release superseded ActivityPub delivery %d: %v", job.ID, err)
	}
}

func nullableError(err error) sql.NullString {
	message := err.Error()
	if len(message) > maxStoredErrorLength {
		message = message[:maxStoredErrorLength]
	}
	return sql.NullString{String: message, Valid: true}
}

func (s *Service) domainBackoff(domain string) (time.Time, bool) {
	s.failedDomainsMu.RLock()
	defer s.failedDomainsMu.RUnlock()

	failure, exists := s.failedDomains[domain]
	if !exists || !time.Now().Before(failure.backoffUntil) {
		return time.Time{}, false
	}
	return failure.backoffUntil, true
}

func (s *Service) recordDomainFailure(domain string) time.Time {
	s.failedDomainsMu.Lock()
	defer s.failedDomainsMu.Unlock()

	failure, exists := s.failedDomains[domain]
	if !exists {
		failure = &domainFailure{}
		s.failedDomains[domain] = failure
	}
	failure.count++

	if failure.count < circuitBreakerFailureThreshold {
		return time.Time{}
	}
	index := failure.count - circuitBreakerFailureThreshold
	if index >= len(circuitBreakerBackoffDurations) {
		index = len(circuitBreakerBackoffDurations) - 1
	}
	backoffDuration := circuitBreakerBackoffDurations[index]
	failure.backoffUntil = time.Now().Add(backoffDuration)
	log.Debugf("Domain %s failed %d times, backing off for %v", domain, failure.count, backoffDuration)
	return failure.backoffUntil
}

func (s *Service) resetDomainFailure(domain string) {
	s.failedDomainsMu.Lock()
	defer s.failedDomainsMu.Unlock()
	delete(s.failedDomains, domain)
}
