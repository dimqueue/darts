package service

import (
	"fmt"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
)

type WordService struct {
	repo repository.Word
}

func NewWordService(repo repository.Word) *WordService {
	return &WordService{
		repo: repo,
	}
}

func (s *WordService) SelectWord(language string) (*model.Word, error) {
	word, err := s.repo.GetRandomWordByLanguage(language)
	if err != nil {
		return nil, fmt.Errorf("failed to select random word: %w", err)
	}

	return word, nil
}

func (s *WordService) GetWordById(wordId int) (*model.Word, error) {
	return s.repo.GetWordById(wordId)
}

// - SelectRandomWord(language string, difficulty string) - select by difficulty
// - SelectWordByCategory(language string, category string) - themed words
// - SelectUnusedWord(userId int, language string) - words user hasn't seen
