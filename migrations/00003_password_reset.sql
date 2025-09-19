-- +goose Up
-- +goose StatementBegin
CREATE TABLE password_resets (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE REFERENCES users(id) ON DELETE CASCADE, -- defines the relationship between the session and the user
    token_hash TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE password_resets;
-- +goose StatementEnd
