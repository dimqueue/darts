package repository

import (
	"context"
	"time"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(ctx context.Context, q Querier, user model.User) (int64, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
}

type Game interface {
	CreateGame(ctx context.Context, game *model.Game) (int64, error)
	GetAllGames(ctx context.Context, userId int64) ([]model.Game, error)
	GetGameById(ctx context.Context, q Querier, gameId int64, forUpdate bool) (*model.Game, error)
	UpdateGameStatus(ctx context.Context, q Querier, gameId int64, status string) error
	GetExpiredGames(ctx context.Context) ([]model.Game, error)
	CreateGuess(ctx context.Context, q Querier, guess *model.Guess) error
	GetAllGuessByGame(ctx context.Context, gameId int64) ([]model.Guess, error)
	CountGuessesByGame(ctx context.Context, q Querier, gameId int64) (int, error)
	GuessExists(ctx context.Context, q Querier, gameId int64, guessWord string) (bool, error)
}

type Word interface {
	GetWordById(ctx context.Context, wordId int64) (*model.Word, error)
	GetRandomWordByLanguage(ctx context.Context, language string) (*model.Word, error)
	GetWordCountByLanguage(ctx context.Context, language string) (int, error)
}

type Profile interface {
	GetProfile(ctx context.Context, userId int64) (*model.UserProfile, error)
	CreateProfile(ctx context.Context, q Querier, profile *model.UserProfile) error
	UpdateProfile(ctx context.Context, userId int64, input model.UpdateProfileInput) error
	GetSettings(ctx context.Context, userId int64) (*model.UserSettings, error)
	CreateSettings(ctx context.Context, q Querier, userId int64) error
	UpdateSettings(ctx context.Context, userId int64, input model.UpdateSettingsInput) error
	GetProfileSummary(ctx context.Context, userId int64) (*model.UserProfileSummary, error)
	GetProfileByUsername(ctx context.Context, username string) (*model.UserProfileSummary, error)
}

type Statistics interface {
	GetStatistics(ctx context.Context, userId int64) (*model.UserStatistics, error)

	CreateGlobalStreaks(ctx context.Context, q Querier, userId int64) error
	GetGlobalStreaks(ctx context.Context, q Querier, userId int64, forUpdate bool) (*model.UserGlobalStreaks, error)
	CreateGlobalStreaksWithData(ctx context.Context, q Querier, streaks *model.UserGlobalStreaks) error
	UpdateGlobalStreaks(ctx context.Context, q Querier, streaks *model.UserGlobalStreaks) error

	GetLanguageStats(ctx context.Context, q Querier, userId int64, language string, forUpdate bool) (*model.UserLanguageStats, error)
	GetAllLanguageStats(ctx context.Context, userId int64) ([]model.UserLanguageStats, error)
	CreateLanguageStats(ctx context.Context, q Querier, stats *model.UserLanguageStats) error
	UpdateLanguageStats(ctx context.Context, q Querier, stats *model.UserLanguageStats) error
}

type Leaderboard interface {
	GetGlobalLeaderboard(ctx context.Context, limit, offset int) ([]model.LeaderboardUser, error)
	GetGlobalLeaderboardByLanguage(ctx context.Context, language string, limit, offset int) ([]model.LeaderboardUser, error)
	GetGlobalLeaderboardCount(ctx context.Context) (int, error)
	GetGlobalLeaderboardByLanguageCount(ctx context.Context, language string) (int, error)
	GetGlobalUserRank(ctx context.Context, userId int64) (*int, error)
	GetGlobalUserRankByLanguage(ctx context.Context, userId int64, language string) (*int, error)

	GetPeriodLeaderboard(ctx context.Context, periodStart time.Time, language *string, limit, offset int) ([]model.LeaderboardUser, error)
	GetPeriodLeaderboardCount(ctx context.Context, periodStart time.Time, language *string) (int, error)
	GetPeriodUserRank(ctx context.Context, userId int64, periodStart time.Time, language *string) (*int, error)
}

type Repository struct {
	Authorization
	Game
	Word
	Profile
	Statistics
	Leaderboard
	TxManager *TransactionManager
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Game:          NewGamePostgres(db),
		Word:          NewWordPostgres(db),
		Profile:       NewProfilePostgres(db),
		Statistics:    NewStatisticsPostgres(db),
		Leaderboard:   NewLeaderboardPostgres(db),
		TxManager:     NewTransactionManager(db),
	}
}
