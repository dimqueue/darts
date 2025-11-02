package migrations

import (
	"embed"

	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

//go:embed postgres/*
var dbMigrations embed.FS

var migrations = migrate.EmbedFileSystemMigrationSource{
	FileSystem: dbMigrations,
	Root:       "postgres",
}

func Up(db *sqlx.DB) error {
	if _, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up); err != nil {
		return err
	}
	return nil
}

func Down(db *sqlx.DB) error {
	if _, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Down); err != nil {
		return err
	}
	return nil
}
