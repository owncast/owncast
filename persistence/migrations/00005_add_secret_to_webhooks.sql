-- +goose Up
-- +goose StatementBegin
ALTER TABLE webhooks ADD COLUMN secret TEXT DEFAULT 'abc123' NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE webhooks DROP COLUMN secret;
-- +goose StatementEnd