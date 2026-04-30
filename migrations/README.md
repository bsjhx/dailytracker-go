# Database Migrations

This directory contains database migration files using [golang-migrate](https://github.com/golang-migrate/migrate).

## How it works

- Migrations run automatically on application startup
- Migrations are **idempotent** - safe to run multiple times
- Migration state is tracked in the `schema_migrations` table

## File naming convention

Files follow the pattern: `{version}_{name}.{direction}.sql`

Example:
- `000001_init_schema.up.sql` - Creates initial tables
- `000001_init_schema.down.sql` - Reverts the migration

## Creating new migrations

1. Create two files with the next version number:
   ```
   000002_add_new_feature.up.sql
   000002_add_new_feature.down.sql
   ```

2. Write the SQL for applying changes in the `.up.sql` file
3. Write the SQL for reverting changes in the `.down.sql` file

## Configuration

Set `MIGRATIONS_PATH` environment variable to override the default path:
```bash
MIGRATIONS_PATH=file:///custom/path/to/migrations
```

## Current migrations

- `000001_init_schema` - Creates `users` and `daily_entries` tables with indexes
