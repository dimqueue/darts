-- +migrate Up

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE email IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_words_language ON words(language);
CREATE INDEX IF NOT EXISTS idx_words_language_active ON words(language, is_active) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_words_language_difficulty ON words(language, difficulty) WHERE is_active = true;

CREATE INDEX IF NOT EXISTS idx_games_user_id ON games(user_id);
CREATE INDEX IF NOT EXISTS idx_games_word_id ON games(word_id);
CREATE INDEX IF NOT EXISTS idx_games_status ON games(status);
CREATE INDEX IF NOT EXISTS idx_games_started_at ON games(started_at DESC);
CREATE INDEX IF NOT EXISTS idx_games_user_status ON games(user_id, status);
CREATE INDEX IF NOT EXISTS idx_games_user_started ON games(user_id, started_at DESC);
CREATE INDEX IF NOT EXISTS idx_games_expires ON games(expires_at) WHERE status = 'in_progress';
CREATE INDEX IF NOT EXISTS idx_games_status_expires_partial ON games(status, expires_at)
    WHERE status = 'in_progress' AND expires_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_games_started_language_status ON games(started_at DESC, language, status);

CREATE INDEX IF NOT EXISTS idx_guesses_game_id ON guesses(game_id);
CREATE INDEX IF NOT EXISTS idx_guesses_created_at ON guesses(created_at);

CREATE INDEX IF NOT EXISTS idx_user_lang_stats_user ON user_language_stats(user_id);
CREATE INDEX IF NOT EXISTS idx_user_lang_stats_language ON user_language_stats(language);
CREATE INDEX IF NOT EXISTS idx_language_stats_user_lang ON user_language_stats(user_id, language);
CREATE INDEX IF NOT EXISTS idx_language_stats_ranking ON user_language_stats(language, total_score DESC, games_won DESC, average_guesses ASC)
    WHERE games_played >= 3;
CREATE INDEX IF NOT EXISTS idx_language_stats_score ON user_language_stats(total_score DESC);

CREATE INDEX IF NOT EXISTS idx_global_streaks_user ON user_global_streaks(user_id);

-- +migrate Down

DROP INDEX IF EXISTS idx_users_email;

DROP INDEX IF EXISTS idx_words_language;
DROP INDEX IF EXISTS idx_words_language_active;
DROP INDEX IF EXISTS idx_words_language_difficulty;

DROP INDEX IF EXISTS idx_games_user_id;
DROP INDEX IF EXISTS idx_games_word_id;
DROP INDEX IF EXISTS idx_games_status;
DROP INDEX IF EXISTS idx_games_started_at;
DROP INDEX IF EXISTS idx_games_user_status;
DROP INDEX IF EXISTS idx_games_user_started;
DROP INDEX IF EXISTS idx_games_expires;
DROP INDEX IF EXISTS idx_games_status_expires_partial;
DROP INDEX IF EXISTS idx_games_started_language_status;

DROP INDEX IF EXISTS idx_guesses_game_id;
DROP INDEX IF EXISTS idx_guesses_created_at;

DROP INDEX IF EXISTS idx_user_lang_stats_user;
DROP INDEX IF EXISTS idx_user_lang_stats_language;
DROP INDEX IF EXISTS idx_language_stats_user_lang;
DROP INDEX IF EXISTS idx_language_stats_ranking;
DROP INDEX IF EXISTS idx_language_stats_score;

DROP INDEX IF EXISTS idx_global_streaks_user;