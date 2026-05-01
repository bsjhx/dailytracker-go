-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create daily_entries table
CREATE TABLE IF NOT EXISTS daily_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    entry_date DATE NOT NULL,
    work_score INTEGER CHECK (work_score BETWEEN 0 AND 5),
    personal_score INTEGER CHECK (personal_score BETWEEN 0 AND 5),
    total INTEGER,
    user_id INTEGER REFERENCES users(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_entry_date ON daily_entries(entry_date DESC);
CREATE INDEX IF NOT EXISTS idx_user_entries ON daily_entries(user_id, entry_date DESC);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_date ON daily_entries(user_id, entry_date);
