package repository

import (
	"database/sql"
	"fmt"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/jmoiron/sqlx"
)

type GamePostgres struct {
	db *sqlx.DB
}

func NewGamePostgres(db *sqlx.DB) *GamePostgres {
	return &GamePostgres{db: db}
}

func (r *GamePostgres) CreateGame(game *model.Game) (int64, error) {
	var id int64

	query := fmt.Sprintf("INSERT INTO %s (user_id,word_id,status,language,expires_at) VALUES ($1,$2,$3,$4,$5) RETURNING id", gamesTable)
	row := r.db.QueryRow(query, game.UserId, game.WordId, game.Status, game.Language, game.ExpiresAt)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *GamePostgres) GetAllGames(userId int64) ([]model.Game, error) {
	games := make([]model.Game, 0)

	query := fmt.Sprintf("SELECT id, user_id, word_id, status, language, started_at, ended_at, expires_at FROM %s WHERE user_id=$1 ORDER BY started_at DESC LIMIT 100", gamesTable)
	err := r.db.Select(&games, query, userId)
	return games, err
}

func (r *GamePostgres) GetGameById(gameId int64) (*model.Game, error) {
	var game model.Game

	query := fmt.Sprintf("SELECT id, user_id, word_id, status, language, started_at, ended_at, expires_at FROM %s WHERE id=$1", gamesTable)
	err := r.db.Get(&game, query, gameId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("game not found")
		}
		return nil, err
	}

	return &game, nil
}

func (r *GamePostgres) UpdateGame(gameId int64) (*model.Game, error) {
	return nil, nil
}

func (r *GamePostgres) UpdateGameStatus(gameId int64, status string) error {
	query := fmt.Sprintf("UPDATE %s SET status=$1, ended_at=NOW() WHERE id=$2", gamesTable)
	_, err := r.db.Exec(query, status, gameId)
	return err
}

func (r *GamePostgres) DeleteGame(gameId int64) (*model.Game, error) {
	return nil, nil
}

func (r *GamePostgres) ExpireGames() (int64, error) {
	query := fmt.Sprintf(`
		UPDATE %s
		SET status = 'lost', ended_at = NOW()
		WHERE status = 'in_progress'
		  AND expires_at IS NOT NULL
		  AND expires_at < NOW()
	`, gamesTable)

	result, err := r.db.Exec(query)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (r *GamePostgres) UpdateGameStatusTx(tx *sqlx.Tx, gameId int64, status string) error {
	query := fmt.Sprintf("UPDATE %s SET status=$1, ended_at=NOW() WHERE id=$2", gamesTable)
	_, err := tx.Exec(query, status, gameId)
	return err
}

func (r *GamePostgres) CreateGuess(guess *model.Guess) error {
	query := fmt.Sprintf("INSERT INTO %s (game_id,guess_word,distance) VALUES ($1,$2,$3) RETURNING id", guessesTable)
	_, err := r.db.Exec(query, guess.GameId, guess.GuessWord, guess.Distance)
	return err
}

func (r *GamePostgres) CreateGuessTx(tx *sqlx.Tx, guess *model.Guess) error {
	query := fmt.Sprintf("INSERT INTO %s (game_id,guess_word,distance) VALUES ($1,$2,$3)", guessesTable)
	_, err := tx.Exec(query, guess.GameId, guess.GuessWord, guess.Distance)
	return err
}

func (r *GamePostgres) GetAllGuessByGame(gameId int64) ([]model.Guess, error) {
	guesses := make([]model.Guess, 0)
	query := fmt.Sprintf("SELECT id, game_id, guess_word, distance, created_at FROM %s WHERE game_id = $1 ORDER BY created_at DESC LIMIT 500", guessesTable)

	err := r.db.Select(&guesses, query, gameId)
	if err != nil {
		return nil, err
	}

	return guesses, nil
}

func (r *GamePostgres) GetGuessById(guessId int64) error {
	return nil
}

func (r *GamePostgres) CountGuessesByGameTx(tx *sqlx.Tx, gameId int64) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE game_id = $1", guessesTable)
	err := tx.Get(&count, query, gameId)
	return count, err
}

func (r *GamePostgres) GetGameByIdForUpdate(tx *sqlx.Tx, gameId int64) (*model.Game, error) {
	var game model.Game
	query := fmt.Sprintf("SELECT id, user_id, word_id, status, language, started_at, ended_at, expires_at FROM %s WHERE id = $1 FOR UPDATE", gamesTable)
	err := tx.Get(&game, query, gameId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("game not found")
		}
		return nil, err
	}
	return &game, nil
}

func (r *GamePostgres) GetExpiredGames() ([]model.Game, error) {
	games := make([]model.Game, 0)
	query := fmt.Sprintf(`
		SELECT id, user_id, word_id, status, language, started_at, ended_at, expires_at
		FROM %s
		WHERE status = 'in_progress'
		  AND expires_at IS NOT NULL
		  AND expires_at < NOW()
	`, gamesTable)
	err := r.db.Select(&games, query)
	return games, err
}

func (r *GamePostgres) GuessExists(gameId int64, guessWord string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE game_id = $1 AND guess_word = $2)", guessesTable)
	err := r.db.Get(&exists, query, gameId, guessWord)
	return exists, err
}

func (r *GamePostgres) GuessExistsTx(tx *sqlx.Tx, gameId int64, guessWord string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE game_id = $1 AND guess_word = $2)", guessesTable)
	err := tx.Get(&exists, query, gameId, guessWord)
	return exists, err
}
