-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN disabled_reason TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN disabled_reason;
-- +goose StatementEnd
