-- +migrate Up

INSERT INTO users (username, name, password_hash, email, is_active) VALUES
    ('player1', 'Alex Johnson', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjZQQX9HqVL6sC1LIaP3Ufv5KXAGSK', 'alex@example.com', true),
    ('player2', 'Maria Garcia', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjZQQX9HqVL6sC1LIaP3Ufv5KXAGSK', 'maria@example.com', true),
    ('player3', 'John Smith', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjZQQX9HqVL6sC1LIaP3Ufv5KXAGSK', 'john@example.com', true),
    ('player4', 'Emma Wilson', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjZQQX9HqVL6sC1LIaP3Ufv5KXAGSK', 'emma@example.com', true),
    ('player5', 'David Brown', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjZQQX9HqVL6sC1LIaP3Ufv5KXAGSK', 'david@example.com', true),
    ('player6', 'Sofia Martinez', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjZQQX9HqVL6sC1LIaP3Ufv5KXAGSK', 'sofia@example.com', true),
    ('player7', 'James Taylor', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjZQQX9HqVL6sC1LIaP3Ufv5KXAGSK', 'james@example.com', true),
    ('player8', 'Olivia Anderson', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjZQQX9HqVL6sC1LIaP3Ufv5KXAGSK', 'olivia@example.com', true),
    ('player9', 'William Lee', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjZQQX9HqVL6sC1LIaP3Ufv5KXAGSK', 'william@example.com', true),
    ('player10', 'Ava Thomas', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjZQQX9HqVL6sC1LIaP3Ufv5KXAGSK', 'ava@example.com', true)
ON CONFLICT (username) DO NOTHING;

INSERT INTO user_profiles (user_id, avatar_url, bio, country_code, timezone)
SELECT id, NULL, 'I love playing word games!',
    CASE (id % 5)
        WHEN 0 THEN 'US'
        WHEN 1 THEN 'GB'
        WHEN 2 THEN 'UA'
        WHEN 3 THEN 'DE'
        WHEN 4 THEN 'FR'
    END,
    'UTC'
FROM users WHERE username LIKE 'player%'
ON CONFLICT (user_id) DO NOTHING;

INSERT INTO user_settings (user_id, preferred_language, theme, sound_enabled, email_notifications, show_profile_public, show_stats_public)
SELECT id,
    CASE WHEN id % 2 = 0 THEN 'en' ELSE 'ua' END,
    CASE (id % 3) WHEN 0 THEN 'light' WHEN 1 THEN 'dark' ELSE 'auto' END,
    true, true, true, true
FROM users WHERE username LIKE 'player%'
ON CONFLICT (user_id) DO NOTHING;

INSERT INTO user_global_streaks (user_id, current_streak, best_streak, last_game_at)
VALUES
    ((SELECT id FROM users WHERE username = 'player1'), 5, 8, NOW() - INTERVAL '1 day'),
    ((SELECT id FROM users WHERE username = 'player2'), 3, 6, NOW() - INTERVAL '1 day'),
    ((SELECT id FROM users WHERE username = 'player3'), 4, 7, NOW() - INTERVAL '2 days'),
    ((SELECT id FROM users WHERE username = 'player4'), 2, 5, NOW() - INTERVAL '2 days'),
    ((SELECT id FROM users WHERE username = 'player5'), 1, 4, NOW() - INTERVAL '3 days'),
    ((SELECT id FROM users WHERE username = 'player6'), 3, 5, NOW() - INTERVAL '3 days'),
    ((SELECT id FROM users WHERE username = 'player7'), 0, 4, NOW() - INTERVAL '4 days'),
    ((SELECT id FROM users WHERE username = 'player8'), 2, 3, NOW() - INTERVAL '4 days'),
    ((SELECT id FROM users WHERE username = 'player9'), 1, 3, NOW() - INTERVAL '5 days'),
    ((SELECT id FROM users WHERE username = 'player10'), 0, 2, NOW() - INTERVAL '5 days')
ON CONFLICT (user_id) DO NOTHING;

INSERT INTO user_language_stats (user_id, language, games_played, games_won, total_guesses, best_streak, current_streak, total_score, fastest_win_seconds, fewest_guesses_win)
VALUES
    ((SELECT id FROM users WHERE username = 'player1'), 'en', 30, 22, 130, 6, 3, 2200, 45, 2),
    ((SELECT id FROM users WHERE username = 'player2'), 'en', 25, 18, 110, 5, 2, 1800, 52, 3),
    ((SELECT id FROM users WHERE username = 'player3'), 'en', 22, 16, 95, 5, 3, 1600, 48, 2),
    ((SELECT id FROM users WHERE username = 'player4'), 'en', 20, 14, 85, 4, 1, 1400, 55, 3),
    ((SELECT id FROM users WHERE username = 'player5'), 'en', 18, 12, 80, 3, 0, 1200, 60, 4),
    ((SELECT id FROM users WHERE username = 'player6'), 'en', 16, 10, 70, 4, 2, 1000, 58, 3),
    ((SELECT id FROM users WHERE username = 'player7'), 'en', 15, 9, 65, 3, 0, 900, 65, 4),
    ((SELECT id FROM users WHERE username = 'player8'), 'en', 14, 8, 60, 2, 1, 800, 70, 4),
    ((SELECT id FROM users WHERE username = 'player9'), 'en', 12, 7, 55, 2, 0, 700, 72, 5),
    ((SELECT id FROM users WHERE username = 'player10'), 'en', 10, 5, 50, 2, 0, 500, 80, 5)
ON CONFLICT (user_id, language) DO NOTHING;

INSERT INTO user_language_stats (user_id, language, games_played, games_won, total_guesses, best_streak, current_streak, total_score, fastest_win_seconds, fewest_guesses_win)
VALUES
    ((SELECT id FROM users WHERE username = 'player1'), 'ua', 20, 13, 90, 4, 2, 1300, 50, 3),
    ((SELECT id FROM users WHERE username = 'player2'), 'ua', 20, 12, 90, 4, 1, 1200, 55, 3),
    ((SELECT id FROM users WHERE username = 'player3'), 'ua', 18, 12, 80, 4, 1, 1200, 52, 3),
    ((SELECT id FROM users WHERE username = 'player4'), 'ua', 18, 11, 80, 3, 1, 1100, 58, 4),
    ((SELECT id FROM users WHERE username = 'player5'), 'ua', 17, 10, 75, 3, 1, 1000, 62, 4),
    ((SELECT id FROM users WHERE username = 'player6'), 'ua', 16, 10, 75, 3, 1, 1000, 60, 4),
    ((SELECT id FROM users WHERE username = 'player7'), 'ua', 15, 9, 70, 3, 0, 900, 68, 5),
    ((SELECT id FROM users WHERE username = 'player8'), 'ua', 14, 8, 65, 2, 1, 800, 73, 5),
    ((SELECT id FROM users WHERE username = 'player9'), 'ua', 13, 7, 60, 2, 1, 700, 75, 5),
    ((SELECT id FROM users WHERE username = 'player10'), 'ua', 10, 5, 50, 1, 0, 500, 85, 6)
ON CONFLICT (user_id, language) DO NOTHING;

INSERT INTO words (word, language, difficulty, is_active) VALUES
    ('apple', 'en', 1, true),
    ('house', 'en', 1, true),
    ('water', 'en', 1, true),
    ('music', 'en', 1, true),
    ('dream', 'en', 2, true),
    ('garden', 'en', 2, true),
    ('planet', 'en', 2, true),
    ('mystery', 'en', 3, true),
    ('adventure', 'en', 3, true),
    ('beautiful', 'en', 3, true),
    ('mustang', 'en', 2, true),
    ('ocean', 'en', 1, true),
    ('knight', 'en', 2, true),
    ('claptrap', 'en', 3, true)
ON CONFLICT (word, language) DO NOTHING;

INSERT INTO words (word, language, difficulty, is_active) VALUES
    ('сонце', 'ua', 1, true),
    ('вода', 'ua', 1, true),
    ('дерево', 'ua', 1, true),
    ('музика', 'ua', 2, true),
    ('мрія', 'ua', 2, true),
    ('щастя', 'ua', 2, true),
    ('природа', 'ua', 2, true),
    ('кохання', 'ua', 3, true),
    ('пригода', 'ua', 3, true),
    ('майбутнє', 'ua', 3, true)
ON CONFLICT (word, language) DO NOTHING;

INSERT INTO games (user_id, word_id, language, status, started_at, ended_at) VALUES
    ((SELECT id FROM users WHERE username = 'player1'), (SELECT id FROM words WHERE word = 'apple'), 'en', 'won', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day' + INTERVAL '5 minutes'),
    ((SELECT id FROM users WHERE username = 'player1'), (SELECT id FROM words WHERE word = 'house'), 'en', 'won', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days' + INTERVAL '7 minutes'),
    ((SELECT id FROM users WHERE username = 'player1'), (SELECT id FROM words WHERE word = 'water'), 'en', 'lost', NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days' + INTERVAL '10 minutes'),
    ((SELECT id FROM users WHERE username = 'player1'), (SELECT id FROM words WHERE word = 'сонце'), 'ua', 'won', NOW() - INTERVAL '4 days', NOW() - INTERVAL '4 days' + INTERVAL '6 minutes'),
    ((SELECT id FROM users WHERE username = 'player1'), (SELECT id FROM words WHERE word = 'вода'), 'ua', 'won', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days' + INTERVAL '8 minutes');

INSERT INTO games (user_id, word_id, language, status, started_at, ended_at) VALUES
    ((SELECT id FROM users WHERE username = 'player2'), (SELECT id FROM words WHERE word = 'music'), 'en', 'won', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day' + INTERVAL '6 minutes'),
    ((SELECT id FROM users WHERE username = 'player2'), (SELECT id FROM words WHERE word = 'dream'), 'en', 'won', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days' + INTERVAL '8 minutes'),
    ((SELECT id FROM users WHERE username = 'player2'), (SELECT id FROM words WHERE word = 'музика'), 'ua', 'lost', NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days' + INTERVAL '12 minutes');

INSERT INTO guesses (game_id, guess_word, distance, created_at)
SELECT g.id, 'fruit', 150, g.started_at + INTERVAL '1 minute'
FROM games g WHERE g.word_id = (SELECT id FROM words WHERE word = 'apple') LIMIT 1;

INSERT INTO guesses (game_id, guess_word, distance, created_at)
SELECT g.id, 'orange', 80, g.started_at + INTERVAL '2 minutes'
FROM games g WHERE g.word_id = (SELECT id FROM words WHERE word = 'apple') LIMIT 1;

INSERT INTO guesses (game_id, guess_word, distance, created_at)
SELECT g.id, 'apple', 0, g.started_at + INTERVAL '3 minutes'
FROM games g WHERE g.word_id = (SELECT id FROM words WHERE word = 'apple') LIMIT 1;

-- +migrate Down

DELETE FROM guesses WHERE game_id IN (SELECT id FROM games WHERE user_id IN (SELECT id FROM users WHERE username LIKE 'player%'));
DELETE FROM games WHERE user_id IN (SELECT id FROM users WHERE username LIKE 'player%');
DELETE FROM words WHERE word IN ('apple', 'house', 'water', 'music', 'dream', 'garden', 'planet', 'mystery', 'adventure', 'beautiful', 'mustang', 'ocean', 'knight', 'claptrap', 'сонце', 'вода', 'дерево', 'музика', 'мрія', 'щастя', 'природа', 'кохання', 'пригода', 'майбутнє');
DELETE FROM user_language_stats WHERE user_id IN (SELECT id FROM users WHERE username LIKE 'player%');
DELETE FROM user_global_streaks WHERE user_id IN (SELECT id FROM users WHERE username LIKE 'player%');
DELETE FROM user_settings WHERE user_id IN (SELECT id FROM users WHERE username LIKE 'player%');
DELETE FROM user_profiles WHERE user_id IN (SELECT id FROM users WHERE username LIKE 'player%');
DELETE FROM users WHERE username LIKE 'player%';