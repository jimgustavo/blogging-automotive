CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title TEXT,
    category TEXT,
    picture TEXT,
    summary TEXT,
    author TEXT,
    editor_data TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
