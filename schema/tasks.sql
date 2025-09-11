CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    task TEXT NOT NULL,
    priority TEXT,
    status bool,
    due_date TIMESTAMP
);
