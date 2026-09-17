# budgetcli — Planned Install Guide (Work in Progress)

**This describes the target setup experience once `install.go` and its related pieces are built. None of the install/config steps below exist yet — see `README.md` for the current, working setup.** This document is written in the same "explain everything" style as the real README will be once this is implemented, so it can be swapped in directly when ready.

---

## What you'll need before starting

1. **Go** — https://go.dev/dl/ — needed once, to build the program.
2. **Postgres** — see options below.
3. **psql** — see options below.

(Same as the current README — these steps don't change.)

### Installing Go

Download and install from https://go.dev/dl/. Confirm with:
```bash
go version
```

### Installing Postgres

#### Option A: Docker (recommended)

Install Docker: https://docs.docker.com/get-docker/. Confirm with:
```bash
docker --version
docker compose version
```

#### Option B: Native Postgres install

Follow the official instructions for your OS: https://www.postgresql.org/download/. Create a database named `budgetcli` and a user for it, and note the username and password. *(If you'd prefer a different database name, or want to use an existing database, see the note in Step 4 below.)*

### Installing psql

Included automatically with a native Postgres install. With Docker, either install a client separately from the link above, or run it from inside the container (shown in Step 5 below).

---

## Step 1: Get the project

```bash
git clone <repo-url>
cd budgetcli
```

## Step 2: Build and install

From the project root:
```bash
cd install
go run install.go
```

This does three things:
1. Builds the `budgetcli` program.
2. Moves it to `~/.local/bin`, a folder your computer already checks when you type a command — so afterward you can just type `budgetcli` from anywhere, no matter what folder you're in.
3. Copies a Docker setup file into `~/.config/budgetcli/`, so it's available even after this project folder is gone.

**If `budgetcli` isn't recognized as a command afterward (Linux/macOS):** your computer doesn't yet know to look in `~/.local/bin`. Add this line to your shell's configuration file (`~/.zshrc` for zsh, `~/.bashrc` for bash):
```bash
export PATH="$HOME/.local/bin:$PATH"
```
Then either restart your terminal, or run `source ~/.zshrc` (or `~/.bashrc`) to apply it immediately.

**Windows:** the installer builds `budgetcli.exe` but does not move it automatically — you'll need to move it yourself to a folder already in your `PATH`.

## Step 3: First-time setup

The first time you run `budgetcli`, it will notice no settings exist yet and walk you through a short setup:

- Whether you're using Docker or a native Postgres install.
- Your Postgres host, port, username, and password. *(The database name is not asked for — `budgetcli` assumes a database named `budgetcli` by default. See the note under Step 4 if you're using a different name.)*

These are saved to a settings file at `~/.config/budgetcli/config.toml` (Windows: an equivalent location in your user folder), so you're never asked again unless you delete that file.

```bash
budgetcli
```
Just follow the prompts.

*(This step doesn't exist yet — right now, connection details still come from a `.env` file inside the project folder. See `README.md`.)*

## Step 4: Starting Postgres

**If you chose Docker during setup:** `budgetcli` starts the Postgres container for you automatically every time it runs — you don't need to do this yourself. The menu appears immediately; there's no need to wait for Postgres to finish starting before choosing what you want to do.

**If you'd like to use a different database name, or an existing database,** open the Docker setup file at `~/.config/budgetcli/docker-compose.yml` and change the `POSTGRES_DB` value before your first run, or point your native Postgres setup at the name of your choosing during Step 3.

**If you chose a native Postgres install during setup:** make sure your Postgres service is running (this varies by OS — see your install's documentation). `budgetcli` will not attempt to start anything for you in this case.

## Step 5: Create the database tables

This still needs to be run once, the same way as today:
```bash
psql "postgres://<user>:<password>@<host>:<port>/budgetcli" -f budgetcli_schema.sql
```
(Replace `budgetcli` in the path above with your own database name if you changed it in Step 4.)

*(Open question, not yet decided: whether this schema file also gets copied into `~/.config/budgetcli/` during install, or whether it's automatically applied by `budgetcli` itself on first run if the tables don't already exist.)*

## Step 6: You're done — delete this folder

Once setup is complete, everything `budgetcli` needs lives outside this project folder:
- The program itself: `~/.local/bin/budgetcli`
- Your settings: `~/.config/budgetcli/config.toml`
- The Docker setup file: `~/.config/budgetcli/docker-compose.yml`
- Your actual bill data (if using Docker): `~/.local/share/budgetcli/pgdata` — this stays in place even if the Docker container is ever removed or rebuilt.

This project folder can be safely deleted.

## Usage

Run `budgetcli` with no arguments to open the interactive menu. Individual commands (`add`, `prep`, `list`, `update`) remain available to run directly as well — they connect to the database the same way the menu does, and will retry briefly and give a clear message if Postgres hasn't finished starting yet.

---

## Still to be designed / built

- [ ] The first-run interactive setup wizard itself (Step 3 above)
- [ ] Final decision on config file format/location (`config.toml` vs. alternatives)
- [ ] Updating the database connection code to read from this config file instead of `.env`
- [ ] Updating `install.go` to copy the Compose file (and possibly the schema file) into `~/.config/budgetcli/`
- [ ] Automatic Postgres startup for Docker mode, plus the shared connection-retry helper used by every command (see `config-toml-migration-plan.md` section 8 for the technical design)
- [ ] Cross-platform config folder handling for Windows (likely via Go's `os.UserConfigDir()`)

Once all of the above is built and tested, this document replaces the relevant sections of `README.md`, and this file goes away.
