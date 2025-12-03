package model

import "time"

type UserStatistics struct {
	UserId            int64      `json:"user_id" db:"user_id"`
	TotalGames        int        `json:"total_games" db:"total_games"`
	TotalWins         int        `json:"total_wins" db:"total_wins"`
	TotalLosses       int        `json:"total_losses" db:"total_losses"`
	CurrentWinStreak  int        `json:"current_win_streak" db:"current_win_streak"`
	BestWinStreak     int        `json:"best_win_streak" db:"best_win_streak"`
	TotalGuesses      int        `json:"total_guesses" db:"total_guesses"`
	AverageGuesses    float64    `json:"average_guesses" db:"average_guesses"`
	FastestWinSeconds *int       `json:"fastest_win_seconds,omitempty" db:"fastest_win_seconds"`
	FewestGuessesWin  *int       `json:"fewest_guesses_win,omitempty" db:"fewest_guesses_win"`
	TotalScore        int        `json:"total_score" db:"total_score"`
	LastGameAt        *time.Time `json:"last_game_at,omitempty" db:"last_game_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

type UserLanguageStats struct {
	Id                int64     `json:"id" db:"id"`
	UserId            int64     `json:"user_id" db:"user_id"`
	Language          string    `json:"language" db:"language"`
	GamesPlayed       int       `json:"games_played" db:"games_played"`
	GamesWon          int       `json:"games_won" db:"games_won"`
	TotalGuesses      int       `json:"total_guesses" db:"total_guesses"`
	AverageGuesses    float64   `json:"average_guesses" db:"average_guesses"`
	CurrentStreak     int       `json:"current_streak" db:"current_streak"`
	BestStreak        int       `json:"best_streak" db:"best_streak"`
	TotalScore        int       `json:"total_score" db:"total_score"`
	FastestWinSeconds *int      `json:"fastest_win_seconds,omitempty" db:"fastest_win_seconds"`
	FewestGuessesWin  *int      `json:"fewest_guesses_win,omitempty" db:"fewest_guesses_win"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

type UserGlobalStreaks struct {
	UserId        int64      `json:"user_id" db:"user_id"`
	CurrentStreak int        `json:"current_streak" db:"current_streak"`
	BestStreak    int        `json:"best_streak" db:"best_streak"`
	LastGameAt    *time.Time `json:"last_game_at,omitempty" db:"last_game_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type StatisticsUpdate struct {
	UserId      int64
	Language    string
	IsWin       bool
	GuessCount  int
	TimeSeconds *int
	ScoreEarned int
}
