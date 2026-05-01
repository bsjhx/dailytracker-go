# Database Migrations

This directory contains database migration files for DailyTracker using [golang-migrate](https://github.com/golang-migrate/migrate).

## Structure

The migrations are organized by database type to support both SQLite (development) and PostgreSQL (production):

```
migrations/
├── sqlite/
│   ├── 000001_init_schema.up.sql
│   └── 000001_init_schema.down.sql
├── postgres/
│   ├── 000001_init_schema.up.sql
│   └── 000001_init_schema.down.sql
└── README.md
```

## How It Works

- Migrations run automatically on application startup
- Migrations are **idempotent** - safe to run multiple times
- Migration state is tracked in the `schema_migrations` table
- The application automatically selects the appropriate migration files based on the `ENV` environment variable:
  - **Development (ENV=dev or not set)**: Uses SQLite migrations from `migrations/sqlite/`
  - **Production (ENV=prod)**: Uses PostgreSQL migrations from `migrations/postgres/`

If the database-specific folder doesn't exist, the application will fall back to the generic `migrations/` folder for backward compatibility.

## File Naming Convention

Files follow the pattern: `{version}_{name}.{direction}.sql`

Example:
- `000001_init_schema.up.sql` - Creates initial tables
- `000001_init_schema.down.sql` - Reverts the migration

## Database-Specific Syntax

The migration files contain database-specific SQL syntax:

**SQLite:**
- `INTEGER PRIMARY KEY AUTOINCREMENT` for auto-incrementing IDs
- `TEXT` for string fields
- `DATETIME` for timestamps
- `INTEGER CHECK` for constrained integer fields

**PostgreSQL:**
- `SERIAL PRIMARY KEY` for auto-incrementing IDs
- `VARCHAR(n)` or `TEXT` for string fields
- `TIMESTAMP` for timestamps
- `SMALLINT CHECK` for constrained integer fields

## Running Migrations

### Development (SQLite)

```bash
# Default - uses SQLite
go run ./cmd/dailytracker

# Or explicitly set ENV
ENV=dev go run ./cmd/dailytracker
```

### Production (PostgreSQL)

```bash
# Set environment to production and provide Postgres credentials
ENV=prod \
DB_URL=localhost:5432 \
DB_USER=dailytracker \
DB_PASSWORD=your_password \
DB_NAME=dailytracker \
go run ./cmd/dailytracker
```

### Docker Compose

```bash
# Development (SQLite)
docker-compose up app-dev

# Production (PostgreSQL)
ENV=prod docker-compose --profile prod up
```

## Environment Variables

### Common
- `ENV` - Environment mode: "dev" (SQLite) or "prod" (PostgreSQL). Default: "dev"
- `MIGRATIONS_PATH` - Override migrations path (optional)

### SQLite (ENV=dev)
- `DB_PATH` - Path to SQLite database file. Default: `./dailytracker.db`

### PostgreSQL (ENV=prod)
- `DB_URL` - PostgreSQL host and port (e.g., `localhost:5432`)
- `DB_USER` - PostgreSQL username
- `DB_PASSWORD` - PostgreSQL password
- `DB_NAME` - PostgreSQL database name

## Creating New Migrations

When adding new migrations:

1. Create files in both `sqlite/` and `postgres/` directories
2. Use sequential numbering: `000002_description.up.sql` and `000002_description.down.sql`
3. Ensure SQL syntax is appropriate for each database type
4. Test migrations on both databases

Example:
```bash
# SQLite migration
touch migrations/sqlite/000002_add_notes_field.up.sql
touch migrations/sqlite/000002_add_notes_field.down.sql

# PostgreSQL migration
touch migrations/postgres/000002_add_notes_field.up.sql
touch migrations/postgres/000002_add_notes_field.down.sql
```

## Current Migrations

- `000001_init_schema` - Creates `users` and `daily_entries` tables with indexes

## Troubleshooting

### "Dirty database version" error

If migrations fail partway through, the database may be marked as "dirty". To fix:

```bash
# For SQLite
sqlite3 dailytracker.db "UPDATE schema_migrations SET dirty = 0;"

# For PostgreSQL
psql -U dailytracker -d dailytracker -c "UPDATE schema_migrations SET dirty = false;"
```

Then manually inspect and fix any partial schema changes before re-running migrations.

### Custom migration path

Override the default migration path:

```bash
MIGRATIONS_PATH=file:///custom/path/to/migrations go run ./cmd/dailytracker
```
