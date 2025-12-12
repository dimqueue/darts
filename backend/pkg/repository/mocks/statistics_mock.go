package mocks

import (
	"context"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
)

type MockStatisticsRepository struct {
	GetStatisticsFn               func(ctx context.Context, userId int64) (*model.UserStatistics, error)
	CreateGlobalStreaksFn         func(ctx context.Context, q repository.Querier, userId int64) error
	GetGlobalStreaksFn            func(ctx context.Context, q repository.Querier, userId int64, forUpdate bool) (*model.UserGlobalStreaks, error)
	CreateGlobalStreaksWithDataFn func(ctx context.Context, q repository.Querier, streaks *model.UserGlobalStreaks) error
	UpdateGlobalStreaksFn         func(ctx context.Context, q repository.Querier, streaks *model.UserGlobalStreaks) error
	GetLanguageStatsFn            func(ctx context.Context, q repository.Querier, userId int64, language string, forUpdate bool) (*model.UserLanguageStats, error)
	GetAllLanguageStatsFn         func(ctx context.Context, userId int64) ([]model.UserLanguageStats, error)
	CreateLanguageStatsFn         func(ctx context.Context, q repository.Querier, stats *model.UserLanguageStats) error
	UpdateLanguageStatsFn         func(ctx context.Context, q repository.Querier, stats *model.UserLanguageStats) error
}

func (m *MockStatisticsRepository) GetStatistics(ctx context.Context, userId int64) (*model.UserStatistics, error) {
	if m.GetStatisticsFn != nil {
		return m.GetStatisticsFn(ctx, userId)
	}
	return nil, nil
}

func (m *MockStatisticsRepository) CreateGlobalStreaks(ctx context.Context, q repository.Querier, userId int64) error {
	if m.CreateGlobalStreaksFn != nil {
		return m.CreateGlobalStreaksFn(ctx, q, userId)
	}
	return nil
}

func (m *MockStatisticsRepository) GetGlobalStreaks(ctx context.Context, q repository.Querier, userId int64, forUpdate bool) (*model.UserGlobalStreaks, error) {
	if m.GetGlobalStreaksFn != nil {
		return m.GetGlobalStreaksFn(ctx, q, userId, forUpdate)
	}
	return nil, nil
}

func (m *MockStatisticsRepository) CreateGlobalStreaksWithData(ctx context.Context, q repository.Querier, streaks *model.UserGlobalStreaks) error {
	if m.CreateGlobalStreaksWithDataFn != nil {
		return m.CreateGlobalStreaksWithDataFn(ctx, q, streaks)
	}
	return nil
}

func (m *MockStatisticsRepository) UpdateGlobalStreaks(ctx context.Context, q repository.Querier, streaks *model.UserGlobalStreaks) error {
	if m.UpdateGlobalStreaksFn != nil {
		return m.UpdateGlobalStreaksFn(ctx, q, streaks)
	}
	return nil
}

func (m *MockStatisticsRepository) GetLanguageStats(ctx context.Context, q repository.Querier, userId int64, language string, forUpdate bool) (*model.UserLanguageStats, error) {
	if m.GetLanguageStatsFn != nil {
		return m.GetLanguageStatsFn(ctx, q, userId, language, forUpdate)
	}
	return nil, nil
}

func (m *MockStatisticsRepository) GetAllLanguageStats(ctx context.Context, userId int64) ([]model.UserLanguageStats, error) {
	if m.GetAllLanguageStatsFn != nil {
		return m.GetAllLanguageStatsFn(ctx, userId)
	}
	return nil, nil
}

func (m *MockStatisticsRepository) CreateLanguageStats(ctx context.Context, q repository.Querier, stats *model.UserLanguageStats) error {
	if m.CreateLanguageStatsFn != nil {
		return m.CreateLanguageStatsFn(ctx, q, stats)
	}
	return nil
}

func (m *MockStatisticsRepository) UpdateLanguageStats(ctx context.Context, q repository.Querier, stats *model.UserLanguageStats) error {
	if m.UpdateLanguageStatsFn != nil {
		return m.UpdateLanguageStatsFn(ctx, q, stats)
	}
	return nil
}
