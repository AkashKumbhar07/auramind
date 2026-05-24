-- +goose Up
CREATE TABLE users (
    id         TEXT PRIMARY KEY,
    email      TEXT NOT NULL UNIQUE,
    username   TEXT NOT NULL UNIQUE,
    password   TEXT NOT NULL,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

-- +goose Down
DROP TABLE IF EXISTS users;
