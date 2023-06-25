# SQL Storage

This package contains the base SQL schema, migrations and queries. It should not need to be imported by any package other than the datstore.

## SQL Migrations

Add migrations to `migrations.go` and use raw SQL make your required changes between schema versions.

## SQL Queries

_sqlc_ generates **type-safe code** from SQL. Here's how it works:

1. You define the schema in `schema.sql`.
1. You write your queries in `query.sql` using regular SQL.
1. You run `sqlc generate` to generate Go code with type-safe interfaces to those queries.
1. You write application code that calls the generated code.

Only those who need to create or update SQL queries will need to have `sqlc` installed on their system. **It is not a dependency required to build the codebase.**

## Install sqlc

### Snap

`sudo snap install sqlc`

### Go install

`go install github.com/kyleconroy/sqlc/cmd/sqlc@latest`

### macOS

`brew install sqlc`

### Download a release

Visit <https://github.com/kyleconroy/sqlc/releases> to download a release for your environment.
