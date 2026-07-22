-- +goose Up
-- +goose StatementBegin
CREATE TABLE ap_delivery_queue (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    inbox TEXT NOT NULL,
    payload BLOB NOT NULL,
    actor_iri TEXT NOT NULL,
    activity_type TEXT NOT NULL,
    coalesce_key TEXT,
    attempts INTEGER NOT NULL DEFAULT 0,
    next_attempt_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    claimed_until TIMESTAMP,
    last_error TEXT,
    failed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revision INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX idx_ap_delivery_queue_due ON ap_delivery_queue (failed_at, next_attempt_at, claimed_until);
CREATE UNIQUE INDEX idx_ap_delivery_queue_coalesce ON ap_delivery_queue (inbox, coalesce_key) WHERE coalesce_key IS NOT NULL AND failed_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ap_delivery_queue;
-- +goose StatementEnd
