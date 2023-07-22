package workerpool

import (
	"net/http"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// workerPoolSize defines the number of concurrent HTTP ActivityPub requests.
var workerPoolSize = runtime.GOMAXPROCS(0)

// Job struct bundling the ActivityPub and the payload in one struct.
type Job struct {
	request *http.Request
}

type WorkerPool struct {
	queue chan Job
}

func New() *WorkerPool {
	wp := &WorkerPool{
		queue: make(chan Job),
	}
	wp.initOutboundWorkerPool()
	return wp
}

var temporaryGlobalInstance *WorkerPool

func Get() *WorkerPool {
	if temporaryGlobalInstance == nil {
		temporaryGlobalInstance = New()
	}
	return temporaryGlobalInstance
}

// InitOutboundWorkerPool starts n go routines that await ActivityPub jobs.
func (wp *WorkerPool) initOutboundWorkerPool() {
	// start workers
	for i := 1; i <= workerPoolSize; i++ {
		go worker(i, queue)
	}
}

// AddToOutboundQueue will queue up an outbound http request.
func (wp *WorkerPool) AddToOutboundQueue(req *http.Request) {
	log.Tracef("Queued request for ActivityPub destination %s", req.RequestURI)
	wp.queue <- Job{req}
}

func (wp *WorkerPool) worker(workerID int, queue <-chan Job) {
	log.Debugf("Started ActivityPub worker %d", workerID)

	for job := range queue {
		if err := wp.sendActivityPubMessageToInbox(job); err != nil {
			log.Errorf("ActivityPub destination %s failed to send Error: %s", job.request.RequestURI, err)
		}
		log.Tracef("Done with ActivityPub destination %s using worker %d", job.request.RequestURI, workerID)
	}
}

func (wp *WorkerPool) sendActivityPubMessageToInbox(job Job) error {
	client := &http.Client{}

	resp, err := client.Do(job.request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
