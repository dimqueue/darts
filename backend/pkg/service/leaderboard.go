package service

import (
	"context"
	"fmt"
	"time"

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

func GetPeriodStart(periodType string) (time.Time, error) {
	now := time.Now()
	switch periodType {
	case "daily":
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), nil
	case "weekly":
		weekStart := now.AddDate(0, 0, -int(now.Weekday())+1)
		if now.Weekday() == time.Sunday {
			weekStart = now.AddDate(0, 0, -6)
		}
		return time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, time.UTC), nil
	case "monthly":
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC), nil
	default:
		return time.Time{}, fmt.Errorf("invalid period type: %s", periodType)
	}
}

func (s *LeaderboardService) GetLeaderboard(ctx context.Context, query model.LeaderboardQuery) (*model.LeaderboardResponse, error) {
	if query.Limit <= 0 {
		query.Limit = config.LeaderboardDefaultLimit
	}
	if query.Limit > config.LeaderboardMaxLimit {
		query.Limit = config.LeaderboardMaxLimit
	}

	var users []model.LeaderboardUser
	var total int
	var err error

	if query.Type == "global" {
		if query.Language != nil {
			users, err = s.leaderboardRepo.GetGlobalLeaderboardByLanguage(ctx, *query.Language, query.Limit, query.Offset)
			if err != nil {
				return nil, err
			}
			total, err = s.leaderboardRepo.GetGlobalLeaderboardByLanguageCount(ctx, *query.Language)
		} else {
			users, err = s.leaderboardRepo.GetGlobalLeaderboard(ctx, query.Limit, query.Offset)
			if err != nil {
				return nil, err
			}
			total, err = s.leaderboardRepo.GetGlobalLeaderboardCount(ctx)
		}
	} else {
		periodStart, err := GetPeriodStart(query.Type)
		if err != nil {
			return nil, err
		}
		users, err = s.leaderboardRepo.GetPeriodLeaderboard(ctx, periodStart, query.Language, query.Limit, query.Offset)
		if err != nil {
			return nil, err
		}
		total, err = s.leaderboardRepo.GetPeriodLeaderboardCount(ctx, periodStart, query.Language)
	}

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

func (s *LeaderboardService) GetLeaderboardWithUserRank(ctx context.Context, userId int64, query model.LeaderboardQuery) (*model.LeaderboardResponse, error) {
	response, err := s.GetLeaderboard(ctx, query)
	if err != nil {
		return nil, err
	}

	userRank, _ := s.GetUserRank(ctx, userId, query)
	response.CurrentUserRank = userRank

	return response, nil
}

func (s *LeaderboardService) GetUserRank(ctx context.Context, userId int64, query model.LeaderboardQuery) (*int, error) {
	if query.Type == "global" {
		if query.Language != nil {
			return s.leaderboardRepo.GetGlobalUserRankByLanguage(ctx, userId, *query.Language)
		}
		return s.leaderboardRepo.GetGlobalUserRank(ctx, userId)
	}

	periodStart, err := GetPeriodStart(query.Type)
	if err != nil {
		return nil, err
	}
	return s.leaderboardRepo.GetPeriodUserRank(ctx, userId, periodStart, query.Language)
}

func (s *LeaderboardService) GetAllUserRanks(ctx context.Context, userId int64) (*model.UserRanks, error) {
	var ranks model.UserRanks

	ranks.GlobalRank, _ = s.leaderboardRepo.GetGlobalUserRank(ctx, userId)

	dayStart, _ := GetPeriodStart("daily")
	ranks.DailyRank, _ = s.leaderboardRepo.GetPeriodUserRank(ctx, userId, dayStart, nil)

	weekStart, _ := GetPeriodStart("weekly")
	ranks.WeeklyRank, _ = s.leaderboardRepo.GetPeriodUserRank(ctx, userId, weekStart, nil)

	monthStart, _ := GetPeriodStart("monthly")
	ranks.MonthlyRank, _ = s.leaderboardRepo.GetPeriodUserRank(ctx, userId, monthStart, nil)

	return &ranks, nil
}
