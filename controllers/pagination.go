package controllers

// PaginatedResponse is a structure for returning a total count with results.
type PaginatedResponse struct {
	Results interface{} `json:"results"`
	Total   int         `json:"total"`
}
