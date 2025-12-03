package repository

import (
	"fmt"
	"time"

	"github.com/dimqueue/darts/pkg/config"
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
	if query.Type == "global" {
		if query.Language != nil {
			return r.getGlobalByLanguage(*query.Language, query.Limit, query.Offset)
		}
		return r.getGlobal(query.Limit, query.Offset)
	}

	periodStart, err := getPeriodStart(query.Type)
	if err != nil {
		return nil, err
	}
	return r.getPeriodLeaderboard(periodStart, query.Language, query.Limit, query.Offset)
}

func (r *LeaderboardPostgres) GetLeaderboardCount(query model.LeaderboardQuery) (int, error) {
	if query.Type == "global" {
		if query.Language != nil {
			return r.getGlobalByLanguageCount(*query.Language)
		}
		return r.getGlobalCount()
	}

	periodStart, err := getPeriodStart(query.Type)
	if err != nil {
		return 0, err
	}
	return r.getPeriodLeaderboardCount(periodStart, query.Language)
}

func (r *LeaderboardPostgres) GetUserRank(userId int64, query model.LeaderboardQuery) (*int, error) {
	if query.Type == "global" {
		if query.Language != nil {
			return r.getUserRankByLanguage(userId, *query.Language)
		}
		return r.getUserRankGlobal(userId)
	}

	periodStart, err := getPeriodStart(query.Type)
	if err != nil {
		return nil, err
	}
	return r.getUserRankPeriod(userId, periodStart, query.Language)
}

func (r *LeaderboardPostgres) getGlobal(limit, offset int) ([]model.LeaderboardUser, error) {
	users := make([]model.LeaderboardUser, 0)

	query := fmt.Sprintf(`
		SELECT rank, user_id, username, name, avatar_url, country_code,
		       total_score, total_wins, total_games, best_win_streak,
		       average_guesses, win_rate
		FROM %s
		ORDER BY rank
		LIMIT $1 OFFSET $2
	`, globalLeaderboardView)

	err := r.db.Select(&users, query, limit, offset)
	return users, err
}

func (r *LeaderboardPostgres) getGlobalByLanguage(language string, limit, offset int) ([]model.LeaderboardUser, error) {
	users := make([]model.LeaderboardUser, 0)

	query := fmt.Sprintf(`
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
		LIMIT $2 OFFSET $3
	`, config.LeaderboardMinGamesGlobal)

	err := r.db.Select(&users, query, language, limit, offset)
	return users, err
}

func (r *LeaderboardPostgres) getGlobalCount() (int, error) {
	var count int
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, globalLeaderboardView)
	err := r.db.Get(&count, query)
	return count, err
}

func (r *LeaderboardPostgres) getGlobalByLanguageCount(language string) (int, error) {
	var count int

	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM user_language_stats uls
		INNER JOIN users u ON uls.user_id = u.id
		WHERE u.is_active = true
		    AND uls.language = $1
		    AND uls.games_played >= %d
	`, config.LeaderboardMinGamesGlobal)

	err := r.db.Get(&count, query, language)
	return count, err
}

func (r *LeaderboardPostgres) getUserRankGlobal(userId int64) (*int, error) {
	var rank *int

	query := fmt.Sprintf(`
		SELECT rank
		FROM %s
		WHERE user_id = $1
	`, globalLeaderboardView)

	err := r.db.Get(&rank, query, userId)
	if err != nil {
		return nil, nil
	}
	return rank, nil
}

func (r *LeaderboardPostgres) getUserRankByLanguage(userId int64, language string) (*int, error) {
	var rank *int

	query := fmt.Sprintf(`
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
		WHERE user_id = $1
	`, config.LeaderboardMinGamesGlobal)

	err := r.db.Get(&rank, query, userId, language)
	if err != nil {
		return nil, nil
	}
	return rank, nil
}

func (r *LeaderboardPostgres) getPeriodLeaderboard(periodStart time.Time, language *string, limit, offset int) ([]model.LeaderboardUser, error) {
	users := make([]model.LeaderboardUser, 0)

	langFilter := ""
	langFilterStreak := ""
	args := []interface{}{periodStart}
	argNum := 2

	if language != nil {
		langFilter = fmt.Sprintf("AND g.language = $%d", argNum)
		langFilterStreak = fmt.Sprintf("AND g.language = $%d", argNum)
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
		),
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
		SELECT rank, user_id, username, name, avatar_url, country_code,
		       total_score, total_wins, total_games, best_win_streak,
		       average_guesses, win_rate
		FROM ranked_users
		ORDER BY rank
		LIMIT $%d OFFSET $%d
	`, langFilterStreak, langFilter, config.LeaderboardMinGamesPeriod, argNum, argNum+1)

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
			WHERE u.is_active = true
			    AND g.started_at >= $1
			    AND g.status != 'in_progress'
			    %s
			GROUP BY u.id
			HAVING COUNT(g.id) >= %d
		) AS eligible_users
	`, langFilter, config.LeaderboardMinGamesPeriod)

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
		WITH game_guesses AS (
		    SELECT game_id, COUNT(*) as guess_count
		    FROM guesses
		    GROUP BY game_id
		),
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
		),
		ranked AS (
		    SELECT user_id, RANK() OVER (ORDER BY total_score DESC, total_wins DESC, average_guesses ASC) as rank
		    FROM period_scores
		)
		SELECT rank FROM ranked WHERE user_id = $%d
	`, langFilter, config.LeaderboardMinGamesPeriod, argNum)

	err := r.db.Get(&rank, query, args...)
	if err != nil {
		return nil, nil
	}
	return &rank, nil
}

func (r *LeaderboardPostgres) GetAllUserRanks(userId int64) (*model.UserRanks, error) {
	var ranks model.UserRanks

	ranks.GlobalRank, _ = r.getUserRankGlobal(userId)

	dayStart, _ := getPeriodStart("daily")
	ranks.DailyRank, _ = r.getUserRankPeriod(userId, dayStart, nil)

	weekStart, _ := getPeriodStart("weekly")
	ranks.WeeklyRank, _ = r.getUserRankPeriod(userId, weekStart, nil)

	monthStart, _ := getPeriodStart("monthly")
	ranks.MonthlyRank, _ = r.getUserRankPeriod(userId, monthStart, nil)

	return &ranks, nil
}

func getPeriodStart(periodType string) (time.Time, error) {
	now := time.Now()
	switch periodType {
	case "daily":
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), nil
	case "weekly":
		weekStart := now.AddDate(0, 0, -int(now.Weekday())+1)
		if now.Weekday() == time.Sunday {
			weekStart = now.AddDate(0, 0, -6)
		}
		return time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, time.UTC), nil
	case "monthly":
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC), nil
	default:
		return time.Time{}, fmt.Errorf("invalid period type: %s", periodType)
	}
}
