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
	token TEXT NOT NULL,
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
			Id:   "create-projects-table",
			Name: "create projects table",
			Execute: func(tx *sqlx.Tx) (sql.Result, error) {
				createTable := `
CREATE TABLE projects(
	id uuid NOT NULL,
	name TEXT NOT NULL,
	path TEXT NOT NULL,
	scm_integration_id uuid NOT NULL,
	created_at TIMESTAMP WITHOUT TIME ZONE,
	updated_at TIMESTAMP WITHOUT TIME ZONE,
	PRIMARY KEY(id),
	UNIQUE(name),
	UNIQUE(path, scm_integration_id),
	CONSTRAINT fk_projects_scm_integration_id FOREIGN KEY(scm_integration_id) REFERENCES scm_integrations(id)
);
CREATE INDEX idx_projects_path ON projects(path);
CREATE INDEX idx_projects_name ON projects(name);`

				return tx.Exec(createTable)
			},
			Rollback: func(tx *sqlx.Tx) (sql.Result, error) {
				dropTable := `DROP TABLE projects;`

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
