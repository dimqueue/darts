package service

import (
	"context"

	"github.com/dimqueue/darts/pkg/connections"
	"github.com/dimqueue/darts/pkg/repository"
	"github.com/jmoiron/sqlx"
)

type MockTxManager struct {
	WithTransactionFn func(ctx context.Context, fn func(*sqlx.Tx) error) error
	DBFn              func() repository.Querier
}

func (m *MockTxManager) WithTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	if m.WithTransactionFn != nil {
		return m.WithTransactionFn(ctx, fn)
	}
	return fn(nil)
}

func (m *MockTxManager) DB() repository.Querier {
	if m.DBFn != nil {
		return m.DBFn()
	}
	return nil
}

// NewGameServiceWithTxManager creates a GameService with a custom TxManager for testing
func NewGameServiceWithTxManager(gameRepo repository.Game, wordService Word, statsService *StatsService, txManager repository.TxManager, computeClient connections.ComputeClient) *GameService {
	return &GameService{
		computeClient: computeClient,
		wordService:   wordService,
		gameRepo:      gameRepo,
		statsService:  statsService,
		txManager:     txManager,
	}
}
