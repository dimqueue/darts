package service

import (
	"github.com/dimqueue/darts/pkg/config"
	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
)

type LeaderboardService struct {
	leaderboardRepo repository.Leaderboard
}

func NewLeaderboardService(leaderboardRepo repository.Leaderboard) *LeaderboardService {
	return &LeaderboardService{
		leaderboardRepo: leaderboardRepo,
	}
}

func (s *LeaderboardService) GetLeaderboard(query model.LeaderboardQuery) (*model.LeaderboardResponse, error) {
	if query.Limit <= 0 {
		query.Limit = config.LeaderboardDefaultLimit
	}
	if query.Limit > config.LeaderboardMaxLimit {
		query.Limit = config.LeaderboardMaxLimit
	}

	users, err := s.leaderboardRepo.GetLeaderboard(query)
	if err != nil {
		return nil, err
	}

	total, err := s.leaderboardRepo.GetLeaderboardCount(query)
	if err != nil {
		return nil, err
	}

	return &model.LeaderboardResponse{
		LeaderboardType: query.Type,
		Language:        query.Language,
		Users:           users,
		Total:           total,
	}, nil
}

func (s *LeaderboardService) GetLeaderboardWithUserRank(userId int64, query model.LeaderboardQuery) (*model.LeaderboardResponse, error) {
	response, err := s.GetLeaderboard(query)
	if err != nil {
		return nil, err
	}

	userRank, _ := s.leaderboardRepo.GetUserRank(userId, query)
	response.CurrentUserRank = userRank

	return response, nil
}

func (s *LeaderboardService) GetUserRank(userId int64, query model.LeaderboardQuery) (*int, error) {
	return s.leaderboardRepo.GetUserRank(userId, query)
}

func (s *LeaderboardService) GetAllUserRanks(userId int64) (*model.UserRanks, error) {
	return s.leaderboardRepo.GetAllUserRanks(userId)
}
