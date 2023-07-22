package inbox

import (
	"runtime"

	"github.com/owncast/owncast/activitypub/apmodels"
	"github.com/owncast/owncast/services/apfederation/apmodels"
	log "github.com/sirupsen/logrus"
)

// workerPoolSize defines the number of concurrent ActivityPub handlers.
var workerPoolSize = runtime.GOMAXPROCS(0)

// Job struct bundling the ActivityPub and the payload in one struct.
type Job struct {
	request apmodels.InboxRequest
}

var queue chan Job

// InitInboxWorkerPool starts n go routines that await ActivityPub jobs.
func (api *APInbox) InitInboxWorkerPool() {
	queue = make(chan Job)

	// start workers
	for i := 1; i <= workerPoolSize; i++ {
		go worker(i, queue)
	}
}

// AddToQueue will queue up an outbound http request.
func (api *APInbox) AddToQueue(req apmodels.InboxRequest) {
	log.Tracef("Queued request for ActivityPub inbox handler")
	queue <- Job{req}
}

func (api *APInbox) worker(workerID int, queue <-chan Job) {
	log.Debugf("Started ActivityPub worker %d", workerID)

	for job := range queue {
		api.handle(job.request)

		log.Tracef("Done with ActivityPub inbox handler using worker %d", workerID)
	}
}
