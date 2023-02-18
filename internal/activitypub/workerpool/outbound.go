package workerpool

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	// ActivityPubWorkerPoolSize defines the number of concurrent HTTP ActivityPub requests.
	ActivityPubWorkerPoolSize = 10
)

// New will initialize the ActivityPub persistence layer service
// with the provided Data.
func New() (*Service, error) {
	s := &Service{}

	return s, nil
}

type Service struct {
	queue chan Job
}

// Job struct bundling the ActivityPub and the payload in one struct.
type Job struct {
	request *http.Request
}

// InitOutboundWorkerPool starts n go routines that await ActivityPub jobs.
func (s *Service) InitOutboundWorkerPool() {
	s.queue = make(chan Job)

	// start workers
	for i := 1; i <= ActivityPubWorkerPoolSize; i++ {
		go s.worker(i, s.queue)
	}
}

// AddToOutboundQueue will queue up an outbound http request.
func (s *Service) AddToOutboundQueue(req *http.Request) {
	log.Tracef("Queued request for ActivityPub destination %s", req.RequestURI)
	s.queue <- Job{req}
}

func (s *Service) worker(workerID int, queue <-chan Job) {
	log.Debugf("Started ActivityPub worker %d", workerID)

	for job := range queue {
		if err := s.sendActivityPubMessageToInbox(job); err != nil {
			log.Errorf("ActivityPub destination %s failed to send Error: %s", job.request.RequestURI, err)
		}
		log.Tracef("Done with ActivityPub destination %s using worker %d", job.request.RequestURI, workerID)
	}
}

func (s *Service) sendActivityPubMessageToInbox(job Job) error {
	client := &http.Client{}

	resp, err := client.Do(job.request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
