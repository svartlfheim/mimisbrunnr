package scm

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/svartlfheim/mimisbrunnr/internal/scm/postgres"
)

type connManager interface {
	GetConnection() (*sqlx.DB, error)
}

func NewPostgresRepository(cm connManager, l zerolog.Logger) *postgres.Repository {
	return postgres.NewRepository(l, cm)
}
