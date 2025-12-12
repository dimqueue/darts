package mocks

import (
	"context"

	"github.com/dimqueue/darts/pkg/model"
)

type MockWordRepository struct {
	GetWordByIdFn             func(ctx context.Context, wordId int64) (*model.Word, error)
	GetRandomWordByLanguageFn func(ctx context.Context, language string) (*model.Word, error)
	GetWordCountByLanguageFn  func(ctx context.Context, language string) (int, error)
}

func (m *MockWordRepository) GetWordById(ctx context.Context, wordId int64) (*model.Word, error) {
	if m.GetWordByIdFn != nil {
		return m.GetWordByIdFn(ctx, wordId)
	}
	return nil, nil
}

func (m *MockWordRepository) GetRandomWordByLanguage(ctx context.Context, language string) (*model.Word, error) {
	if m.GetRandomWordByLanguageFn != nil {
		return m.GetRandomWordByLanguageFn(ctx, language)
	}
	return nil, nil
}

func (m *MockWordRepository) GetWordCountByLanguage(ctx context.Context, language string) (int, error) {
	if m.GetWordCountByLanguageFn != nil {
		return m.GetWordCountByLanguageFn(ctx, language)
	}
	return 0, nil
}
