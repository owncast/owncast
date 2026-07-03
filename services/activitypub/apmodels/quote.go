package apmodels

import (
	"net/url"

	"code.superseriousbusiness.org/activity/streams"
	"code.superseriousbusiness.org/activity/streams/vocab"
)

// MakeNoteQuotable marks a note as quotable by anyone by attaching an
// interaction policy that automatically approves quote requests (FEP-044f).
// Remote servers read this policy to decide whether to offer their users a
// quote option for the post.
func MakeNoteQuotable(note vocab.ActivityStreamsNote) vocab.ActivityStreamsNote {
	public, _ := url.Parse(PUBLIC)

	automaticApproval := streams.NewGoToSocialAutomaticApprovalProperty()
	automaticApproval.AppendIRI(public)

	canQuote := streams.NewGoToSocialCanQuote()
	canQuote.SetGoToSocialAutomaticApproval(automaticApproval)

	canQuoteProp := streams.NewGoToSocialCanQuoteProperty()
	canQuoteProp.AppendGoToSocialCanQuote(canQuote)

	policy := streams.NewGoToSocialInteractionPolicy()
	policy.SetGoToSocialCanQuote(canQuoteProp)

	policyProp := streams.NewGoToSocialInteractionPolicyProperty()
	policyProp.AppendGoToSocialInteractionPolicy(policy)
	note.SetGoToSocialInteractionPolicy(policyProp)

	return note
}

// MakeQuoteAuthorization creates a QuoteAuthorization "stamp" (FEP-044f)
// approving quotePostIRI (the remote quote post) quoting quotedPostIRI (one
// of our posts). The stamp must be served at its IRI so third-party servers
// can fetch it to verify the quote was authorized.
func MakeQuoteAuthorization(id, attributedToIRI, quotePostIRI, quotedPostIRI *url.URL) vocab.GoToSocialQuoteAuthorization {
	authorization := streams.NewGoToSocialQuoteAuthorization()

	idProp := streams.NewJSONLDIdProperty()
	idProp.Set(id)
	authorization.SetJSONLDId(idProp)

	attributedTo := streams.NewActivityStreamsAttributedToProperty()
	attributedTo.AppendIRI(attributedToIRI)
	authorization.SetActivityStreamsAttributedTo(attributedTo)

	interactingObject := streams.NewGoToSocialInteractingObjectProperty()
	interactingObject.AppendIRI(quotePostIRI)
	authorization.SetGoToSocialInteractingObject(interactingObject)

	interactionTarget := streams.NewGoToSocialInteractionTargetProperty()
	interactionTarget.AppendIRI(quotedPostIRI)
	authorization.SetGoToSocialInteractionTarget(interactionTarget)

	return authorization
}
