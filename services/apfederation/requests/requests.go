package requests

import "github.com/owncast/owncast/services/apfederation/workerpool"

type Requests struct {
	outboundWorkerPool *workerpool.WorkerPool
}

func New() *Requests {
	return &Requests{}
}

var temporaryGlobalInstance *Requests

func Get() *Requests {
	if temporaryGlobalInstance == nil {
		temporaryGlobalInstance = New()
	}
	return temporaryGlobalInstance
}
