# SQL Queries

sqlc generates **type-safe code** from SQL. Here's how it works:

1. You define the schema in `schema.sql`.
1. You write your queries in `query.sql` using regular SQL.
1. You run `sqlc generate` to generate Go code with type-safe interfaces to those queries.
1. You write application code that calls the generated code.

Only those who need to create or update SQL queries will need to have `sqlc` installed on their system. **It is not a dependency required to build the codebase.**

## Using sqlc

Run from the repository root:

```bash
make sqlc
```

Or manually:

```bash
./bin/sqlc generate
```

## Managing Tool Versions

To upgrade sqlc to a specific version, edit `tools/go.mod` and run:

```bash
cd tools && go get github.com/sqlc-dev/sqlc@latest && go mod tidy
```

Then reinstall:

```bash
rm ./bin/sqlc && make sqlc
```
