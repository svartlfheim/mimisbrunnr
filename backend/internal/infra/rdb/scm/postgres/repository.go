package scmpostgres

import (
	"fmt"
	"reflect"
	"strings"
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

func (pI postgresSCMIntegration) ToDomainModel() *models.SCMIntegration {
	return models.NewSCMIntegration(
		uuid.MustParse(pI.ID),
		pI.Name,
		models.SCMIntegrationType(pI.Type),
		pI.Endpoint,
		pI.Token,
		pI.CreatedAt,
		pI.UpdatedAt,
	)
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

	return res.ToDomainModel(), nil
}

func (r *Repository) Find(id uuid.UUID) (*models.SCMIntegration, error) {
	conn, err := r.cm.GetConnection()

	if err != nil {
		return nil, err
	}

	// There is a unique key on name, so limit 1 is just covering our arse
	// And ensures that the code below will fetch the only row in the result
	q := fmt.Sprintf("SELECT * from %s WHERE id = $1 LIMIT 1", PostgresSCMIntegrationsTableName)
	res := postgresSCMIntegration{}
	err = conn.Get(&res, q, id.String())

	// An error is returned if the result set is empty
	// This is apparrently the only documented error that is returned
	// It doesn't have any specific type to check either
	if err != nil {
		return nil, nil
	}

	return res.ToDomainModel(), nil
}

func (r *Repository) Count() (int, error) {
	conn, err := r.cm.GetConnection()

	if err != nil {
		return -1, err
	}

	q := fmt.Sprintf("SELECT COUNT(1) from %s", PostgresSCMIntegrationsTableName)
	row := conn.QueryRow(q)

	var count int

	if err = row.Scan(&count); err != nil {
		return -1, err
	}

	return count, nil
}

func (r *Repository) Paginate(page int, limit int) ([]*models.SCMIntegration, error) {
	conn, err := r.cm.GetConnection()

	if err != nil {
		return []*models.SCMIntegration{}, err
	}

	offset := (page - 1) * limit
	q := fmt.Sprintf("SELECT * from %s ORDER BY name OFFSET $1 LIMIT $2", PostgresSCMIntegrationsTableName)
	rows, err := conn.Queryx(q, offset, limit)

	if err != nil {
		return []*models.SCMIntegration{}, err
	}

	dbResults := []postgresSCMIntegration{}
	for rows.Next() {
		res := postgresSCMIntegration{}

		if err := rows.StructScan(&res); err != nil {
			return []*models.SCMIntegration{}, err
		}

		dbResults = append(dbResults, res)
	}

	results := []*models.SCMIntegration{}

	for _, dbRes := range dbResults {
		results = append(results, dbRes.ToDomainModel())
	}

	return results, nil
}

func dbColumnByStructFieldName(f string) (string, error) {
	rv := reflect.TypeOf(postgresSCMIntegration{})
	frv, found := rv.FieldByName(f)

	if !found {
		return "", fmt.Errorf("field %s could not be mapped to db column", f)
	}

	columnValue, found := frv.Tag.Lookup("db")

	if !found {
		return "", fmt.Errorf("field %s could not be mapped to db column", f)
	}

	return columnValue, nil
}

func buildPatchQuery(id uuid.UUID, cs *models.ChangeSet) (string, []interface{}, error) {
	columns := []string{}
	values := []interface{}{}

	cs.RegisterChange("UpdatedAt", time.Now())

	for k, val := range cs.Changes {
		if k == "CreatedAt" {
			// this should never be changed
			continue
		}
		colName, err := dbColumnByStructFieldName(k)

		if err != nil {
			return "", values, err
		}

		setColumn := fmt.Sprintf(
			"%s = $%d",
			colName,
			(len(columns) + 1), // this will be inserted to the next index 1
		)
		columns = append(columns, setColumn)
		values = append(values, val)
	}

	values = append(values, id.String())

	return fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = $%d",
		PostgresSCMIntegrationsTableName,
		strings.Join(columns, ","),
		(len(values)),
	), values, nil
}

func (r *Repository) Patch(id uuid.UUID, cs *models.ChangeSet) (*models.SCMIntegration, error) {
	if cs.IsEmpty() {
		return r.Find(id)
	}

	conn, err := r.cm.GetConnection()

	if err != nil {
		return nil, err
	}

	q, args, err := buildPatchQuery(id, cs)

	if err != nil {
		return nil, err
	}

	if _, err := conn.Exec(q, args...); err != nil {
		return nil, err
	}

	return r.Find(id)
}

func (r *Repository) Delete(id uuid.UUID) error {
	conn, err := r.cm.GetConnection()

	if err != nil {
		return err
	}

	q := fmt.Sprintf("DELETE from %s WHERE id = $1", PostgresSCMIntegrationsTableName)
	_, err = conn.Exec(q, id.String())

	return err
}
