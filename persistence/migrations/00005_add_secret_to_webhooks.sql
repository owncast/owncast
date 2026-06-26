-- +goose Up
ALTER TABLE webhooks ADD COLUMN secret TEXT DEFAULT 'abc123' NOT NULL;

-- +goose Down
ALTER TABLE webhooks DROP COLUMN secret;