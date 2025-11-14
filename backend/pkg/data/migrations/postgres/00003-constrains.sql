-- +migrate Up

ALTER TABLE guesses ADD CONSTRAINT unique_game_word UNIQUE (game_id, guess_word);

-- +migrate Down