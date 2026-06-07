package models

// BaseAPIResponse is a simple response to API requests.
type BaseAPIResponse struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode,omitempty"`
	Success   bool   `json:"success"`
}
