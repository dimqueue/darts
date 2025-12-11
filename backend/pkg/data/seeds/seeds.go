//go:build dev

package seeds

import (
	"embed"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

//go:embed postgres/*
var seedFiles embed.FS

func Run(db *sqlx.DB) error {
	slog.Info("Loading seed data (dev mode)...")

	content, err := seedFiles.ReadFile("postgres/sample-data.sql")
	if err != nil {
		return err
	}

	if _, err := db.Exec(string(content)); err != nil {
		return err
	}

	slog.Info("Seed data loaded successfully")
	return nil
}
