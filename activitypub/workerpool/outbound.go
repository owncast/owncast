package workerpool

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Job struct bundling the ActivityPub and the payload in one struct.
type Job struct {
	request *http.Request
}

var queue chan Job

// InitOutboundWorkerPool starts n go routines that await ActivityPub jobs.
func InitOutboundWorkerPool(workerPoolSize int) {
	queue = make(chan Job, workerPoolSize)

	// start workers
	for i := 1; i <= workerPoolSize; i++ {
		go worker(i, queue)
	}
}

// AddToOutboundQueue will queue up an outbound http request.
func AddToOutboundQueue(req *http.Request) {
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
		}
		log.Tracef("Done with ActivityPub destination %s using worker %d", job.request.RequestURI, workerID)
	}
}

func sendActivityPubMessageToInbox(job Job) error {
	client := &http.Client{}

	resp, err := client.Do(job.request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
