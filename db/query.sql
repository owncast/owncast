-- Queries added to query.sql must be compiled into Go code with sqlc. Read README.md for details.

-- Federation related queries.

-- name: GetFollowerCount :one
-- Featured-streams follows (another Owncast server following us so it can show
-- our live status in its directory) are excluded: they are a directory
-- relationship, not a fan follow, so they must not inflate the follower count.
SELECT count(*) FROM ap_followers WHERE approved_at is not null AND directory IS NOT 1;

-- name: GetLocalPostCount :one
-- Only posts (Notes) count: the outbox table also stores other
-- dereferenceable objects, such as QuoteAuthorization stamps, that are not
-- posts and must not inflate the public post count.
SELECT count(*) FROM ap_outbox WHERE type = 'Note';

-- name: GetFederationFollowersWithOffset :many
-- Excludes featured-streams (Owncast-server) follows so they don't show up in
-- the public or admin followers list; they are tracked as a directory
-- relationship, not surfaced as followers.
-- rowid breaks created_at ties (second resolution) so pagination order is
-- deterministic; without it, same-second follows shift between pages.
SELECT iri, inbox, shared_inbox, name, username, image, created_at FROM ap_followers WHERE approved_at is not null AND directory IS NOT 1 ORDER BY created_at DESC, rowid DESC LIMIT ? OFFSET ?;

-- name: GetRejectedAndBlockedFollowers :many
SELECT iri, name, username, image, created_at, disabled_at FROM ap_followers WHERE disabled_at is not null;

-- name: GetFederationFollowerApprovalRequests :many
-- Regular (fan) follow approval requests only. Featured-streams (Owncast
-- server) requests are excluded here and surfaced separately via
-- GetPendingFeaturedFollowRequests so they can be approved from the featured
-- streams admin instead of the followers admin.
SELECT iri, inbox, shared_inbox, name, username, image, created_at FROM ap_followers WHERE approved_at IS null AND disabled_at is null AND directory IS NOT 1;

-- name: GetPendingFeaturedFollowRequests :many
-- Pending requests from other Owncast servers asking to feature this server's
-- stream in their directory. These always require explicit approval.
SELECT iri, inbox, shared_inbox, name, username, image, created_at FROM ap_followers WHERE approved_at IS null AND disabled_at is null AND directory IS 1 ORDER BY created_at DESC;

-- name: ApproveFederationFollower :exec
UPDATE ap_followers SET approved_at = ?, disabled_at = null WHERE iri = ?;

-- name: RejectFederationFollower :exec
UPDATE ap_followers SET approved_at = null, disabled_at = ? WHERE iri = ?;

-- name: GetFollowerByIRI :one
SELECT iri, inbox, shared_inbox, name, username, image, request, request_object, created_at, approved_at, disabled_at, directory FROM ap_followers WHERE iri = ?;

-- name: GetOutboxWithOffset :many
-- Only posts (Notes) are listed in the public outbox collection. Other stored
-- objects, such as QuoteAuthorization stamps, stay fetchable by IRI but are
-- not part of the collection.
SELECT value FROM ap_outbox WHERE type = 'Note' LIMIT ? OFFSET ?;


-- name: GetObjectFromOutboxByIRI :one
SELECT value, live_notification, created_at FROM ap_outbox WHERE iri = ?;

-- name: GetNoteFromOutboxByIRI :one
-- Like GetObjectFromOutboxByIRI but only matches posts (Notes). Used to
-- decide whether an object may be quoted: stamps and directory pings are
-- stored in the same table and are not quotable.
SELECT value, live_notification, created_at FROM ap_outbox WHERE iri = ? AND type = 'Note';

-- name: RemoveFollowerByIRI :exec
DELETE FROM ap_followers WHERE iri = ?;

-- name: AddFollower :exec
INSERT INTO ap_followers(iri, inbox, shared_inbox, request, request_object, name, username, image, approved_at, directory) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: AddToOutbox :exec
INSERT INTO ap_outbox(iri, value, type, live_notification) values(?, ?, ?, ?);

-- name: UpsertOutboxObject :exec
INSERT INTO ap_outbox(iri, value, type, live_notification) VALUES(?, ?, ?, ?)
ON CONFLICT(iri) DO UPDATE SET value = excluded.value, type = excluded.type, live_notification = excluded.live_notification;

-- name: QueueActivityPubDelivery :one
INSERT INTO ap_delivery_queue (
    inbox,
    payload,
    actor_iri,
    activity_type,
    coalesce_key,
    next_attempt_at
) VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT(inbox, coalesce_key) WHERE coalesce_key IS NOT NULL AND failed_at IS NULL
DO UPDATE SET
    payload = excluded.payload,
    actor_iri = excluded.actor_iri,
    activity_type = excluded.activity_type,
    next_attempt_at = excluded.next_attempt_at,
    attempts = 0,
    last_error = NULL,
    revision = ap_delivery_queue.revision + 1
RETURNING id;

-- name: ClaimActivityPubDelivery :one
UPDATE ap_delivery_queue
SET claimed_until = @claimed_until,
    attempts = attempts + 1
WHERE id = (
    SELECT candidate.id
    FROM ap_delivery_queue AS candidate
    WHERE candidate.failed_at IS NULL
      AND candidate.next_attempt_at <= @now
      AND (candidate.claimed_until IS NULL OR candidate.claimed_until <= @now)
    ORDER BY candidate.next_attempt_at, candidate.id
    LIMIT 1
)
RETURNING id, inbox, payload, actor_iri, activity_type, coalesce_key, attempts, revision;

-- name: CompleteActivityPubDelivery :execrows
DELETE FROM ap_delivery_queue WHERE id = ? AND revision = ?;

-- name: RetryActivityPubDelivery :execrows
UPDATE ap_delivery_queue
SET claimed_until = NULL,
    next_attempt_at = ?,
    last_error = ?
WHERE id = ? AND revision = ?;

-- name: DeferActivityPubDelivery :execrows
UPDATE ap_delivery_queue
SET claimed_until = NULL,
    next_attempt_at = ?,
    attempts = MAX(attempts - 1, 0)
WHERE id = ? AND revision = ?;

-- name: FailActivityPubDelivery :execrows
UPDATE ap_delivery_queue
SET claimed_until = NULL,
    last_error = ?,
    failed_at = ?
WHERE id = ? AND revision = ?;

-- name: ReleaseSupersededActivityPubDelivery :exec
UPDATE ap_delivery_queue SET claimed_until = NULL WHERE id = ? AND revision != ?;

-- name: AddToAcceptedActivities :exec
INSERT INTO ap_accepted_activities(iri, actor, type, timestamp) values(?, ?, ?, ?);

-- name: GetInboundActivityCount :one
SELECT count(*) FROM ap_accepted_activities;

-- name: GetInboundActivitiesWithOffset :many
SELECT iri, actor, type, timestamp FROM ap_accepted_activities ORDER BY timestamp DESC LIMIT ? OFFSET ?;

-- name: DoesInboundActivityExist :one
SELECT count(*) FROM ap_accepted_activities WHERE iri = ? AND actor = ? AND TYPE = ?;

-- name: UpdateFollowerByIRI :exec
UPDATE ap_followers SET inbox = ?, shared_inbox = ?, name = ?, username = ?, image = ? WHERE iri = ?;

-- name: GetFollowersToValidate :many
SELECT iri, inbox, shared_inbox, name, username, image, first_validation_failure_at
FROM ap_followers
WHERE approved_at IS NOT NULL AND disabled_at IS NULL
ORDER BY last_validated_at ASC NULLS FIRST
LIMIT ?;

-- name: UpdateFollowerValidationSuccess :exec
UPDATE ap_followers
SET last_validated_at = ?, first_validation_failure_at = NULL
WHERE iri = ?;

-- name: UpdateFollowerValidationFailure :exec
UPDATE ap_followers
SET last_validated_at = @last_validated_at, first_validation_failure_at = COALESCE(first_validation_failure_at, @last_validated_at)
WHERE iri = @iri;

-- name: GetUniqueDeliveryInboxes :many
SELECT COALESCE(shared_inbox, inbox) as delivery_inbox FROM ap_followers WHERE approved_at is not null GROUP BY delivery_inbox;

-- name: GetUniqueDirectoryDeliveryInboxes :many
-- Approved directory followers only. The Offer/Leave stream pings are delivered
-- here, not to fan followers, who only need the go-live Create/Note.
SELECT COALESCE(shared_inbox, inbox) as delivery_inbox FROM ap_followers WHERE approved_at is not null AND directory IS 1 GROUP BY delivery_inbox;

-- name: GetApprovedDirectoryFollowers :many
-- Approved directories that are featuring/listing this server. Shown in the
-- admin so the operator can review and remove them.
SELECT iri, inbox, shared_inbox, name, username, image, created_at FROM ap_followers WHERE approved_at IS NOT NULL AND disabled_at IS NULL AND directory IS 1 ORDER BY created_at DESC;

-- name: BanIPAddress :exec
INSERT INTO ip_bans(ip_address, notes) values(?, ?);

-- name: RemoveIPAddressBan :exec
DELETE FROM ip_bans WHERE ip_address = ?;

-- name: IsIPAddressBlocked :one
SELECT count(*) FROM ip_bans WHERE ip_address = ?;

-- name: GetIPAddressBans :many
SELECT * FROM ip_bans;

-- name: AddNotification :exec
INSERT INTO notifications (channel, destination) VALUES(?, ?);

-- name: GetNotificationDestinationsForChannel :many
SELECT destination FROM notifications WHERE channel = ?;

-- name: RemoveNotificationDestinationForChannel :exec
DELETE FROM notifications WHERE channel = ? AND destination = ?;

-- name: AddAuthForUser :exec
INSERT INTO auth(user_id, auth_key, type, provider, profile_url, handle, is_public) values(?, ?, ?, ?, ?, ?, ?);

-- name: GetUserByAuth :one
SELECT users.id, display_name, display_color, users.created_at, disabled_at, disabled_reason, previous_names, namechanged_at, authenticated_at, scopes FROM auth, users WHERE auth_key = ? AND auth.type = ? AND users.id = auth.user_id;

-- name: GetUserByPluginAuth :one
-- Plugin-auth lookup keyed on (provider, auth_key). The type literal mirrors
-- models.PluginAuth and pins the lookup to the plugin namespace, so a plugin
-- slug can never resolve a built-in identity.
SELECT users.id, display_name, display_color, users.created_at, disabled_at, disabled_reason, previous_names, namechanged_at, authenticated_at, scopes FROM auth, users WHERE auth.type = 'plugin.auth' AND auth.provider = ? AND auth.auth_key = ? AND users.id = auth.user_id;

-- name: CountUserAuthByProvider :one
SELECT count(*) FROM auth WHERE user_id = ? AND type = ? AND provider = ?;

-- name: GetAuthForUsers :many
-- External auth identities (type + provider) for the given users, so the admin
-- user list can show how each authenticated user signed in.
SELECT user_id, type, provider FROM auth WHERE user_id IN (sqlc.slice('user_ids'));

-- name: AddAccessTokenForUser :exec
INSERT INTO user_access_tokens(token, user_id) values(?, ?);

-- name: GetUserByAccessToken :one
SELECT users.id, display_name, display_color, users.created_at, disabled_at, disabled_reason, previous_names, namechanged_at, authenticated_at, scopes, users.type = 'API' AS is_bot FROM users, user_access_tokens WHERE token = ? AND users.id = user_id;

-- name: GetUserByID :one
SELECT id, display_name, display_color, created_at, disabled_at, disabled_reason, previous_names, namechanged_at, authenticated_at, scopes, type = 'API' AS is_bot FROM users WHERE id = ?;

-- name: GetUsers :many
SELECT id, display_name, display_color, created_at, disabled_at, disabled_reason, previous_names, namechanged_at, authenticated_at, scopes, type = 'API' AS is_bot FROM users ORDER BY created_at DESC;

-- name: GetUserDisplayNameByToken :one
SELECT display_name FROM users JOIN user_access_tokens ON users.id = user_access_tokens.user_id WHERE token = ? AND users.disabled_at IS NULL;

-- name: SetAccessTokenToOwner :exec
UPDATE user_access_tokens SET user_id = ? WHERE token = ?;

-- name: SetUserAsAuthenticated :exec
UPDATE users SET authenticated_at = CURRENT_TIMESTAMP WHERE id = ?;

-- name: GetMessagesFromUser :many
SELECT id, body, hidden_at, timestamp FROM messages WHERE eventType = 'CHAT' AND user_id = ? ORDER BY TIMESTAMP DESC;

-- name: IsDisplayNameAvailable :one
SELECT count(*) FROM users WHERE display_name = ? AND ( type='API' OR authenticated_at IS NOT NULL ) AND disabled_at IS NULL;

-- name: ChangeDisplayName :exec
UPDATE users SET display_name = ?, previous_names = previous_names || ?, namechanged_at = ? WHERE id = ?;

-- name: ChangeDisplayColor :exec
UPDATE users SET display_color = ? WHERE id = ?;

-- name: GetUsersPaginated :many
-- A page of users of every type (chat viewers, authenticated/plugin users, and
-- API integrations), newest first, filtered to display names containing
-- @search and an optional @status: '' or 'all' = every user; otherwise
-- 'active' (not banned) or 'bots' (API users). An empty @search matches every
-- user (LIKE '%%'). The complete banned and moderator sets come from
-- GetDisabledUsers and GetModeratorUsers instead, since those need every match
-- rather than a page. GetUsersPaginatedAsc is the oldest-first counterpart;
-- the two exist as separate queries because sqlc can't bind a sort direction
-- into ORDER BY.
SELECT id, display_name, display_color, created_at, disabled_at, disabled_reason, previous_names, namechanged_at, authenticated_at, scopes, type = 'API' AS is_bot
FROM users
WHERE display_name LIKE '%' || @search || '%'
  AND (
    @status = '' OR @status = 'all'
    OR (@status = 'active' AND disabled_at IS NULL)
    OR (@status = 'bots' AND type = 'API')
  )
ORDER BY created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: GetUsersPaginatedAsc :many
-- Oldest-first counterpart to GetUsersPaginated. Same @search/@status filter,
-- only the created_at ordering differs.
SELECT id, display_name, display_color, created_at, disabled_at, disabled_reason, previous_names, namechanged_at, authenticated_at, scopes, type = 'API' AS is_bot
FROM users
WHERE display_name LIKE '%' || @search || '%'
  AND (
    @status = '' OR @status = 'all'
    OR (@status = 'active' AND disabled_at IS NULL)
    OR (@status = 'bots' AND type = 'API')
  )
ORDER BY created_at ASC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountUsers :one
-- Total number of users matching the same @search/@status filter as GetUsersPaginated.
SELECT count(*) FROM users
WHERE display_name LIKE '%' || @search || '%'
  AND (
    @status = '' OR @status = 'all'
    OR (@status = 'active' AND disabled_at IS NULL)
    OR (@status = 'bots' AND type = 'API')
  );

-- name: DeleteUserAccessTokens :exec
DELETE FROM user_access_tokens WHERE user_id = ?;

-- name: DeleteUserAuth :exec
DELETE FROM auth WHERE user_id = ?;

-- name: DeleteUserMessages :exec
DELETE FROM messages WHERE user_id = ?;

-- name: DeleteUserByID :execrows
DELETE FROM users WHERE id = ?;

-- Federated servers queries

-- name: GetFederatedServers :many
SELECT id, iri, name, logo_url, is_online, stream_title, stream_description, stream_tags, thumbnail_url, last_seen_online, last_status_update, added_at, followed_at, pending, username, display_name, summary, accepted_at, rejected_at, follow_status FROM federated_servers ORDER BY added_at DESC;

-- name: GetFederatedServer :one
SELECT id, iri, name, logo_url, is_online, stream_title, stream_description, stream_tags, thumbnail_url, last_seen_online, last_status_update, added_at, followed_at, pending, username, display_name, summary, accepted_at, rejected_at, follow_status FROM federated_servers WHERE iri = ?;

-- name: AddFederatedServer :exec
INSERT INTO federated_servers(iri, name, logo_url, followed_at, pending, username, follow_status) values(?, ?, ?, ?, ?, ?, ?);

-- name: UpdateFederatedServerStatus :exec
UPDATE federated_servers SET is_online = ?, stream_title = ?, stream_description = ?, stream_tags = ?, thumbnail_url = ?, last_status_update = ? WHERE iri = ?;

-- name: UpdateFederatedServerOnlineStatus :exec
UPDATE federated_servers SET is_online = ?, last_seen_online = ?, last_status_update = ? WHERE iri = ?;

-- name: RemoveFederatedServer :exec
DELETE FROM federated_servers WHERE id = ?;

-- name: UpdateFederatedServerFollowStatus :exec
UPDATE federated_servers SET follow_status = ?, pending = ?, accepted_at = ?, rejected_at = ? WHERE iri = ?;

-- name: UpdateFederatedServerMetadata :exec
UPDATE federated_servers SET name = ?, display_name = ?, summary = ?, logo_url = ? WHERE iri = ?;

-- name: GetPendingFederatedServers :many
SELECT id, iri, name, logo_url, is_online, stream_title, stream_description, stream_tags, thumbnail_url, last_seen_online, last_status_update, added_at, followed_at, pending, username, display_name, summary, accepted_at, rejected_at, follow_status FROM federated_servers WHERE pending = true ORDER BY added_at DESC;

-- Scheduled stream events (schedule feature).

-- name: AddStreamEventSeries :exec
INSERT INTO stream_event_series(id, name, description, reminder_message, recurrence, duration_minutes) VALUES(?, ?, ?, ?, ?, ?);

-- name: GetStreamEventSeries :one
SELECT * FROM stream_event_series WHERE id = ?;

-- name: GetAllStreamEventSeries :many
SELECT * FROM stream_event_series ORDER BY created_at;

-- name: GetActiveStreamEventSeries :many
SELECT * FROM stream_event_series WHERE active = TRUE ORDER BY created_at;

-- name: UpdateStreamEventSeries :exec
UPDATE stream_event_series SET name = ?, description = ?, reminder_message = ?, recurrence = ?, duration_minutes = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?;

-- name: SetStreamEventSeriesActive :exec
UPDATE stream_event_series SET active = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?;

-- name: DeleteStreamEventSeries :exec
DELETE FROM stream_event_series WHERE id = ?;

-- name: AddStreamEvent :execrows
-- The targeted conflict clause keeps materialization idempotent:
-- re-expanding a series hits the UNIQUE(series_id, original_start) index
-- and skips existing rows, while any other conflict (a generated primary
-- key colliding with an existing row) still errors instead of silently
-- reporting the slot as pre-existing.
INSERT INTO stream_events(id, series_id, original_start, name, description, reminder_message, start_time, duration_minutes, timezone) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(series_id, original_start) DO NOTHING;

-- name: GetStreamEvent :one
SELECT * FROM stream_events WHERE id = ?;

-- name: GetStreamEventsInRange :many
SELECT * FROM stream_events WHERE start_time >= ? AND start_time < ? ORDER BY start_time;

-- name: GetStreamEventsForSeries :many
SELECT * FROM stream_events WHERE series_id = ? ORDER BY start_time;

-- name: UpdateStreamEventDetails :exec
UPDATE stream_events SET name = ?, description = ?, reminder_message = ?, duration_minutes = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?;

-- name: CancelStreamEvent :exec
-- The row is kept: it holds the federation state needed to announce the
-- cancellation, and its original_start keeps the materializer from
-- re-creating the slot.
UPDATE stream_events SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP WHERE id = ?;

-- name: MoveStreamEvent :exec
-- original_start is deliberately untouched: it is the identity the
-- materializer keys on, so the vacated slot is not re-inserted. A moved
-- future event enters a new warning window.
UPDATE stream_events SET start_time = ?, webhook_warning_sent_at = NULL, updated_at = CURRENT_TIMESTAMP WHERE id = ?;

-- name: DeleteStreamEvent :exec
DELETE FROM stream_events WHERE id = ?;

-- name: DeleteUnfederatedFutureStreamEventsForSeries :exec
-- Series edits regenerate only rows nobody has seen: future occurrences that
-- never federated. Announced rows stay and get Update/Delete activities.
-- Cancelled rows also stay, whatever their federation state: deleting one
-- would let the materializer resurrect the slot as a fresh scheduled row.
DELETE FROM stream_events WHERE series_id = ? AND federated_at IS NULL AND start_time > ? AND status != 'cancelled';

-- name: GetCurrentOrUpcomingStreamEvents :many
-- Events that are still running (start + duration in the future) or have
-- not started yet. Powers the status endpoint's next-event answer, so a
-- stream that starts late keeps its chat window open until the event's
-- scheduled end passes.
SELECT * FROM stream_events WHERE status = 'scheduled' AND datetime(start_time, '+' || duration_minutes || ' minutes') > datetime(?) ORDER BY start_time LIMIT ?;

-- name: GetNextUpcomingStreamEvents :many
SELECT * FROM stream_events WHERE status = 'scheduled' AND start_time > ? ORDER BY start_time LIMIT ?;

-- name: GetStreamEventsToFederate :many
SELECT * FROM stream_events WHERE federated_at IS NULL AND status = 'scheduled' AND start_time > ? ORDER BY start_time;

-- name: SetStreamEventFederatedAt :exec
UPDATE stream_events SET federated_at = ? WHERE id = ?;

-- name: GetStreamEventsNeedingReminder1 :many
SELECT * FROM stream_events WHERE reminder_1_sent_at IS NULL AND status = 'scheduled' AND start_time > ? AND start_time <= ? ORDER BY start_time;

-- name: GetStreamEventsNeedingReminder2 :many
SELECT * FROM stream_events WHERE reminder_2_sent_at IS NULL AND status = 'scheduled' AND start_time > ? AND start_time <= ? ORDER BY start_time;

-- name: SetStreamEventReminder1SentAt :exec
UPDATE stream_events SET reminder_1_sent_at = ? WHERE id = ?;

-- name: SetStreamEventReminder2SentAt :exec
UPDATE stream_events SET reminder_2_sent_at = ? WHERE id = ?;

-- name: GetStreamEventsNeedingWebhookWarning :many
SELECT * FROM stream_events WHERE webhook_warning_sent_at IS NULL AND status = 'scheduled' AND start_time > ? AND start_time <= ? ORDER BY start_time;

-- name: GetStreamEventsNeedingWebhookStart :many
SELECT * FROM stream_events WHERE webhook_started_sent_at IS NULL AND status = 'scheduled' AND start_time <= ? ORDER BY start_time;

-- name: GetStreamEventsNeedingWebhookEnd :many
SELECT * FROM stream_events WHERE webhook_started_sent_at IS NOT NULL AND webhook_ended_sent_at IS NULL AND datetime(start_time, '+' || duration_minutes || ' minutes') <= datetime(?) ORDER BY start_time;

-- name: SetStreamEventWebhookWarningSentAt :exec
UPDATE stream_events SET webhook_warning_sent_at = ? WHERE id = ?;

-- name: SetStreamEventWebhookStartedSentAt :exec
UPDATE stream_events SET webhook_started_sent_at = ? WHERE id = ?;

-- name: SetStreamEventWebhookEndedSentAt :exec
UPDATE stream_events SET webhook_ended_sent_at = ? WHERE id = ?;
