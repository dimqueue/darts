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
	UpdateGameEndStats(ctx context.Context, q repository.Querier, update model.StatisticsUpdate) error
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

	logrus.WithField("interval", w.interval.String()).Info("game expiry worker started")

	for {
		select {
		case <-ctx.Done():
			logrus.Info("game expiry worker stopped")
			return
		case <-ticker.C:
			w.expireGames(ctx)
		}
	}
}

func (w *GameExpiryWorker) expireGames(ctx context.Context) {
	games, err := w.gameRepo.GetExpiredGames(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to get expired games")
		return
	}

	if len(games) == 0 {
		return
	}

	logrus.WithField("count", len(games)).Debug("found expired games")

	expiredCount := 0
	for _, game := range games {
		err := w.expireGame(ctx, game)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"game_id": game.Id,
				"user_id": game.UserId,
			}).WithError(err).Error("failed to expire game")
			continue
		}
		expiredCount++
	}

	if expiredCount > 0 {
		logrus.WithField("count", expiredCount).Info("expired games processed")
	}
}

func (w *GameExpiryWorker) expireGame(ctx context.Context, game model.Game) error {
	return w.txManager.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		lockedGame, err := w.gameRepo.GetGameById(ctx, tx, game.Id, true)
		if err != nil {
			return err
		}

		if lockedGame.Status != "in_progress" {
			return nil
		}

		if err := w.gameRepo.UpdateGameStatus(ctx, tx, game.Id, "lost"); err != nil {
			return err
		}

		guessCount, err := w.gameRepo.CountGuessesByGame(ctx, tx, game.Id)
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

		if err := w.statsUpdater.UpdateGameEndStats(ctx, tx, statsUpdate); err != nil {
			return err
		}

		return nil
	})
}
