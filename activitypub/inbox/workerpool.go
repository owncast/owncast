package inbox

import (
	"github.com/owncast/owncast/activitypub/apmodels"
	log "github.com/sirupsen/logrus"
)

const (
	// InboxWorkerPoolSize defines the number of concurrent ActivityPub handlers.
	InboxWorkerPoolSize = 10
)

// Job struct bundling the ActivityPub and the payload in one struct.
type Job struct {
	request apmodels.InboxRequest
}

var queue chan Job

// InitInboxWorkerPool starts n go routines that await ActivityPub jobs.
func InitInboxWorkerPool() {
	queue = make(chan Job)

	// start workers
	for i := 1; i <= InboxWorkerPoolSize; i++ {
		go worker(i, queue)
	}
}

// AddToQueue will queue up an outbound http request.
func AddToQueue(req apmodels.InboxRequest) {
	log.Tracef("Queued request for ActivityPub inbox handler")
	queue <- Job{req}
}

func worker(workerID int, queue <-chan Job) {
	log.Debugf("Started ActivityPub worker %d", workerID)

	for job := range queue {
		handle(job.request)

		log.Tracef("Done with ActivityPub inbox handler using worker %d", workerID)
	}
}
