-- Schema update to query.sql must be referenced in queries located in query.sql
-- and compiled into code with sqlc. Read README.md for details.

CREATE TABLE ap_followers (
		"iri" TEXT NOT NULL,
		"inbox" TEXT NOT NULL,
		"name" TEXT,
		"username" TEXT NOT NULL,
		"image" TEXT,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (iri));
		CREATE INDEX iri ON ap_followers (iri);

CREATE TABLE IF NOT EXISTS ap_outbox (
		"id" TEXT NOT NULL,
		"iri" TEXT NOT NULL,
		"value" BLOB,
		"type" TEXT NOT NULL,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id));
		CREATE INDEX id ON ap_outbox (id);
		CREATE INDEX iri ON ap_outbox (iri);
		CREATE INDEX type ON ap_outbox (type);
		