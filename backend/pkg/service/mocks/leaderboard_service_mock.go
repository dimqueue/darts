package mocks

import (
	"context"

	"github.com/dimqueue/darts/pkg/model"
)

type MockLeaderboardService struct {
	GetLeaderboardFn             func(ctx context.Context, query model.LeaderboardQuery) (*model.LeaderboardResponse, error)
	GetLeaderboardWithUserRankFn func(ctx context.Context, userId int64, query model.LeaderboardQuery) (*model.LeaderboardResponse, error)
	GetUserRankFn                func(ctx context.Context, userId int64, query model.LeaderboardQuery) (*int, error)
	GetAllUserRanksFn            func(ctx context.Context, userId int64) (*model.UserRanks, error)
}

func (m *MockLeaderboardService) GetLeaderboard(ctx context.Context, query model.LeaderboardQuery) (*model.LeaderboardResponse, error) {
	if m.GetLeaderboardFn != nil {
		return m.GetLeaderboardFn(ctx, query)
	}
	return &model.LeaderboardResponse{
		LeaderboardType: query.Type,
		Users:           []model.LeaderboardUser{},
		Total:           0,
	}, nil
}

func (m *MockLeaderboardService) GetLeaderboardWithUserRank(ctx context.Context, userId int64, query model.LeaderboardQuery) (*model.LeaderboardResponse, error) {
	if m.GetLeaderboardWithUserRankFn != nil {
		return m.GetLeaderboardWithUserRankFn(ctx, userId, query)
	}
	return &model.LeaderboardResponse{
		LeaderboardType: query.Type,
		Users:           []model.LeaderboardUser{},
		Total:           0,
	}, nil
}

func (m *MockLeaderboardService) GetUserRank(ctx context.Context, userId int64, query model.LeaderboardQuery) (*int, error) {
	if m.GetUserRankFn != nil {
		return m.GetUserRankFn(ctx, userId, query)
	}
	return nil, nil
}

func (m *MockLeaderboardService) GetAllUserRanks(ctx context.Context, userId int64) (*model.UserRanks, error) {
	if m.GetAllUserRanksFn != nil {
		return m.GetAllUserRanksFn(ctx, userId)
	}
	return &model.UserRanks{}, nil
}
