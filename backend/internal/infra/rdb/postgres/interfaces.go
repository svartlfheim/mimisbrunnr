package postgres

import "github.com/jmoiron/sqlx"

type connManager interface {
	GetConnection() (*sqlx.DB, error)
}
