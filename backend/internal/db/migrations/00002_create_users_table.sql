-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    nickname      VARCHAR(50) NOT NULL,
    avatar        VARCHAR(20) NOT NULL DEFAULT '🚗',
    games_played  INT NOT NULL DEFAULT 0,
    games_won     INT NOT NULL DEFAULT 0,
    total_score   INT NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (lower(email));
CREATE INDEX IF NOT EXISTS idx_users_nickname ON users (lower(nickname));

-- +goose Down
DROP TABLE IF EXISTS users;
