package inbox

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/go-fed/activity/streams/vocab"
)

func handleUndoInboxRequest(c context.Context, activity vocab.ActivityStreamsUndo) error {
	// Determine if this is an undo of a follow, favorite, announce, etc.
	o := activity.GetActivityStreamsObject()
	for iter := o.Begin(); iter != o.End(); iter = iter.Next() {
		if iter.IsActivityStreamsFollow() {
			// This is an Unfollow request
			if err := handleUnfollowRequest(c, activity); err != nil {
				return err
			}
		} else {
			log.Traceln("Undo", iter.GetType().GetTypeName(), "ignored")
			return nil
		}
	}

	return nil
}
