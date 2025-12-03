package service

import (
	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
	"github.com/jmoiron/sqlx"
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

func (s *StatsService) InitializeStats(tx *sqlx.Tx, userId int64) error {
	return s.statsRepo.CreateGlobalStreaks(tx, userId)
}

func (s *StatsService) UpdateGameEndStats(tx *sqlx.Tx, update model.StatisticsUpdate) error {
	if err := s.statsRepo.UpdateLanguageStats(tx, update); err != nil {
		return err
	}

	if err := s.statsRepo.UpdateGlobalStreaksAfterGame(tx, update); err != nil {
		return err
	}

	return nil
}

func (s *StatsService) GetStatistics(userId int64) (*model.UserStatistics, error) {
	return s.statsRepo.GetStatistics(userId)
}

func (s *StatsService) GetLanguageStats(userId int64, language string) (*model.UserLanguageStats, error) {
	return s.statsRepo.GetLanguageStats(userId, language)
}

func (s *StatsService) GetAllLanguageStats(userId int64) ([]model.UserLanguageStats, error) {
	return s.statsRepo.GetAllLanguageStats(userId)
}
