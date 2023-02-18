package middleware

import (
	"net/http"
	"strconv"
)

// PaginatedHandlerFunc is a handler for endpoints that require pagination.
type PaginatedHandlerFunc func(int, int, http.ResponseWriter, *http.Request)

// HandlePagination is a middleware handler that pulls pagination values
// and passes them along.
func HandlePagination(handler PaginatedHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Default 50 items per page
		limitString := r.URL.Query().Get("limit")
		if limitString == "" {
			limitString = "50"
		}
		limit, err := strconv.Atoi(limitString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Default first page 0
		offsetString := r.URL.Query().Get("offset")
		if offsetString == "" {
			offsetString = "0"
		}
		offset, err := strconv.Atoi(offsetString)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		handler(offset, limit, w, r)
	}
}
