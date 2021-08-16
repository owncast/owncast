package models

import "net/http"

type InboxRequest struct {
	Request         *http.Request
	ForLocalAccount string
	Body            []byte
}
