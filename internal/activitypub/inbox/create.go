package inbox

import (
	"context"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/pkg/errors"
)

func handleCreateRequest(c context.Context, activity vocab.ActivityStreamsCreate) error {
	iri := activity.GetJSONLDId().GetIRI().String()
	return errors.New("not handling create request of: " + iri)
}
