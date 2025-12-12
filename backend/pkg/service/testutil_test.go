package service

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type MockTxManager struct {
	WithTransactionFn func(ctx context.Context, fn func(*sqlx.Tx) error) error
}

func (m *MockTxManager) WithTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	if m.WithTransactionFn != nil {
		return m.WithTransactionFn(ctx, fn)
	}
	return fn(nil)
}
