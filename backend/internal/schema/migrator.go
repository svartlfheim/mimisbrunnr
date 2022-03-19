package schema

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/svartlfheim/gomigrator"
)

type connectionManager interface {
	GetConnection() (*sqlx.DB, error)
}

var postgresMigrations gomigrator.MigrationList = gomigrator.NewMigrationList(
	[]gomigrator.Migration{
		{
			Id:   "create-scm_integrations-table",
			Name: "create scm_integrations table",
			Execute: func(tx *sqlx.Tx) (sql.Result, error) {
				createTable := `
CREATE TABLE scm_integrations(
	id uuid NOT NULL,
	name TEXT NOT NULL,
	type TEXT NOT NULL,
	endpoint TEXT NOT NULL,
	created_at TIMESTAMP WITHOUT TIME ZONE,
	updated_at TIMESTAMP WITHOUT TIME ZONE,
	PRIMARY KEY(id),
	UNIQUE(name)
);
CREATE INDEX idx_scm_integrations_name ON scm_integrations(name);
CREATE INDEX idx_scm_integrations_type ON scm_integrations(type);`

				return tx.Exec(createTable)
			},
			Rollback: func(tx *sqlx.Tx) (sql.Result, error) {
				dropTable := `DROP TABLE scm_integrations;`

				return tx.Exec(dropTable)
			},
		},

		{
			Id:   "create-scm_integrations_access_tokens-table",
			Name: "create scm_integrations_access_tokens table",
			Execute: func(tx *sqlx.Tx) (sql.Result, error) {
				createTable := `
CREATE TABLE scm_integrations_access_tokens(
	id uuid NOT NULL,
	name TEXT NOT NULL,
	token TEXT NOT NULL,
	active BOOLEAN,
	integration_id uuid NOT NULL,
	created_at TIMESTAMP WITHOUT TIME ZONE,
	updated_at TIMESTAMP WITHOUT TIME ZONE,
	PRIMARY KEY(id),
	UNIQUE(name, integration_id),
	CONSTRAINT fk_scm_integrations_access_tokens_to_scm_integrations FOREIGN KEY(integration_id) REFERENCES scm_integrations(id)
);

CREATE INDEX idx_scm_integrations_access_tokens_name ON scm_integrations_access_tokens(name);
CREATE INDEX idx_scm_integrations_access_tokens_active ON scm_integrations_access_tokens(active);
CREATE INDEX idx_scm_integrations_access_tokens_integration_id ON scm_integrations_access_tokens(integration_id);`

				return tx.Exec(createTable)
			},
			Rollback: func(tx *sqlx.Tx) (sql.Result, error) {
				dropTable := `DROP TABLE scm_integrations_access_tokens;`

				return tx.Exec(dropTable)
			},
		},
	},
)

type hasSchema interface {
	GetRDBSchema() string
}

func NewMigrator(cm connectionManager, cfg hasSchema, l zerolog.Logger) (*gomigrator.Migrator, error) {
	conn, err := cm.GetConnection()

	if err != nil {
		return nil, err
	}

	return gomigrator.NewMigrator(conn, postgresMigrations, gomigrator.Opts{
		Schema:  cfg.GetRDBSchema(),
		Applyer: "some-name-for-now",
	}, l)
}
