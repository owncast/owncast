-- +goose Up
-- +goose StatementBegin

-- Linked-identity fields for the auth table.
--
-- An auth row is one external identity a user proved they own (IndieAuth,
-- Fediverse, or a viewer-auth plugin). Until now a row carried only the auth
-- key (token) and a coarse type. To surface a verified identity publicly we
-- need more: which provider it is (first-class, not parsed from a token
-- prefix), where it links, a human label, and whether the user consents to
-- showing it.
--
--   token -> auth_key  the value matched at login. For plugin auth this is now
--                      the raw external id, with no "<slug>:" prefix.
--   provider           first-class provider id. For built-ins it equals the
--                      type; for plugins it is the plugin slug. `type` is kept
--                      as the security namespace, so a plugin whose slug is
--                      "fediverse" still lands in type='plugin.auth' and can
--                      never collide with the built-in fediverse provider.
--   profile_url        public, clickable link. Optional.
--   handle             human label, e.g. @me@host. Optional.
--   is_public          consent. Off by default; the user opts an identity in.
ALTER TABLE auth RENAME COLUMN token TO auth_key;
ALTER TABLE auth ADD COLUMN provider TEXT NOT NULL DEFAULT '';
ALTER TABLE auth ADD COLUMN profile_url TEXT;
ALTER TABLE auth ADD COLUMN handle TEXT;
ALTER TABLE auth ADD COLUMN is_public BOOLEAN NOT NULL DEFAULT FALSE;

-- Backfill existing rows. Only IndieAuth and Fediverse identities exist today
-- (plugin auth has not shipped), so provider == type for every current row.
UPDATE auth SET provider = type;
-- IndieAuth's key is the user's website, which already is the public profile.
UPDATE auth SET profile_url = auth_key WHERE type = 'indieauth';
-- Fediverse's key is the @me@host handle, which is the human label.
UPDATE auth SET handle = auth_key WHERE type = 'fediverse';

DROP INDEX IF EXISTS idx_auth_token;
CREATE INDEX IF NOT EXISTS idx_auth_auth_key ON auth (auth_key);

-- A plugin identity (provider=slug, auth_key=external id) must map to exactly
-- one user, so a concurrent double-registration can't mint two users for the
-- same external identity. Partial to type='plugin.auth': it is the formalized
-- path and has no existing rows (plugin auth has not shipped), so the unique
-- index can never fail to build on legacy data. Built-in IndieAuth/Fediverse
-- linking is guarded by a find-first check in their handlers.
CREATE UNIQUE INDEX IF NOT EXISTS idx_auth_plugin_identity ON auth (provider, auth_key) WHERE type = 'plugin.auth';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_auth_plugin_identity;
DROP INDEX IF EXISTS idx_auth_auth_key;
ALTER TABLE auth DROP COLUMN is_public;
ALTER TABLE auth DROP COLUMN handle;
ALTER TABLE auth DROP COLUMN profile_url;
ALTER TABLE auth DROP COLUMN provider;
ALTER TABLE auth RENAME COLUMN auth_key TO token;
CREATE INDEX IF NOT EXISTS idx_auth_token ON auth (token);
-- +goose StatementEnd
