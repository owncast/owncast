package apmodels

import "net/http"

// InboxRequest represents an inbound request to the ActivityPub inbox.
type InboxRequest struct {
	Request         *http.Request
	ForLocalAccount string
	Body            []byte
}
