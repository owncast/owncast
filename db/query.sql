-- Queries added to query.sql must be compiled into Go code with sqlc. Read README.md for details.

-- Federation related queries.

-- name: GetFollowerCount :one
SELECT count(*) FROM ap_followers WHERE approved_at is not null;

-- name: GetLocalPostCount :one
SELECT count(*) FROM ap_outbox;

-- name: GetFederationFollowersWithOffset :many
SELECT iri, inbox, shared_inbox, name, username, image, created_at FROM ap_followers WHERE approved_at is not null ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- name: GetRejectedAndBlockedFollowers :many
SELECT iri, name, username, image, created_at, disabled_at FROM ap_followers WHERE disabled_at is not null;

-- name: GetFederationFollowerApprovalRequests :many
SELECT iri, inbox, shared_inbox, name, username, image, created_at FROM ap_followers WHERE approved_at IS null AND disabled_at is null;

-- name: ApproveFederationFollower :exec
UPDATE ap_followers SET approved_at = ?, disabled_at = null WHERE iri = ?;

-- name: RejectFederationFollower :exec
UPDATE ap_followers SET approved_at = null, disabled_at = ? WHERE iri = ?;

-- name: GetFollowerByIRI :one
SELECT iri, inbox, shared_inbox, name, username, image, request, request_object, created_at, approved_at, disabled_at, owncast_server FROM ap_followers WHERE iri = ?;

-- name: GetOutboxWithOffset :many
SELECT value FROM ap_outbox LIMIT ? OFFSET ?;


-- name: GetObjectFromOutboxByIRI :one
SELECT value, live_notification, created_at FROM ap_outbox WHERE iri = ?;

-- name: RemoveFollowerByIRI :exec
DELETE FROM ap_followers WHERE iri = ?;

-- name: AddFollower :exec
INSERT INTO ap_followers(iri, inbox, shared_inbox, request, request_object, name, username, image, approved_at, owncast_server) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

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
