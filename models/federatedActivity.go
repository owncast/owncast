package models

import "time"

type FederatedActivity struct {
	IRI       string    `json:"iri"`
	ActorIRI  string    `json:"actorIRI"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}
