package requests

import (
	"encoding/json"
	"net/url"

	"code.superseriousbusiness.org/activity/streams"
	"code.superseriousbusiness.org/activity/streams/vocab"
	"github.com/pkg/errors"

	"github.com/owncast/owncast/services/activitypub/apmodels"
	"github.com/owncast/owncast/services/activitypub/crypto"
	"github.com/owncast/owncast/services/activitypub/workerpool"

	"github.com/teris-io/shortid"
)

// SendQuoteRequestAccept accepts a remote QuoteRequest (FEP-044f) by
// delivering an Accept whose result carries the IRI of the stored
// QuoteAuthorization stamp. The requesting server flips its pending quote to
// approved and passes the stamp IRI along so other servers can verify it.
func SendQuoteRequestAccept(wp *workerpool.Service, inbox *url.URL, originalQuoteRequest vocab.GoToSocialQuoteRequest, stampIRI *url.URL, fromLocalAccountName string, builder *apmodels.Builder, signer *crypto.Signer) error {
	accept := streams.NewActivityStreamsAccept()
	setQuoteResponseProperties(accept, originalQuoteRequest, fromLocalAccountName, builder)

	result := streams.NewActivityStreamsResultProperty()
	result.AppendIRI(stampIRI)
	accept.SetActivityStreamsResult(result)

	return signAndQueueQuoteResponse(wp, inbox, accept, fromLocalAccountName, builder, signer)
}

// SendQuoteRequestReject declines a remote QuoteRequest, clearing the pending
// quote on the requesting server.
func SendQuoteRequestReject(wp *workerpool.Service, inbox *url.URL, originalQuoteRequest vocab.GoToSocialQuoteRequest, fromLocalAccountName string, builder *apmodels.Builder, signer *crypto.Signer) error {
	reject := streams.NewActivityStreamsReject()
	setQuoteResponseProperties(reject, originalQuoteRequest, fromLocalAccountName, builder)

	return signAndQueueQuoteResponse(wp, inbox, reject, fromLocalAccountName, builder, signer)
}

// quoteResponse is the intersection of Accept and Reject used when responding
// to a QuoteRequest.
type quoteResponse interface {
	SetJSONLDId(vocab.JSONLDIdProperty)
	SetActivityStreamsActor(vocab.ActivityStreamsActorProperty)
	SetActivityStreamsObject(vocab.ActivityStreamsObjectProperty)
}

// setQuoteResponseProperties fills the id, actor and (embedded original
// QuoteRequest) object properties shared by both response activities.
func setQuoteResponseProperties(response quoteResponse, originalQuoteRequest vocab.GoToSocialQuoteRequest, fromAccountName string, builder *apmodels.Builder) {
	responseID := builder.MakeLocalIRIForResource(shortid.MustGenerate())
	idProperty := streams.NewJSONLDIdProperty()
	idProperty.SetIRI(responseID)
	response.SetJSONLDId(idProperty)

	actorID := builder.MakeLocalIRIForAccount(fromAccountName)
	response.SetActivityStreamsActor(apmodels.MakeActorPropertyWithID(actorID))

	// Embed the original request so the receiving server can match it by id
	// even if it does not track our activity IRIs.
	object := streams.NewActivityStreamsObjectProperty()
	object.AppendGoToSocialQuoteRequest(originalQuoteRequest)
	response.SetActivityStreamsObject(object)
}

func signAndQueueQuoteResponse(wp *workerpool.Service, inbox *url.URL, response vocab.Type, fromLocalAccountName string, builder *apmodels.Builder, signer *crypto.Signer) error {
	// SSRF protection: reject non-HTTPS schemes and internal/loopback hosts.
	if err := validateRemoteInbox(inbox); err != nil {
		return err
	}

	localAccountIRI := builder.MakeLocalIRIForAccount(fromLocalAccountName)

	jsonmap, err := streams.Serialize(response)
	if err != nil {
		return errors.Wrap(err, "unable to serialize quote request response")
	}
	b, err := json.Marshal(jsonmap)
	if err != nil {
		return errors.Wrap(err, "unable to marshal quote request response")
	}

	req, err := signer.CreateSignedRequest(b, inbox, localAccountIRI)
	if err != nil {
		return err
	}

	wp.AddToOutboundQueue(req)

	return nil
}
