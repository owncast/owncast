// Package scheduler is the application's central home for recurring
// background work: a registry of named jobs running on fixed intervals,
// replacing the ad-hoc time.NewTicker loops that used to be scattered
// across services.
//
// The scheduling engine is gocron (github.com/go-co-op/gocron/v2), a
// widely-used, actively-maintained library, wrapped in this thin service so
// consumers depend on one small Owncast-shaped API instead of a third-party
// surface, and so every job gets the same conventions: overlap protection
// (a run that outlives its interval causes the next tick to be skipped,
// never stacked), panic recovery, and a shutdown that waits for in-flight
// runs.
//
// Jobs are expected to be idempotent sweeps that read their own state and
// tolerate rerunning after a restart, the same contract the old tickers
// had. Interval-based on purpose: every recurring internal task Owncast
// has is "every N". gocron supports crontab and wall-clock schedules too,
// so a RegisterCron variant can join Register if a real calendar-time need
// ever appears.
package scheduler

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	log "github.com/sirupsen/logrus"
)

// Job is one recurring unit of background work.
type Job struct {
	// Name uniquely identifies the job in logs and status listings.
	Name string
	// Interval is how often the job runs. Must be positive.
	Interval time.Duration
	// RunAtStart also runs the job once immediately when the scheduler
	// starts instead of waiting a full interval for the first run.
	RunAtStart bool
	// Run does the work. It should honor ctx, which is cancelled when the
	// scheduler stops.
	Run func(ctx context.Context)
}

// JobStatus is a point-in-time view of a registered job, for logging and
// future admin or debug surfaces.
type JobStatus struct {
	Name      string    `json:"name"`
	LastRunAt time.Time `json:"lastRunAt"`
	NextRunAt time.Time `json:"nextRunAt"`
}

// Service runs registered jobs until stopped. Construct with New in
// main.go and inject into anything that needs to register work.
type Service struct {
	engine gocron.Scheduler
	ctx    context.Context
	cancel context.CancelFunc

	stoppedMutex sync.Mutex
	stopped      bool
}

// New constructs the scheduler.
func New() (*Service, error) {
	engine, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Service{
		engine: engine,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Register adds a job. Jobs registered before Start begin running at
// Start; jobs registered after Start (but before Stop) are picked up
// immediately. Names must be unique and intervals positive.
func (s *Service) Register(job Job) error {
	s.stoppedMutex.Lock()
	stopped := s.stopped
	s.stoppedMutex.Unlock()
	if stopped {
		return errors.New("scheduler is stopped and cannot register " + job.Name)
	}
	if job.Name == "" {
		return errors.New("scheduler job needs a name")
	}
	if job.Interval <= 0 {
		return errors.New("scheduler job " + job.Name + " needs a positive interval")
	}
	if job.Run == nil {
		return errors.New("scheduler job " + job.Name + " needs a Run function")
	}
	for _, existing := range s.engine.Jobs() {
		if existing.Name() == job.Name {
			return errors.New("scheduler job " + job.Name + " is already registered")
		}
	}

	options := []gocron.JobOption{
		gocron.WithName(job.Name),
		// A run that outlives its interval causes following ticks to be
		// skipped rather than queued or run concurrently.
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	}
	if job.RunAtStart {
		options = append(options, gocron.WithStartAt(gocron.WithStartImmediately()))
	}

	run := job.Run
	name := job.Name
	_, err := s.engine.NewJob(
		gocron.DurationJob(job.Interval),
		gocron.NewTask(func() {
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("Scheduler job %s panicked: %v", name, r)
				}
			}()
			if s.ctx.Err() != nil {
				return
			}
			run(s.ctx)
		}),
		options...,
	)
	return err
}

// Start begins running every registered job.
func (s *Service) Start() {
	s.engine.Start()
	log.Infof("Started scheduler with %d job(s)", len(s.engine.Jobs()))
}

// Stop cancels the job context and shuts the engine down, waiting for
// in-flight runs to finish. A stopped Service is permanently done: it
// cannot be restarted, and Register calls after Stop return an error.
func (s *Service) Stop() {
	s.stoppedMutex.Lock()
	s.stopped = true
	s.stoppedMutex.Unlock()

	s.cancel()
	if err := s.engine.Shutdown(); err != nil {
		log.Errorf("error shutting down scheduler: %v", err)
		return
	}
	log.Infoln("Stopped scheduler")
}

// Statuses lists every registered job.
func (s *Service) Statuses() []JobStatus {
	jobs := s.engine.Jobs()
	statuses := make([]JobStatus, 0, len(jobs))
	for _, job := range jobs {
		status := JobStatus{Name: job.Name()}
		if lastRun, err := job.LastRunStartedAt(); err == nil {
			status.LastRunAt = lastRun
		}
		if nextRun, err := job.NextRun(); err == nil {
			status.NextRunAt = nextRun
		}
		statuses = append(statuses, status)
	}
	return statuses
}
