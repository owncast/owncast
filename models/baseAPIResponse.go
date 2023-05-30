package models

// BaseAPIResponse is a simple response to API requests.
type BaseAPIResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
