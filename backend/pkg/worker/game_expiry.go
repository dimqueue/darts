package worker

import (
	"context"
	"time"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/dimqueue/darts/pkg/repository"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// StatsUpdater interface to avoid circular dependency with service package
type StatsUpdater interface {
	UpdateGameEndStats(q repository.Querier, update model.StatisticsUpdate) error
}

type GameExpiryWorker struct {
	gameRepo     repository.Game
	statsUpdater StatsUpdater
	txManager    *repository.TransactionManager
	interval     time.Duration
}

func NewGameExpiryWorker(gameRepo repository.Game, statsUpdater StatsUpdater, txManager *repository.TransactionManager, interval time.Duration) *GameExpiryWorker {
	return &GameExpiryWorker{
		gameRepo:     gameRepo,
		statsUpdater: statsUpdater,
		txManager:    txManager,
		interval:     interval,
	}
}

func (w *GameExpiryWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	logrus.Infof("Game expiry worker started (interval: %s)", w.interval)

	for {
		select {
		case <-ctx.Done():
			logrus.Info("Game expiry worker stopped")
			return
		case <-ticker.C:
			w.expireGames()
		}
	}
}

func (w *GameExpiryWorker) expireGames() {
	games, err := w.gameRepo.GetExpiredGames()
	if err != nil {
		logrus.Errorf("Failed to get expired games: %v", err)
		return
	}

	if len(games) == 0 {
		return
	}

	expiredCount := 0
	for _, game := range games {
		err := w.expireGame(game)
		if err != nil {
			logrus.Errorf("Failed to expire game %d: %v", game.Id, err)
			continue
		}
		expiredCount++
	}

	if expiredCount > 0 {
		logrus.Infof("Expired %d games with statistics updated", expiredCount)
	}
}

func (w *GameExpiryWorker) expireGame(game model.Game) error {
	return w.txManager.WithTransaction(func(tx *sqlx.Tx) error {
		lockedGame, err := w.gameRepo.GetGameById(tx, game.Id, true)
		if err != nil {
			return err
		}

		if lockedGame.Status != "in_progress" {
			return nil
		}

		if err := w.gameRepo.UpdateGameStatus(tx, game.Id, "lost"); err != nil {
			return err
		}

		guessCount, err := w.gameRepo.CountGuessesByGame(tx, game.Id)
		if err != nil {
			return err
		}

		statsUpdate := model.StatisticsUpdate{
			UserId:      game.UserId,
			Language:    game.Language,
			IsWin:       false,
			GuessCount:  guessCount,
			ScoreEarned: 0,
		}

		if err := w.statsUpdater.UpdateGameEndStats(tx, statsUpdate); err != nil {
			return err
		}

		return nil
	})
}
