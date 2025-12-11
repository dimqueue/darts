-- +migrate Up

CREATE TABLE "users"
(
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(100)        NOT NULL,
    username      VARCHAR(50) UNIQUE  NOT NULL,
    password_hash VARCHAR(255)        NOT NULL,
    email         VARCHAR(255),
    is_active     BOOLEAN NOT NULL DEFAULT true,
    last_login_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE "words"
(
    id         BIGSERIAL PRIMARY KEY,
    word       VARCHAR(100) NOT NULL,
    language   VARCHAR(3)   NOT NULL,
    difficulty INTEGER NOT NULL DEFAULT 1 CHECK (difficulty BETWEEN 1 AND 5),
    is_active  BOOLEAN NOT NULL DEFAULT true,
    UNIQUE(word, language)
);

CREATE TABLE "games"
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    word_id    BIGINT      NOT NULL REFERENCES words (id) ON DELETE CASCADE,
    status     VARCHAR(20) NOT NULL CHECK (status IN ('in_progress', 'won', 'lost', 'abandoned')),
    language   VARCHAR(3)  NOT NULL,
    started_at TIMESTAMPTZ DEFAULT NOW(),
    ended_at   TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    CONSTRAINT check_game_completion CHECK (
        (status = 'in_progress' AND ended_at IS NULL) OR
        (status IN ('won', 'lost', 'abandoned') AND ended_at IS NOT NULL)
    )
);

CREATE TABLE "guesses"
(
    id         BIGSERIAL PRIMARY KEY,
    game_id    BIGINT       NOT NULL REFERENCES games (id) ON DELETE CASCADE,
    guess_word VARCHAR(100) NOT NULL,
    distance   INTEGER      NOT NULL CHECK (distance >= 0),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (game_id, guess_word)
);

CREATE TABLE "user_profiles" (
    user_id       BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    avatar_url    VARCHAR(500),
    bio           TEXT,
    country_code  VARCHAR(2),
    timezone      VARCHAR(50) DEFAULT 'UTC',
    date_of_birth DATE,
    created_at    TIMESTAMPTZ DEFAULT NOW(),
    updated_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE "user_settings" (
    user_id               BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    preferred_language    VARCHAR(3) NOT NULL DEFAULT 'en',
    theme                 VARCHAR(20) NOT NULL DEFAULT 'light',
    sound_enabled         BOOLEAN NOT NULL DEFAULT true,
    email_notifications   BOOLEAN NOT NULL DEFAULT true,
    show_profile_public   BOOLEAN NOT NULL DEFAULT true,
    show_stats_public     BOOLEAN NOT NULL DEFAULT true,
    created_at            TIMESTAMPTZ DEFAULT NOW(),
    updated_at            TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE "user_language_stats" (
    id                  BIGSERIAL PRIMARY KEY,
    user_id             BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    language            VARCHAR(3) NOT NULL,
    games_played        INTEGER NOT NULL DEFAULT 0,
    games_won           INTEGER NOT NULL DEFAULT 0,
    total_guesses       INTEGER NOT NULL DEFAULT 0,
    average_guesses     DECIMAL(5,2) GENERATED ALWAYS AS (
                            CASE WHEN games_played > 0
                                 THEN ROUND(total_guesses::DECIMAL / games_played, 2)
                                 ELSE 0
                            END
                        ) STORED,
    current_streak      INTEGER NOT NULL DEFAULT 0,
    best_streak         INTEGER NOT NULL DEFAULT 0,
    total_score         INTEGER NOT NULL DEFAULT 0,
    fastest_win_seconds INTEGER,
    fewest_guesses_win  INTEGER,
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, language)
);

CREATE TABLE "user_global_streaks" (
    user_id         BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current_streak  INTEGER NOT NULL DEFAULT 0,
    best_streak     INTEGER NOT NULL DEFAULT 0,
    last_game_at    TIMESTAMPTZ,
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE VIEW "user_statistics" AS
SELECT
    u.id as user_id,
    COALESCE(SUM(uls.games_played), 0)::INTEGER as total_games,
    COALESCE(SUM(uls.games_won), 0)::INTEGER as total_wins,
    (COALESCE(SUM(uls.games_played), 0) - COALESCE(SUM(uls.games_won), 0))::INTEGER as total_losses,
    COALESCE(SUM(uls.total_guesses), 0)::INTEGER as total_guesses,
    COALESCE(ROUND(SUM(uls.total_guesses)::DECIMAL / NULLIF(SUM(uls.games_played), 0), 2), 0)::DECIMAL(5,2) as average_guesses,
    COALESCE(SUM(uls.total_score), 0)::INTEGER as total_score,
    MIN(uls.fastest_win_seconds) as fastest_win_seconds,
    MIN(uls.fewest_guesses_win) as fewest_guesses_win,
    COALESCE(ugs.current_streak, 0) as current_win_streak,
    COALESCE(ugs.best_streak, 0) as best_win_streak,
    ugs.last_game_at,
    GREATEST(MAX(uls.updated_at), ugs.updated_at) as updated_at
FROM users u
LEFT JOIN user_language_stats uls ON u.id = uls.user_id
LEFT JOIN user_global_streaks ugs ON u.id = ugs.user_id
GROUP BY u.id, ugs.current_streak, ugs.best_streak, ugs.last_game_at, ugs.updated_at;

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $update_modified_column$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$update_modified_column$ LANGUAGE plpgsql;
-- +migrate StatementEnd

CREATE TRIGGER update_language_stats_modtime
    BEFORE UPDATE ON user_language_stats
    FOR EACH ROW EXECUTE FUNCTION update_modified_column();

CREATE TRIGGER update_global_streaks_modtime
    BEFORE UPDATE ON user_global_streaks
    FOR EACH ROW EXECUTE FUNCTION update_modified_column();

-- Views
CREATE VIEW global_leaderboard_view AS
SELECT
    RANK() OVER (ORDER BY us.total_score DESC, us.total_wins DESC,
                 COALESCE(us.average_guesses, 999) ASC) as rank,
    u.id as user_id,
    u.username,
    u.name,
    up.avatar_url,
    up.country_code,
    us.total_score,
    us.total_wins as total_wins,
    us.total_games as total_games,
    us.best_win_streak,
    COALESCE(us.average_guesses, 0) as average_guesses,
    CASE WHEN us.total_games > 0
         THEN ROUND(us.total_wins::DECIMAL / us.total_games * 100, 2)
         ELSE 0 END as win_rate
FROM users u
INNER JOIN user_statistics us ON u.id = us.user_id
LEFT JOIN user_profiles up ON u.id = up.user_id
WHERE u.is_active = true
  AND us.total_games >= 3;

CREATE VIEW user_profile_summary AS
SELECT
    u.id,
    u.username,
    u.name,
    u.email,
    u.created_at as member_since,
    u.last_login_at,
    up.avatar_url,
    up.bio,
    up.country_code,
    COALESCE(uset.preferred_language, 'en') as preferred_language,
    COALESCE(uset.theme, 'light') as theme,
    COALESCE(uset.show_profile_public, true) as show_profile_public,
    COALESCE(uset.show_stats_public, true) as show_stats_public,
    COALESCE(us.total_games, 0) as total_games,
    COALESCE(us.total_wins, 0) as total_wins,
    COALESCE(us.total_losses, 0) as total_losses,
    COALESCE(us.current_win_streak, 0) as current_win_streak,
    COALESCE(us.best_win_streak, 0) as best_win_streak,
    COALESCE(us.average_guesses, 0) as average_guesses,
    COALESCE(us.total_score, 0) as total_score,
    us.last_game_at
FROM users u
LEFT JOIN user_profiles up ON u.id = up.user_id
LEFT JOIN user_settings uset ON u.id = uset.user_id
LEFT JOIN user_statistics us ON u.id = us.user_id;

-- +migrate Down

DROP VIEW IF EXISTS user_profile_summary;
DROP VIEW IF EXISTS global_leaderboard_view;
DROP VIEW IF EXISTS user_statistics;

DROP TRIGGER IF EXISTS update_language_stats_modtime ON user_language_stats;
DROP TRIGGER IF EXISTS update_global_streaks_modtime ON user_global_streaks;
DROP FUNCTION IF EXISTS update_modified_column();

DROP TABLE IF EXISTS user_global_streaks;
DROP TABLE IF EXISTS user_language_stats;
DROP TABLE IF EXISTS user_settings;
DROP TABLE IF EXISTS user_profiles;
DROP TABLE IF EXISTS guesses;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS words;
DROP TABLE IF EXISTS users;