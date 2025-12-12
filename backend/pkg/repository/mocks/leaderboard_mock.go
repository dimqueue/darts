package mocks

import (
	"context"
	"time"

	"github.com/dimqueue/darts/pkg/model"
)

type MockLeaderboardRepository struct {
	GetGlobalLeaderboardFn                func(ctx context.Context, limit, offset int) ([]model.LeaderboardUser, error)
	GetGlobalLeaderboardByLanguageFn      func(ctx context.Context, language string, limit, offset int) ([]model.LeaderboardUser, error)
	GetGlobalLeaderboardCountFn           func(ctx context.Context) (int, error)
	GetGlobalLeaderboardByLanguageCountFn func(ctx context.Context, language string) (int, error)
	GetGlobalUserRankFn                   func(ctx context.Context, userId int64) (*int, error)
	GetGlobalUserRankByLanguageFn         func(ctx context.Context, userId int64, language string) (*int, error)
	GetPeriodLeaderboardFn                func(ctx context.Context, periodStart time.Time, language *string, limit, offset int) ([]model.LeaderboardUser, error)
	GetPeriodLeaderboardCountFn           func(ctx context.Context, periodStart time.Time, language *string) (int, error)
	GetPeriodUserRankFn                   func(ctx context.Context, userId int64, periodStart time.Time, language *string) (*int, error)
}

func (m *MockLeaderboardRepository) GetGlobalLeaderboard(ctx context.Context, limit, offset int) ([]model.LeaderboardUser, error) {
	if m.GetGlobalLeaderboardFn != nil {
		return m.GetGlobalLeaderboardFn(ctx, limit, offset)
	}
	return nil, nil
}

func (m *MockLeaderboardRepository) GetGlobalLeaderboardByLanguage(ctx context.Context, language string, limit, offset int) ([]model.LeaderboardUser, error) {
	if m.GetGlobalLeaderboardByLanguageFn != nil {
		return m.GetGlobalLeaderboardByLanguageFn(ctx, language, limit, offset)
	}
	return nil, nil
}

func (m *MockLeaderboardRepository) GetGlobalLeaderboardCount(ctx context.Context) (int, error) {
	if m.GetGlobalLeaderboardCountFn != nil {
		return m.GetGlobalLeaderboardCountFn(ctx)
	}
	return 0, nil
}

func (m *MockLeaderboardRepository) GetGlobalLeaderboardByLanguageCount(ctx context.Context, language string) (int, error) {
	if m.GetGlobalLeaderboardByLanguageCountFn != nil {
		return m.GetGlobalLeaderboardByLanguageCountFn(ctx, language)
	}
	return 0, nil
}

func (m *MockLeaderboardRepository) GetGlobalUserRank(ctx context.Context, userId int64) (*int, error) {
	if m.GetGlobalUserRankFn != nil {
		return m.GetGlobalUserRankFn(ctx, userId)
	}
	return nil, nil
}

func (m *MockLeaderboardRepository) GetGlobalUserRankByLanguage(ctx context.Context, userId int64, language string) (*int, error) {
	if m.GetGlobalUserRankByLanguageFn != nil {
		return m.GetGlobalUserRankByLanguageFn(ctx, userId, language)
	}
	return nil, nil
}

func (m *MockLeaderboardRepository) GetPeriodLeaderboard(ctx context.Context, periodStart time.Time, language *string, limit, offset int) ([]model.LeaderboardUser, error) {
	if m.GetPeriodLeaderboardFn != nil {
		return m.GetPeriodLeaderboardFn(ctx, periodStart, language, limit, offset)
	}
	return nil, nil
}

func (m *MockLeaderboardRepository) GetPeriodLeaderboardCount(ctx context.Context, periodStart time.Time, language *string) (int, error) {
	if m.GetPeriodLeaderboardCountFn != nil {
		return m.GetPeriodLeaderboardCountFn(ctx, periodStart, language)
	}
	return 0, nil
}

func (m *MockLeaderboardRepository) GetPeriodUserRank(ctx context.Context, userId int64, periodStart time.Time, language *string) (*int, error) {
	if m.GetPeriodUserRankFn != nil {
		return m.GetPeriodUserRankFn(ctx, userId, periodStart, language)
	}
	return nil, nil
}
