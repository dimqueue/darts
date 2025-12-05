//go:build !dev

package seeds

import "github.com/jmoiron/sqlx"

func Run(db *sqlx.DB) error {
	return nil
}
