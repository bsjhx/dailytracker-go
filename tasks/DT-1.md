# [DT-1][Technical] Migrate to Postgres

Remove SQLite at all. Use Postgres for both development and production.
Use envs to configure connection string.
Update documentation and deployment scripts accordingly.
Docker setup should also be updated to use Postgres.
In docker compose for local development, use official Postgres image and set up a volume for data persistence. Refactor docker compose for prod.