package postgres

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
)

type connManager interface {
	GetConnection() (*sqlx.DB, error)
}

const PostgresSCMIntegrationsTableName string = "scm_integrations"

type postgresSCMIntegration struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Type      string    `db:"type"`
	Token     string    `db:"token"`
	Endpoint  string    `db:"endpoint"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Repository struct {
	logger zerolog.Logger
	cm     connManager
}

func NewRepository(l zerolog.Logger, cm connManager) *Repository {
	return &Repository{
		logger: l,
		cm:     cm,
	}
}

func (r *Repository) toPostgresSCMIntegration(gh *models.SCMIntegration) *postgresSCMIntegration {
	return &postgresSCMIntegration{
		ID:        gh.GetID().String(),
		Name:      gh.GetName(),
		Type:      string(gh.GetType()),
		Token:     gh.GetToken(),
		Endpoint:  gh.GetEndpoint(),
		CreatedAt: gh.GetCreationTime(),
		UpdatedAt: gh.GetLastUpdatedTime(),
	}
}

func (r *Repository) Create(gh *models.SCMIntegration) error {
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
(id, name, type, token, endpoint, created_at, updated_at) 
VALUES 
(:id, :name, :type, :token, :endpoint, :created_at, :updated_at);`,
		PostgresSCMIntegrationsTableName)

	if _, err := tx.NamedExec(insert, dbGh); err != nil {
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


func (r *Repository) FindByName(name string) (*models.SCMIntegration, error) {
	conn, err := r.cm.GetConnection()

	if err != nil {
		return nil, err
	}

	// There is a unique key on name, so limit 1 is just covering our arse
	// And ensures that the code below will fetch the only row in the result
	q := fmt.Sprintf("SELECT * from %s WHERE name = $1 LIMIT 1", PostgresSCMIntegrationsTableName)
	res := postgresSCMIntegration{}
	err = conn.Get(&res, q, name)

	// An error is returned if the result set is empty
	// This is apparrently the only documented error that is returned
	// It doesn't have any specific type to check either
	if err != nil {
		return nil, nil
	}

	model := models.NewSCMIntegration(
		uuid.MustParse(res.ID),
		res.Name,
		models.SCMIntegrationType(res.Type),
		res.Endpoint,
		res.Token,
		res.CreatedAt,
		res.UpdatedAt,
	)


	return model, nil
}
