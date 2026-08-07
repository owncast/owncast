-- +goose Up
-- +goose StatementBegin

-- Scheduled streams: recurring series and their concrete occurrences.
--
-- A series row stores an iCalendar RFC 5545 recurrence value (DTSTART with
-- an IANA TZID plus an RRULE). The scheduler expands it in Go and
-- materializes concrete occurrence rows on a rolling horizon. One-off
-- events are occurrence rows with no series. Occurrence rows persist after
-- they pass: they carry the federation state needed to send ActivityPub
-- Update/Delete activities for announced events, and feed the viewer
-- calendar's past view.

CREATE TABLE IF NOT EXISTS stream_event_series (
    "id" TEXT NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL DEFAULT '',
    "recurrence" TEXT NOT NULL,
    "duration_minutes" INTEGER NOT NULL DEFAULT 60,
    "active" BOOLEAN NOT NULL DEFAULT TRUE,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS stream_events (
    "id" TEXT NOT NULL PRIMARY KEY,
    "series_id" TEXT REFERENCES stream_event_series(id),
    "original_start" TIMESTAMP,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL DEFAULT '',
    "start_time" TIMESTAMP NOT NULL,
    "duration_minutes" INTEGER NOT NULL DEFAULT 60,
    "timezone" TEXT NOT NULL DEFAULT 'UTC',
    "status" TEXT NOT NULL DEFAULT 'scheduled',
    "federated_at" TIMESTAMP,
    "reminder_sent_at" TIMESTAMP,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- original_start is the instant the recurrence rule produced, immutable even
-- if the occurrence is later moved. The unique index makes materialization an
-- idempotent INSERT OR IGNORE. SQLite treats NULLs as distinct here, so
-- one-off events (NULL series_id and original_start) never collide.
CREATE UNIQUE INDEX IF NOT EXISTS stream_events_series_original_start ON stream_events (series_id, original_start);
CREATE INDEX IF NOT EXISTS stream_events_start_time ON stream_events (start_time);
CREATE INDEX IF NOT EXISTS stream_events_federated_at ON stream_events (federated_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS stream_events_federated_at;
DROP INDEX IF EXISTS stream_events_start_time;
DROP INDEX IF EXISTS stream_events_series_original_start;
DROP TABLE IF EXISTS stream_events;
DROP TABLE IF EXISTS stream_event_series;
-- +goose StatementEnd
