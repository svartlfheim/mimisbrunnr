package scm

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

const postgresSCMIntegrationsTableName string = "scm_integrations"
const postgresSCMIntegrationsAccessTokensTableName string = "scm_integrations_access_tokens"

type postgresSCMIntegration struct {
	ID string `db:"id"`
	Name string `db:"name"`
	Type string `db:"type"`
	Endpoint string `db:"endpoint"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type postgresSCMIntegrationAccessToken struct {
	ID string `db:"id"`
	Name string `db:"name"`
	Token string `db:"token"`
	Active bool `db:"active"`
	SCMIntegrationID string `db:"scm_integration_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type PostgresRepository struct {
	logger zerolog.Logger
	cm connManager
}

func (r *PostgresRepository) toPostgresSCMIntegration(gh *SCMIntegration) *postgresSCMIntegration {
	return &postgresSCMIntegration{
		ID: gh.GetID().String(),
		Name: gh.GetName(),
		Type: string(gh.GetType()),
		Endpoint: gh.GetEndpoint(),
		CreatedAt: gh.GetCreationTime(),
		UpdatedAt: gh.GetLastUpdatedTime(),
	}
}

func (r *PostgresRepository) toPostgresSCMIntegrationAccessToken(token *AccessToken, gh *SCMIntegration) *postgresSCMIntegrationAccessToken {
	return &postgresSCMIntegrationAccessToken{
		ID: token.GetID().String(),
		Name: token.GetName(),
		Token: token.GetToken(),
		Active: token.IsActive(),
		SCMIntegrationID: gh.GetID().String(),
		CreatedAt: token.GetCreationTime(),
		UpdatedAt: token.GetLastUpdatedTime(),
	}
}

func (r *PostgresRepository) doCreateSCMIntegrationAccessToken(tx *sqlx.Tx, creds *AccessToken, gh *SCMIntegration) error {
	dbGhC := r.toPostgresSCMIntegrationAccessToken(creds, gh)

	insert := fmt.Sprintf(`
INSERT INTO %s 
(id, name, token, active, scm_integration_id, created_at, updated_at) 
VALUES 
(:id, :name, :token, :active, :scm_integration_id, :created_at, :updated_at);`,
postgresSCMIntegrationsAccessTokensTableName)

	_, err := tx.NamedExec(insert, dbGhC)

	return err
}

func (r *PostgresRepository) Create(gh *SCMIntegration, cred *AccessToken) error {
	conn, err := r.cm.GetConnection()

	if err != nil {
		return err
	}

	tx, err := conn.Beginx()

	if err != nil {
		return err
	}

	dbGh := r.toPostgresSCMIntegration(gh)

	insert := fmt.Sprintf(`
INSERT INTO %s 
(id, name, type, endpoint, created_at, updated_at) 
VALUES 
(:id, :name, :type, :endpoint, :created_at, :updated_at);`,
postgresSCMIntegrationsTableName)


	if _, err := tx.NamedExec(insert, dbGh); err != nil {
		return err
	}

	if err := r.doCreateSCMIntegrationAccessToken(tx, cred, gh); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil  {
			return err
		}

		return err
	}

	return nil
}