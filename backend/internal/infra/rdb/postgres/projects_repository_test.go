package postgres_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/svartlfheim/mimisbrunnr/internal/infra/rdb/postgres"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/test/integration"
	postgresmocks "github.com/svartlfheim/mimisbrunnr/test/mocks/infra/rdb/postgres"
	zerologmocks "github.com/svartlfheim/mimisbrunnr/test/mocks/zerolog"
)

func buildProjectsRepo(t *testing.T) (*postgres.ProjectsRepository, *sqlx.DB, *zerologmocks.Zerologger) {
	l := zerologmocks.NewLogger()

	conn := integration.GetDatabaseConnectionOrFail(t)

	cm := &postgresmocks.ConnManager{}
	cm.EXPECT().GetConnection().Return(conn, nil)

	return postgres.NewProjectsRepository(l.Logger, cm), conn, l
}

func TestIntegration_Projects_Create(t *testing.T) {
	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ResetMigrationsOrFail(t)

	r, conn, _ := buildProjectsRepo(t)

	seed := seedSCMIntegrations(t, conn, 1)
	i := seed[0]

	createdAt := time.Now()
	updatedAt := createdAt.Add(1 * time.Hour)
	pID := "2b56d418-b1a9-4335-a2cb-575ea9f0c66d"
	p := &models.Project{
		ID:             uuid.MustParse(pID),
		Name:           "my-first-project",
		Path:           "myorg/myrepo",
		SCMIntegration: i,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}

	err := r.Create(p)

	assert.Nil(t, err)
	integration.AssertRowCount(
		t,
		fmt.Sprintf("SELECT * FROM %s WHERE id='%s';", postgres.ProjectsTableName, pID),
		1,
		conn,
		fmt.Sprintf("row was not inserted to '%s' with ID: %s", postgres.ProjectsTableName, pID),
	)
}

func TestIntegration_Projects_Find(t *testing.T) {
	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ResetMigrationsOrFail(t)

	r, conn, _ := buildProjectsRepo(t)

	seed := seedSCMIntegrations(t, conn, 1)
	i := seed[0]

	pSeeds := seedProjects(t, i, conn, 1)
	p := pSeeds[0]

	res, err := r.Find(p.GetID())

	require.Nil(t, err)
	require.NotNil(t, res)

	assert.Equal(t, p.ID, res.ID)
	assert.Equal(t, p.Name, res.Name)
	assert.Equal(t, p.Path, res.Path)
	assert.Equal(t, p.SCMIntegration.ID, i.ID)
}

func TestIntegration_Projects_Find_not_found(t *testing.T) {
	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ResetMigrationsOrFail(t)

	r, conn, _ := buildProjectsRepo(t)

	seed := seedSCMIntegrations(t, conn, 1)
	i := seed[0]
	pSeeds := seedProjects(t, i, conn, 1)
	p := pSeeds[0]

	res, err := r.Find(p.GetID())

	require.Nil(t, err)
	require.NotNil(t, res)
	assert.Equal(t, res.GetID(), p.GetID())

	res, err = r.Find(uuid.MustParse("bf832431-545c-4afc-90e9-87406c1ab0e9"))
	assert.Nil(t, err)
	assert.Nil(t, res)
}

func TestIntegration_Projects_FindByName(t *testing.T) {
	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ResetMigrationsOrFail(t)

	r, conn, _ := buildProjectsRepo(t)

	seed := seedSCMIntegrations(t, conn, 1)
	i := seed[0]

	pSeeds := seedProjects(t, i, conn, 1)
	p := pSeeds[0]

	res, err := r.FindByName(p.GetName())

	require.Nil(t, err)
	require.NotNil(t, res)

	assert.Equal(t, p.ID, res.ID)
	assert.Equal(t, p.Name, res.Name)
	assert.Equal(t, p.Path, res.Path)
	assert.Equal(t, p.SCMIntegration.ID, i.ID)
}

func TestIntegration_Projects_FindByName_not_found(t *testing.T) {
	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ResetMigrationsOrFail(t)

	r, conn, _ := buildProjectsRepo(t)

	seed := seedSCMIntegrations(t, conn, 1)
	i := seed[0]

	pSeeds := seedProjects(t, i, conn, 1)
	p := pSeeds[0]

	res, err := r.FindByName(p.GetName())

	require.Nil(t, err)
	assert.Equal(t, res.GetID(), p.ID)

	res, err = r.FindByName("not-created")
	require.Nil(t, err)
	assert.Nil(t, res)
}

func TestIntegration_Projects_FindByPathAndIntegrationID(t *testing.T) {
	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ResetMigrationsOrFail(t)

	r, conn, _ := buildProjectsRepo(t)

	seed := seedSCMIntegrations(t, conn, 2)

	firstIntegrationSeeds := seedProjects(t, seed[0], conn, 3)
	secondIntegrationSeeds := seedProjects(t, seed[1], conn, 3)

	integration.AssertRowCount(
		t,
		fmt.Sprintf("SELECT * FROM %s;", postgres.ProjectsTableName),
		6,
		conn,
		"not enough rows after seed",
	)

	toFind := firstIntegrationSeeds[0]
	res, err := r.FindByPathAndIntegrationID(toFind.GetPath(), seed[0].GetID())

	require.Nil(t, err)
	require.NotNil(t, res)

	assert.Equal(t, toFind.ID, res.ID)
	assert.Equal(t, toFind.Name, res.Name)
	assert.Equal(t, toFind.Path, res.Path)
	assert.Equal(t, toFind.SCMIntegration.ID, seed[0].ID)
	assert.Equal(t, toFind.CreatedAt.Unix(), res.CreatedAt.Unix())
	assert.Equal(t, toFind.UpdatedAt.Unix(), res.UpdatedAt.Unix())

	toFind = secondIntegrationSeeds[0]
	res, err = r.FindByPathAndIntegrationID(toFind.GetPath(), seed[1].GetID())

	require.Nil(t, err)
	require.NotNil(t, res)

	assert.Equal(t, toFind.ID, res.ID)
	assert.Equal(t, toFind.Name, res.Name)
	assert.Equal(t, toFind.Path, res.Path)
	assert.Equal(t, toFind.SCMIntegration.ID, seed[1].ID)
	assert.Equal(t, toFind.CreatedAt.Unix(), res.CreatedAt.Unix())
	assert.Equal(t, toFind.UpdatedAt.Unix(), res.UpdatedAt.Unix())
}
