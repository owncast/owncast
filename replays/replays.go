package replays

import (
	"context"

	"github.com/owncast/owncast/core/data"
	log "github.com/sirupsen/logrus"
)

// Setup will setup the replay package.
func Setup() {
	fixUnfinishedStreams()
}

// fixUnfinishedStreams will find streams with no end time and attempt to
// give them an end time based on the last segment assigned to that stream.
func fixUnfinishedStreams() {
	if err := data.GetDatastore().GetQueries().FixUnfinishedStreams(context.Background()); err != nil {
		log.Warnln(err)
	}
}
