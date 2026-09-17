# Migration Plan: .env → config.toml + install.go

Working document. Describes what needs to be built and in what order. Nothing described here exists yet — no `install.go`, no config.toml support. Today's actual working setup is `.env`, documented in `README.md`.

---

## 1. Config format and location

- Move from `.env` (key=value, read via `os.Getenv`) to a `config.toml` file.
- Location: a per-user config directory, resolved via Go's `os.UserConfigDir()` — this returns the correct OS-specific path automatically (`~/.config` on Linux, `~/Library/Application Support` on macOS, `%AppData%` on Windows), so no manual OS-detection is needed.
- Final path would be something like `<UserConfigDir>/budgetcli/config.toml`.
- Package to review: a TOML parsing library, since Go's standard library has no built-in TOML support (unlike JSON). `github.com/BurntSushi/toml` is the long-standing, widely used option — worth reading its README for basic struct-tag usage (mapping TOML keys to Go struct fields, similar in spirit to how `pgx` maps SQL columns to struct fields via `Scan`).

## 2. Changes needed in `internal/database`

- `Connect()` currently builds its connection string from `os.Getenv(...)` calls — this is the real, existing code today.
- New responsibility to add: check whether `config.toml` exists at the resolved path.
  - If it exists, read and parse it into a struct (host/port/user/password/dbname fields), build the connection string from that.
  - If it doesn't exist, this is the trigger for first-run setup (see section 3) — `Connect()` itself likely shouldn't own the prompting logic; better to have `main.go` or a dedicated setup step check for the config's existence *before* ever calling `Connect()`, and only call `Connect()` once a config is known to exist.
- Decide whether `.env` support is dropped entirely or kept as a fallback during a transition period. Recommend dropping it once config.toml is confirmed working, to avoid maintaining two parallel paths.

## 3. First-run setup wizard (does not exist yet)

- New logic, likely its own function in a new file (e.g. `cmd/setup.go`) — not part of any existing command, since it needs to run *before* any command that touches the database.
- Trigger condition: config file doesn't exist at the resolved path.
- Uses `huh` (already a dependency) to prompt for host/port/user/password/dbname — same shape as any existing form in this project.
- On completion: create the config directory if it doesn't exist (`os.MkdirAll`), then write the collected values out as TOML (via the chosen TOML library's encoding side) to the resolved path.
- Where this check would happen: most natural spot is early in `main.go`, before `cmd.Execute()` is called — check for config, run setup if missing, then proceed as normal either way.

## 4. A future `install.go` — proposed order of operations (nothing here exists yet)

This would be a new, standalone helper script, separate from the main `budgetcli` program (its own file/folder, its own `package main`, run manually with `go run install.go` — not something that gets compiled into `budgetcli` itself). Proposed steps, in sequence:

1. Build the `budgetcli` binary (`go build`).
2. Determine a destination directory for the binary — a location already in most users' `PATH` without needing admin/sudo rights, such as `~/.local/bin`.
3. Determine the config directory via `os.UserConfigDir()` + `/budgetcli`.
4. Create that config directory if it doesn't exist (`os.MkdirAll`).
5. Copy `docker-compose.prod.yml` from the repo into that config directory.
   - Copying a file isn't a single stdlib call — needs an explicit read-source/write-destination step, or the `io.Copy` pattern (worth reviewing `io` and `os.Open`/`os.Create` together for this).
6. (Open question, not yet decided) — whether `budgetcli_schema.sql` also gets copied here, for a possible future "auto-create tables if missing" feature, or stays a repo-only, manually-run file.
7. Move the built binary into the destination directory from step 2.
8. Print final instructions for the user (PATH reminder if needed, note that first run will prompt for setup).

## 5. Packages worth reviewing before writing any of this

- `github.com/BurntSushi/toml` (or equivalent) — TOML read/write.
- `os.UserConfigDir()` — stdlib, cross-platform config path resolution.
- `os.MkdirAll` — stdlib, already familiar from other parts of this project's design discussions.
- `io.Copy` alongside `os.Open` / `os.Create` — for the file-copy step described in section 4.
- No new third-party packages needed beyond the TOML library — everything else is stdlib.

## 6. Still undecided (blockers before implementation)

- Where `pgdata` (the actual Postgres data files) lives once relocated out of the repo/devcontainer structure — affects what path gets written into the Compose file's volume line.
- Whether `.env` is removed entirely once this ships, or kept temporarily as a fallback.
- Whether the schema file gets copied alongside the Compose file, or stays repo-only.

## 7. Suggested build order (once ready to implement)

1. Settle the `pgdata` location decision (blocks the Compose file's final contents).
2. Add TOML read/write support and the config struct.
3. Update `Connect()` to read from config.toml.
4. Build the first-run setup wizard.
5. Build `install.go` for the first time, with the config-directory and Compose-file-copy steps included from the start.
6. Update `README.md` (merge in the relevant parts of `README-wip.md`), then retire `README-wip.md`.

---

## 8. Startup / connection-readiness design## 8. Startup / connection-readiness design (settled)

User-facing behavior described in `README-wip.md` (Step 4, Usage). Technical summary for implementation:

- `pgxpool.New(...)` stays in `main.go`, unchanged — confirmed lazy (pgx v5 docs: does not connect until first use), so no startup delay regardless of Postgres readiness.
- On `mode == "docker"` (read from `config.toml`), kick off `docker compose -f <path> up -d` via `os/exec` at startup, don't wait for or check its result.
- Build one shared helper (e.g. `queryWithRetry` / `execWithRetry`) wrapping `DB.Query`/`DB.Exec`: on a failed first attempt, retry a small fixed number of times with a short pause, then surface a clear error if still failing. Every command uses this instead of calling `DB.Query`/`DB.Exec` directly.
- Commands remain independently callable outside the menu (`budgetcli prep`, etc.) — the shared retry helper is what makes this safe.

## 9. Database name 9. Database name — hardcoded default (settled)

User-facing behavior described in `README-wip.md` (Installing Postgres, Step 4). Technical summary for implementation:

- First-run wizard does not prompt for `dbname` — assume `budgetcli` by default.
- `docker-compose.prod.yml`'s `POSTGRES_DB: ${DB_NAME}` becomes a hardcoded `POSTGRES_DB: budgetcli`.
- Override path (different name / existing database) is a manual edit to the Compose file or native setup, not a wizard prompt — documented in `README-wip.md`.
