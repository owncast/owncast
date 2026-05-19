package inbox

import (
	"context"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/pkg/errors"
)

func (s *Service) handleCreateRequest(c context.Context, activity vocab.ActivityStreamsCreate) error {
	iri, err := apmodels.GetIRIStringFromJSONLDIdProperty(activity.GetJSONLDId())
	if err != nil {
		return errors.Wrap(err, "create activity is missing IRI")
	}
	return errors.New("not handling create request of: " + iri)
}
