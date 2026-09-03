-- +goose Up
-- +goose StatementBegin

ALTER TABLE stream_events ADD COLUMN federation_reminder_1_sent_at TIMESTAMP;
ALTER TABLE stream_events ADD COLUMN federation_reminder_2_sent_at TIMESTAMP;
ALTER TABLE stream_events ADD COLUMN federation_update_pending BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE stream_events ADD COLUMN federation_version INTEGER NOT NULL DEFAULT 0;
ALTER TABLE ap_delivery_queue ADD COLUMN ordering_key TEXT;
ALTER TABLE ap_delivery_queue ADD COLUMN coalesce_version INTEGER NOT NULL DEFAULT 0;
ALTER TABLE ap_delivery_queue ADD COLUMN blocks_following BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE ap_outbox ADD COLUMN coalesce_version INTEGER NOT NULL DEFAULT 0;
CREATE INDEX idx_ap_delivery_queue_ordering ON ap_delivery_queue(inbox, ordering_key, id);
CREATE TABLE stream_event_federation_deletes (
    event_id TEXT NOT NULL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE stream_event_federation_deletes;
DROP INDEX idx_ap_delivery_queue_ordering;
ALTER TABLE ap_delivery_queue DROP COLUMN blocks_following;
ALTER TABLE ap_outbox DROP COLUMN coalesce_version;
ALTER TABLE ap_delivery_queue DROP COLUMN coalesce_version;
ALTER TABLE ap_delivery_queue DROP COLUMN ordering_key;
ALTER TABLE stream_events DROP COLUMN federation_version;
ALTER TABLE stream_events DROP COLUMN federation_update_pending;
ALTER TABLE stream_events DROP COLUMN federation_reminder_2_sent_at;
ALTER TABLE stream_events DROP COLUMN federation_reminder_1_sent_at;

-- +goose StatementEnd
