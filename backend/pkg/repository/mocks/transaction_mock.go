package mocks

import (
	"context"
	"database/sql"

	"github.com/dimqueue/darts/pkg/repository"
	"github.com/jmoiron/sqlx"
)

type MockQuerier struct {
	GetContextFn      func(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContextFn   func(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContextFn     func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContextFn func(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func (m *MockQuerier) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if m.GetContextFn != nil {
		return m.GetContextFn(ctx, dest, query, args...)
	}
	return nil
}

func (m *MockQuerier) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if m.SelectContextFn != nil {
		return m.SelectContextFn(ctx, dest, query, args...)
	}
	return nil
}

func (m *MockQuerier) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if m.ExecContextFn != nil {
		return m.ExecContextFn(ctx, query, args...)
	}
	return nil, nil
}

func (m *MockQuerier) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if m.QueryRowContextFn != nil {
		return m.QueryRowContextFn(ctx, query, args...)
	}
	return nil
}

type MockTransactionManager struct {
	WithTransactionFn func(ctx context.Context, fn func(*sqlx.Tx) error) error
	DBFn              func() repository.Querier
}

func (m *MockTransactionManager) WithTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	if m.WithTransactionFn != nil {
		return m.WithTransactionFn(ctx, fn)
	}
	return fn(nil)
}

func (m *MockTransactionManager) DB() repository.Querier {
	if m.DBFn != nil {
		return m.DBFn()
	}
	return &MockQuerier{}
}
