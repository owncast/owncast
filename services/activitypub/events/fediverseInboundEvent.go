package events

// FediverseActor is the remote account responsible for an inbound activity.
type FediverseActor struct {
	Name   string `json:"name"`
	Handle string `json:"handle"`
	URL    string `json:"url,omitempty"`
	Image  string `json:"image,omitempty"`
}

// FediverseTarget identifies the local object acted on by a remote account.
type FediverseTarget struct {
	URL string `json:"url"`
}

// FediverseEngagementEvent is the payload for inbound like, repost, and quote events.
type FediverseEngagementEvent struct {
	Actor  FediverseActor   `json:"actor"`
	Target *FediverseTarget `json:"target,omitempty"`
}

// FediverseAttachment is media attached to an inbound post.
type FediverseAttachment struct {
	URL       string `json:"url"`
	MediaType string `json:"mediaType"`
	Alt       string `json:"alt,omitempty"`
}

// FediverseInboundPostEvent is the payload for inbound mention and reply events.
type FediverseInboundPostEvent struct {
	Actor       FediverseActor        `json:"actor"`
	Content     string                `json:"content"`
	ContentText string                `json:"contentText"`
	URL         string                `json:"url"`
	PostedAt    string                `json:"postedAt"`
	InReplyTo   string                `json:"inReplyTo,omitempty"`
	Attachments []FediverseAttachment `json:"attachments,omitempty"`
	Language    string                `json:"language,omitempty"`
}
