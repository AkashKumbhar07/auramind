-- +goose Up
CREATE TABLE alerts (
    id         TEXT PRIMARY KEY,
    user_id    TEXT NOT NULL,
    symbol     TEXT NOT NULL,
    condition  TEXT NOT NULL,
    value      REAL NOT NULL,
    active     INTEGER NOT NULL DEFAULT 1,
    created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_alerts_user ON alerts(user_id);
CREATE INDEX idx_alerts_active ON alerts(active);

-- +goose Down
DROP TABLE IF EXISTS alerts;
