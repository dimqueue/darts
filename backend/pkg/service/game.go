package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dimqueue/darts/pkg/config"
	"github.com/dimqueue/darts/pkg/connections"
	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
	"github.com/jmoiron/sqlx"
)

type GameService struct {
	computeClient connections.ComputeClient
	wordService   Word
	gameRepo      repository.Game
	statsService  *StatsService
	txManager     *repository.TransactionManager
}

func NewGameService(gameRepo repository.Game, wordService Word, statsService *StatsService, txManager *repository.TransactionManager, computeClient connections.ComputeClient) *GameService {
	return &GameService{
		computeClient: computeClient,
		wordService:   wordService,
		gameRepo:      gameRepo,
		statsService:  statsService,
		txManager:     txManager,
	}
}

func (s *GameService) CreateGame(userId int64, lang string) (int64, error) {
	selectedWord, err := s.wordService.SelectWord(lang)
	if err != nil {
		return 0, fmt.Errorf("failed to select word: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.GameComputeTimeout)
	defer cancel()

	req := &connections.StartGameRequest{
		Language:   lang,
		SecretWord: selectedWord.Word,
	}

	resp, err := s.computeClient.StartGame(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("failed to start game on compute client: %w", err)
	}

	expiresAt := time.Now().Add(config.GameTTL)

	game := model.Game{
		Language:  lang,
		Status:    "in_progress",
		UserId:    userId,
		WordId:    selectedWord.Id,
		ExpiresAt: &expiresAt,
	}

	fmt.Printf("Game started - Word: %s, Calculation time: %.3fs, Hint: %s\n", selectedWord.Word, resp.CalculationTime, resp.HintWord)
	return s.gameRepo.CreateGame(&game)
}

func (s *GameService) GetAllGames(userId int64) ([]model.Game, error) {
	return s.gameRepo.GetAllGames(userId)
}

func (s *GameService) GetGameById(userId, gameId int64) (*model.Game, error) {
	game, err := s.gameRepo.GetGameById(gameId)
	if err != nil {
		return nil, err
	}

	if game.UserId != userId {
		return nil, fmt.Errorf("unauthorized: access denied")
	}

	return game, nil
}

func (s *GameService) UpdateGame(gameId int64) (*model.Game, error) {
	return nil, nil
}

func (s *GameService) UpdateGameStatus(gameId int64, status string) error {
	return s.gameRepo.UpdateGameStatus(gameId, status)
}

func (s *GameService) DeleteGame(gameId int64) (*model.Game, error) {
	return nil, nil
}

func (s *GameService) MakeGuess(userId, gameId int64, guess string) (int, error) {
	game, err := s.GetGameById(userId, gameId)
	if err != nil {
		return 0, err
	}

	if game.Status != "in_progress" {
		return 0, errors.New("game is not active")
	}

	exists, err := s.gameRepo.GuessExists(gameId, guess)
	if err != nil {
		return 0, fmt.Errorf("failed to check guess: %w", err)
	}
	if exists {
		return 0, errors.New("word already guessed")
	}

	word, err := s.wordService.GetWordById(game.WordId)
	if err != nil {
		return 0, err
	}

	var distance int
	isWinningGuess := guess == word.Word

	if isWinningGuess {
		distance = 1
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), config.GameComputeTimeout)
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
		distance = resp.Distance
	}

	err = s.txManager.WithTransaction(func(tx *sqlx.Tx) error {
		lockedGame, err := s.gameRepo.GetGameByIdForUpdate(tx, gameId)
		if err != nil {
			return fmt.Errorf("failed to lock game: %w", err)
		}

		if lockedGame.Status != "in_progress" {
			return errors.New("game is no longer active")
		}

		existsTx, err := s.gameRepo.GuessExistsTx(tx, gameId, guess)
		if err != nil {
			return fmt.Errorf("failed to check guess in transaction: %w", err)
		}
		if existsTx {
			return errors.New("word already guessed")
		}

		newGuess := model.Guess{
			GameId:    gameId,
			GuessWord: guess,
			Distance:  distance,
		}
		if err := s.gameRepo.CreateGuessTx(tx, &newGuess); err != nil {
			return fmt.Errorf("failed to create guess: %w", err)
		}

		if isWinningGuess {
			if err := s.gameRepo.UpdateGameStatusTx(tx, gameId, "won"); err != nil {
				return fmt.Errorf("failed to update game status: %w", err)
			}

			guessCount, err := s.gameRepo.CountGuessesByGameTx(tx, gameId)
			if err != nil {
				return fmt.Errorf("failed to count guesses: %w", err)
			}

			elapsed := int(time.Since(lockedGame.StartedAt).Seconds())

			statsUpdate := model.StatisticsUpdate{
				UserId:      userId,
				Language:    lockedGame.Language,
				IsWin:       true,
				GuessCount:  guessCount,
				TimeSeconds: &elapsed,
				ScoreEarned: 100,
			}

			if err := s.statsService.UpdateGameEndStats(tx, statsUpdate); err != nil {
				return fmt.Errorf("failed to update statistics: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return distance, nil
}

func (s *GameService) GetAllGuessByGame(userId, gameId int64) ([]model.Guess, error) {
	_, err := s.GetGameById(userId, gameId)
	if err != nil {
		return nil, err
	}

	return s.gameRepo.GetAllGuessByGame(gameId)
}
