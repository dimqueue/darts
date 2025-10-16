-- +migrate Up
CREATE TABLE "users"
(
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(100)       NOT NULL,
    username      VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255)       NOT NULL,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "words"
(
    id   SERIAL PRIMARY KEY,
    word VARCHAR(100) NOT NULL
);

CREATE TABLE "games"
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER     NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    word_id    INTEGER     NOT NULL REFERENCES words (id) ON DELETE CASCADE,
    status     VARCHAR(20) NOT NULL CHECK (status IN ('in_progress', 'won', 'lost')),
    started_at TIMESTAMP DEFAULT NOW(),
    ended_at   TIMESTAMP
);

CREATE TABLE guesses
(
    id         SERIAL PRIMARY KEY,
    game_id    INTEGER      NOT NULL REFERENCES games (id) ON DELETE CASCADE,
    guess_word VARCHAR(100) NOT NULL,
    distance   INTEGER      NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
-- +migrate Down

DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "words";
DROP TABLE IF EXISTS "games";
DROP TABLE IF EXISTS "guesses";