package postgres_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/svartlfheim/mimisbrunnr/internal/infra/rdb/postgres"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
)

func seedProjects(t *testing.T, m *models.SCMIntegration, conn *sqlx.DB, count int) []*models.Project {

	createdAt := time.Now()
	updatedAt := createdAt.Add(1 * time.Hour)

	tx, err := conn.Beginx()

	require.Nil(t, err)

	rows := []*models.Project{}
	for i := 0; i < count; i++ {
		nameIndex := strconv.Itoa(i)

		if i < 10 {
			nameIndex = "0" + nameIndex
		}
		r := models.NewProject(
			uuid.New(),
			fmt.Sprintf("%s-my-project-%s", m.GetName(), nameIndex),
			fmt.Sprintf("my/path/index-%s", nameIndex),
			m,
			createdAt,
			updatedAt,
		)

		rows = append(rows, r)

		e := `
INSERT INTO %s 
(id, name, path, scm_integration_id, created_at, updated_at) 
VALUES 
($1, $2, $3, $4, $5, $6);
		`
		_, err := tx.Exec(
			fmt.Sprintf(e, postgres.ProjectsTableName),
			r.GetID().String(),
			r.GetName(),
			r.GetPath(),
			r.GetSCMIntegration().GetID().String(),
			r.GetCreationTime(),
			r.GetLastUpdatedTime(),
		)

		require.Nil(t, err, "project seeding error")
	}

	err = tx.Commit()

	require.Nil(t, err, "project seeding error (commit)")

	return rows
}

func seedSCMIntegrations(t *testing.T, conn *sqlx.DB, count int) []*models.SCMIntegration {

	createdAt := time.Now()
	updatedAt := createdAt.Add(1 * time.Hour)

	tx, err := conn.Beginx()

	require.Nil(t, err)

	rows := []*models.SCMIntegration{}
	for i := 0; i < count; i++ {
		nameIndex := strconv.Itoa(i)

		if i < 10 {
			nameIndex = "0" + nameIndex
		}
		r := models.NewSCMIntegration(
			uuid.New(),
			fmt.Sprintf("my-integration-%s", nameIndex),
			getRandomSCMIntegrationType(),
			"http://fake.example.local",
			"mytoken",
			createdAt,
			updatedAt,
		)

		rows = append(rows, r)

		e := `
INSERT INTO %s 
(id, name, type, token, endpoint, created_at, updated_at) 
VALUES 
($1, $2, $3, $4, $5, $6, $7);
		`
		_, err := tx.Exec(
			fmt.Sprintf(e, postgres.SCMIntegrationsTableName),
			r.GetID().String(),
			r.GetName(),
			string(r.GetType()),
			r.GetToken(),
			r.GetEndpoint(),
			r.GetCreationTime(),
			r.GetLastUpdatedTime(),
		)

		require.Nil(t, err, "integration seeding error")
	}

	err = tx.Commit()

	require.Nil(t, err, "integration seeding error (commit)")

	return rows
}
