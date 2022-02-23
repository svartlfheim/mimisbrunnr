package scm

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type connManager interface {
	GetConnection() (*sqlx.DB, error)
}

type repository interface {
	Create(*SCMIntegration, *AccessToken) error
}

func NewPostgresRepository(cm connManager, l zerolog.Logger) *PostgresRepository {
	return &PostgresRepository{
		cm: cm,
		logger: l,
	}
}
