# Database

Two pieces work together here:

- **Schema** lives in `persistence/migrations/` as numbered [goose](https://github.com/pressly/goose) SQL migrations. These are the single source of truth for what the tables look like. They run automatically on startup.
- **Queries** live in `db/query.sql`. [sqlc](https://sqlc.dev) reads them (plus the schema from `persistence/migrations/`, see `sqlc.yaml`) and generates the type-safe Go in `db/query.sql.go` and `db/models.go`.

There is no `schema.sql`. sqlc derives the schema from the migrations.

## Adding or changing a column (or any schema change)

1. **Write a migration.** Create the next numbered file in `persistence/migrations/`, e.g. `00005_add_widget_color.sql`. Copy the structure of an existing one. Put your `ALTER TABLE` / `CREATE TABLE` between the goose markers, and write the reverse in the `Down` section:

   ```sql
   -- +goose Up
   -- +goose StatementBegin
   ALTER TABLE users ADD COLUMN widget_color TEXT NOT NULL DEFAULT '';
   -- +goose StatementEnd

   -- +goose Down
   -- +goose StatementBegin
   ALTER TABLE users DROP COLUMN widget_color;
   -- +goose StatementEnd
   ```

   Rules: never edit a migration that has already shipped (add a new one instead), keep the numbers sequential, and make statements idempotent where practical (`IF NOT EXISTS`).

2. **Update queries** in `db/query.sql` if you need to read or write the new column.

3. **Regenerate** the Go code:

   ```bash
   make sqlc
   ```

4. **Build.** `go build ./...` — the migration applies on next startup; no manual step.

Do not hand-write raw SQL in Go for new work. Add it to `db/query.sql` and regenerate.

### Tables that are not yet sqlc-managed

Some older repositories still use hand-written SQL and are **not** wired into sqlc — for example `persistence/webhookrepository/`. For those tables, `make sqlc` generates nothing and there are no generated files to edit. The schema change still goes in a goose migration, but you then edit the raw SQL by hand in that repository (the `INSERT`, the `SELECT`, and the matching struct fields).

Watch for `SELECT *` with a positional `rows.Scan(...)`: adding a column changes how many values the row returns, so you must add the new column to the `Scan` argument list (in column order — a column added via `ALTER TABLE` lands last) or the scan fails at runtime. `GetWebhooks()` in the webhook repository is exactly this shape.

## sqlc

Only contributors who change SQL need sqlc installed; it is **not** required to build the project. `make sqlc` installs it into `./bin` automatically.

To upgrade sqlc, edit `tools/go.mod`, then:

```bash
cd tools && go get github.com/sqlc-dev/sqlc@latest && go mod tidy
rm ../bin/sqlc && cd .. && make sqlc
```

## Legacy installs

Installs that predate goose tracked schema state in a `config.version` row advanced by a switch statement in `persistence/legacymigrations`. That package is **frozen** — no new cases. On startup, old installs are caught up to the goose baseline and then goose takes over. All new schema work goes in `persistence/migrations/` only. See `persistence/migrations/migrations.go`.
