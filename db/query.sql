-- Queries added to query.sql must be compiled into Go code with sqlc. Read README.md for details.

-- Federation related queries.

-- name: GetFollowerCount :one
SElECT count(*) FROM ap_followers WHERE approved_at is not null;

-- name: GetLocalPostCount :one
SElECT count(*) FROM ap_outbox;

-- name: GetFederationFollowersWithOffset :many
SELECT iri, inbox, name, username, image, created_at FROM ap_followers WHERE approved_at is not null ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: GetRejectedAndBlockedFollowers :many
SELECT iri, name, username, image, created_at, disabled_at FROM ap_followers WHERE disabled_at is not null;

-- name: GetFederationFollowerApprovalRequests :many
SELECT iri, inbox, name, username, image, created_at FROM ap_followers WHERE approved_at IS null AND disabled_at is null;

-- name: ApproveFederationFollower :exec
UPDATE ap_followers SET approved_at = $1, disabled_at = null WHERE iri = $2;

-- name: RejectFederationFollower :exec
UPDATE ap_followers SET approved_at = null, disabled_at = $1 WHERE iri = $2;

-- name: GetFollowerByIRI :one
SELECT iri, inbox, name, username, image, request, request_object, created_at, approved_at, disabled_at FROM ap_followers WHERE iri = $1;

-- name: GetOutboxWithOffset :many
SELECT value FROM ap_outbox LIMIT $1 OFFSET $2;


-- name: GetObjectFromOutboxByIRI :one
SELECT value, live_notification, created_at FROM ap_outbox WHERE iri = $1;

-- name: RemoveFollowerByIRI :exec
DELETE FROM ap_followers WHERE iri = $1;

-- name: AddFollower :exec
INSERT INTO ap_followers(iri, inbox, request, request_object, name, username, image, approved_at) values($1, $2, $3, $4, $5, $6, $7, $8);

-- name: AddToOutbox :exec
INSERT INTO ap_outbox(iri, value, type, live_notification) values($1, $2, $3, $4);

-- name: AddToAcceptedActivities :exec
INSERT INTO ap_accepted_activities(iri, actor, type, timestamp) values($1, $2, $3, $4);

-- name: GetInboundActivityCount :one
SELECT count(*) FROM ap_accepted_activities;

-- name: GetInboundActivitiesWithOffset :many
SELECT iri, actor, type, timestamp FROM ap_accepted_activities ORDER BY timestamp DESC LIMIT $1 OFFSET $2;

-- name: DoesInboundActivityExist :one
SELECT count(*) FROM ap_accepted_activities WHERE iri = $1 AND actor = $2 AND TYPE = $3;

-- name: UpdateFollowerByIRI :exec
UPDATE ap_followers SET inbox = $1, name = $2, username = $3, image = $4 WHERE iri = $5;

-- name: BanIPAddress :exec
INSERT INTO ip_bans(ip_address, notes) values($1, $2);

-- name: RemoveIPAddressBan :exec
DELETE FROM ip_bans WHERE ip_address = $1;

-- name: IsIPAddressBlocked :one
SELECT count(*) FROM ip_bans WHERE ip_address = $1;

-- name: GetIPAddressBans :many
SELECT * FROM ip_bans;
-- name: AddNotification :exec
INSERT INTO notifications (channel, destination) VALUES($1, $2);

-- name: GetNotificationDestinationsForChannel :many
SELECT destination FROM notifications WHERE channel = $1;

-- name: RemoveNotificationDestinationForChannel :exec
DELETE FROM notifications WHERE channel = $1 AND destination = $2;
-- name: AddAuthForUser :exec
INSERT INTO auth(user_id, token, type) values($1, $2, $3);

-- name: GetUserByAuth :one
SELECT users.id, display_name, display_color, users.created_at, disabled_at, previous_names, namechanged_at, authenticated_at, scopes FROM auth, users WHERE token = $1 AND auth.type = $2 AND users.id = auth.user_id;

-- name: AddAccessTokenForUser :exec
INSERT INTO user_access_tokens(token, user_id) values($1, $2);

-- name: GetUserByAccessToken :one
SELECT users.id, display_name, display_color, users.created_at, disabled_at, previous_names, namechanged_at, authenticated_at, scopes FROM users, user_access_tokens WHERE token = $1 AND users.id = user_id;

-- name: GetUserDisplayNameByToken :one
SELECT display_name FROM users, user_access_tokens WHERE token = $1 AND users.id = user_id AND disabled_at = NULL;

-- name: SetAccessTokenToOwner :exec
UPDATE user_access_tokens SET user_id = $1 WHERE token = $2;

-- name: SetUserAsAuthenticated :exec
UPDATE users SET authenticated_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: GetMessagesFromUser :many
SELECT id, body, hidden_at, timestamp FROM messages WHERE eventType = 'CHAT' AND user_id = $1 ORDER BY TIMESTAMP DESC;

-- name: IsDisplayNameAvailable :one
SELECT count(*) FROM users WHERE display_name = $1 AND authenticated_at is not null AND disabled_at is NULL;

-- name: ChangeDisplayName :exec
UPDATE users SET display_name = $1, previous_names = previous_names || $2, namechanged_at = $3 WHERE id = $4;

-- name: ChangeDisplayColor :exec
UPDATE users SET display_color = $1 WHERE id = $2;

-- Recording and clip related queries.

-- name: GetStreams :many
SELECT id, stream_title, start_time, end_time FROM streams ORDER BY start_time DESC;

-- name: GetStreamById :one
SELECT id, stream_title, start_time, end_time FROM streams WHERE id = $1 LIMIT 1;

-- name: GetOutputConfigurationsForStreamId :many
SELECT id, stream_id, variant_id, name, segment_duration, bitrate, framerate, resolution_width, resolution_height FROM video_segment_output_configuration WHERE stream_id = $1;

-- name: GetOutputConfigurationForId :one
SELECT id, stream_id, variant_id, name, segment_duration, bitrate, framerate, resolution_width, resolution_height FROM video_segment_output_configuration WHERE id = $1;

-- name: GetSegmentsForOutputId :many
SELECT id, stream_id, output_configuration_id, path, timestamp FROM video_segments WHERE output_configuration_id = $1 ORDER BY timestamp ASC;

-- name: GetSegmentsForOutputIdAndWindow :many
SELECT id, stream_id, output_configuration_id, path, relative_timestamp, timestamp FROM video_segments WHERE output_configuration_id = $1 AND (cast ( relative_timestamp as int ) - ( relative_timestamp < cast ( relative_timestamp as int ))) >= @start_seconds::REAL AND (cast ( relative_timestamp as int ) + ( relative_timestamp > cast ( relative_timestamp as int ))) <= @end_seconds::REAL ORDER BY relative_timestamp ASC;

-- name: InsertStream :exec
INSERT INTO streams (id, stream_title, start_time, end_time) VALUES($1, $2, $3, $4);

-- name: InsertOutputConfiguration :exec
INSERT INTO video_segment_output_configuration (id, variant_id, stream_id, name, segment_duration, bitrate, framerate, resolution_width, resolution_height, timestamp) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: InsertSegment :exec
INSERT INTO video_segments (id, stream_id, output_configuration_id, path, relative_timestamp, timestamp) VALUES($1, $2, $3, $4, $5, $6);

-- name: SetStreamEnded :exec
UPDATE streams SET end_time = $1 WHERE id = $2;

-- name: InsertClip :exec
INSERT INTO replay_clips (id, stream_id, clip_title, relative_start_time, relative_end_time, timestamp) VALUES($1, $2, $3, $4, $5, $6);

-- name: GetAllClips :many
SELECT rc.id AS id, rc.clip_title, rc.stream_id, rc.relative_start_time, rc.relative_end_time, (rc.relative_end_time - rc.relative_start_time) AS duration_seconds, rc.timestamp, s.stream_title AS stream_title
	FROM replay_clips rc
	JOIN streams s ON rc.stream_id = s.id
	ORDER BY timestamp DESC;

-- name: GetAllClipsForStream :many
SELECT rc.id AS clip_id, rc.stream_id, rc.clipped_by, rc.clip_title, rc.relative_start_time, rc.relative_end_time, rc.timestamp,
	s.id AS stream_id, s.stream_title AS stream_title
	FROM replay_clips rc
	JOIN streams s ON rc.stream_id = s.id
	WHERE rc.stream_id = $1
	ORDER BY timestamp DESC;

-- name: GetClip :one
SELECT id AS clip_id, stream_id, clipped_by, clip_title, timestamp AS clip_timestamp, relative_start_time, relative_end_time FROM replay_clips WHERE id = $1;

-- name: GetFinalSegmentForStream :one
SELECT id, stream_id, output_configuration_id, path, relative_timestamp, timestamp FROM video_segments WHERE stream_id = $1 ORDER BY relative_timestamp DESC LIMIT 1;
