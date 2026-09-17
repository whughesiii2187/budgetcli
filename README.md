# budgetcli

A terminal tool for tracking your bills, figuring out what's due each month, and marking things as paid. It stores everything in a Postgres database.

This guide assumes no prior experience setting any of this up. Follow the steps in order.

---

## What you'll need before starting

1. **Go** — the programming language this tool is written in. You need this to build the program from its source code.
2. **Postgres** — the database that stores your bills. You have two choices for how to run it, explained below.
3. **psql** — a command-line tool for talking to Postgres directly (used once, to set up the database tables).

### Installing Go

Download and install it from the official site: https://go.dev/dl/

Follow the installer instructions for your operating system. Once installed, confirm it worked by opening a terminal and running:
```bash
go version
```
You should see a version number printed back.

### Installing Postgres

You have two options. **If you're not sure which to pick, use Option A (Docker).** It's more self-contained and easier to remove later if you change your mind.

#### Option A: Postgres via Docker

Docker lets you run Postgres in an isolated container without installing it directly on your computer.

1. Install Docker: https://docs.docker.com/get-docker/ (this includes Docker Compose, which is also needed)
2. Follow their installer for your operating system.
3. Confirm it worked:
   ```bash
   docker --version
   docker compose version
   ```

You don't need to do anything else with Docker yet — the steps further down in this guide explain how it gets used.

#### Option B: Installing Postgres directly

If you'd rather not use Docker, install Postgres directly using your operating system's official instructions:
- Official downloads and guides: https://www.postgresql.org/download/

You'll need to create a database and a database user yourself as part of that install process. Make note of the username, password, and database name you choose — you'll need them shortly.

### Installing psql

`psql` is Postgres's command-line client. It's used to run the one-time setup command that creates this project's tables.

- If you installed Postgres directly (Option B above), `psql` is usually included automatically.
- If you're using Docker (Option A), you can either install `psql` separately on your own machine (see the Postgres download page above, which includes client-only installers), or run `psql` from inside the Docker container itself — this is explained in Step 3 below.

---

## Step 1: Get the project

If you have `git` installed:
```bash
git clone <repo-url>
cd budgetcli
```
Otherwise, download the project as a ZIP from wherever you got this repository, and extract it.

## Step 2: Set up your `.env` file

This project needs to know how to connect to your Postgres database — the address, the username, the password, and so on. Rather than typing those in every time, they're stored in a file called `.env`.

1. In the project folder, find the file named `.env.example`. Make a copy of it and rename the copy to `.env`.
   ```bash
   cp .env.example .env
   ```
2. Open `.env` in any text editor. You'll see something like:
   ```
   DB_HOST=
   DB_PORT=
   DB_USER=
   DB_PASSWORD=
   DB_NAME=
   ```
3. Fill in each value:
   - `DB_HOST` — `localhost` if using Docker or a local Postgres install
   - `DB_PORT` — `5432` unless you specifically changed it
   - `DB_USER` — the username you set up (Docker: whatever you choose in Step 3 below; native install: the one you created yourself)
   - `DB_PASSWORD` — the matching password
   - `DB_NAME` — the name of the database (Docker: whatever you choose in Step 3 below; native install: the one you created yourself)

**Do not share this file or upload it anywhere public** — it will contain your database password. It's already set up to be ignored by `git` so it won't accidentally get uploaded if you use version control.

## Step 3: Start Postgres

### If using Docker

A ready-made configuration file, `docker-compose.prod.yml`, is included in this project. It already reads the username, password, and database name from your `.env` file, so make sure Step 2 is done first.

Start it with:
```bash
docker compose -f docker-compose.prod.yml up -d
```

This will download Postgres (the first time only) and start it running in the background. To confirm it's running:
```bash
docker ps
```
You should see a container listed.

If you need to run `psql` and don't have it installed locally, you can run it from inside this same container:
```bash
docker compose -f docker-compose.prod.yml exec db psql -U <DB_USER> -d <DB_NAME>
```
(Replace `<DB_USER>` and `<DB_NAME>` with the values from your `.env` file. Type `\q` to exit when done.)

### If using a native Postgres install

Make sure the Postgres service is running (this varies by operating system — check the documentation from your install).

## Step 4: Create the database tables

This project includes a file, `budgetcli_schema.sql`, containing the instructions to set up all the necessary tables. Run it once:

```bash
psql "postgres://<DB_USER>:<DB_PASSWORD>@<DB_HOST>:<DB_PORT>/<DB_NAME>" -f budgetcli_schema.sql
```

Replace each `<...>` with the matching value from your `.env` file.

If successful, you won't see any errors — just some brief output confirming tables were created.

## Step 5: Build the program

From the project's root folder:
```bash
go build -o budgetcli .
```

This creates a program file named `budgetcli` in the current folder.

## Step 6: Run it

From that same folder:
```bash
./budgetcli
```

This opens the main menu. You can also run a specific action directly:
```bash
./budgetcli add      # Add a new bill
./budgetcli prep     # Prepare bills for an upcoming month
./budgetcli list     # View upcoming bills
./budgetcli update   # Update payment amounts, mark paid, or modify a bill
```

**Keep this project folder around.** As it stands today, the program still reads your database settings from the `.env` file inside this folder every time it runs — so this folder needs to stay in place, and you'll run `./budgetcli` from inside it (or reference its full path) each time.

---

## Project structure

```
budgetcli/
├── cmd/                     # Individual commands (add, prep, list, update, menu)
├── internal/
│   ├── database/            # Handles connecting to Postgres
│   └── models/              # Data structures used throughout the program
├── main.go                  # Program entry point
├── budgetcli_schema.sql     # Database table definitions
├── docker-compose.prod.yml  # Runs Postgres via Docker
├── .env.example             # Template for your own .env file
└── .env                     # Your personal database settings (not included — you create this)
```
