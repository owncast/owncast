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

CREATE TABLE IF NOT EXISTS messages (
		"id" string NOT NULL,
		"user_id" TEXT,
		"body" TEXT,
		"eventType" TEXT,
		"hidden_at" DATE,
		"timestamp" DATE,
    "title" TEXT,
    "subtitle" TEXT,
    "image" TEXT,
    "link" TEXT,
		PRIMARY KEY (id)
	);CREATE INDEX index ON messages (id, user_id, hidden_at, timestamp);
	CREATE INDEX id ON messages (id);
	CREATE INDEX user_id ON messages (user_id);
	CREATE INDEX hidden_at ON messages (hidden_at);
	CREATE INDEX timestamp ON messages (timestamp);

-- Record the high level details of each stream.
CREATE TABLE IF NOT EXISTS streams (
	"id" string NOT NULL PRIMARY KEY,
	"stream_title" TEXT,
	"start_time" DATE,
	"end_time" DATE,
	PRIMARY KEY (id)
);
CREATE INDEX streams_id ON streams (id);
CREATE INDEX streams_start_time ON streams (start_time);
CREATE INDEX streams_start_end_time ON streams (start_time,end_time);

-- Record the output configuration of a stream.
CREATE TABLE IF NOT EXISTS video_segment_output_configuration (
	"id" string NOT NULL PRIMARY KEY,
	"variant_id" string NOT NULL,
	"name" string NOT NULL,
	"stream_id" string NOT NULL,
	"segment_duration" INTEGER NOT NULL,
	"bitrate" INTEGER NOT NULL,
	"framerate" INTEGER NOT NULL,
	"resolution_width" INTEGER,
	"resolution_height" INTEGER,
	"timestamp" DATE,
	PRIMARY KEY (id)
);
	CREATE INDEX video_segment_output_configuration_stream_id ON video_segment_output_configuration (stream_id);

-- Support querying all segments for a single stream as well
-- as segments for a time window.
CREATE TABLE IF NOT EXISTS video_segments (
	"id" string NOT NULL PRIMARY KEY,
	"stream_id" string NOT NULL,
	"output_configuration_id" string NOT NULL,
	"path" TEXT NOT NULL,
	"relative_timestamp" REAL NOT NULL,
	"timestamp" DATE,
	PRIMARY KEY (id)
);
	CREATE INDEX video_segments_stream_id ON video_segments (stream_id);
	CREATE INDEX video_segments_stream_id_timestamp ON video_segments (stream_id,timestamp);

-- Record the details of a replayable clip.
CREATE TABLE IF NOT EXISTS replay_clips (
	"id" string NOT NULL,
	"stream_id" string NOT NULL,
	"clipped_by" string,
	"clip_title" TEXT,
	"relative_start_time" REAL,
	"relative_end_time" REAL,
	"timestamp" DATE,
	PRIMARY KEY (id),
  FOREIGN KEY(stream_id) REFERENCES streams(id)
);
CREATE INDEX clip_id ON replay_clips (id);
CREATE INDEX clip_stream_id ON replay_clips (stream_id);
CREATE INDEX clip_start_end_time ON replay_clips (start_time,end_time);
