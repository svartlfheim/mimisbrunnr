package postgres

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
)

type connManager interface {
	GetConnection() (*sqlx.DB, error)
}

const PostgresSCMIntegrationsTableName string = "scm_integrations"
const PostgresSCMIntegrationsAccessTokensTableName string = "scm_integrations_access_tokens"

type postgresSCMIntegration struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Type      string    `db:"type"`
	Endpoint  string    `db:"endpoint"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type postgresSCMIntegrationAccessToken struct {
	ID               string    `db:"id"`
	Name             string    `db:"name"`
	Token            string    `db:"token"`
	Active           bool      `db:"active"`
	SCMIntegrationID string    `db:"integration_id"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

type Repository struct {
	logger zerolog.Logger
	cm     connManager
}

func NewRepository(l zerolog.Logger, cm connManager) *Repository {
	return &Repository{
		logger: l,
		cm: cm,
	}
}

func (r *Repository) toPostgresSCMIntegration(gh *models.SCMIntegration) *postgresSCMIntegration {
	return &postgresSCMIntegration{
		ID:        gh.GetID().String(),
		Name:      gh.GetName(),
		Type:      string(gh.GetType()),
		Endpoint:  gh.GetEndpoint(),
		CreatedAt: gh.GetCreationTime(),
		UpdatedAt: gh.GetLastUpdatedTime(),
	}
}

func (r *Repository) toPostgresSCMIntegrationAccessToken(token *models.SCMAccessToken, gh *models.SCMIntegration) *postgresSCMIntegrationAccessToken {
	return &postgresSCMIntegrationAccessToken{
		ID:               token.GetID().String(),
		Name:             token.GetName(),
		Token:            token.GetToken(),
		Active:           token.IsActive(),
		SCMIntegrationID: gh.GetID().String(),
		CreatedAt:        token.GetCreationTime(),
		UpdatedAt:        token.GetLastUpdatedTime(),
	}
}

func (r *Repository) doCreateSCMIntegrationAccessToken(tx *sqlx.Tx, creds *models.SCMAccessToken, gh *models.SCMIntegration) error {
	dbGhC := r.toPostgresSCMIntegrationAccessToken(creds, gh)

	insert := fmt.Sprintf(`
INSERT INTO %s 
(id, name, token, active, integration_id, created_at, updated_at) 
VALUES 
(:id, :name, :token, :active, :integration_id, :created_at, :updated_at);`,
		PostgresSCMIntegrationsAccessTokensTableName)

	_, err := tx.NamedExec(insert, dbGhC)

	return err
}

func (r *Repository) Create(gh *models.SCMIntegration, cred *models.SCMAccessToken) error {
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
		PostgresSCMIntegrationsTableName)

	if _, err := tx.NamedExec(insert, dbGh); err != nil {
		return err
	}

	if err := r.doCreateSCMIntegrationAccessToken(tx, cred, gh); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	return nil
}
