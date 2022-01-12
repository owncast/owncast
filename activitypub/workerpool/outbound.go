package workerpool

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	// ActivityPubWorkerPoolSize defines the number of concurrent HTTP ActivityPub requests.
	ActivityPubWorkerPoolSize = 10
)

// Job struct bundling the ActivityPub and the payload in one struct.
type Job struct {
	request *http.Request
}

var queue chan Job

// InitOutboundWorkerPool starts n go routines that await ActivityPub jobs.
func InitOutboundWorkerPool() {
	queue = make(chan Job)

	// start workers
	for i := 1; i <= ActivityPubWorkerPoolSize; i++ {
		go worker(i, queue)
	}
}

// AddToOutboundQueue will queue up an outbound http request.
func AddToOutboundQueue(req *http.Request) {
	log.Tracef("Queued request for ActivityPub destination %s", req.RequestURI)
	queue <- Job{req}
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
	// req, err := http.NewRequest("POST", job.inbox.String(), bytes.NewReader(job.payload))
	// if err != nil {
	// 	return err
	// }

	// req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(job.request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
