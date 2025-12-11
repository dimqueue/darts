package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/dimqueue/darts/pkg/model"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(ctx context.Context, q Querier, user model.User) (int64, error) {
	var id int64

	query := fmt.Sprintf("INSERT INTO %s (name,username,password_hash) VALUES ($1,$2,$3) RETURNING id", usersTable)
	row := q.QueryRowContext(ctx, query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, ErrDuplicateKey
		}
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT id, password_hash FROM %s WHERE username=$1", usersTable)
	err := r.db.GetContext(ctx, &user, query, username)
	return user, err
}
