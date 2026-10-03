-- +goose Up
-- +goose StatementBegin

-- Add a canonical operator-set priority to each federated server.
-- Lower values appear first (priority 1 = highest priority recommendation).
-- Pending and rejected servers keep their default of 0; priority is only
-- assigned when a server transitions to "accepted" so pending rows do not
-- consume positions.
ALTER TABLE federated_servers ADD COLUMN priority INTEGER NOT NULL DEFAULT 0;

-- +goose StatementEnd

-- +goose StatementBegin

-- Assign an initial priority to already-accepted servers based on insertion
-- order (id) so they have a deterministic starting order. Operators can
-- reorder using the move-up/move-down admin controls once the server is running.
UPDATE federated_servers
SET priority = (
    SELECT COUNT(*) + 1
    FROM federated_servers AS fs2
    WHERE fs2.follow_status = 'accepted' AND fs2.id < federated_servers.id
)
WHERE follow_status = 'accepted';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE federated_servers DROP COLUMN priority;
-- +goose StatementEnd
