# Federation

This document describes how Owncast federates with other ActivityPub servers, following [FEP-67ff](https://codeberg.org/fediverse/fep/src/branch/main/fep/67ff/fep-67ff.md).

## Supported Federation Protocols and Standards

- [ActivityPub](https://www.w3.org/TR/activitypub/) (Server-to-Server)
- [WebFinger](https://webfinger.net/)
- [HTTP Signatures](https://datatracker.ietf.org/doc/html/draft-cavage-http-signatures)
- [NodeInfo](https://nodeinfo.diaspora.software/)

## Supported FEPs

- [FEP-67ff: FEDERATION.md](https://codeberg.org/fediverse/fep/src/branch/main/fep/67ff/fep-67ff.md)
- [FEP-f1d5: NodeInfo in Fediverse Software](https://codeberg.org/fediverse/fep/src/branch/main/fep/f1d5/fep-f1d5.md)
- [FEP-044f: Consent-respecting quote posts](https://codeberg.org/fediverse/fep/src/branch/main/fep/044f/fep-044f.md)

## Actor

Owncast servers present as a single `Service` actor (not `Person`). Each Owncast instance has one actor representing the live stream.

The actor profile includes:
- `icon`: Server logo/avatar
- `image`: Profile banner image
- `summary`: Server description (HTML)
- `attachment`: Array of `PropertyValue` objects containing social links and metadata
- `tag`: Hashtags describing the server content
- `manuallyApprovesFollowers`: `true` if the server operates in private mode
- `discoverable`: Always `true`

Owncast actors do not follow other actors. The `following` collection endpoint returns 404.

## Activities

| Activity | Object | Send | Receive | Notes |
|----------|--------|:----:|:-------:|-------|
| `Create` | `Note` | Yes | Yes | Sent on go-live and manual fediverse posts. Qualifying inbound notes produce plugin mention or reply events only. |
| `Update` | `Service` | Yes | No | Sent when server profile changes. |
| `Update` | `Person` | No | Yes | Updates cached follower profile information. |
| `Follow` | - | No | Yes | Queued for approval if private mode is enabled. |
| `Accept` | `Follow` | Yes | No | Sent automatically unless in private mode. |
| `Undo` | `Follow` | No | Yes | Removes the follower. |
| `Like` | `Note` | No | Yes | Optionally displayed in live chat. |
| `Announce` | `Note` | No | Yes | Optionally displayed in live chat. |
| `QuoteRequest` | `Note` | No | Yes | A remote user asking to quote one of our posts (FEP-044f). |
| `Accept` | `QuoteRequest` | Yes | No | Sent when the quoted post exists and quoting is allowed. Carries the `QuoteAuthorization` stamp IRI in `result`. |
| `Reject` | `QuoteRequest` | Yes | No | Sent for unknown posts, while in private mode, or when quoting is disabled. |

Owncast is a broadcast-only service. It does not follow other actors. Inbound posts are not displayed as a social timeline or federated onward.

## Notes

Posts from Owncast are `Note` objects with:

- `content`: HTML-formatted text with linked hashtags
- `attachment`: `Image` object with stream thumbnail (for go-live posts)
- `tag`: Array containing `Hashtag` and `Mention` objects
- `sensitive`: `true` if the stream is marked NSFW
- `interactionPolicy`: Declares that anyone may quote the post (`canQuote` with `automaticApproval` set to the public collection), attached to public posts unless the operator disables quoting

### Go-Live Announcements

When a stream starts, Owncast sends a `Create(Note)` containing:
- The configured go-live message (or default announcement)
- A thumbnail image of the stream
- Hashtags configured for the server
- A link to watch the stream

### Hashtags

Hashtags use `toot:Hashtag` type and link to `https://directory.owncast.online/tags/{tag}` for discovery across Owncast instances.

## Quote Posts

Owncast implements the quoted-server side of [FEP-044f](https://codeberg.org/fediverse/fep/src/branch/main/fep/044f/fep-044f.md): remote users can quote Owncast posts with the consent flow used by Mastodon 4.5+ and GoToSocial.

- Public posts carry an `interactionPolicy` advertising that quoting is automatically approved for anyone.
- An inbound `QuoteRequest` for a post this server authored is answered with an `Accept` whose `result` is the IRI of a stored `QuoteAuthorization` stamp. The stamp is dereferenceable at that IRI so third-party servers can verify the quote was approved.
- A `QuoteRequest` for an unknown object, any request while private mode is enabled, or any request while the operator has disabled quoting receives a `Reject`.
- Quoting can be turned off by the operator in the admin under the Social configuration ("Allow quotes"). It is enabled by default.
- Owncast does not author quote posts and does not revoke previously issued stamps.

## Plugin Integration

Inbound Fediverse activities can produce internal plugin events. These subscriptions are not external HTTP webhooks. A plugin must declare the `fediverse.inbound` permission to subscribe to `fediverse.follow`, `fediverse.like`, `fediverse.repost`, `fediverse.mention`, `fediverse.reply`, `fediverse.quote`, or `fediverse.activity`.

All inbox requests must pass HTTP signature verification and the activity actor must share an origin with the verified signing key owner. After those checks, `fediverse.activity` receives the raw JSON object in addition to any specialized event the activity produces.

Inbound `Create(Note)` activities have narrower requirements for mention and reply events. The `Create` must contain exactly one actor and exactly one `Note`. The note must have an IRI and exactly one `attributedTo` value that matches the activity actor. The resolved actor must also match the activity actor.

- A note whose first `inReplyTo` IRI identifies a locally authored post produces `fediverse.reply`.
- Any other qualifying note produces `fediverse.mention` only when the `Create` or `Note` addresses the local actor in `to` or `cc`.

These inbound notes are delivered to plugins. Owncast does not add them to a social timeline or federate them onward.

An accepted `QuoteRequest` produces `fediverse.quote` after Owncast signs and queues the acceptance. Its payload contains the remote actor and the target, which is the locally authored quoted post. Rejected quote requests do not produce this event.

## Addressing

Public posts are addressed to:
```json
{
  "to": ["https://www.w3.org/ns/activitystreams#Public"],
  "cc": ["{actor}/followers"]
}
```

Owncast supports sending direct messages to specific actors (used for replying to engagement). These include the recipient in `cc` with a corresponding `Mention` tag for Mastodon compatibility.

## HTTP Signatures

**Outbound**: All POST requests are signed using RSA-SHA256. Signed headers: `(request-target)`, `host`, `date`, `digest`.

**Inbound**: All POST requests to the inbox must have a valid signature. Unsigned requests are rejected.

**GET requests** to actor, outbox, and followers endpoints do not require signatures.

## WebFinger

Actor discovery via `/.well-known/webfinger?resource=acct:{username}@{domain}`

Response links include:
- `self` (application/activity+json): Actor document
- `http://webfinger.net/rel/profile-page`: Web profile
- `http://webfinger.net/rel/avatar`: Profile image
- `alternate` (application/x-mpegURL): HLS stream URL

The HLS stream link allows clients to discover the live stream URL directly from WebFinger.

## NodeInfo

Available at `/.well-known/nodeinfo` with NodeInfo 2.0 at `/nodeinfo/2.0`.

Additional endpoints:
- `/.well-known/x-nodeinfo2`: Extended format
- `/api/v1/instance`: Mastodon-compatible instance information
- `/.well-known/host-meta`: WebFinger discovery

## Private Mode

Owncast can operate in private mode where:
- `manuallyApprovesFollowers` is `true` on the actor
- Follow requests are queued for admin approval
- `Accept` is only sent after manual approval

## Blocking

Owncast supports blocking at the domain and individual actor level. Blocked entities:
- Cannot deliver to the inbox (requests rejected)
- Do not receive activities from the server

## Interoperability Notes

- Owncast uses `Service` actor type, not `Person`
- There is no `following` collection (Owncast does not follow accounts)
- Owncast does not display inbound `Create(Note)` activities as social timeline posts or federate them onward. Qualifying mentions and replies are delivered only to plugins
- Engagement (`Like`, `Announce`) can be displayed in the live chat if enabled
- The `sensitive` flag indicates NSFW content for the entire stream

## Additional Resources

- [Owncast Documentation](https://owncast.online/docs/)
- [ActivityPub Specification](https://www.w3.org/TR/activitypub/)
- [Fediverse Enhancement Proposals](https://codeberg.org/fediverse/fep)
