# budgetcli

A terminal tool for tracking your bills, figuring out what's due each month, and marking things as paid. It stores everything in a Postgres database.

This guide assumes no prior experience setting any of this up. Follow the steps in order. Where Windows and Linux/macOS differ, both are shown.

---

## What you'll need before starting

1. **Go** — the programming language this tool is written in. You need this to build the program from its source code.
2. **Postgres** — the database that stores your bills. You have two choices for how to run it, explained below.
3. **A way to run SQL commands against Postgres** — either `psql` (a command-line tool) or **pgAdmin** (a graphical tool that comes bundled with the Windows Postgres installer, and works just as well for this guide's purposes).

### Installing Go

Download and install it from the official site: https://go.dev/dl/

Follow the installer instructions for your operating system. Once installed, confirm it worked by opening a terminal (Windows: Command Prompt or PowerShell) and running:
```bash
go version
```
You should see a version number printed back.

### Installing Postgres

You have two options. **If you're not sure which to pick, use Option A (Docker).** It's more self-contained and easier to remove later if you change your mind. **If you're on Windows and don't already use Docker, Option B (native install) is often the more familiar path.**

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

Download and run the official installer for your OS: https://www.postgresql.org/download/

**During the Windows installer:**
- When asked, set a password for the default `postgres` superuser account. **Write this password down** — you'll need it shortly.
- At the end of the installer, it will offer to launch **Stack Builder**. **You can skip this — click Cancel or close it.** Stack Builder installs extra add-ons that aren't needed for this project.
- The installer also installs **pgAdmin 4** (a graphical tool for managing Postgres) and a **SQL Shell (psql)** shortcut — both are used later in this guide.

**Important — what username and database name do you actually have?** Unless you specifically created a new one, the installer sets up:
- A default **username**: `postgres`
- A default **database**: also named `postgres`

You can either use these defaults as-is (simplest — just use `postgres` for both `DB_USER` and `DB_NAME` in the next step), or create a new, separate database and user of your own using pgAdmin. This guide assumes you're using the defaults unless you know you created your own.

### About psql / pgAdmin

You'll need one of these later to run a setup script against the database:
- **psql**: installed automatically with the Windows installer above, available from your Start Menu as "SQL Shell (psql)". On macOS/Linux with a native install, it's usually already on your `PATH`.
- **pgAdmin**: also installed automatically on Windows. To run SQL with it: open pgAdmin, connect to your server (it will prompt for the postgres password you set), right-click your database in the left-hand tree, and choose **Query Tool**. You can paste and run SQL commands there.

Either tool works for the steps below — use whichever you're more comfortable with.

---

## Step 1: Get the project

**If you have `git` installed:**
```bash
git clone <repo-url>
cd budgetcli
```

**If you don't have `git`:** on the project's GitHub page, click the green **`<> Code`** button, then **Download ZIP**. Extract the ZIP file anywhere on your computer — the extracted folder is your project folder. You'll need to know exactly where this folder is for the next steps.

## Step 2: Set up your `.env` file

This project needs to know how to connect to your Postgres database — the address, the username, the password, and so on. Rather than typing those in every time, they're stored in a file called `.env`.

1. In the project folder, find the file named `.env.example`. Make a copy of it and rename the copy to `.env`.
   - Command line: `cp .env.example .env` (Windows PowerShell: `copy .env.example .env`)
   - Or just do this in your file explorer: copy the file, paste it, rename the copy to `.env` (including removing any `.txt` your OS might try to add).
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
   - `DB_USER` — `postgres`, unless you created your own user (see the note under "Installing Postgres directly" above)
   - `DB_PASSWORD` — the password you set for the `postgres` user during install
   - `DB_NAME` — `postgres`, unless you created your own database

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

### If using a native Postgres install

The Postgres installer already set Postgres up to run automatically as a background service — there's nothing you need to start manually. You can confirm it's running:
- **Windows:** open the Start Menu, search "Services", find "postgresql-x64-\<version\>" in the list — it should say "Running".
- **macOS/Linux:** this varies by install method; check your installer's documentation.

## Step 4: Create the database tables

This project includes a file, `budgetcli_schema.sql`. It contains the raw SQL commands needed to create all the tables this program uses — it is not a set of instructions to read, it's a script to *run*, once, using `psql` or pgAdmin.

**Using psql (command line):**
```bash
psql -U <DB_USER> -d <DB_NAME> -h <DB_HOST> -f budgetcli_schema.sql
```
Replace each `<...>` with the matching value from your `.env` file. You'll be prompted for your password. **This command must be run from inside the project folder** (see Step 5 below for how to get your terminal there), since it refers to `budgetcli_schema.sql` by its plain filename.

**Using pgAdmin (graphical):**
1. Open pgAdmin, connect to your server.
2. Right-click your database (e.g. `postgres`) in the left-hand tree → **Query Tool**.
3. Open `budgetcli_schema.sql` in a text editor, copy its entire contents, and paste them into the Query Tool.
4. Click the "Execute/Run" button (a play-button icon, or press F5).

Either way, if successful, you won't see any errors — the tables have been created.

## Step 5: Build the program

First, you need your terminal (Command Prompt, PowerShell, or a terminal app) to be *inside* the project folder — not just anywhere on your computer.

- **Windows:** open the project folder in File Explorer, click in the address bar at the top, type `cmd`, and press Enter — this opens a Command Prompt already inside that folder. (Alternative: right-click inside the folder and look for "Open in Terminal" or "Open PowerShell window here".)
- **macOS/Linux:** open a terminal, then use `cd` to navigate into the project folder, e.g. `cd ~/Downloads/budgetcli`.

Once your terminal is in the project folder, build the program:

**Windows:**
```powershell
go build -o budgetcli.exe .
```

**macOS/Linux:**
```bash
go build -o budgetcli .
```

This creates a program file (`budgetcli.exe` on Windows, `budgetcli` on macOS/Linux) in that same folder.

## Step 6: Run it

Still from the same terminal, in the same folder:

**Windows:**
```powershell
.\budgetcli.exe
```

**macOS/Linux:**
```bash
./budgetcli
```

**Windows note:** double-clicking the `.exe` file directly in File Explorer will not work correctly for this program — it needs to be run from the command line, using the command above, not opened like a normal application.

This opens the main menu. You can also run a specific action directly:
```powershell
.\budgetcli.exe add      # Add a new bill
.\budgetcli.exe prep     # Prepare bills for an upcoming month
.\budgetcli.exe list     # View upcoming bills
.\budgetcli.exe update   # Update payment amounts, mark paid, or modify a bill
```
(macOS/Linux: same commands, using `./budgetcli` instead of `.\budgetcli.exe`.)

**Keep this project folder around.** As it stands today, the program still reads your database settings from the `.env` file inside this folder every time it runs — so this folder needs to stay in place, and you'll run the program from inside it (or reference its full path) each time.

---

## Project structure

```
budgetcli/
├── cmd/                     # Individual commands (add, prep, list, update, menu)
├── internal/
│   ├── database/            # Handles connecting to Postgres
│   └── models/              # Data structures used throughout the program
├── main.go                  # Program entry point
├── budgetcli_schema.sql     # Database table definitions (run once, via psql or pgAdmin)
├── docker-compose.prod.yml  # Runs Postgres via Docker
├── .env.example             # Template for your own .env file
└── .env                     # Your personal database settings (not included — you create this)
```
