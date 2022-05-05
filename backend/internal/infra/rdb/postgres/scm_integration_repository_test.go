package postgres_test

import (
	"fmt"
	"math/rand"
	"strconv"
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

func getRandomSCMIntegrationType() models.SCMIntegrationType {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(models.AvailableSCMIntegrationTypes()))
	pick := models.AvailableSCMIntegrationTypes()[randomIndex]

	return pick
}

func buildSCMIntegrationsRepo(t *testing.T) (*postgres.SCMIntegrationsRepository, *sqlx.DB, *zerologmocks.Zerologger) {
	l := zerologmocks.NewLogger()

	conn := integration.GetDatabaseConnectionOrFail(t)

	cm := &postgresmocks.ConnManager{}
	cm.EXPECT().GetConnection().Return(conn, nil)

	return postgres.NewSCMIntegrationsRepository(l.Logger, cm), conn, l
}

func TestIntegration_SCMIntegrations_Create(t *testing.T) {
	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ResetMigrationsOrFail(t)

	r, conn, _ := buildSCMIntegrationsRepo(t)

	createdAt := time.Now()
	updatedAt := createdAt.Add(1 * time.Hour)

	scmIntegrationID := "d0170c68-38cb-40b7-a6db-4b70210c60d7"
	scmIntegration := &models.SCMIntegration{
		ID:        uuid.MustParse(scmIntegrationID),
		Name:      "my-first-integration",
		Type:      models.GithubType,
		Endpoint:  "http://fake.example.local",
		Token:     "mytoken",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	err := r.Create(scmIntegration)

	assert.Nil(t, err)
	integration.AssertRowCount(
		t,
		fmt.Sprintf("SELECT * FROM %s WHERE id='%s';", postgres.SCMIntegrationsTableName, scmIntegrationID),
		1,
		conn,
		fmt.Sprintf("row was not inserted to '%s' with ID: %s", postgres.SCMIntegrationsTableName, scmIntegrationID),
	)
}

func TestIntegration_SCMIntegrations_CreateThenFind(t *testing.T) {
	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ResetMigrationsOrFail(t)

	r, conn, _ := buildSCMIntegrationsRepo(t)

	createdAt := time.Now()
	updatedAt := createdAt.Add(1 * time.Hour)

	scmIntegrationID := "d0170c68-38cb-40b7-a6db-4b70210c60d7"
	scmIntegration := &models.SCMIntegration{
		ID:        uuid.MustParse(scmIntegrationID),
		Name:      "my-first-integration",
		Type:      models.GithubType,
		Endpoint:  "http://fake.example.local",
		Token:     "mytoken",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	err := r.Create(scmIntegration)

	assert.Nil(t, err)
	integration.AssertRowCount(
		t,
		fmt.Sprintf("SELECT * FROM %s WHERE id='%s';", postgres.SCMIntegrationsTableName, scmIntegrationID),
		1,
		conn,
		fmt.Sprintf("row was not inserted to '%s' with ID: %s", postgres.SCMIntegrationsTableName, scmIntegrationID),
	)

	m, err := r.Find(scmIntegration.ID)

	require.Nil(t, err)

	assert.Equal(t, scmIntegration.ID, m.ID)
	assert.Equal(t, scmIntegration.Name, m.Name)
	assert.Equal(t, scmIntegration.Type, m.Type)
	assert.Equal(t, scmIntegration.Endpoint, m.Endpoint)
	assert.Equal(t, scmIntegration.Token, m.Token)
	assert.Equal(t, scmIntegration.CreatedAt.Unix(), m.CreatedAt.Unix())
	assert.Equal(t, scmIntegration.UpdatedAt.Unix(), m.UpdatedAt.Unix())
}

func TestIntegration_SCMIntegrations_Find_not_found(t *testing.T) {
	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ResetMigrationsOrFail(t)

	r, conn, _ := buildSCMIntegrationsRepo(t)

	createdAt := time.Now()
	updatedAt := createdAt.Add(1 * time.Hour)

	scmIntegrationID := "d0170c68-38cb-40b7-a6db-4b70210c60d7"
	scmIntegration := &models.SCMIntegration{
		ID:        uuid.MustParse(scmIntegrationID),
		Name:      "my-first-integration",
		Type:      models.GithubType,
		Endpoint:  "http://fake.example.local",
		Token:     "mytoken",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	err := r.Create(scmIntegration)

	assert.Nil(t, err)
	integration.AssertRowCount(
		t,
		fmt.Sprintf("SELECT * FROM %s WHERE id='%s';", postgres.SCMIntegrationsTableName, scmIntegrationID),
		1,
		conn,
		fmt.Sprintf("row was not inserted to '%s' with ID: %s", postgres.SCMIntegrationsTableName, scmIntegrationID),
	)

	m, err := r.Find(scmIntegration.ID)

	require.Nil(t, err)
	assert.Equal(t, scmIntegration.ID, m.ID)

	m, err = r.Find(uuid.MustParse("bf832431-545c-4afc-90e9-87406c1ab0e9"))
	require.Nil(t, err)
	assert.Nil(t, m)
}

func TestIntegration_SCMIntegrations_CreateThenFindByName(t *testing.T) {
	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ResetMigrationsOrFail(t)

	r, conn, _ := buildSCMIntegrationsRepo(t)

	createdAt := time.Now()
	updatedAt := createdAt.Add(1 * time.Hour)

	scmIntegrationID := "d0170c68-38cb-40b7-a6db-4b70210c60d7"
	scmIntegration := &models.SCMIntegration{
		ID:        uuid.MustParse(scmIntegrationID),
		Name:      "my-first-integration",
		Type:      models.GithubType,
		Endpoint:  "http://fake.example.local",
		Token:     "mytoken",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	err := r.Create(scmIntegration)

	assert.Nil(t, err)
	integration.AssertRowCount(
		t,
		fmt.Sprintf("SELECT * FROM %s WHERE id='%s';", postgres.SCMIntegrationsTableName, scmIntegrationID),
		1,
		conn,
		fmt.Sprintf("row was not inserted to '%s' with ID: %s", postgres.SCMIntegrationsTableName, scmIntegrationID),
	)

	m, err := r.FindByName("my-first-integration")

	require.Nil(t, err)

	assert.Equal(t, scmIntegration.ID, m.ID)
	assert.Equal(t, scmIntegration.Name, m.Name)
	assert.Equal(t, scmIntegration.Type, m.Type)
	assert.Equal(t, scmIntegration.Endpoint, m.Endpoint)
	assert.Equal(t, scmIntegration.Token, m.Token)
	assert.Equal(t, scmIntegration.CreatedAt.Unix(), m.CreatedAt.Unix())
	assert.Equal(t, scmIntegration.UpdatedAt.Unix(), m.UpdatedAt.Unix())
}

func TestIntegration_SCMIntegrations_FindByName_not_found(t *testing.T) {
	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ResetMigrationsOrFail(t)

	r, conn, _ := buildSCMIntegrationsRepo(t)

	createdAt := time.Now()
	updatedAt := createdAt.Add(1 * time.Hour)

	scmIntegrationID := "d0170c68-38cb-40b7-a6db-4b70210c60d7"
	scmIntegration := &models.SCMIntegration{
		ID:        uuid.MustParse(scmIntegrationID),
		Name:      "my-first-integration",
		Type:      models.GithubType,
		Endpoint:  "http://fake.example.local",
		Token:     "mytoken",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	err := r.Create(scmIntegration)

	assert.Nil(t, err)
	integration.AssertRowCount(
		t,
		fmt.Sprintf("SELECT * FROM %s WHERE id='%s';", postgres.SCMIntegrationsTableName, scmIntegrationID),
		1,
		conn,
		fmt.Sprintf("row was not inserted to '%s' with ID: %s", postgres.SCMIntegrationsTableName, scmIntegrationID),
	)

	m, err := r.FindByName("my-first-integration")

	require.Nil(t, err)
	assert.Equal(t, scmIntegration.ID, m.ID)

	m, err = r.FindByName("not-created")
	require.Nil(t, err)
	assert.Nil(t, m)
}

func TestIntegration_SCMIntegrations_Count_no_records(t *testing.T) {

	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ResetMigrationsOrFail(t)

	r, _, _ := buildSCMIntegrationsRepo(t)

	total, err := r.Count()

	assert.Nil(t, err)
	assert.Equal(t, 0, total)
}

func TestIntegration_SCMIntegrations_Count_some_records(t *testing.T) {

	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ResetMigrationsOrFail(t)

	r, _, _ := buildSCMIntegrationsRepo(t)

	createdAt := time.Now()
	updatedAt := createdAt.Add(1 * time.Hour)

	for i := 0; i < 37; i++ {
		scmIntegration := &models.SCMIntegration{
			ID:        uuid.New(),
			Name:      fmt.Sprintf("my-integration-%d", i),
			Type:      getRandomSCMIntegrationType(),
			Endpoint:  "http://fake.example.local",
			Token:     "mytoken",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		err := r.Create(scmIntegration)

		require.Nil(t, err, "error seeding data tocount")
	}

	total, err := r.Count()

	assert.Nil(t, err)
	assert.Equal(t, 37, total)
}

func TestIntegration_SCMIntegrations_Paginate_single_page(t *testing.T) {

	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ResetMigrationsOrFail(t)

	r, _, _ := buildSCMIntegrationsRepo(t)

	createdAt := time.Now()
	updatedAt := createdAt.Add(1 * time.Hour)

	for i := 0; i < 50; i++ {
		scmIntegration := &models.SCMIntegration{
			ID:        uuid.New(),
			Name:      fmt.Sprintf("my-integration-%d", i),
			Type:      getRandomSCMIntegrationType(),
			Endpoint:  "http://fake.example.local",
			Token:     "mytoken",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		err := r.Create(scmIntegration)

		require.Nil(t, err, "error seeding data tocount")
	}

	items, err := r.Paginate(1, 50)

	assert.Nil(t, err)
	assert.Equal(t, 50, len(items))
}

func TestIntegration_SCMIntegrations_Paginate_multiple_pages(t *testing.T) {

	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ResetMigrationsOrFail(t)

	r, _, _ := buildSCMIntegrationsRepo(t)

	createdAt := time.Now()
	updatedAt := createdAt.Add(1 * time.Hour)

	for i := 0; i < 50; i++ {
		nameIndex := strconv.Itoa(i)

		if i < 10 {
			nameIndex = "0" + nameIndex
		}

		scmIntegration := &models.SCMIntegration{
			ID:        uuid.New(),
			Name:      fmt.Sprintf("my-integration-%s", nameIndex),
			Type:      getRandomSCMIntegrationType(),
			Endpoint:  "http://fake.example.local",
			Token:     "mytoken",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		err := r.Create(scmIntegration)

		require.Nil(t, err, "error seeding data tocount")
	}

	items, err := r.Paginate(1, 25)

	assert.Nil(t, err)
	assert.Equal(t, 25, len(items))
	assert.Equal(t, items[0].GetName(), "my-integration-00")
	assert.Equal(t, items[24].GetName(), "my-integration-24")

	items, err = r.Paginate(2, 25)

	assert.Nil(t, err)
	assert.Equal(t, 25, len(items))
	assert.Equal(t, items[0].GetName(), "my-integration-25")
	assert.Equal(t, items[24].GetName(), "my-integration-49")
}

func TestIntegration_SCMIntegrations_Paginate_multiple_pages_random_cutoff(t *testing.T) {

	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ResetMigrationsOrFail(t)

	r, _, _ := buildSCMIntegrationsRepo(t)

	createdAt := time.Now()
	updatedAt := createdAt.Add(1 * time.Hour)

	for i := 0; i < 43; i++ {
		nameIndex := strconv.Itoa(i)

		if i < 10 {
			nameIndex = "0" + nameIndex
		}

		scmIntegration := &models.SCMIntegration{
			ID:        uuid.New(),
			Name:      fmt.Sprintf("my-integration-%s", nameIndex),
			Type:      getRandomSCMIntegrationType(),
			Endpoint:  "http://fake.example.local",
			Token:     "mytoken",
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		err := r.Create(scmIntegration)

		require.Nil(t, err, "error seeding data tocount")
	}

	items, err := r.Paginate(1, 25)

	assert.Nil(t, err)
	assert.Equal(t, 25, len(items))
	assert.Equal(t, items[0].GetName(), "my-integration-00")
	assert.Equal(t, items[24].GetName(), "my-integration-24")

	items, err = r.Paginate(2, 25)

	assert.Nil(t, err)
	assert.Equal(t, 18, len(items))
	assert.Equal(t, items[0].GetName(), "my-integration-25")
	assert.Equal(t, items[17].GetName(), "my-integration-42")
}
