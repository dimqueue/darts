package repository

import (
	"time"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(q Querier, user model.User) (int64, error)
	GetUser(username, password string) (model.User, error)
}

type Game interface {
	CreateGame(game *model.Game) (int64, error)
	GetAllGames(userId int64) ([]model.Game, error)
	GetGameById(q Querier, gameId int64, forUpdate bool) (*model.Game, error)
	UpdateGameStatus(q Querier, gameId int64, status string) error
	GetExpiredGames() ([]model.Game, error)
	CreateGuess(q Querier, guess *model.Guess) error
	GetAllGuessByGame(gameId int64) ([]model.Guess, error)
	CountGuessesByGame(q Querier, gameId int64) (int, error)
	GuessExists(q Querier, gameId int64, guessWord string) (bool, error)
}

type Word interface {
	GetWordById(wordId int64) (*model.Word, error)
	GetRandomWordByLanguage(language string) (*model.Word, error)
	GetWordCountByLanguage(language string) (int, error)
}

type Profile interface {
	GetProfile(userId int64) (*model.UserProfile, error)
	CreateProfile(q Querier, profile *model.UserProfile) error
	UpdateProfile(userId int64, input model.UpdateProfileInput) error
	GetSettings(userId int64) (*model.UserSettings, error)
	CreateSettings(q Querier, userId int64) error
	UpdateSettings(userId int64, input model.UpdateSettingsInput) error
	GetProfileSummary(userId int64) (*model.UserProfileSummary, error)
	GetProfileByUsername(username string) (*model.UserProfileSummary, error)
}

type Statistics interface {
	GetStatistics(userId int64) (*model.UserStatistics, error)

	CreateGlobalStreaks(q Querier, userId int64) error
	GetGlobalStreaks(q Querier, userId int64, forUpdate bool) (*model.UserGlobalStreaks, error)
	CreateGlobalStreaksWithData(q Querier, streaks *model.UserGlobalStreaks) error
	UpdateGlobalStreaks(q Querier, streaks *model.UserGlobalStreaks) error

	GetLanguageStats(q Querier, userId int64, language string, forUpdate bool) (*model.UserLanguageStats, error)
	GetAllLanguageStats(userId int64) ([]model.UserLanguageStats, error)
	CreateLanguageStats(q Querier, stats *model.UserLanguageStats) error
	UpdateLanguageStats(q Querier, stats *model.UserLanguageStats) error
}

type Leaderboard interface {
	GetGlobalLeaderboard(limit, offset int) ([]model.LeaderboardUser, error)
	GetGlobalLeaderboardByLanguage(language string, limit, offset int) ([]model.LeaderboardUser, error)
	GetGlobalLeaderboardCount() (int, error)
	GetGlobalLeaderboardByLanguageCount(language string) (int, error)
	GetGlobalUserRank(userId int64) (*int, error)
	GetGlobalUserRankByLanguage(userId int64, language string) (*int, error)

	GetPeriodLeaderboard(periodStart time.Time, language *string, limit, offset int) ([]model.LeaderboardUser, error)
	GetPeriodLeaderboardCount(periodStart time.Time, language *string) (int, error)
	GetPeriodUserRank(userId int64, periodStart time.Time, language *string) (*int, error)
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
