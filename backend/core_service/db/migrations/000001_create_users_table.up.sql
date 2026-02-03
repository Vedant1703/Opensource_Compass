CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    github_username TEXT UNIQUE NOT NULL,
    experience_level TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

