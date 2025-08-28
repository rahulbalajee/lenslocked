CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE REFERENCES users(id) ON DELETE CASCADE, -- defines the relationship between the session and the user
    token_hash TEXT UNIQUE NOT NULL
);
