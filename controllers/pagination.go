package controllers

// PaginatedResponse is a structure for returning a total count with results.
type PaginatedResponse struct {
	Total   int         `json:"total"`
	Results interface{} `json:"results"`
}
