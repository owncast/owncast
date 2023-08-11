package data

import "database/sql"

func createRecordingTables(db *sql.DB) {
	createSegmentsTableSQL := `CREATE TABLE IF NOT EXISTS video_segments (
		"id" string NOT NULL,
		"stream_id" string NOT NULL,
		"output_configuration_id" string NOT NULL,
		"path" TEXT NOT NULL,
		"relative_timestamp" REAL NOT NULL,
		"timestamp" DATETIME,
		PRIMARY KEY (id)
	);CREATE INDEX video_segments_stream_id ON video_segments (stream_id);CREATE INDEX video_segments_stream_id_timestamp ON video_segments (stream_id,timestamp);`

	createVideoOutputConfigsTableSQL := `CREATE TABLE IF NOT EXISTS video_segment_output_configuration (
		"id" string NOT NULL,
		"variant_id" string NOT NULL,
		"name" string NOT NULL,
		"stream_id" string NOT NULL,
		"segment_duration" INTEGER NOT NULL,
		"bitrate" INTEGER NOT NULL,
		"framerate" INTEGER NOT NULL,
		"resolution_width" INTEGER,
		"resolution_height" INTEGER,
		"timestamp" DATETIME,
		PRIMARY KEY (id)
	);CREATE INDEX video_segment_output_configuration_stream_id ON video_segment_output_configuration (stream_id);`

	createVideoStreamsTableSQL := `CREATE TABLE IF NOT EXISTS streams (
		"id" string NOT NULL,
		"stream_title" TEXT,
		"start_time" DATETIME,
		"end_time" DATETIME,
		PRIMARY KEY (id)
	);
	CREATE INDEX streams_id ON streams (id);
	CREATE INDEX streams_start_time ON streams (start_time);
	CREATE INDEX streams_start_end_time ON streams (start_time,end_time);
	`

	createClipsTableSQL := `CREATE TABLE IF NOT EXISTS replay_clips (
		"id" string NOT NULL,
		"stream_id" string NOT NULL,
		"clipped_by" string,
		"clip_title" TEXT,
		"relative_start_time" REAL,
		"relative_end_time" REAL,
		"timestamp" DATETIME,
		PRIMARY KEY (id),
		FOREIGN KEY(stream_id) REFERENCES streams(id)
	);
	CREATE INDEX clip_id ON replay_clips (id);
	CREATE INDEX clip_stream_id ON replay_clips (stream_id);
	CREATE INDEX clip_start_end_time ON replay_clips (start_time,end_time);
	`

	MustExec(createSegmentsTableSQL, db)
	MustExec(createVideoOutputConfigsTableSQL, db)
	MustExec(createVideoStreamsTableSQL, db)
	MustExec(createClipsTableSQL, db)
}
