package postgres_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/scm/postgres"
	"github.com/svartlfheim/mimisbrunnr/test/integration"
	scmpostgresmocks "github.com/svartlfheim/mimisbrunnr/test/mocks/scm/postgres"
	zerologmocks "github.com/svartlfheim/mimisbrunnr/test/mocks/zerolog"
)

func buildRepo(t *testing.T) (*postgres.Repository, *sqlx.DB, *zerologmocks.Zerologger) {
	l := zerologmocks.NewLogger()

	conn := integration.GetDatabaseConnectionOrFail(t)

	cm := &scmpostgresmocks.ConnManager{}
	cm.EXPECT().GetConnection().Return(conn, nil)

	return postgres.NewRepository(l.Logger, cm), conn, l
}

func TestIntegration_Create(t *testing.T) {
	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ResetMigrationsOrFail(t)

	r, conn, _ := buildRepo(t)

	createdAt := time.Now()
	updatedAt := createdAt.Add(1 * time.Hour)

	scmIntegrationID := "d0170c68-38cb-40b7-a6db-4b70210c60d7"
	scmIntegration := &models.SCMIntegration{
		ID: uuid.MustParse(scmIntegrationID),
		Name: "my-first-integration",
		Type: models.GithubType,
		Endpoint: "http://fake.example.local",
		AccessTokens: []*models.SCMAccessToken{},
		CreatedAt:createdAt,
		UpdatedAt: updatedAt,
	}

	tokenID := "7c544640-a3eb-4428-9cf2-093eb849a53f"
	token := &models.SCMAccessToken{
		ID: uuid.MustParse(tokenID),
		Name: "my-token",
		Active: true,
		Token: "some-token",
		CreatedAt:createdAt,
		UpdatedAt: updatedAt,
	}

	err := r.Create(scmIntegration, token)

	assert.Nil(t, err)
	integration.AssertRowCount(
		t, 
		fmt.Sprintf("SELECT * FROM %s WHERE id='%s';", postgres.PostgresSCMIntegrationsTableName, scmIntegrationID), 
		1, 
		conn, 
		fmt.Sprintf("row was not inserted to '%s' with ID: %s", postgres.PostgresSCMIntegrationsTableName, scmIntegrationID),
	)

	integration.AssertRowCount(
		t, 
		fmt.Sprintf("SELECT * FROM %s WHERE id='%s';", postgres.PostgresSCMIntegrationsAccessTokensTableName, tokenID), 
		1, 
		conn, 
		fmt.Sprintf("row was not inserted to '%s' with ID: %s", postgres.PostgresSCMIntegrationsAccessTokensTableName, tokenID),
	)
}