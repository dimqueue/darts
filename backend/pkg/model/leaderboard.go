package model

import "time"

type LeaderboardUser struct {
	Rank           int     `json:"rank" db:"rank"`
	UserId         int64   `json:"user_id" db:"user_id"`
	Username       string  `json:"username" db:"username"`
	Name           string  `json:"name" db:"name"`
	AvatarURL      *string `json:"avatar_url,omitempty" db:"avatar_url"`
	CountryCode    *string `json:"country_code,omitempty" db:"country_code"`
	TotalScore     int     `json:"total_score" db:"total_score"`
	TotalWins      int     `json:"total_wins" db:"total_wins"`
	TotalGames     int     `json:"total_games" db:"total_games"`
	BestWinStreak  int     `json:"best_win_streak" db:"best_win_streak"`
	AverageGuesses float64 `json:"average_guesses" db:"average_guesses"`
	WinRate        float64 `json:"win_rate" db:"win_rate"`
}

type LeaderboardResponse struct {
	LeaderboardType string            `json:"leaderboard_type"`
	Language        *string           `json:"language,omitempty"`
	PeriodStart     *time.Time        `json:"period_start,omitempty"`
	PeriodEnd       *time.Time        `json:"period_end,omitempty"`
	Users           []LeaderboardUser `json:"users"`
	Total           int               `json:"total"`
	CurrentUserRank *int              `json:"current_user_rank,omitempty"`
}

type LeaderboardQuery struct {
	Type     string  // "global", "daily", "weekly", "monthly"
	Language *string // Optional - filter by language for any type
	Limit    int     // Number of entries to return (default: 50, max: 100)
	Offset   int     // Pagination offset
}

type UserRanks struct {
	GlobalRank  *int `json:"global_rank" db:"global_rank"`
	DailyRank   *int `json:"daily_rank" db:"daily_rank"`
	WeeklyRank  *int `json:"weekly_rank" db:"weekly_rank"`
	MonthlyRank *int `json:"monthly_rank" db:"monthly_rank"`
}
