# Config Repository

The configuration repository represents all the getters, setters and storage logic for user-defined configuration values. This includes things like the server name, enabled/disabled flags, etc. See `keys.go` to see the full list of keys that are used for accessing these values in the database.

## Migrations

Add migrations to `migrations.go` and you can use the datastore and config repository to make your required changes between datastore versions.
