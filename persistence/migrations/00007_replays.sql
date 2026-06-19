-- +goose Up
-- +goose StatementBegin
-- Record the high level details of each stream so it can be replayed later.
CREATE TABLE IF NOT EXISTS streams (
	"id" TEXT NOT NULL PRIMARY KEY,
	"stream_title" TEXT,
	"start_time" DATETIME,
	"end_time" DATETIME
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS streams_start_time ON streams (start_time);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS streams_start_end_time ON streams (start_time, end_time);
-- +goose StatementEnd

-- +goose StatementBegin
-- Record the output (variant) configuration that was used for a stream so the
-- replay playlists can be regenerated with matching parameters.
CREATE TABLE IF NOT EXISTS video_segment_output_configuration (
	"id" TEXT NOT NULL PRIMARY KEY,
	"variant_id" TEXT NOT NULL,
	"name" TEXT NOT NULL,
	"stream_id" TEXT NOT NULL,
	"segment_duration" INTEGER NOT NULL,
	"bitrate" INTEGER NOT NULL,
	"framerate" INTEGER NOT NULL,
	"resolution_width" INTEGER,
	"resolution_height" INTEGER,
	"timestamp" DATETIME
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS video_segment_output_configuration_stream_id ON video_segment_output_configuration (stream_id);
-- +goose StatementEnd

-- +goose StatementBegin
-- Every HLS segment written for a recorded stream, so a stream (or a window
-- of it, for clips) can be reconstructed into a playlist.
CREATE TABLE IF NOT EXISTS video_segments (
	"id" TEXT NOT NULL PRIMARY KEY,
	"stream_id" TEXT NOT NULL,
	"output_configuration_id" TEXT NOT NULL,
	"path" TEXT NOT NULL,
	"relative_timestamp" REAL NOT NULL,
	"timestamp" DATETIME
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS video_segments_stream_id ON video_segments (stream_id);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS video_segments_stream_id_timestamp ON video_segments (stream_id, timestamp);
-- +goose StatementEnd

-- +goose StatementBegin
-- A replayable clip: a named window of a recorded stream.
CREATE TABLE IF NOT EXISTS replay_clips (
	"id" TEXT NOT NULL PRIMARY KEY,
	"stream_id" TEXT NOT NULL,
	"clipped_by" TEXT,
	"clip_title" TEXT,
	"relative_start_time" REAL,
	"relative_end_time" REAL,
	"timestamp" DATETIME,
	FOREIGN KEY(stream_id) REFERENCES streams(id)
);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS clip_stream_id ON replay_clips (stream_id);
-- +goose StatementEnd
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS clip_start_end_time ON replay_clips (relative_start_time, relative_end_time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS replay_clips;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS video_segments;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS video_segment_output_configuration;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS streams;
-- +goose StatementEnd
