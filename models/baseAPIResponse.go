package models

// BaseAPIResponse is a simple response to API requests.
type BaseAPIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
