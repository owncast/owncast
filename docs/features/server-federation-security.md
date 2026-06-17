# Featured Streams: Security & Trust Model

Featured streams let an Owncast server display a small directory of other
Owncast servers and their live status. Because ActivityPub is an open protocol,
any server can send anything and assert any identity. This document describes
what is and isn't guaranteed, so operators understand the trust boundary.

## What is guaranteed

- **The directory is admin-curated.** Entries in your featured directory are
  created only by an admin adding a server URL. No inbound ActivityPub activity
  can insert an entry — inbound `Offer`/`Accept`/`Leave` only *update* rows that
  already exist. A remote server cannot put itself on your list.
- **The link is the URL you entered, and it is immutable.** The link target for
  an entry is the `https` URL the admin typed (non-`https` is rejected). No
  inbound handler ever changes it, so a remote server cannot repoint your
  directory entry somewhere else after the fact.
- **Status and metadata are bound to that host.** Inbound activities are
  HTTP-signature verified, and the activity's `actor` must share the host of the
  verified signing key (a server cannot sign with its own key while claiming to
  be another). The featured-follow also pins the resolved actor to the host of
  the URL you entered. Together these ensure an entry's live status and metadata
  genuinely come from the server at that URL.
- **The real hostname is shown.** The public stream card displays the entry's
  hostname (from the immutable URL), so a remote server cannot fully masquerade
  behind a spoofed display name.
- **Being featured requires approval.** Another server asking to feature your
  stream appears under "Requests to feature your stream" and must be approved by
  an admin, regardless of whether your server otherwise accepts follows
  automatically.

## What is NOT guaranteed (inherent to an open protocol)

- **"Is this really Owncast?" is self-asserted.** Adding a server validates its
  `nodeinfo` (software name `owncast`, ActivityPub enabled), but any server can
  serve a `nodeinfo` claiming that. Treat this as a sanity check, not proof. The
  real trust anchor is your own curation plus the guarantees above.
- **A featured server controls its own name, logo, title and thumbnail.** These
  are supplied by the remote server and can change over time (including
  bait-and-switch). The link, however, always points to the host you vetted, and
  that host is shown in the UI.

## Practical guidance for operators

- Only feature servers you trust; the directory is a curation you stand behind.
- The shown hostname is the source of truth for where a card leads — judge an
  entry by its hostname, not its display name.
- If a featured server starts misrepresenting itself, unfeature it.

## Hardening implemented

- HTTP-signature verification binds an activity's `actor` to the signing key's
  host (`services/activitypub/inbox/worker.go`).
- IRI/actor/key resolution refuses internal/loopback addresses and caps response
  size; redirects are re-validated (`services/activitypub/resolvers/resolve.go`,
  `utils/tlsconfig.go`).
- Replayed (stale-`Date`) signed requests are rejected within a tolerance window
  (`services/activitypub/inbox/worker.go`).
- Remote-supplied metadata is length/count clamped before storage
  (`services/activitypub/inbox/offer.go`, `accept.go`).
- The featured-follow pins the resolved actor to the entered URL's host
  (`services/activitypub/outbox/follow.go`).
- Remote URLs rendered as links are validated to be `http(s)`
  (`web/components/admin/FederatedServers/FeatureRequests.tsx`), and the public
  card surfaces the real hostname (`web/components/ui/StreamCard/StreamCard.tsx`).
