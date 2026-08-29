# Federation

This document describes how Owncast federates with other ActivityPub servers, following [FEP-67ff](https://codeberg.org/fediverse/fep/src/branch/main/fep/67ff/fep-67ff.md). It is intended to help operators and developers of other Fediverse software interoperate with Owncast.

Owncast is a self-hosted live streaming server. Each instance federates as a single actor representing its live stream, broadcasting go-live announcements and posts to its followers. It does not present a social timeline and does not follow individual fediverse accounts. Servers can, however, list one another's live streams through the [featured-streams](#owncast-live-stream-status) protocol.

## Supported federation protocols and standards

- [ActivityPub](https://www.w3.org/TR/activitypub/) (Server-to-Server)
- [WebFinger](https://webfinger.net/)
- [HTTP Signatures](https://datatracker.ietf.org/doc/html/draft-cavage-http-signatures)
- [NodeInfo](https://nodeinfo.diaspora.software/)

## Supported FEPs

- [FEP-67ff: FEDERATION.md](https://codeberg.org/fediverse/fep/src/branch/main/fep/67ff/fep-67ff.md)
- [FEP-f1d5: NodeInfo in Fediverse Software](https://codeberg.org/fediverse/fep/src/branch/main/fep/f1d5/fep-f1d5.md)
- [FEP-044f: Consent-respecting quote posts](https://codeberg.org/fediverse/fep/src/branch/main/fep/044f/fep-044f.md)

## ActivityPub in Owncast

### Actor

Each Owncast instance presents as a single `Service` actor (not `Person`) representing the live stream. The actor profile includes:

- `icon`: Server logo/avatar
- `image`: Profile banner image (currently the same asset as `icon`)
- `summary`: Server description, emitted verbatim as configured
- `attachment`: Array of `PropertyValue` objects containing social links and metadata
- `tag`: Hashtags describing the server content
- `manuallyApprovesFollowers`: `true` when the server operates in private mode
- `discoverable`: Always `true`

Owncast does not expose a `following` collection; that endpoint returns `404`. It does not follow individual fediverse accounts, but it does follow other servers as part of the [featured-streams](#owncast-live-stream-status) directory protocol described below.

### Activities

| Activity       | Object         | Send | Receive | Notes                                                                                                                                                                       |
| -------------- | -------------- | :--: | :-----: | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Create`       | `Note`         | Yes  |   Yes   | Sent on go-live and manual fediverse posts. Qualifying inbound notes produce plugin mention or reply events only.                                                           |
| `Update`       | `Service`      | Yes  |   No    | Sent when the server profile changes.                                                                                                                                       |
| `Update`       | `Person`       |  No  |   Yes   | Updates cached follower profile information.                                                                                                                                |
| `Follow`       | -              | Yes  |   Yes   | Received from followers (queued for approval in private mode; directory follows always require approval). Sent to other servers to feature their stream (featured-streams). |
| `Accept`       | `Follow`       | Yes  |   Yes   | Sent to accept an inbound follow. Received when a server accepts our featured-streams follow.                                                                               |
| `Reject`       | `Follow`       | Yes  |   Yes   | Sent to reject or remove a follower. Received when a server rejects our featured-streams follow.                                                                            |
| `Undo`         | `Follow`       | Yes  |   Yes   | Received to remove a follower. Sent to stop featuring a server.                                                                                                             |
| `Like`         | `Note`         |  No  |   Yes   | Optionally displayed in live chat.                                                                                                                                          |
| `Announce`     | `Note`         |  No  |   Yes   | Optionally displayed in live chat.                                                                                                                                          |
| `Offer`        | -              | Yes  |   Yes   | Featured-streams live signal (see below).                                                                                                                                   |
| `Leave`        | -              | Yes  |   Yes   | Featured-streams offline signal (see below).                                                                                                                                |
| `QuoteRequest` | `Note`         |  No  |   Yes   | A remote user asking to quote one of our posts (FEP-044f).                                                                                                                  |
| `Accept`       | `QuoteRequest` | Yes  |   No    | Sent when the quoted post exists and quoting is allowed. Carries the `QuoteAuthorization` stamp IRI in `result`.                                                            |
| `Reject`       | `QuoteRequest` | Yes  |   No    | Sent for unknown posts, while in private mode, or when quoting is disabled.                                                                                                 |

Owncast does not display inbound `Create(Note)` activities as a social timeline and does not federate them onward. Qualifying mentions and replies are delivered only to server plugins.

### Owncast live stream status

Owncast servers can list one another's live streams in a lightweight directory built on standard ActivityPub activities plus a small JSON-LD extension in the `https://owncast.online/ns#` namespace. A directory follows the servers it features; each featured server then signals its live and offline transitions to its directory followers.

Using this mechanism, a standalone, custom, directory can display a list of currently live streams, their status, and their metadata. The protocol is open, so third-party ActivityPub implementations can act as directories that list Owncast streams.

**Following (listing a stream).** A directory sends a `Follow` whose `actor` is the directory and whose `object` is the target server's actor, carrying the extension property `https://owncast.online/ns#directory: true`. This marker distinguishes a directory listing from an ordinary personal Fediverse follow. A directory follow always requires operator approval, even on otherwise public instances. The target replies with `Accept(Follow)` once approved, or `Reject(Follow)`; it may later send `Reject(Follow)` to drop the listing, and the directory may send `Undo(Follow)` to stop featuring.

**Live and offline signals.** Once the follow is accepted, the Owncast server delivers stream-status activities to its directory followers only (never personal account followers): `Offer` when the stream goes live (at go-live and periodically on a timer) and `Leave` when it goes offline. Both use the sending server's base URL (`scheme://host`) as the activity `object` and attach stream metadata as JSON-LD extension properties:

| Property (in `https://owncast.online/ns#`) | Meaning                                         |
| ------------------------------------------ | ----------------------------------------------- |
| `streamStatus`                             | `live` or `offline`                             |
| `streamTitle`                              | Current stream title                            |
| `streamDescription`                        | Stream or server description                    |
| `serverName`                               | Server display name                             |
| `logoUrl`                                  | Server logo URL                                 |
| `thumbnailUrl`                             | Live preview thumbnail URL (present while live) |
| `streamTags`                               | Array of stream tags                            |

Received metadata is length- and count-clamped before storage.

**Interoperating with other software.** This protocol is open, but the two directions differ today. Owncast accepts a directory `Follow` (carrying the `ns#directory` marker) from any actor and will deliver `Offer`/`Leave` to it, so a third-party ActivityPub implementation can act as a directory that features Owncast streams. The reverse is currently Owncast-only: before following a server in order to feature it, Owncast validates that the server's NodeInfo reports `software.name == "owncast"` and `metadata.federation.featuredStreams >= 1`, so it will not feature a non-Owncast server.

### Notes

Posts from Owncast are `Note` objects with:

- `content`: HTML-formatted text with linked hashtags
- `attachment`: `Image` object with the stream preview, included on go-live posts when one is available
- `tag`: Array containing `Hashtag` and `Mention` objects
- `sensitive`: `true` on go-live posts when the stream is marked NSFW
- `interactionPolicy`: Declares that anyone may quote the post (`canQuote` with `automaticApproval` set to the public collection), attached to public posts unless quoting is disabled or the server is in private mode

When a stream starts, Owncast sends a `Create(Note)` go-live announcement containing the configured go-live message (defaulting to `I've gone live!`), the server's configured hashtags, a link to watch the stream, and the current stream preview image when one is available. If the operator clears the go-live message, no announcement is sent.

Hashtags use the `Hashtag` type (Mastodon/`toot` vocabulary) and link to `https://owncast.directory/tags/{tag}` for discovery across Owncast instances.

### Quote posts

Owncast implements the quoted-server side of [FEP-044f](https://codeberg.org/fediverse/fep/src/branch/main/fep/044f/fep-044f.md): remote users can quote Owncast posts with a consent flow.

- Public posts carry an `interactionPolicy` advertising that quoting is automatically approved for anyone.
- An inbound `QuoteRequest` for a post this server authored is answered with an `Accept` whose `result` is the IRI of a stored `QuoteAuthorization` stamp. The stamp is dereferenceable at that IRI so third-party servers can verify the quote was approved.
- A `QuoteRequest` for an unknown object, any request while private mode is enabled, or any request while quoting is disabled receives a `Reject`.
- Quoting is enabled by default and can be turned off by the operator.
- Owncast does not author quote posts and does not revoke previously issued stamps.

### Addressing

Public posts are addressed to:

```json
{
	"to": ["https://www.w3.org/ns/activitystreams#Public"],
	"cc": ["{actor}/followers"]
}
```

Owncast can send direct messages to specific actors (used for replying to engagement). These address the recipient in `to` with a corresponding `Mention` tag for Mastodon compatibility.

### HTTP Signatures

- **Outbound**: All `POST` requests are signed using RSA-SHA256. Signed headers are `(request-target)`, `host`, `date`, and `digest`. A fresh signature is generated for every delivery attempt.
- **Inbound**: Activities are accepted over HTTP (`202 Accepted`) and verified asynchronously. Only activities carrying a valid HTTP signature are processed; unsigned or invalidly signed activities are discarded without being acted on.
- **GET requests** to the actor, outbox, and followers endpoints do not require signatures.

### Private mode

Owncast can operate in private mode, where:

- `manuallyApprovesFollowers` is `true` on the actor
- Follow requests are queued for operator approval
- `Accept` is only sent after manual approval

### Blocking

Owncast supports blocking at both the domain and individual actor level. Inbound activities from a blocked domain or a disabled actor are discarded without being processed, and a rejected or blocked follower is removed from future outbound delivery.

### Interoperability notes

- Owncast uses the `Service` actor type, not `Person`.
- Owncast does not follow individual fediverse accounts and exposes no `following` collection, but servers follow one another to build the [featured-streams](#owncast-live-stream-status) directory.
- Inbound `Create(Note)` activities are not shown as timeline posts or federated onward. Qualifying mentions and replies are delivered only to plugins.
- Engagement (`Like`, `Announce`) can be displayed in the live chat if enabled.
- The `sensitive` flag indicates NSFW content for the entire stream.

## WebFinger

Actor discovery is available via `/.well-known/webfinger?resource=acct:{username}@{domain}`.

Response links include:

- `self` (`application/activity+json`): Actor document
- `http://webfinger.net/rel/profile-page`: Web profile
- `http://webfinger.net/rel/avatar`: Profile image
- `alternate` (`application/x-mpegURL`): HLS stream URL, allowing clients to discover the live stream directly from WebFinger

## NodeInfo

NodeInfo discovery is available at `/.well-known/nodeinfo`, with NodeInfo 2.0 served at `/nodeinfo/2.0`. Additional related endpoints:

- `/.well-known/x-nodeinfo2`: Extended format
- `/api/v1/instance`: Mastodon-compatible instance information
- `/.well-known/host-meta`: WebFinger discovery

## Additional documentation

- [Owncast Documentation](https://owncast.online/docs/)
- [ActivityPub Specification](https://www.w3.org/TR/activitypub/)
- [Fediverse Enhancement Proposals](https://codeberg.org/fediverse/fep)
