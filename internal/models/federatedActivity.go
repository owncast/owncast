package models

import "time"

// FederatedActivity is an internal representation of an activity that was
// accepted and stored.
type FederatedActivity struct {
	IRI       string    `json:"iri"`
	ActorIRI  string    `json:"actorIRI"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}
