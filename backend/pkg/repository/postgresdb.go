package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable         = "users"
	gamesTable         = "games"
	guessesTable       = "guesses"
	wordsTable         = "words"
	profilesTable      = "user_profiles"
	settingsTable      = "user_settings"
	languageStatsTable = "user_language_stats"
	globalStreaksTable = "user_global_streaks"
)

const (
	globalLeaderboardView  = "global_leaderboard_view"
	userProfileSummaryView = "user_profile_summary"
	userStatisticsView     = "user_statistics"
)

type Config struct {
	Host            string
	Port            string
	Username        string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
