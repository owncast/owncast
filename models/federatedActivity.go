package models

import "time"

// FederatedActivity is an internal representation of an activity that was
// accepted and stored.
type FederatedActivity struct {
	Timestamp time.Time `json:"timestamp"`
	IRI       string    `json:"iri"`
	ActorIRI  string    `json:"actorIRI"`
	Type      string    `json:"type"`
}
