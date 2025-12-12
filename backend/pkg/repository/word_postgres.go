package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/jmoiron/sqlx"
)

type WordPostgres struct {
	db *sqlx.DB
}

func NewWordPostgres(db *sqlx.DB) *WordPostgres {
	return &WordPostgres{db: db}
}

func (r *WordPostgres) GetWordById(ctx context.Context, wordId int64) (*model.Word, error) {
	var word model.Word

	query := fmt.Sprintf("SELECT id, word, language, difficulty, is_active FROM %s WHERE id=$1", wordsTable)
	err := r.db.GetContext(ctx, &word, query, wordId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("word not found")
		}
		return nil, err
	}

	return &word, nil
}

func (r *WordPostgres) GetRandomWordByLanguage(ctx context.Context, language string) (*model.Word, error) {
	var word model.Word

	query := fmt.Sprintf("SELECT id, word, language, difficulty, is_active FROM %s WHERE language=$1 AND is_active=true ORDER BY RANDOM() LIMIT 1", wordsTable)
	err := r.db.GetContext(ctx, &word, query, language)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no words found for language: %s", language)
		}
		return nil, err
	}

	return &word, nil
}

func (r *WordPostgres) GetWordCountByLanguage(ctx context.Context, language string) (int, error) {
	var count int

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE language=$1 AND is_active=true", wordsTable)
	err := r.db.GetContext(ctx, &count, query, language)

	if err != nil {
		return 0, err
	}

	return count, nil
}
