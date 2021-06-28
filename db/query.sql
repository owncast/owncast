-- Queries added to query.sql must be compiled into Go code with sqlc. Read README.md for details.

-- Federation related queries.

-- name: GetFollowerCount :one
SElECT count(*) FROM ap_followers;

-- name: GetLocalPostCount :one
SElECT count(*) FROM ap_outbox;

-- name: GetFederationFollowers :many
SELECT iri, inbox, name, username, image FROM ap_followers WHERE approved_at is not null;

-- name: GetFederationFollowerApprovalRequests :many
SELECT iri, inbox, name, username, image FROM ap_followers WHERE approved_at = null;

-- name: ApproveFederationFollower :exec
UPDATE ap_followers SET approved_at = $1 WHERE iri = $2;

-- name: GetOutbox :many
SELECT value FROM ap_outbox;

-- name: GetObjectFromOutboxByID :one
SELECT value FROM ap_outbox WHERE id = $1;

-- name: GetObjectFromOutboxByIRI :one
SELECT value FROM ap_outbox WHERE iri = $1;

-- name: RemoveFollowerByIRI :exec
DELETE FROM ap_followers WHERE iri = $1;

-- name: AddFollower :exec
INSERT INTO ap_followers(iri, inbox, name, username, image, approved_at) values($1, $2, $3, $4, $5, $6);

-- name: AddToOutbox :exec
INSERT INTO ap_outbox(id, iri, value, type) values($1, $2, $3, $4);

-- name: UpdateFollowerByIRI :exec
UPDATE ap_followers SET inbox = $1, name = $2, username = $3, image = $4 WHERE iri = $5;