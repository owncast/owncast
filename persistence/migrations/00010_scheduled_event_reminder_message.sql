-- +goose Up
-- +goose StatementBegin

ALTER TABLE stream_event_series ADD COLUMN reminder_message TEXT NOT NULL DEFAULT '';
ALTER TABLE stream_events ADD COLUMN reminder_message TEXT NOT NULL DEFAULT '';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE stream_events DROP COLUMN reminder_message;
ALTER TABLE stream_event_series DROP COLUMN reminder_message;

-- +goose StatementEnd
