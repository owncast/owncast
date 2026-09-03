-- +goose Up
-- +goose StatementBegin

ALTER TABLE stream_events ADD COLUMN webhook_warning_sent_at TIMESTAMP;
ALTER TABLE stream_events ADD COLUMN webhook_started_sent_at TIMESTAMP;
ALTER TABLE stream_events ADD COLUMN webhook_ended_sent_at TIMESTAMP;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE stream_events DROP COLUMN webhook_ended_sent_at;
ALTER TABLE stream_events DROP COLUMN webhook_started_sent_at;
ALTER TABLE stream_events DROP COLUMN webhook_warning_sent_at;

-- +goose StatementEnd
