package repository

import (
	"database/sql"
	"fmt"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/jmoiron/sqlx"
)

type StatisticsPostgres struct {
	db *sqlx.DB
}

func NewStatisticsPostgres(db *sqlx.DB) *StatisticsPostgres {
	return &StatisticsPostgres{db: db}
}

func (r *StatisticsPostgres) GetStatistics(userId int64) (*model.UserStatistics, error) {
	var stats model.UserStatistics

	query := fmt.Sprintf(`
		SELECT user_id, total_games, total_wins, total_losses,
		       current_win_streak, best_win_streak, total_guesses, average_guesses,
		       fastest_win_seconds, fewest_guesses_win, total_score, last_game_at, updated_at
		FROM %s
		WHERE user_id = $1
	`, statisticsTable)

	err := r.db.Get(&stats, query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("statistics not found")
		}
		return nil, err
	}

	return &stats, nil
}

func (r *StatisticsPostgres) CreateGlobalStreaks(tx *sqlx.Tx, userId int64) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (user_id)
		VALUES ($1)
	`, globalStreaksTable)

	_, err := tx.Exec(query, userId)
	return err
}

func (r *StatisticsPostgres) UpdateGlobalStreaksAfterGame(tx *sqlx.Tx, update model.StatisticsUpdate) error {
	var exists bool
	checkQuery := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE user_id = $1)`, globalStreaksTable)
	err := tx.Get(&exists, checkQuery, update.UserId)
	if err != nil {
		return fmt.Errorf("failed to check global streaks existence: %w", err)
	}

	if !exists {
		currentStreak := 0
		if update.IsWin {
			currentStreak = 1
		}
		insertQuery := fmt.Sprintf(`
			INSERT INTO %s (user_id, current_streak, best_streak, last_game_at)
			VALUES ($1, $2, $2, NOW())
		`, globalStreaksTable)
		_, err = tx.Exec(insertQuery, update.UserId, currentStreak)
		return err
	}

	var currentStreaks model.UserGlobalStreaks
	query := fmt.Sprintf(`
		SELECT user_id, current_streak, best_streak
		FROM %s
		WHERE user_id = $1
		FOR UPDATE
	`, globalStreaksTable)

	err = tx.Get(&currentStreaks, query, update.UserId)
	if err != nil {
		return fmt.Errorf("failed to get current global streaks: %w", err)
	}

	newCurrentStreak := currentStreaks.CurrentStreak
	newBestStreak := currentStreaks.BestStreak

	if update.IsWin {
		newCurrentStreak++
		if newCurrentStreak > newBestStreak {
			newBestStreak = newCurrentStreak
		}
	} else {
		newCurrentStreak = 0
	}

	updateQuery := fmt.Sprintf(`
		UPDATE %s
		SET current_streak = $2,
		    best_streak = $3,
		    last_game_at = NOW()
		WHERE user_id = $1
	`, globalStreaksTable)

	_, err = tx.Exec(updateQuery,
		update.UserId,
		newCurrentStreak,
		newBestStreak,
	)

	return err
}

func (r *StatisticsPostgres) GetLanguageStats(userId int64, language string) (*model.UserLanguageStats, error) {
	var stats model.UserLanguageStats

	query := fmt.Sprintf(`
		SELECT id, user_id, language, games_played, games_won,
		       total_guesses, average_guesses, best_streak, current_streak, total_score, updated_at
		FROM %s
		WHERE user_id = $1 AND language = $2
	`, languageStatsTable)

	err := r.db.Get(&stats, query, userId, language)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not an error, just no stats yet
		}
		return nil, err
	}

	return &stats, nil
}

func (r *StatisticsPostgres) GetAllLanguageStats(userId int64) ([]model.UserLanguageStats, error) {
	stats := []model.UserLanguageStats{}

	query := fmt.Sprintf(`
		SELECT id, user_id, language, games_played, games_won,
		       total_guesses, average_guesses, best_streak, current_streak, total_score, updated_at
		FROM %s
		WHERE user_id = $1
		ORDER BY games_played DESC
	`, languageStatsTable)

	err := r.db.Select(&stats, query, userId)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *StatisticsPostgres) UpdateLanguageStats(tx *sqlx.Tx, update model.StatisticsUpdate) error {
	var exists bool
	checkQuery := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE user_id = $1 AND language = $2)`, languageStatsTable)
	err := tx.Get(&exists, checkQuery, update.UserId, update.Language)
	if err != nil {
		return fmt.Errorf("failed to check language stats existence: %w", err)
	}

	gamesWon := 0
	if update.IsWin {
		gamesWon = 1
	}

	if !exists {
		insertQuery := fmt.Sprintf(`
			INSERT INTO %s (user_id, language, games_played, games_won, total_guesses, best_streak, current_streak, total_score, fastest_win_seconds, fewest_guesses_win)
			VALUES ($1, $2, 1, $3, $4, $5, $5, $6, $7, $8)
		`, languageStatsTable)

		currentStreak := 0
		var fastestWin *int
		var fewestGuesses *int
		if update.IsWin {
			currentStreak = 1
			fastestWin = update.TimeSeconds
			fewestGuesses = &update.GuessCount
		}

		_, err = tx.Exec(insertQuery,
			update.UserId,
			update.Language,
			gamesWon,
			update.GuessCount,
			currentStreak,
			update.ScoreEarned,
			fastestWin,
			fewestGuesses,
		)
		return err
	}

	var currentStats model.UserLanguageStats
	selectQuery := fmt.Sprintf(`
		SELECT id, user_id, language, games_played, games_won, total_guesses, best_streak, current_streak, total_score, fastest_win_seconds, fewest_guesses_win
		FROM %s
		WHERE user_id = $1 AND language = $2
		FOR UPDATE
	`, languageStatsTable)

	err = tx.Get(&currentStats, selectQuery, update.UserId, update.Language)
	if err != nil {
		return fmt.Errorf("failed to get current language stats: %w", err)
	}

	newGamesPlayed := currentStats.GamesPlayed + 1
	newGamesWon := currentStats.GamesWon + gamesWon
	newTotalGuesses := currentStats.TotalGuesses + update.GuessCount
	newTotalScore := currentStats.TotalScore + update.ScoreEarned
	newCurrentStreak := currentStats.CurrentStreak
	newBestStreak := currentStats.BestStreak

	if update.IsWin {
		newCurrentStreak++
		if newCurrentStreak > newBestStreak {
			newBestStreak = newCurrentStreak
		}
	} else {
		newCurrentStreak = 0
	}

	newFastestWin := currentStats.FastestWinSeconds
	if update.IsWin && update.TimeSeconds != nil {
		if newFastestWin == nil || *update.TimeSeconds < *newFastestWin {
			newFastestWin = update.TimeSeconds
		}
	}

	newFewestGuesses := currentStats.FewestGuessesWin
	if update.IsWin {
		if newFewestGuesses == nil || update.GuessCount < *newFewestGuesses {
			newFewestGuesses = &update.GuessCount
		}
	}

	updateQuery := fmt.Sprintf(`
		UPDATE %s
		SET games_played = $3,
		    games_won = $4,
		    total_guesses = $5,
		    best_streak = $6,
		    current_streak = $7,
		    total_score = $8,
		    fastest_win_seconds = $9,
		    fewest_guesses_win = $10
		WHERE user_id = $1 AND language = $2
	`, languageStatsTable)

	_, err = tx.Exec(updateQuery,
		update.UserId,
		update.Language,
		newGamesPlayed,
		newGamesWon,
		newTotalGuesses,
		newBestStreak,
		newCurrentStreak,
		newTotalScore,
		newFastestWin,
		newFewestGuesses,
	)

	return err
}
