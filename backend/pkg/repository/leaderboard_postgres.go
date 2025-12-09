package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/dimqueue/darts/pkg/config"
	"github.com/dimqueue/darts/pkg/model"
	"github.com/jmoiron/sqlx"
)

const (
	globalLeaderboardByLanguageQuery = `
		WITH ranked_users AS (
			SELECT
			    RANK() OVER (ORDER BY uls.total_score DESC, uls.games_won DESC, uls.average_guesses ASC) as rank,
			    u.id as user_id,
			    u.username,
			    u.name,
			    up.avatar_url,
			    up.country_code,
			    uls.total_score,
			    uls.games_won as total_wins,
			    uls.games_played as total_games,
			    uls.best_streak as best_win_streak,
			    uls.average_guesses,
			    CASE WHEN uls.games_played > 0
			        THEN ROUND((uls.games_won::DECIMAL / uls.games_played * 100), 2)
			        ELSE 0
			    END as win_rate
			FROM user_language_stats uls
			INNER JOIN users u ON uls.user_id = u.id
			LEFT JOIN user_profiles up ON u.id = up.user_id
			WHERE u.is_active = true
			    AND uls.language = $1
			    AND uls.games_played >= %d
		)
		SELECT rank, user_id, username, name, avatar_url, country_code,
		       total_score, total_wins, total_games, best_win_streak,
		       average_guesses, win_rate
		FROM ranked_users
		ORDER BY rank
		LIMIT $2 OFFSET $3`

	globalUserRankByLanguageQuery = `
		WITH ranked_users AS (
		    SELECT
		        user_id,
		        RANK() OVER (ORDER BY total_score DESC, games_won DESC, average_guesses ASC) as rank
		    FROM user_language_stats
		    WHERE language = $2
		        AND games_played >= %d
		)
		SELECT rank
		FROM ranked_users
		WHERE user_id = $1`

	periodStreaksCTE = `
		period_streaks AS (
			SELECT user_id, COALESCE(MAX(streak_len), 0) as best_period_streak
			FROM (
				SELECT user_id, COUNT(*) as streak_len
				FROM (
					SELECT
						g.user_id,
						g.status,
						SUM(CASE WHEN g.status != 'won' THEN 1 ELSE 0 END)
							OVER (PARTITION BY g.user_id ORDER BY g.started_at) as loss_group
					FROM games g
					WHERE g.started_at >= $1
						AND g.status != 'in_progress'
						%s
				) grouped
				WHERE status = 'won'
				GROUP BY user_id, loss_group
			) streaks
			GROUP BY user_id
		)`

	periodStatsCTE = `
		period_stats AS (
			SELECT
			    u.id as user_id,
			    u.username,
			    u.name,
			    up.avatar_url,
			    up.country_code,
			    SUM(CASE WHEN g.status = 'won' THEN 100 ELSE 0 END) as total_score,
			    COUNT(CASE WHEN g.status = 'won' THEN 1 END) as total_wins,
			    COUNT(CASE WHEN g.status != 'in_progress' THEN 1 END) as total_games,
			    COALESCE(
			        ROUND(SUM(CASE WHEN g.status != 'in_progress' THEN COALESCE(gg.guess_count, 0) ELSE 0 END)::DECIMAL /
			              NULLIF(COUNT(CASE WHEN g.status != 'in_progress' THEN 1 END), 0), 2),
			        0
			    ) as average_guesses,
			    CASE WHEN COUNT(CASE WHEN g.status != 'in_progress' THEN 1 END) > 0
			        THEN ROUND((COUNT(CASE WHEN g.status = 'won' THEN 1 END)::DECIMAL /
			                    COUNT(CASE WHEN g.status != 'in_progress' THEN 1 END) * 100), 2)
			        ELSE 0
			    END as win_rate
			FROM users u
			INNER JOIN games g ON u.id = g.user_id
			LEFT JOIN game_guesses gg ON g.id = gg.game_id
			LEFT JOIN user_profiles up ON u.id = up.user_id
			WHERE u.is_active = true
			    AND g.started_at >= $1
			    %s
			GROUP BY u.id, u.username, u.name, up.avatar_url, up.country_code
			HAVING COUNT(CASE WHEN g.status != 'in_progress' THEN 1 END) >= %d
		)`

	periodScoresCTE = `
		period_scores AS (
		    SELECT
		        u.id as user_id,
		        SUM(CASE WHEN g.status = 'won' THEN 100 ELSE 0 END) as total_score,
		        COUNT(CASE WHEN g.status = 'won' THEN 1 END) as total_wins,
		        COALESCE(
		            ROUND(SUM(CASE WHEN g.status != 'in_progress' THEN COALESCE(gg.guess_count, 0) ELSE 0 END)::DECIMAL /
		                  NULLIF(COUNT(CASE WHEN g.status != 'in_progress' THEN 1 END), 0), 2),
		            0
		        ) as average_guesses
		    FROM users u
		    INNER JOIN games g ON u.id = g.user_id
		    LEFT JOIN game_guesses gg ON g.id = gg.game_id
		    WHERE u.is_active = true
		        AND g.started_at >= $1
		        AND g.status != 'in_progress'
		        %s
		    GROUP BY u.id
		    HAVING COUNT(g.id) >= %d
		)`

	gameGuessesCTE = `
		game_guesses AS (
			SELECT game_id, COUNT(*) as guess_count
			FROM guesses
			GROUP BY game_id
		)`

	leaderboardSelectFields = `rank, user_id, username, name, avatar_url, country_code,
		       total_score, total_wins, total_games, best_win_streak,
		       average_guesses, win_rate`
)

type LeaderboardPostgres struct {
	db *sqlx.DB
}

func NewLeaderboardPostgres(db *sqlx.DB) *LeaderboardPostgres {
	return &LeaderboardPostgres{db: db}
}

func buildLanguageFilter(periodStart time.Time, language *string) (langFilter string, args []interface{}, nextArgNum int) {
	args = []interface{}{periodStart}
	nextArgNum = 2

	if language != nil {
		langFilter = fmt.Sprintf("AND g.language = $%d", nextArgNum)
		args = append(args, *language)
		nextArgNum++
	}
	return
}

func (r *LeaderboardPostgres) GetGlobalLeaderboard(ctx context.Context, limit, offset int) ([]model.LeaderboardUser, error) {
	users := make([]model.LeaderboardUser, 0)

	query := fmt.Sprintf(`
		SELECT rank, user_id, username, name, avatar_url, country_code,
		       total_score, total_wins, total_games, best_win_streak,
		       average_guesses, win_rate
		FROM %s
		ORDER BY rank
		LIMIT $1 OFFSET $2
	`, globalLeaderboardView)

	err := r.db.SelectContext(ctx, &users, query, limit, offset)
	return users, err
}

func (r *LeaderboardPostgres) GetGlobalLeaderboardByLanguage(ctx context.Context, language string, limit, offset int) ([]model.LeaderboardUser, error) {
	users := make([]model.LeaderboardUser, 0)
	query := fmt.Sprintf(globalLeaderboardByLanguageQuery, config.LeaderboardMinGamesGlobal)
	err := r.db.SelectContext(ctx, &users, query, language, limit, offset)
	return users, err
}

func (r *LeaderboardPostgres) GetGlobalLeaderboardCount(ctx context.Context) (int, error) {
	var count int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, globalLeaderboardView)
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}

func (r *LeaderboardPostgres) GetGlobalLeaderboardByLanguageCount(ctx context.Context, language string) (int, error) {
	var count int

	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM user_language_stats uls
		INNER JOIN users u ON uls.user_id = u.id
		WHERE u.is_active = true
		    AND uls.language = $1
		    AND uls.games_played >= %d
	`, config.LeaderboardMinGamesGlobal)

	err := r.db.GetContext(ctx, &count, query, language)
	return count, err
}

func (r *LeaderboardPostgres) GetGlobalUserRank(ctx context.Context, userId int64) (*int, error) {
	var rank *int

	query := fmt.Sprintf(`
		SELECT rank
		FROM %s
		WHERE user_id = $1
	`, globalLeaderboardView)

	err := r.db.GetContext(ctx, &rank, query, userId)
	if err != nil {
		return nil, nil
	}
	return rank, nil
}

func (r *LeaderboardPostgres) GetGlobalUserRankByLanguage(ctx context.Context, userId int64, language string) (*int, error) {
	var rank *int
	query := fmt.Sprintf(globalUserRankByLanguageQuery, config.LeaderboardMinGamesGlobal)
	err := r.db.GetContext(ctx, &rank, query, userId, language)
	if err != nil {
		return nil, nil
	}
	return rank, nil
}

func (r *LeaderboardPostgres) GetPeriodLeaderboard(ctx context.Context, periodStart time.Time, language *string, limit, offset int) ([]model.LeaderboardUser, error) {
	users := make([]model.LeaderboardUser, 0)
	langFilter, args, argNum := buildLanguageFilter(periodStart, language)

	args = append(args, limit, offset)

	streaksCTE := fmt.Sprintf(periodStreaksCTE, langFilter)
	statsCTE := fmt.Sprintf(periodStatsCTE, langFilter, config.LeaderboardMinGamesPeriod)

	query := fmt.Sprintf(`
		WITH %s,
		%s,
		%s,
		ranked_users AS (
			SELECT
			    RANK() OVER (ORDER BY ps.total_score DESC, ps.total_wins DESC, ps.average_guesses ASC) as rank,
			    ps.user_id, ps.username, ps.name, ps.avatar_url, ps.country_code,
			    ps.total_score, ps.total_wins, ps.total_games,
			    COALESCE(pst.best_period_streak, 0) as best_win_streak,
			    ps.average_guesses, ps.win_rate
			FROM period_stats ps
			LEFT JOIN period_streaks pst ON ps.user_id = pst.user_id
		)
		SELECT %s
		FROM ranked_users
		ORDER BY rank
		LIMIT $%d OFFSET $%d
	`, gameGuessesCTE, streaksCTE, statsCTE, leaderboardSelectFields, argNum, argNum+1)

	err := r.db.SelectContext(ctx, &users, query, args...)
	return users, err
}

func (r *LeaderboardPostgres) GetPeriodLeaderboardCount(ctx context.Context, periodStart time.Time, language *string) (int, error) {
	var count int

	langFilter := ""
	args := []interface{}{periodStart}

	if language != nil {
		langFilter = "AND g.language = $2"
		args = append(args, *language)
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) FROM (
			SELECT u.id
			FROM users u
			INNER JOIN games g ON u.id = g.user_id
			WHERE u.is_active = true
			    AND g.started_at >= $1
			    AND g.status != 'in_progress'
			    %s
			GROUP BY u.id
			HAVING COUNT(g.id) >= %d
		) AS eligible_users
	`, langFilter, config.LeaderboardMinGamesPeriod)

	err := r.db.GetContext(ctx, &count, query, args...)
	return count, err
}

func (r *LeaderboardPostgres) GetPeriodUserRank(ctx context.Context, userId int64, periodStart time.Time, language *string) (*int, error) {
	var rank int
	langFilter, args, argNum := buildLanguageFilter(periodStart, language)

	args = append(args, userId)

	scoresCTE := fmt.Sprintf(periodScoresCTE, langFilter, config.LeaderboardMinGamesPeriod)

	query := fmt.Sprintf(`
		WITH %s,
		%s,
		ranked AS (
		    SELECT user_id, RANK() OVER (ORDER BY total_score DESC, total_wins DESC, average_guesses ASC) as rank
		    FROM period_scores
		)
		SELECT rank FROM ranked WHERE user_id = $%d
	`, gameGuessesCTE, scoresCTE, argNum)

	err := r.db.GetContext(ctx, &rank, query, args...)
	if err != nil {
		return nil, nil
	}
	return &rank, nil
}
