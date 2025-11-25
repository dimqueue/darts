package service

import (
	"context"
	"errors"
	"time"

	"github.com/dimqueue/darts/pkg/connections"
	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
)

type GuessService struct {
	computeClient connections.ComputeClient
	guessRepo     repository.Guess
	gameService   Game
	wordService   Word
}

func NewGuessService(guessRepo repository.Guess, gameService Game, wordService Word, computeClient connections.ComputeClient) *GuessService {
	return &GuessService{
		guessRepo:     guessRepo,
		gameService:   gameService,
		wordService:   wordService,
		computeClient: computeClient,
	}
}

func (s *GuessService) CreateGuess(userId, gameId int, guess string) (int, error) {
	game, err := s.gameService.GetGameById(userId, gameId)
	if err != nil {
		return 0, err
	}

	if game.Status != "in_progress" {
		return 0, errors.New("game is not active")
	}

	word, err := s.wordService.GetWordById(game.WordId)
	if err != nil {
		return 0, err
	}

	var distance int

	if guess == word.Word {
		distance = 1

		newGuess := model.Guess{
			GameId:    gameId,
			GuessWord: guess,
			Distance:  distance,
		}
		err = s.guessRepo.CreateGuess(&newGuess)
		if err != nil {
			return 0, err
		}

		err = s.gameService.UpdateGameStatus(gameId, "won")
		if err != nil {
			return 0, err
		}

		return distance, nil
	}

	// Create context with timeout for compute service call
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &connections.GuessRequest{
		SecretWord: word.Word,
		Guess:      guess,
		Language:   game.Language,
	}

	resp, err := s.computeClient.MakeGuess(ctx, req)
	if err != nil {
		return 0, err
	}

	newGuess := model.Guess{
		GameId:    gameId,
		GuessWord: guess,
		Distance:  resp.Distance,
	}
	err = s.guessRepo.CreateGuess(&newGuess)
	if err != nil {
		return 0, err
	}

	return resp.Distance, nil
}

func (s *GuessService) GetAllGuessByGame(userId, gameId int) ([]model.Guess, error) {
	_, err := s.gameService.GetGameById(userId, gameId)
	if err != nil {
		return nil, err
	}

	guesses, err := s.guessRepo.GetAllGuessByGame(gameId)
	if err != nil {
		return nil, err
	}

	return guesses, nil
}

func (s *GuessService) GetGuessById(i int) error {
	return nil
}
