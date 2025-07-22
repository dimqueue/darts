-- +migrate Up
CREATE TABLE "users"
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(100)       NOT NULL,
    username      VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255)       NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down

DROP TABLE IF EXISTS "users";