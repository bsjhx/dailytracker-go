-- Drop indexes
DROP INDEX IF EXISTS idx_user_date;
DROP INDEX IF EXISTS idx_user_entries;
DROP INDEX IF EXISTS idx_entry_date;

-- Drop tables
DROP TABLE IF EXISTS daily_entries;
DROP TABLE IF EXISTS users;
