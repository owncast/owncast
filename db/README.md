# SQL Queries

sqlc generates **type-safe code** from SQL. Here's how it works:

1. You define the schema in `schema.sql`.
1. You write your queries in `query.sql` using regular SQL.
1. You run `sqlc generate` to generate Go code with type-safe interfaces to those queries.
1. You write application code that calls the generated code.

Only those who need to create or update SQL queries will need to have `sqlc` installed on their system. **It is not a dependency required to build the codebase.**

## Using sqlc

This project uses Go 1.24+'s native tool management. The `sqlc` tool is already defined in `go.mod`.

After cloning the repository, ensure tools are downloaded:

```bash
go mod download
```

Then run sqlc:

```bash
go tool sqlc generate
```

## Managing Tool Versions

To upgrade sqlc to a specific version:

```bash
go get -tool github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

To check the current version:

```bash
go list -m github.com/sqlc-dev/sqlc
```
