package repository

import (
	"fmt"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/jmoiron/sqlx"
)

type GuessPostgres struct {
	db *sqlx.DB
}

func NewGuessPostgres(db *sqlx.DB) *GuessPostgres {
	return &GuessPostgres{db: db}
}

func (r *GuessPostgres) CreateGuess(guess *model.Guess) error {

	query := fmt.Sprintf("INSERT INTO %s (game_id,guess_word,distance) VALUES ($1,$2,$3) RETURNING id", guessesTable)
	_, err := r.db.Exec(query, guess.GameId, guess.GuessWord, guess.Distance)
	if err != nil {
		return err
	}
	return nil
}

func (r *GuessPostgres) GetAllGuessByGame(gameId int) ([]model.Guess, error) {
	var guesses []model.Guess
	query := fmt.Sprintf("SELECT id, game_id, guess_word, distance, created_at FROM %s WHERE game_id = $1 ORDER BY created_at DESC", guessesTable)

	err := r.db.Select(&guesses, query, gameId)
	if err != nil {
		return nil, err
	}

	return guesses, nil
}

func (r *GuessPostgres) GetGuessById(i int) error {
	return nil
}
