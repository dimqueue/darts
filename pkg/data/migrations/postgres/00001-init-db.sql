-- +migrate Up
CREATE TABLE "user"
(
    id    SERIAL PRIMARY KEY,
    email TEXT
);

-- +migrate Down

DROP TABLE IF EXISTS "user";