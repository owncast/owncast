package controllers

import (
	"net/http"
)

// GetFollowers will handle an API request to fetch the list of followers (non-activitypub response).
func (s *Service) GetFollowers(offset int, limit int, w http.ResponseWriter, r *http.Request) {
	followers, total, err := s.Follower.GetFederationFollowers(limit, offset)
	if err != nil {
		s.WriteSimpleResponse(w, false, "unable to fetch followers")
		return
	}

	response := PaginatedResponse{
		Total:   total,
		Results: followers,
	}

	s.WriteResponse(w, response)
}
