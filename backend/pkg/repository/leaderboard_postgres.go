package repository

import (
	"fmt"
	"time"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/jmoiron/sqlx"
)

type LeaderboardPostgres struct {
	db *sqlx.DB
}

func NewLeaderboardPostgres(db *sqlx.DB) *LeaderboardPostgres {
	return &LeaderboardPostgres{db: db}
}

func (r *LeaderboardPostgres) GetLeaderboard(query model.LeaderboardQuery) ([]model.LeaderboardUser, error) {
	switch query.Type {
	case "global":
		if query.Language != nil {
			return r.getGlobalByLanguage(*query.Language, query.Limit, query.Offset)
		}
		return r.getGlobal(query.Limit, query.Offset)
	case "daily":
		return r.getPeriodLeaderboard(getDayStart(), query.Language, query.Limit, query.Offset)
	case "weekly":
		return r.getPeriodLeaderboard(getWeekStart(), query.Language, query.Limit, query.Offset)
	case "monthly":
		return r.getPeriodLeaderboard(getMonthStart(), query.Language, query.Limit, query.Offset)
	default:
		return nil, fmt.Errorf("invalid leaderboard type: %s", query.Type)
	}
}

func (r *LeaderboardPostgres) GetLeaderboardCount(query model.LeaderboardQuery) (int, error) {
	switch query.Type {
	case "global":
		if query.Language != nil {
			return r.getGlobalByLanguageCount(*query.Language)
		}
		return r.getGlobalCount()
	case "daily":
		return r.getPeriodLeaderboardCount(getDayStart(), query.Language)
	case "weekly":
		return r.getPeriodLeaderboardCount(getWeekStart(), query.Language)
	case "monthly":
		return r.getPeriodLeaderboardCount(getMonthStart(), query.Language)
	default:
		return 0, fmt.Errorf("invalid leaderboard type: %s", query.Type)
	}
}

func (r *LeaderboardPostgres) GetUserRank(userId int64, query model.LeaderboardQuery) (*int, error) {
	switch query.Type {
	case "global":
		if query.Language != nil {
			return r.getUserRankByLanguage(userId, *query.Language)
		}
		return r.getUserRankGlobal(userId)
	case "daily":
		return r.getUserRankPeriod(userId, getDayStart(), query.Language)
	case "weekly":
		return r.getUserRankPeriod(userId, getWeekStart(), query.Language)
	case "monthly":
		return r.getUserRankPeriod(userId, getMonthStart(), query.Language)
	default:
		return nil, fmt.Errorf("invalid leaderboard type: %s", query.Type)
	}
}

func (r *LeaderboardPostgres) getGlobal(limit, offset int) ([]model.LeaderboardUser, error) {
	users := make([]model.LeaderboardUser, 0)

	query := `
		SELECT rank, user_id, username, name, avatar_url, country_code,
		       total_score, total_wins, total_games, best_win_streak,
		       average_guesses, win_rate
		FROM global_leaderboard_view
		ORDER BY rank
		LIMIT $1 OFFSET $2
	`

	err := r.db.Select(&users, query, limit, offset)
	return users, err
}

func (r *LeaderboardPostgres) getGlobalByLanguage(language string, limit, offset int) ([]model.LeaderboardUser, error) {
	users := make([]model.LeaderboardUser, 0)

	query := `
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
			LEFT JOIN user_settings us ON u.id = us.user_id
			WHERE u.is_active = true
			    AND uls.language = $1
			    AND uls.games_played >= 3
			    AND COALESCE(us.show_stats_public, true) = true
		)
		SELECT rank, user_id, username, name, avatar_url, country_code,
		       total_score, total_wins, total_games, best_win_streak,
		       average_guesses, win_rate
		FROM ranked_users
		ORDER BY rank
		LIMIT $2 OFFSET $3
	`

	err := r.db.Select(&users, query, language, limit, offset)
	return users, err
}

func (r *LeaderboardPostgres) getGlobalCount() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM global_leaderboard_view`
	err := r.db.Get(&count, query)
	return count, err
}

func (r *LeaderboardPostgres) getGlobalByLanguageCount(language string) (int, error) {
	var count int

	query := `
		SELECT COUNT(*)
		FROM user_language_stats uls
		INNER JOIN users u ON uls.user_id = u.id
		LEFT JOIN user_settings us ON u.id = us.user_id
		WHERE u.is_active = true
		    AND uls.language = $1
		    AND uls.games_played >= 3
		    AND COALESCE(us.show_stats_public, true) = true
	`

	err := r.db.Get(&count, query, language)
	return count, err
}

func (r *LeaderboardPostgres) getUserRankGlobal(userId int64) (*int, error) {
	var rank *int

	query := `
		SELECT rank
		FROM global_leaderboard_view
		WHERE user_id = $1
	`

	err := r.db.Get(&rank, query, userId)
	if err != nil {
		return nil, nil
	}
	return rank, nil
}

func (r *LeaderboardPostgres) getUserRankByLanguage(userId int64, language string) (*int, error) {
	var rank *int

	query := `
		WITH ranked_users AS (
		    SELECT
		        user_id,
		        RANK() OVER (ORDER BY total_score DESC, games_won DESC, average_guesses ASC) as rank
		    FROM user_language_stats
		    WHERE language = $2
		        AND games_played >= 3
		)
		SELECT rank
		FROM ranked_users
		WHERE user_id = $1
	`

	err := r.db.Get(&rank, query, userId, language)
	if err != nil {
		return nil, nil
	}
	return rank, nil
}

func (r *LeaderboardPostgres) getPeriodLeaderboard(periodStart time.Time, language *string, limit, offset int) ([]model.LeaderboardUser, error) {
	users := make([]model.LeaderboardUser, 0)

	langFilter := ""
	args := []interface{}{periodStart}
	argNum := 2

	if language != nil {
		langFilter = fmt.Sprintf("AND g.language = $%d", argNum)
		args = append(args, *language)
		argNum++
	}

	args = append(args, limit, offset)

	query := fmt.Sprintf(`
		WITH game_guesses AS (
			SELECT game_id, COUNT(*) as guess_count
			FROM guesses
			GROUP BY game_id
		),
		period_stats AS (
			SELECT
			    u.id as user_id,
			    u.username,
			    u.name,
			    up.avatar_url,
			    up.country_code,
			    SUM(CASE WHEN g.status = 'won' THEN 100 ELSE 0 END) as total_score,
			    COUNT(CASE WHEN g.status = 'won' THEN 1 END) as total_wins,
			    COUNT(g.id) as total_games,
			    0 as best_win_streak,
			    COALESCE(
			        ROUND(SUM(COALESCE(gg.guess_count, 0))::DECIMAL / NULLIF(COUNT(g.id), 0), 2),
			        0
			    ) as average_guesses,
			    CASE WHEN COUNT(g.id) > 0
			        THEN ROUND((COUNT(CASE WHEN g.status = 'won' THEN 1 END)::DECIMAL / COUNT(g.id) * 100), 2)
			        ELSE 0
			    END as win_rate
			FROM users u
			INNER JOIN games g ON u.id = g.user_id
			LEFT JOIN game_guesses gg ON g.id = gg.game_id
			LEFT JOIN user_profiles up ON u.id = up.user_id
			LEFT JOIN user_settings us ON u.id = us.user_id
			WHERE u.is_active = true
			    AND g.started_at >= $1
			    %s
			    AND COALESCE(us.show_stats_public, true) = true
			GROUP BY u.id, u.username, u.name, up.avatar_url, up.country_code
			HAVING COUNT(g.id) >= 3
		),
		ranked_users AS (
			SELECT
			    RANK() OVER (ORDER BY total_score DESC, total_wins DESC, average_guesses ASC) as rank,
			    user_id, username, name, avatar_url, country_code,
			    total_score, total_wins, total_games, best_win_streak,
			    average_guesses, win_rate
			FROM period_stats
		)
		SELECT rank, user_id, username, name, avatar_url, country_code,
		       total_score, total_wins, total_games, best_win_streak,
		       average_guesses, win_rate
		FROM ranked_users
		ORDER BY rank
		LIMIT $%d OFFSET $%d
	`, langFilter, argNum, argNum+1)

	err := r.db.Select(&users, query, args...)
	return users, err
}

func (r *LeaderboardPostgres) getPeriodLeaderboardCount(periodStart time.Time, language *string) (int, error) {
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
			LEFT JOIN user_settings us ON u.id = us.user_id
			WHERE u.is_active = true
			    AND g.started_at >= $1
			    %s
			    AND COALESCE(us.show_stats_public, true) = true
			GROUP BY u.id
			HAVING COUNT(g.id) >= 3
		) AS eligible_users
	`, langFilter)

	err := r.db.Get(&count, query, args...)
	return count, err
}

func (r *LeaderboardPostgres) getUserRankPeriod(userId int64, periodStart time.Time, language *string) (*int, error) {
	var rank int

	langFilter := ""
	args := []interface{}{periodStart}
	argNum := 2

	if language != nil {
		langFilter = fmt.Sprintf("AND g.language = $%d", argNum)
		args = append(args, *language)
		argNum++
	}

	args = append(args, userId)

	query := fmt.Sprintf(`
		WITH period_scores AS (
		    SELECT
		        u.id as user_id,
		        SUM(CASE WHEN g.status = 'won' THEN 100 ELSE 0 END) as total_score,
		        COUNT(CASE WHEN g.status = 'won' THEN 1 END) as total_wins
		    FROM users u
		    INNER JOIN games g ON u.id = g.user_id
		    LEFT JOIN user_settings us ON u.id = us.user_id
		    WHERE u.is_active = true
		        AND g.started_at >= $1
		        %s
		        AND COALESCE(us.show_stats_public, true) = true
		    GROUP BY u.id
		    HAVING COUNT(g.id) >= 3
		),
		ranked AS (
		    SELECT user_id, RANK() OVER (ORDER BY total_score DESC, total_wins DESC) as rank
		    FROM period_scores
		)
		SELECT rank FROM ranked WHERE user_id = $%d
	`, langFilter, argNum)

	err := r.db.Get(&rank, query, args...)
	if err != nil {
		return nil, nil
	}
	return &rank, nil
}

func (r *LeaderboardPostgres) GetAllUserRanks(userId int64) (*model.UserRanks, error) {
	dayStart := getDayStart()
	weekStart := getWeekStart()
	monthStart := getMonthStart()

	query := `
		WITH global_rank AS (
			SELECT rank FROM global_leaderboard_view WHERE user_id = $1
		),
		daily_scores AS (
			SELECT
				u.id as user_id,
				SUM(CASE WHEN g.status = 'won' THEN 100 ELSE 0 END) as total_score,
				COUNT(CASE WHEN g.status = 'won' THEN 1 END) as total_wins
			FROM users u
			INNER JOIN games g ON u.id = g.user_id
			LEFT JOIN user_settings us ON u.id = us.user_id
			WHERE u.is_active = true
				AND g.started_at >= $2
				AND COALESCE(us.show_stats_public, true) = true
			GROUP BY u.id
			HAVING COUNT(g.id) >= 3
		),
		daily_ranked AS (
			SELECT user_id, RANK() OVER (ORDER BY total_score DESC, total_wins DESC) as rank
			FROM daily_scores
		),
		weekly_scores AS (
			SELECT
				u.id as user_id,
				SUM(CASE WHEN g.status = 'won' THEN 100 ELSE 0 END) as total_score,
				COUNT(CASE WHEN g.status = 'won' THEN 1 END) as total_wins
			FROM users u
			INNER JOIN games g ON u.id = g.user_id
			LEFT JOIN user_settings us ON u.id = us.user_id
			WHERE u.is_active = true
				AND g.started_at >= $3
				AND COALESCE(us.show_stats_public, true) = true
			GROUP BY u.id
			HAVING COUNT(g.id) >= 3
		),
		weekly_ranked AS (
			SELECT user_id, RANK() OVER (ORDER BY total_score DESC, total_wins DESC) as rank
			FROM weekly_scores
		),
		monthly_scores AS (
			SELECT
				u.id as user_id,
				SUM(CASE WHEN g.status = 'won' THEN 100 ELSE 0 END) as total_score,
				COUNT(CASE WHEN g.status = 'won' THEN 1 END) as total_wins
			FROM users u
			INNER JOIN games g ON u.id = g.user_id
			LEFT JOIN user_settings us ON u.id = us.user_id
			WHERE u.is_active = true
				AND g.started_at >= $4
				AND COALESCE(us.show_stats_public, true) = true
			GROUP BY u.id
			HAVING COUNT(g.id) >= 3
		),
		monthly_ranked AS (
			SELECT user_id, RANK() OVER (ORDER BY total_score DESC, total_wins DESC) as rank
			FROM monthly_scores
		)
		SELECT
			(SELECT rank FROM global_rank) as global_rank,
			(SELECT rank FROM daily_ranked WHERE user_id = $1) as daily_rank,
			(SELECT rank FROM weekly_ranked WHERE user_id = $1) as weekly_rank,
			(SELECT rank FROM monthly_ranked WHERE user_id = $1) as monthly_rank
	`

	var ranks model.UserRanks
	err := r.db.Get(&ranks, query, userId, dayStart, weekStart, monthStart)
	if err != nil {
		return &model.UserRanks{}, nil
	}

	return &ranks, nil
}

func getDayStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func getWeekStart() time.Time {
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday())+1)
	if now.Weekday() == time.Sunday {
		weekStart = now.AddDate(0, 0, -6)
	}
	return time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, time.UTC)
}

func getMonthStart() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
}
