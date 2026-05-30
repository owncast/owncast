-- +goose Up
-- +goose StatementBegin

-- Featured streams / mini-directory schema additions.
--
-- Tracks which remote Owncast servers we follow so they can advertise live
-- stream status to this server via ActivityPub (Offer/Leave). Also flags
-- ap_followers rows as Owncast peers so we can treat them specially.

CREATE TABLE IF NOT EXISTS federated_servers (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "iri" TEXT NOT NULL UNIQUE,
    "name" TEXT,
    "logo_url" TEXT,
    "is_online" BOOLEAN DEFAULT FALSE,
    "stream_title" TEXT,
    "stream_description" TEXT,
    "stream_tags" TEXT,
    "thumbnail_url" TEXT,
    "last_seen_online" TIMESTAMP,
    "last_status_update" TIMESTAMP,
    "added_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "followed_at" TIMESTAMP,
    "pending" BOOLEAN DEFAULT TRUE,
    "username" TEXT,
    "display_name" TEXT,
    "summary" TEXT,
    "accepted_at" TIMESTAMP,
    "rejected_at" TIMESTAMP,
    "follow_status" TEXT DEFAULT 'pending'
);
CREATE INDEX IF NOT EXISTS federated_servers_iri ON federated_servers (iri);
CREATE INDEX IF NOT EXISTS federated_servers_is_online ON federated_servers (is_online);
CREATE INDEX IF NOT EXISTS federated_servers_last_seen ON federated_servers (last_seen_online);

-- +goose StatementEnd

-- +goose StatementBegin
-- Tag federated followers that are Owncast peers so we can ping/Offer to
-- them without crossing wires with regular Mastodon-style followers.
ALTER TABLE ap_followers ADD COLUMN owncast_server BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS federated_servers_last_seen;
DROP INDEX IF EXISTS federated_servers_is_online;
DROP INDEX IF EXISTS federated_servers_iri;
DROP TABLE IF EXISTS federated_servers;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE ap_followers DROP COLUMN owncast_server;
-- +goose StatementEnd
