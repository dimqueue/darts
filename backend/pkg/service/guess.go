package service

import (
	"errors"

	"github.com/dimqueue/darts/pkg/connections"
	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
)

type GuessService struct {
	computeClient *connections.ComputeClientService
	guessRepo     repository.Guess
	gameRepo      repository.Game
	wordRepo      repository.Word
}

func NewGuessService(guessRepo repository.Guess, gameRepo repository.Game, wordRepo repository.Word, computeClient *connections.ComputeClientService) *GuessService {
	return &GuessService{
		guessRepo:     guessRepo,
		gameRepo:      gameRepo,
		wordRepo:      wordRepo,
		computeClient: computeClient,
	}
}

func (s *GuessService) CreateGuess(userId, gameId int, guess string) (int, error) {
	game, err := s.gameRepo.GetGameById(gameId)
	if err != nil {
		return 0, err
	}

	if game.UserId != userId {
		return 0, errors.New("unauthorized: cannot add guess to another user's game")
	}

	if game.Status != "in_progress" {
		return 0, errors.New("game is not active")
	}

	word, err := s.wordRepo.GetWordById(game.WordId)
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

		err = s.gameRepo.UpdateGameStatus(gameId, "won")
		if err != nil {
			return 0, err
		}

		return distance, nil
	}

	resp, err := s.computeClient.MakeGuess(word.Word, guess, game.Language)
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
	game, err := s.gameRepo.GetGameById(gameId)
	if err != nil {
		return nil, err
	}

	if game.UserId != userId {
		return nil, errors.New("unauthorized: cannot view guesses for another user's game")
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
