package test

import (
	"time"

	log "github.com/sirupsen/logrus"
)

var timestamp time.Time

func Mark() {
	now := time.Now()
	if !timestamp.IsZero() {
		delta := now.Sub(timestamp)
		log.Println(delta.Milliseconds(), "ms")
	}

	timestamp = now
}
