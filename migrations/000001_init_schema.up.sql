-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create daily_entries table
CREATE TABLE IF NOT EXISTS daily_entries (
    id SERIAL PRIMARY KEY,
    entry_date DATE NOT NULL,
    work_score SMALLINT CHECK (work_score BETWEEN 0 AND 5),
    personal_score SMALLINT CHECK (personal_score BETWEEN 0 AND 5),
    total INTEGER,
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_entry_date ON daily_entries(entry_date DESC);
CREATE INDEX IF NOT EXISTS idx_user_entries ON daily_entries(user_id, entry_date DESC);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_date ON daily_entries(user_id, entry_date);
