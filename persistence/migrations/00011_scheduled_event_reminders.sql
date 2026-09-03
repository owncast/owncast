-- +goose Up
-- +goose StatementBegin

ALTER TABLE stream_events ADD COLUMN reminder_1_sent_at TIMESTAMP;
ALTER TABLE stream_events ADD COLUMN reminder_2_sent_at TIMESTAMP;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE stream_events DROP COLUMN reminder_2_sent_at;
ALTER TABLE stream_events DROP COLUMN reminder_1_sent_at;

-- +goose StatementEnd
