package mocks

import (
	"context"

	"github.com/dimqueue/darts/pkg/model"
)

type MockWordService struct {
	SelectWordFn  func(ctx context.Context, language string) (*model.Word, error)
	GetWordByIdFn func(ctx context.Context, wordId int64) (*model.Word, error)
}

func (m *MockWordService) SelectWord(ctx context.Context, language string) (*model.Word, error) {
	if m.SelectWordFn != nil {
		return m.SelectWordFn(ctx, language)
	}
	return &model.Word{Id: 1, Word: "testword", Language: language}, nil
}

func (m *MockWordService) GetWordById(ctx context.Context, wordId int64) (*model.Word, error) {
	if m.GetWordByIdFn != nil {
		return m.GetWordByIdFn(ctx, wordId)
	}
	return &model.Word{Id: wordId, Word: "testword", Language: "en"}, nil
}
