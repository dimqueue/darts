package service

import (
	"context"
	"fmt"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
)

type StatsService struct {
	statsRepo repository.Statistics
	txManager *repository.TransactionManager
}

func NewStatsService(statsRepo repository.Statistics, txManager *repository.TransactionManager) *StatsService {
	return &StatsService{
		statsRepo: statsRepo,
		txManager: txManager,
	}
}

func (s *StatsService) InitializeStats(ctx context.Context, q repository.Querier, userId int64) error {
	return s.statsRepo.CreateGlobalStreaks(ctx, q, userId)
}

func (s *StatsService) CalculateNewStreaks(current *model.UserGlobalStreaks, userId int64, isWin bool) *model.UserGlobalStreaks {
	result := &model.UserGlobalStreaks{
		UserId:        userId,
		CurrentStreak: 0,
		BestStreak:    0,
	}

	if current != nil {
		result.CurrentStreak = current.CurrentStreak
		result.BestStreak = current.BestStreak
	}

	if isWin {
		result.CurrentStreak++
		if result.CurrentStreak > result.BestStreak {
			result.BestStreak = result.CurrentStreak
		}
	} else {
		result.CurrentStreak = 0
	}

	return result
}

func (s *StatsService) CalculateNewLanguageStats(
	current *model.UserLanguageStats,
	userId int64,
	language string,
	isWin bool,
	guessCount int,
	timeSeconds *int,
	scoreEarned int,
) *model.UserLanguageStats {
	result := &model.UserLanguageStats{
		UserId:            userId,
		Language:          language,
		GamesPlayed:       1,
		GamesWon:          0,
		TotalGuesses:      guessCount,
		CurrentStreak:     0,
		BestStreak:        0,
		TotalScore:        scoreEarned,
		FastestWinSeconds: nil,
		FewestGuessesWin:  nil,
	}

	if current != nil {
		result.GamesPlayed = current.GamesPlayed + 1
		result.GamesWon = current.GamesWon
		result.TotalGuesses = current.TotalGuesses + guessCount
		result.CurrentStreak = current.CurrentStreak
		result.BestStreak = current.BestStreak
		result.TotalScore = current.TotalScore + scoreEarned
		result.FastestWinSeconds = current.FastestWinSeconds
		result.FewestGuessesWin = current.FewestGuessesWin
	}

	if isWin {
		result.GamesWon++
		result.CurrentStreak++
		if result.CurrentStreak > result.BestStreak {
			result.BestStreak = result.CurrentStreak
		}

		if timeSeconds != nil {
			if result.FastestWinSeconds == nil || *timeSeconds < *result.FastestWinSeconds {
				result.FastestWinSeconds = timeSeconds
			}
		}

		if result.FewestGuessesWin == nil || guessCount < *result.FewestGuessesWin {
			result.FewestGuessesWin = &guessCount
		}
	} else {
		result.CurrentStreak = 0
	}

	return result
}

func (s *StatsService) UpdateStatsAfterGame(ctx context.Context, q repository.Querier, userId int64, language string, isWin bool, guessCount int, timeSeconds *int, scoreEarned int) error {
	currentStreaks, err := s.statsRepo.GetGlobalStreaks(ctx, q, userId, true)
	if err != nil {
		return fmt.Errorf("failed to get global streaks: %w", err)
	}

	newStreaks := s.CalculateNewStreaks(currentStreaks, userId, isWin)

	if currentStreaks == nil {
		if err := s.statsRepo.CreateGlobalStreaksWithData(ctx, q, newStreaks); err != nil {
			return fmt.Errorf("failed to create global streaks: %w", err)
		}
	} else {
		if err := s.statsRepo.UpdateGlobalStreaks(ctx, q, newStreaks); err != nil {
			return fmt.Errorf("failed to update global streaks: %w", err)
		}
	}

	currentLangStats, err := s.statsRepo.GetLanguageStats(ctx, q, userId, language, true)
	if err != nil {
		return fmt.Errorf("failed to get language stats: %w", err)
	}

	newLangStats := s.CalculateNewLanguageStats(currentLangStats, userId, language, isWin, guessCount, timeSeconds, scoreEarned)

	if currentLangStats == nil {
		if err := s.statsRepo.CreateLanguageStats(ctx, q, newLangStats); err != nil {
			return fmt.Errorf("failed to create language stats: %w", err)
		}
	} else {
		if err := s.statsRepo.UpdateLanguageStats(ctx, q, newLangStats); err != nil {
			return fmt.Errorf("failed to update language stats: %w", err)
		}
	}

	return nil
}

func (s *StatsService) UpdateGameEndStats(ctx context.Context, q repository.Querier, update model.StatisticsUpdate) error {
	return s.UpdateStatsAfterGame(
		ctx,
		q,
		update.UserId,
		update.Language,
		update.IsWin,
		update.GuessCount,
		update.TimeSeconds,
		update.ScoreEarned,
	)
}

func (s *StatsService) GetStatistics(ctx context.Context, userId int64) (*model.UserStatistics, error) {
	return s.statsRepo.GetStatistics(ctx, userId)
}

func (s *StatsService) GetLanguageStats(ctx context.Context, userId int64, language string) (*model.UserLanguageStats, error) {
	return s.statsRepo.GetLanguageStats(ctx, s.txManager.DB(), userId, language, false)
}

func (s *StatsService) GetAllLanguageStats(ctx context.Context, userId int64) ([]model.UserLanguageStats, error) {
	return s.statsRepo.GetAllLanguageStats(ctx, userId)
}
