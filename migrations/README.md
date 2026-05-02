# Database Migrations

This directory contains database migration files for DailyTracker using [golang-migrate](https://github.com/golang-migrate/migrate).

## Structure

```
migrations/
├── 000001_init_schema.up.sql
├── 000001_init_schema.down.sql
└── README.md
```

## How It Works

- Migrations run automatically on application startup
- Migrations are **idempotent** - safe to run multiple times
- Migration state is tracked in the `schema_migrations` table
- The application uses PostgreSQL for all environments (development and production)

## File Naming Convention

Files follow the pattern: `{version}_{name}.{direction}.sql`

Example:
- `000001_init_schema.up.sql` - Creates initial tables
- `000001_init_schema.down.sql` - Reverts the migration

## PostgreSQL Syntax

The migration files use PostgreSQL-specific SQL syntax:
- `SERIAL PRIMARY KEY` for auto-incrementing IDs
- `VARCHAR(n)` or `TEXT` for string fields
- `TIMESTAMP` for timestamps
- `SMALLINT CHECK` for constrained integer fields

## Running Migrations

### Local Development

```bash
# Start PostgreSQL and run migrations
go run ./cmd/dailytracker
```

### Docker Compose

```bash
# Starts both PostgreSQL and the application
docker-compose up
```

## Environment Variables

### PostgreSQL Configuration (Required)
- `DB_URL` - PostgreSQL host and port (e.g., `localhost:5432`)
- `DB_USER` - PostgreSQL username
- `DB_PASSWORD` - PostgreSQL password
- `DB_NAME` - PostgreSQL database name

### Optional
- `MIGRATIONS_PATH` - Override migrations path (default: `file://migrations`)

## Creating New Migrations

When adding new migrations:

1. Create new files in the `migrations/` directory
2. Use sequential numbering: `000002_description.up.sql` and `000002_description.down.sql`
3. Use PostgreSQL syntax

Example:
```bash
touch migrations/000002_add_notes_field.up.sql
touch migrations/000002_add_notes_field.down.sql
```

## Current Migrations

- `000001_init_schema` - Creates `users` and `daily_entries` tables with indexes

## Troubleshooting

### "Dirty database version" error

If migrations fail partway through, the database may be marked as "dirty". To fix:

```bash
# Connect to PostgreSQL
psql -U dailytracker -d dailytracker

# Fix dirty flag
UPDATE schema_migrations SET dirty = false;
```

Then manually inspect and fix any partial schema changes before re-running migrations.

### Custom migration path

Override the default migration path:

```bash
MIGRATIONS_PATH=file:///custom/path/to/migrations go run ./cmd/dailytracker
```

## Connecting to PostgreSQL

### Using Docker Compose

```bash
# Connect to the PostgreSQL container
docker exec -it dailytracker-postgres psql -U dailytracker -d dailytracker
```

### Local PostgreSQL

```bash
psql -U dailytracker -d dailytracker -h localhost
```
