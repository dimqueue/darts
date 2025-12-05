//go:build dev

package seeds

import (
	"embed"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

//go:embed postgres/*
var seedFiles embed.FS

func Run(db *sqlx.DB) error {
	logrus.Info("Loading seed data (dev mode)...")

	content, err := seedFiles.ReadFile("postgres/sample-data.sql")
	if err != nil {
		return err
	}

	if _, err := db.Exec(string(content)); err != nil {
		return err
	}

	logrus.Info("Seed data loaded successfully")
	return nil
}
