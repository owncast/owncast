# Automated upgrade test

Proves that a data directory created by the latest published Owncast release survives an upgrade to the code in this checkout.

`run.sh`:

1. Downloads the latest GitHub release asset for the current platform (Linux amd64/arm64 and macOS Intel/Apple Silicon are supported).
2. Runs it in a scratch directory on non-default ports (web 8098, RTMP 1998) and writes recognizable config values through the admin API.
3. Stops it, builds this checkout, and starts the new binary on the same data directory.
4. Asserts the upgraded server becomes ready, reports the dev version, kept the previously written config, and migrated the database to the newest goose schema version.

Any failed assertion exits nonzero. Server logs are left in `scratch/` for inspection.

Requires `curl`, `jq`, `unzip`, `sqlite3`, `go`, and `ffmpeg`. Runs in CI via `.github/workflows/upgrade-tests.yml` (weekly, on demand, and on PRs touching this suite).

This intentionally tests a single hop: latest release to local HEAD, which is the upgrade every user performs when a release ships. It does not replay the full historical release chain.
