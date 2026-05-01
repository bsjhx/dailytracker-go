# [DT-1][Technical] Migrate to Postgres (local = sqlite, prod = Postres)

App on start should check if it's running in production or development mode. Must be based on env variable `ENV=prod|dev`), and initialize database connection accordingly:
- for local mode use SQLite - as it is now
- for production use Postgres which is hosted on VPS, so url, username and password should be read from env variables (`DB_URL`, `DB_USER`, `DB_PASSWORD`)

Also update migration logic to work with both databases - it must provide separate folder and files for each database, so that we can have different migration files for local and production if needed.

As we dont have any rele1vant data, we do not migrate anything.