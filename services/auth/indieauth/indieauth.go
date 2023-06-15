package indieauth

import "sync"

type IndieAuthClient struct {
	pendingAuthRequests map[string]*Request
	lock                sync.Mutex
}

type IndieAuthServer struct {
	pendingServerAuthRequests map[string]ServerAuthRequest
}

var temporaryGlobalClientInstance *IndieAuthClient

// GetIndieAuthClient returns the temporary global instance of IndieAuthClient.
// Remove this after dependency injection is implemented.
func GetIndieAuthClient() *IndieAuthClient {
	if temporaryGlobalClientInstance == nil {
		temporaryGlobalClientInstance = NewIndieAuthClient()
	}

	return temporaryGlobalClientInstance
}

// NewIndieAuthClient creates a new IndieAuth client instance.
func NewIndieAuthClient() *IndieAuthClient {
	i := &IndieAuthClient{
		pendingAuthRequests: make(map[string]*Request),
	}
	go i.setupExpiredRequestPruner()

	return i
}

var temporaryGlobalServerInstance *IndieAuthServer

// GetIndieAuthServer returns the temporary global instance of IndieAuthServer.
// Remove this after dependency injection is implemented.
func GetIndieAuthServer() *IndieAuthServer {
	if temporaryGlobalServerInstance == nil {
		temporaryGlobalServerInstance = NewIndieAuthServer()
	}

	return temporaryGlobalServerInstance
}

// NewIndieAuthServer creates a new IndieAuth client instance.
func NewIndieAuthServer() *IndieAuthServer {
	i := &IndieAuthServer{
		pendingServerAuthRequests: make(map[string]ServerAuthRequest),
	}

	return i
}
