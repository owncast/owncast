-- Queries added to query.sql must be compiled into Go code with sqlc. Read README.md for details.

-- Federation related queries.

-- name: GetFollowerCount :one
-- Featured-streams follows (another Owncast server following us so it can show
-- our live status in its directory) are excluded: they are a directory
-- relationship, not a fan follow, so they must not inflate the follower count.
SELECT count(*) FROM ap_followers WHERE approved_at is not null AND directory IS NOT 1;

-- name: GetLocalPostCount :one
SELECT count(*) FROM ap_outbox;

-- name: GetFederationFollowersWithOffset :many
-- Excludes featured-streams (Owncast-server) follows so they don't show up in
-- the public or admin followers list; they are tracked as a directory
-- relationship, not surfaced as followers.
SELECT iri, inbox, shared_inbox, name, username, image, created_at FROM ap_followers WHERE approved_at is not null AND directory IS NOT 1 ORDER BY created_at DESC LIMIT ? OFFSET ?;

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
SELECT value FROM ap_outbox LIMIT ? OFFSET ?;


-- name: GetObjectFromOutboxByIRI :one
SELECT value, live_notification, created_at FROM ap_outbox WHERE iri = ?;

-- name: RemoveFollowerByIRI :exec
DELETE FROM ap_followers WHERE iri = ?;

-- name: AddFollower :exec
INSERT INTO ap_followers(iri, inbox, shared_inbox, request, request_object, name, username, image, approved_at, directory) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: AddToOutbox :exec
INSERT INTO ap_outbox(iri, value, type, live_notification) values(?, ?, ?, ?);

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
INSERT INTO auth(user_id, token, type) values(?, ?, ?);

-- name: GetUserByAuth :one
SELECT users.id, display_name, display_color, users.created_at, disabled_at, previous_names, namechanged_at, authenticated_at, scopes FROM auth, users WHERE token = ? AND auth.type = ? AND users.id = auth.user_id;

-- name: CountUserAuthByTypeAndTokenPrefix :one
SELECT count(*) FROM auth WHERE user_id = ? AND type = ? AND token LIKE ?;

-- name: GetAuthForUsers :many
-- External auth identities (type + token) for the given users, so the admin
-- user list can show how each authenticated user signed in.
SELECT user_id, type, token FROM auth WHERE user_id IN (sqlc.slice('user_ids'));

-- name: AddAccessTokenForUser :exec
INSERT INTO user_access_tokens(token, user_id) values(?, ?);

-- name: GetUserByAccessToken :one
SELECT users.id, display_name, display_color, users.created_at, disabled_at, previous_names, namechanged_at, authenticated_at, scopes, users.type = 'API' AS is_bot FROM users, user_access_tokens WHERE token = ? AND users.id = user_id;

-- name: GetUserByID :one
SELECT id, display_name, display_color, created_at, disabled_at, previous_names, namechanged_at, authenticated_at, scopes, type = 'API' AS is_bot FROM users WHERE id = ?;

-- name: GetUsers :many
SELECT id, display_name, display_color, created_at, disabled_at, previous_names, namechanged_at, authenticated_at, scopes, type = 'API' AS is_bot FROM users ORDER BY created_at DESC;

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
-- 'active' (not banned), 'banned' (disabled), 'moderators', or 'bots' (API
-- users). An empty @search matches every user (LIKE '%%').
SELECT id, display_name, display_color, created_at, disabled_at, previous_names, namechanged_at, authenticated_at, scopes, type = 'API' AS is_bot
FROM users
WHERE display_name LIKE '%' || @search || '%'
  AND (
    @status = '' OR @status = 'all'
    OR (@status = 'active' AND disabled_at IS NULL)
    OR (@status = 'banned' AND disabled_at IS NOT NULL)
    OR (@status = 'bots' AND type = 'API')
    OR (@status = 'moderators' AND scopes LIKE '%MODERATOR%')
  )
ORDER BY created_at DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: CountUsers :one
-- Total number of users matching the same @search/@status filter as GetUsersPaginated.
SELECT count(*) FROM users
WHERE display_name LIKE '%' || @search || '%'
  AND (
    @status = '' OR @status = 'all'
    OR (@status = 'active' AND disabled_at IS NULL)
    OR (@status = 'banned' AND disabled_at IS NOT NULL)
    OR (@status = 'bots' AND type = 'API')
    OR (@status = 'moderators' AND scopes LIKE '%MODERATOR%')
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
