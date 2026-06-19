-- +goose Up
-- +goose StatementBegin
-- Rename the follower flag from owncast_server to directory.
--
-- The flag marks followers that asked to list this server's stream in their
-- directory (set when their Follow carries the ns#directory marker). The old
-- name conflated "is an Owncast peer" with "is a directory"; the flag has only
-- ever meant the latter, so name it for what it is.
ALTER TABLE ap_followers RENAME COLUMN owncast_server TO directory;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE ap_followers RENAME COLUMN directory TO owncast_server;
-- +goose StatementEnd
