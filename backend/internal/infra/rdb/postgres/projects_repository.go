package postgres

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

type ProjectsRepository struct {
	logger zerolog.Logger
	cm     connManager
}

func (r *ProjectsRepository) Create(p *models.Project) error {
	conn, err := r.cm.GetConnection()

	if err != nil {
		return err
	}

	tx, err := conn.Beginx()

	if err != nil {
		return err
	}

	dbP := toDBProject(p)

	insert := fmt.Sprintf(`
INSERT INTO %s 
(id, name, path, scm_integration_id, created_at, updated_at) 
VALUES 
(:id, :name, :path, :scm_integration_id, :created_at, :updated_at);`,
		ProjectsTableName)

	if _, err := tx.NamedExec(insert, dbP); err != nil {
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

func (r *ProjectsRepository) selectWithIntegration() string {
	return fmt.Sprintf(`
	SELECT
		p.id,
		p.name,
		p.path,
		p.created_at,
		p.updated_at,

		scm.id AS scmId,
		scm.name AS scmName,
		scm.type AS scmType,
		scm.token AS scmToken,
		scm.endpoint AS scmEndpoint,
		scm.created_at AS scmCreatedAt,
		scm.updated_at AS scmUpdatedAt
	FROM
		%s p
	INNER JOIN %s scm
		ON p.scm_integration_id = scm.id
	`, ProjectsTableName, SCMIntegrationsTableName)
}

func (r *ProjectsRepository) hydrateRowWithIntegrationToProject(rows *sqlx.Rows) (*models.Project, error) {
	var pID, pName, pPath, iID, iName, iType, iToken, iEndpoint string
	var pCreatedAt, pUpdatedAt, iCreatedAt, iUpdatedAt time.Time

	if err := rows.Scan(&pID, &pName, &pPath, &pCreatedAt, &pUpdatedAt, &iID, &iName, &iType, &iToken, &iEndpoint, &iCreatedAt, &iUpdatedAt); err != nil {
		return nil, err
	}

	dbI := scmIntegration{
		ID:        iID,
		Name:      iName,
		Type:      iType,
		Token:     iToken,
		Endpoint:  iEndpoint,
		CreatedAt: iCreatedAt,
		UpdatedAt: iUpdatedAt,
	}

	dbP := project{
		ID:        pID,
		Name:      pName,
		Path:      pPath,
		CreatedAt: pCreatedAt,
		UpdatedAt: pUpdatedAt,
	}

	return dbP.ToDomainModel(dbI.ToDomainModel()), nil
}

func (r *ProjectsRepository) hydrateRowsWithIntegrationToProject(rows *sqlx.Rows) ([]*models.Project, error) {
	integrations := map[string]*models.SCMIntegration{}
	results := []*models.Project{}
	for rows.Next() {
		var pID, pName, pPath, iID, iName, iType, iToken, iEndpoint string
		var pCreatedAt, pUpdatedAt, iCreatedAt, iUpdatedAt time.Time

		if err := rows.Scan(&pID, &pName, &pPath, &pCreatedAt, &pUpdatedAt, &iID, &iName, &iType, &iToken, &iEndpoint, &iCreatedAt, &iUpdatedAt); err != nil {
			return []*models.Project{}, err
		}

		var i *models.SCMIntegration

		if inMap, found := integrations[iID]; found {
			i = inMap
		} else {
			dbI := scmIntegration{
				ID:        iID,
				Name:      iName,
				Type:      iType,
				Token:     iToken,
				Endpoint:  iEndpoint,
				CreatedAt: iCreatedAt,
				UpdatedAt: iUpdatedAt,
			}
			i = dbI.ToDomainModel()
			integrations[iID] = i
		}

		dbP := project{
			ID:        pID,
			Name:      pName,
			Path:      pPath,
			CreatedAt: pCreatedAt,
			UpdatedAt: pUpdatedAt,
		}

		results = append(results, dbP.ToDomainModel(i))
	}

	return results, nil
}

func (r *ProjectsRepository) Find(id uuid.UUID) (*models.Project, error) {
	conn, err := r.cm.GetConnection()

	if err != nil {
		return nil, err
	}

	q := r.selectWithIntegration() + `
WHERE 
	p.id = $1 
LIMIT 1`

	rows, err := conn.Queryx(q, id.String())

	if err != nil {
		return nil, err
	}
	// Should only be 1, see limit
	found := rows.Next()

	if !found {
		return nil, nil
	}

	return r.hydrateRowWithIntegrationToProject(rows)
}

func (r *ProjectsRepository) FindByName(name string) (*models.Project, error) {
	conn, err := r.cm.GetConnection()

	if err != nil {
		return nil, err
	}

	q := r.selectWithIntegration() + `
WHERE 
	p.name = $1 
LIMIT 1`

	rows, err := conn.Queryx(q, name)

	if err != nil {
		return nil, err
	}
	// Should only be 1, see limit
	found := rows.Next()

	if !found {
		return nil, nil
	}

	return r.hydrateRowWithIntegrationToProject(rows)
}

func (r *ProjectsRepository) FindByPathAndIntegrationID(path string, integrationID uuid.UUID) (*models.Project, error) {
	conn, err := r.cm.GetConnection()

	if err != nil {
		return nil, err
	}

	q := r.selectWithIntegration() + `
WHERE 
	p.path = $1 AND
	p.scm_integration_id= $2
LIMIT 1`

	rows, err := conn.Queryx(q, path, integrationID.String())

	if err != nil {
		return nil, err
	}
	// Should only be 1, see limit
	found := rows.Next()

	if !found {
		return nil, nil
	}

	return r.hydrateRowWithIntegrationToProject(rows)
}

func (r *ProjectsRepository) Count() (int, error) {
	conn, err := r.cm.GetConnection()

	if err != nil {
		return -1, err
	}

	q := fmt.Sprintf("SELECT COUNT(1) from %s", ProjectsTableName)
	row := conn.QueryRow(q)

	var count int

	if err = row.Scan(&count); err != nil {
		return -1, err
	}

	return count, nil
}

func (r *ProjectsRepository) Paginate(page int, limit int) ([]*models.Project, error) {
	conn, err := r.cm.GetConnection()

	if err != nil {
		return []*models.Project{}, err
	}

	offset := (page - 1) * limit
	q := r.selectWithIntegration() + `
ORDER BY p.name ASC
OFFSET $1
LIMIT $2
`
	rows, err := conn.Queryx(q, offset, limit)

	if err != nil {
		return []*models.Project{}, err
	}

	return r.hydrateRowsWithIntegrationToProject(rows)
}

func (r *ProjectsRepository) dbColumnByStructFieldName(f string) (string, error) {
	rv := reflect.TypeOf(project{})
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

func (r *ProjectsRepository) buildPatchQuery(id uuid.UUID, cs *models.ChangeSet) (string, []interface{}, error) {
	columns := []string{}
	values := []interface{}{}

	cs.RegisterChange("UpdatedAt", time.Now().UTC())

	for k, val := range cs.Changes {
		if k == "CreatedAt" {
			// this should never be changed
			continue
		}
		colName, err := r.dbColumnByStructFieldName(k)

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
		ProjectsTableName,
		strings.Join(columns, ","),
		(len(values)),
	), values, nil
}

func (r *ProjectsRepository) Patch(id uuid.UUID, cs *models.ChangeSet) (*models.Project, error) {
	if cs.IsEmpty() {
		return r.Find(id)
	}

	conn, err := r.cm.GetConnection()

	if err != nil {
		return nil, err
	}

	q, args, err := r.buildPatchQuery(id, cs)

	if err != nil {
		return nil, err
	}

	if _, err := conn.Exec(q, args...); err != nil {
		return nil, err
	}

	return r.Find(id)
}

func (r *ProjectsRepository) Delete(id uuid.UUID) error {
	conn, err := r.cm.GetConnection()

	if err != nil {
		return err
	}

	q := fmt.Sprintf("DELETE from %s WHERE id = $1", ProjectsTableName)
	_, err = conn.Exec(q, id.String())

	return err
}

func NewProjectsRepository(l zerolog.Logger, cm connManager) *ProjectsRepository {
	return &ProjectsRepository{
		logger: l,
		cm:     cm,
	}
}
