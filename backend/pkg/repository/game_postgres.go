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

func (r *GamePostgres) GetGameById(q Querier, gameId int64, forUpdate bool) (*model.Game, error) {
	var game model.Game

	query := fmt.Sprintf("SELECT id, user_id, word_id, status, language, started_at, ended_at, expires_at FROM %s WHERE id=$1", gamesTable)
	if forUpdate {
		query += " FOR UPDATE"
	}

	err := q.Get(&game, query, gameId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("game not found")
		}
		return nil, err
	}

	return &game, nil
}

func (r *GamePostgres) UpdateGameStatus(q Querier, gameId int64, status string) error {
	query := fmt.Sprintf("UPDATE %s SET status=$1, ended_at=NOW() WHERE id=$2", gamesTable)
	_, err := q.Exec(query, status, gameId)
	return err
}

func (r *GamePostgres) CreateGuess(q Querier, guess *model.Guess) error {
	query := fmt.Sprintf("INSERT INTO %s (game_id,guess_word,distance) VALUES ($1,$2,$3)", guessesTable)
	_, err := q.Exec(query, guess.GameId, guess.GuessWord, guess.Distance)
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

func (r *GamePostgres) CountGuessesByGame(q Querier, gameId int64) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE game_id = $1", guessesTable)
	err := q.Get(&count, query, gameId)
	return count, err
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

func (r *GamePostgres) GuessExists(q Querier, gameId int64, guessWord string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE game_id = $1 AND guess_word = $2)", guessesTable)
	err := q.Get(&exists, query, gameId, guessWord)
	return exists, err
}
