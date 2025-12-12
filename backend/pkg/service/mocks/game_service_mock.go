package mocks

import (
	"context"

	"github.com/dimqueue/darts/pkg/model"
)

type MockGameService struct {
	CreateGameFn        func(ctx context.Context, userId int64, lang string) (int64, error)
	GetAllGamesFn       func(ctx context.Context, userId int64) ([]model.Game, error)
	GetActiveGameFn     func(ctx context.Context, userId int64) (*model.Game, error)
	GetGameByIdFn       func(ctx context.Context, userId, gameId int64) (*model.Game, error)
	UpdateGameStatusFn  func(ctx context.Context, gameId int64, status string) error
	AbandonGameFn       func(ctx context.Context, userId, gameId int64) error
	MakeGuessFn         func(ctx context.Context, userId, gameId int64, guess string) (int, error)
	GetAllGuessByGameFn func(ctx context.Context, userId, gameId int64) ([]model.Guess, error)
}

func (m *MockGameService) CreateGame(ctx context.Context, userId int64, lang string) (int64, error) {
	if m.CreateGameFn != nil {
		return m.CreateGameFn(ctx, userId, lang)
	}
	return 1, nil
}

func (m *MockGameService) GetAllGames(ctx context.Context, userId int64) ([]model.Game, error) {
	if m.GetAllGamesFn != nil {
		return m.GetAllGamesFn(ctx, userId)
	}
	return []model.Game{}, nil
}

func (m *MockGameService) GetActiveGame(ctx context.Context, userId int64) (*model.Game, error) {
	if m.GetActiveGameFn != nil {
		return m.GetActiveGameFn(ctx, userId)
	}
	return nil, nil
}

func (m *MockGameService) GetGameById(ctx context.Context, userId, gameId int64) (*model.Game, error) {
	if m.GetGameByIdFn != nil {
		return m.GetGameByIdFn(ctx, userId, gameId)
	}
	return nil, nil
}

func (m *MockGameService) UpdateGameStatus(ctx context.Context, gameId int64, status string) error {
	if m.UpdateGameStatusFn != nil {
		return m.UpdateGameStatusFn(ctx, gameId, status)
	}
	return nil
}

func (m *MockGameService) AbandonGame(ctx context.Context, userId, gameId int64) error {
	if m.AbandonGameFn != nil {
		return m.AbandonGameFn(ctx, userId, gameId)
	}
	return nil
}

func (m *MockGameService) MakeGuess(ctx context.Context, userId, gameId int64, guess string) (int, error) {
	if m.MakeGuessFn != nil {
		return m.MakeGuessFn(ctx, userId, gameId, guess)
	}
	return 100, nil
}

func (m *MockGameService) GetAllGuessByGame(ctx context.Context, userId, gameId int64) ([]model.Guess, error) {
	if m.GetAllGuessByGameFn != nil {
		return m.GetAllGuessByGameFn(ctx, userId, gameId)
	}
	return []model.Guess{}, nil
}
