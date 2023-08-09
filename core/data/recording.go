package data

import "database/sql"

func createRecordingTables(db *sql.DB) {
	createSegmentsTableSQL := `CREATE TABLE IF NOT EXISTS video_segments (
		"id" string NOT NULL,
		"stream_id" string NOT NULL,
		"output_configuration_id" string NOT NULL,
		"path" TEXT NOT NULL,
		"timestamp" DATE DEFAULT CURRENT_TIMESTAMP NOT NULL,
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
		"timestamp" DATE DEFAULT CURRENT_TIMESTAMP NOT NULL,
		PRIMARY KEY (id)
	);CREATE INDEX video_segment_output_configuration_stream_id ON video_segment_output_configuration (stream_id);`

	createVideoStreamsTableSQL := `CREATE TABLE IF NOT EXISTS streams (
		"id" string NOT NULL,
		"stream_title" TEXT,
		"start_time" DATE NOT NULL,
		"end_time" DATE,
		PRIMARY KEY (id)
	);
	CREATE INDEX streams_id ON streams (id);
	CREATE INDEX streams_start_time ON streams (start_time);
	CREATE INDEX streams_start_end_time ON streams (start_time,end_time);
	`

	MustExec(createSegmentsTableSQL, db)
	MustExec(createVideoOutputConfigsTableSQL, db)
	MustExec(createVideoStreamsTableSQL, db)
}
