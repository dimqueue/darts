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
	`, userStatisticsView)

	if err := r.db.Get(&stats, query, userId); err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *StatisticsPostgres) CreateGlobalStreaks(q Querier, userId int64) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (user_id)
		VALUES ($1)
	`, globalStreaksTable)

	_, err := q.Exec(query, userId)
	return err
}

func (r *StatisticsPostgres) GetGlobalStreaks(q Querier, userId int64, forUpdate bool) (*model.UserGlobalStreaks, error) {
	var streaks model.UserGlobalStreaks

	query := fmt.Sprintf(`
		SELECT user_id, current_streak, best_streak, last_game_at, updated_at
		FROM %s
		WHERE user_id = $1
	`, globalStreaksTable)

	if forUpdate {
		query += " FOR UPDATE"
	}

	err := q.Get(&streaks, query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &streaks, nil
}

func (r *StatisticsPostgres) CreateGlobalStreaksWithData(q Querier, streaks *model.UserGlobalStreaks) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, current_streak, best_streak, last_game_at)
		VALUES ($1, $2, $3, NOW())
	`, globalStreaksTable)

	_, err := q.Exec(query, streaks.UserId, streaks.CurrentStreak, streaks.BestStreak)
	return err
}

func (r *StatisticsPostgres) UpdateGlobalStreaks(q Querier, streaks *model.UserGlobalStreaks) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET current_streak = $2,
		    best_streak = $3,
		    last_game_at = NOW()
		WHERE user_id = $1
	`, globalStreaksTable)

	_, err := q.Exec(query, streaks.UserId, streaks.CurrentStreak, streaks.BestStreak)
	return err
}

func (r *StatisticsPostgres) GetLanguageStats(q Querier, userId int64, language string, forUpdate bool) (*model.UserLanguageStats, error) {
	var stats model.UserLanguageStats

	query := fmt.Sprintf(`
		SELECT id, user_id, language, games_played, games_won,
		       total_guesses, average_guesses, best_streak, current_streak, total_score,
		       fastest_win_seconds, fewest_guesses_win, updated_at
		FROM %s
		WHERE user_id = $1 AND language = $2
	`, languageStatsTable)

	if forUpdate {
		query += " FOR UPDATE"
	}

	err := q.Get(&stats, query, userId, language)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
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

func (r *StatisticsPostgres) CreateLanguageStats(q Querier, stats *model.UserLanguageStats) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, language, games_played, games_won, total_guesses,
		                best_streak, current_streak, total_score, fastest_win_seconds, fewest_guesses_win)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, languageStatsTable)

	_, err := q.Exec(query,
		stats.UserId,
		stats.Language,
		stats.GamesPlayed,
		stats.GamesWon,
		stats.TotalGuesses,
		stats.BestStreak,
		stats.CurrentStreak,
		stats.TotalScore,
		stats.FastestWinSeconds,
		stats.FewestGuessesWin,
	)
	return err
}

func (r *StatisticsPostgres) UpdateLanguageStats(q Querier, stats *model.UserLanguageStats) error {
	query := fmt.Sprintf(`
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

	_, err := q.Exec(query,
		stats.UserId,
		stats.Language,
		stats.GamesPlayed,
		stats.GamesWon,
		stats.TotalGuesses,
		stats.BestStreak,
		stats.CurrentStreak,
		stats.TotalScore,
		stats.FastestWinSeconds,
		stats.FewestGuessesWin,
	)
	return err
}
