-- +goose Up
-- +goose StatementBegin

-- Baseline schema migration.
--
-- This migration represents the consolidated database schema as of the
-- transition from the legacy hand-rolled migration system (previously tracked
-- via the config.version key) to goose.
--
-- For existing installs at legacy schema version 9, the catch-up path in
-- persistence/legacymigrations runs first to bring them to v9; this baseline
-- is then a sequence of no-ops (every statement is IF NOT EXISTS) and goose
-- simply records it as applied.
--
-- For fresh installs, this migration creates every table and index the
-- application expects.
--
-- Table ordering matches the runtime setup order of the previous system so
-- that deterministic index-name collisions (see notes below) resolve to the
-- same table they resolve to in existing installs.

-- Key-value store used for application config values.
-- NOTE: column types were historically declared as the non-standard
-- "string" affinity. They're written here with standard SQLite types
-- (TEXT) so the sqlc SQLite parser can infer Go types correctly; the
-- stored values are identical in practice because SQLite's type
-- affinity system is tolerant of either declaration.
CREATE TABLE IF NOT EXISTS datastore (
    "key" TEXT NOT NULL PRIMARY KEY,
    "value" BLOB,
    "timestamp" DATE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS webhooks (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "url" TEXT NOT NULL,
    "events" TEXT NOT NULL,
    "timestamp" DATETIME DEFAULT CURRENT_TIMESTAMP,
    "last_used" DATETIME
);

CREATE TABLE IF NOT EXISTS users (
    "id" TEXT NOT NULL,
    "display_name" TEXT NOT NULL,
    "display_color" INTEGER NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "disabled_at" TIMESTAMP,
    "previous_names" TEXT DEFAULT '',
    "namechanged_at" TIMESTAMP,
    "authenticated_at" TIMESTAMP,
    "scopes" TEXT,
    "type" TEXT DEFAULT 'STANDARD',
    "last_used" DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS idx_user_id ON users (id);
CREATE INDEX IF NOT EXISTS idx_user_id_disabled ON users (id, disabled_at);
CREATE INDEX IF NOT EXISTS idx_user_disabled_at ON users (disabled_at);

CREATE TABLE IF NOT EXISTS user_access_tokens (
    "token" TEXT NOT NULL PRIMARY KEY,
    "user_id" TEXT NOT NULL,
    "timestamp" DATE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS ap_followers (
    "iri" TEXT NOT NULL,
    "inbox" TEXT NOT NULL,
    "shared_inbox" TEXT,
    "name" TEXT,
    "username" TEXT NOT NULL,
    "image" TEXT,
    "request" TEXT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "approved_at" TIMESTAMP,
    "disabled_at" TIMESTAMP,
    "request_object" BLOB,
    "last_validated_at" TIMESTAMP,
    "first_validation_failure_at" TIMESTAMP,
    PRIMARY KEY (iri)
);
CREATE INDEX IF NOT EXISTS idx_iri ON ap_followers (iri);
CREATE INDEX IF NOT EXISTS idx_ap_followers_iri ON ap_followers (iri);
CREATE INDEX IF NOT EXISTS idx_approved_at ON ap_followers (approved_at);

CREATE TABLE IF NOT EXISTS ap_outbox (
    "iri" TEXT NOT NULL,
    "value" BLOB,
    "type" TEXT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "live_notification" BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (iri)
);
CREATE INDEX IF NOT EXISTS idx_iri ON ap_outbox (iri);
CREATE INDEX IF NOT EXISTS idx_ap_outbox_iri ON ap_outbox (iri);
CREATE INDEX IF NOT EXISTS idx_type ON ap_outbox (type);
CREATE INDEX IF NOT EXISTS idx_live_notification ON ap_outbox (live_notification);

CREATE TABLE IF NOT EXISTS ap_accepted_activities (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "iri" TEXT NOT NULL,
    "actor" TEXT NOT NULL,
    "type" TEXT NOT NULL,
    "timestamp" TIMESTAMP NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_iri_actor_index ON ap_accepted_activities (iri, actor);

CREATE TABLE IF NOT EXISTS notifications (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "channel" TEXT NOT NULL,
    "destination" TEXT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_channel ON notifications (channel);

CREATE TABLE IF NOT EXISTS messages (
    "id" TEXT NOT NULL,
    "user_id" TEXT,
    "body" TEXT,
    "eventType" TEXT,
    "hidden_at" DATETIME,
    "timestamp" DATETIME,
    "title" TEXT,
    "subtitle" TEXT,
    "image" TEXT,
    "link" TEXT,
    PRIMARY KEY (id)
);
CREATE INDEX IF NOT EXISTS user_id_hidden_at_timestamp ON messages (id, user_id, hidden_at, timestamp);
CREATE INDEX IF NOT EXISTS idx_id ON messages (id);
CREATE INDEX IF NOT EXISTS idx_user_id ON messages (user_id);
CREATE INDEX IF NOT EXISTS idx_messages_user_id ON messages (user_id);
CREATE INDEX IF NOT EXISTS idx_hidden_at ON messages (hidden_at);
CREATE INDEX IF NOT EXISTS idx_timestamp ON messages (timestamp);
CREATE INDEX IF NOT EXISTS idx_messages_hidden_at_timestamp ON messages (hidden_at, timestamp);

CREATE TABLE IF NOT EXISTS auth (
    "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "user_id" TEXT NOT NULL,
    "token" TEXT NOT NULL,
    "type" TEXT NOT NULL,
    "timestamp" DATE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE INDEX IF NOT EXISTS idx_auth_token ON auth (token);

CREATE TABLE IF NOT EXISTS ip_bans (
    "ip_address" TEXT NOT NULL PRIMARY KEY,
    "notes" TEXT,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- The baseline migration intentionally has no Down. Reverting the entire
-- schema is not a meaningful operation for this application.
SELECT 'baseline migration has no down step';
-- +goose StatementEnd
