-- Schema update to query.sql must be referenced in queries located in query.sql
-- and compiled into code with sqlc. Read README.md for details.

CREATE TABLE IF NOT EXISTS ap_followers (
		"iri" TEXT NOT NULL,
		"inbox" TEXT NOT NULL,
		"name" TEXT,
		"username" TEXT NOT NULL,
		"image" TEXT,
    "request" TEXT NOT NULL,
    "request_object" BLOB,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"approved_at" TIMESTAMP,
    "disabled_at" TIMESTAMP,
		PRIMARY KEY (iri));
		CREATE INDEX iri_index ON ap_followers (iri);
    CREATE INDEX approved_at_index ON ap_followers (approved_at);


CREATE TABLE IF NOT EXISTS ap_outbox (
		"iri" TEXT NOT NULL,
		"value" BLOB,
		"type" TEXT NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "live_notification" BOOLEAN DEFAULT FALSE,
		PRIMARY KEY (iri));
		CREATE INDEX iri ON ap_outbox (iri);
		CREATE INDEX type ON ap_outbox (type);
    CREATE INDEX live_notification ON ap_outbox (live_notification);

CREATE TABLE IF NOT EXISTS ap_accepted_activities (
    "id" INTEGER NOT NULL PRIMARY KEY,
		"iri" TEXT NOT NULL,
    "actor" TEXT NOT NULL,
    "type" TEXT NOT NULL,
		"timestamp" TIMESTAMP NOT NULL
	);
	CREATE INDEX iri_actor_index ON ap_accepted_activities (iri,actor);

  CREATE TABLE IF NOT EXISTS ip_bans (
    "ip_address" TEXT NOT NULL PRIMARY KEY,
    "notes" TEXT,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );

CREATE TABLE IF NOT EXISTS notifications (
    "id" INTEGER NOT NULL PRIMARY KEY,
		"channel" TEXT NOT NULL,
		"destination" TEXT NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
		CREATE INDEX channel_index ON notifications (channel);

CREATE TABLE IF NOT EXISTS users (
		"id" TEXT,
		"display_name" TEXT NOT NULL,
		"display_color" INTEGER NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"disabled_at" TIMESTAMP,
		"previous_names" TEXT DEFAULT '',
		"namechanged_at" TIMESTAMP,
		"scopes" TEXT,
    "authenticated_at" TIMESTAMP,
		"type" TEXT DEFAULT 'STANDARD',
		"last_used" DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	);

CREATE TABLE IF NOT EXISTS user_access_tokens (
  "token" TEXT NOT NULL PRIMARY KEY,
  "user_id" TEXT NOT NULL,
  "timestamp" DATE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS auth (
    "id" INTEGER NOT NULL PRIMARY KEY,
		"user_id" TEXT NOT NULL,
    "token" TEXT NOT NULL,
    "type" TEXT NOT NULL,
		"timestamp" DATE DEFAULT CURRENT_TIMESTAMP NOT NULL
	);
  CREATE INDEX auth_token ON auth (token);
