package mocks

import (
	"context"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
)

type MockGameRepository struct {
	CreateGameFn         func(ctx context.Context, game *model.Game) (int64, error)
	GetAllGamesFn        func(ctx context.Context, userId int64) ([]model.Game, error)
	GetActiveGameFn      func(ctx context.Context, userId int64) (*model.Game, error)
	GetGameByIdFn        func(ctx context.Context, q repository.Querier, gameId int64, forUpdate bool) (*model.Game, error)
	UpdateGameStatusFn   func(ctx context.Context, q repository.Querier, gameId int64, status string) error
	GetExpiredGamesFn    func(ctx context.Context) ([]model.Game, error)
	CreateGuessFn        func(ctx context.Context, q repository.Querier, guess *model.Guess) error
	GetAllGuessByGameFn  func(ctx context.Context, gameId int64) ([]model.Guess, error)
	CountGuessesByGameFn func(ctx context.Context, q repository.Querier, gameId int64) (int, error)
	GuessExistsFn        func(ctx context.Context, q repository.Querier, gameId int64, guessWord string) (bool, error)
}

func (m *MockGameRepository) CreateGame(ctx context.Context, game *model.Game) (int64, error) {
	if m.CreateGameFn != nil {
		return m.CreateGameFn(ctx, game)
	}
	return 0, nil
}

func (m *MockGameRepository) GetAllGames(ctx context.Context, userId int64) ([]model.Game, error) {
	if m.GetAllGamesFn != nil {
		return m.GetAllGamesFn(ctx, userId)
	}
	return nil, nil
}

func (m *MockGameRepository) GetActiveGame(ctx context.Context, userId int64) (*model.Game, error) {
	if m.GetActiveGameFn != nil {
		return m.GetActiveGameFn(ctx, userId)
	}
	return nil, nil
}

func (m *MockGameRepository) GetGameById(ctx context.Context, q repository.Querier, gameId int64, forUpdate bool) (*model.Game, error) {
	if m.GetGameByIdFn != nil {
		return m.GetGameByIdFn(ctx, q, gameId, forUpdate)
	}
	return nil, nil
}

func (m *MockGameRepository) UpdateGameStatus(ctx context.Context, q repository.Querier, gameId int64, status string) error {
	if m.UpdateGameStatusFn != nil {
		return m.UpdateGameStatusFn(ctx, q, gameId, status)
	}
	return nil
}

func (m *MockGameRepository) GetExpiredGames(ctx context.Context) ([]model.Game, error) {
	if m.GetExpiredGamesFn != nil {
		return m.GetExpiredGamesFn(ctx)
	}
	return nil, nil
}

func (m *MockGameRepository) CreateGuess(ctx context.Context, q repository.Querier, guess *model.Guess) error {
	if m.CreateGuessFn != nil {
		return m.CreateGuessFn(ctx, q, guess)
	}
	return nil
}

func (m *MockGameRepository) GetAllGuessByGame(ctx context.Context, gameId int64) ([]model.Guess, error) {
	if m.GetAllGuessByGameFn != nil {
		return m.GetAllGuessByGameFn(ctx, gameId)
	}
	return nil, nil
}

func (m *MockGameRepository) CountGuessesByGame(ctx context.Context, q repository.Querier, gameId int64) (int, error) {
	if m.CountGuessesByGameFn != nil {
		return m.CountGuessesByGameFn(ctx, q, gameId)
	}
	return 0, nil
}

func (m *MockGameRepository) GuessExists(ctx context.Context, q repository.Querier, gameId int64, guessWord string) (bool, error) {
	if m.GuessExistsFn != nil {
		return m.GuessExistsFn(ctx, q, gameId, guessWord)
	}
	return false, nil
}
