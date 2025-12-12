package service

import (
	"context"
	"errors"
	"testing"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository/mocks"
)

func TestSelectWord_Success(t *testing.T) {
	mockRepo := &mocks.MockWordRepository{
		GetRandomWordByLanguageFn: func(ctx context.Context, language string) (*model.Word, error) {
			if language != "en" {
				t.Errorf("expected language 'en', got: %s", language)
			}
			return &model.Word{
				Id:         1,
				Word:       "example",
				Language:   language,
				Difficulty: 1,
				IsActive:   true,
			}, nil
		},
	}

	service := NewWordService(mockRepo)

	word, err := service.SelectWord(context.Background(), "en")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if word.Word != "example" {
		t.Errorf("expected word 'example', got: %s", word.Word)
	}
	if word.Language != "en" {
		t.Errorf("expected language 'en', got: %s", word.Language)
	}
}

func TestSelectWord_DifferentLanguages(t *testing.T) {
	languages := []string{"en", "ua"}

	for _, lang := range languages {
		t.Run(lang, func(t *testing.T) {
			mockRepo := &mocks.MockWordRepository{
				GetRandomWordByLanguageFn: func(ctx context.Context, language string) (*model.Word, error) {
					return &model.Word{
						Id:       1,
						Word:     "test",
						Language: language,
					}, nil
				},
			}

			service := NewWordService(mockRepo)

			word, err := service.SelectWord(context.Background(), lang)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if word.Language != lang {
				t.Errorf("expected language '%s', got: %s", lang, word.Language)
			}
		})
	}
}

func TestSelectWord_RepoError(t *testing.T) {
	mockRepo := &mocks.MockWordRepository{
		GetRandomWordByLanguageFn: func(ctx context.Context, language string) (*model.Word, error) {
			return nil, errors.New("database error")
		},
	}

	service := NewWordService(mockRepo)

	_, err := service.SelectWord(context.Background(), "en")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestSelectWord_NoWordsAvailable(t *testing.T) {
	mockRepo := &mocks.MockWordRepository{
		GetRandomWordByLanguageFn: func(ctx context.Context, language string) (*model.Word, error) {
			return nil, errors.New("no words found")
		},
	}

	service := NewWordService(mockRepo)

	_, err := service.SelectWord(context.Background(), "xyz")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetWordById_Success(t *testing.T) {
	mockRepo := &mocks.MockWordRepository{
		GetWordByIdFn: func(ctx context.Context, wordId int64) (*model.Word, error) {
			if wordId != 42 {
				t.Errorf("expected wordId 42, got: %d", wordId)
			}
			return &model.Word{
				Id:         wordId,
				Word:       "hello",
				Language:   "en",
				Difficulty: 2,
				IsActive:   true,
			}, nil
		},
	}

	service := NewWordService(mockRepo)

	word, err := service.GetWordById(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if word.Id != 42 {
		t.Errorf("expected id 42, got: %d", word.Id)
	}
	if word.Word != "hello" {
		t.Errorf("expected word 'hello', got: %s", word.Word)
	}
}

func TestGetWordById_NotFound(t *testing.T) {
	mockRepo := &mocks.MockWordRepository{
		GetWordByIdFn: func(ctx context.Context, wordId int64) (*model.Word, error) {
			return nil, errors.New("word not found")
		},
	}

	service := NewWordService(mockRepo)

	_, err := service.GetWordById(context.Background(), 999)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetWordById_RepoError(t *testing.T) {
	mockRepo := &mocks.MockWordRepository{
		GetWordByIdFn: func(ctx context.Context, wordId int64) (*model.Word, error) {
			return nil, errors.New("database connection failed")
		},
	}

	service := NewWordService(mockRepo)

	_, err := service.GetWordById(context.Background(), 1)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
