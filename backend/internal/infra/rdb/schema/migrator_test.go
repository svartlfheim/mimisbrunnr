package schema_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/svartlfheim/gomigrator"
	"github.com/svartlfheim/mimisbrunnr/internal/infra/rdb/schema"
	"github.com/svartlfheim/mimisbrunnr/test/integration"
	schemamocks "github.com/svartlfheim/mimisbrunnr/test/mocks/infra/rdb/schema"
	zerologmocks "github.com/svartlfheim/mimisbrunnr/test/mocks/zerolog"
)

func TestIntegration_AllMigrationsAreAppliedSuccessfully(t *testing.T) {
	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ClearDBOrFail(t)

	conn := integration.GetDatabaseMigrationsConnectionOrFail(t)
	defer conn.Close()

	l := zerologmocks.NewLogger()
	cm := &schemamocks.ConnectionManager{}
	cm.EXPECT().GetConnection().Return(conn, nil)
	cfg := &schemamocks.HasSchema{}
	cfg.EXPECT().GetRDBSchema().Return(integration.GetenvOrFail(t, "RDB_SCHEMA"))

	m, err := schema.NewMigrator(cm, cfg, l.Logger)

	require.Nil(t, err)

	err = m.Up(gomigrator.MigrateToLatest)

	require.Nil(t, err)

	integration.AssertTableExists("scm_integrations", conn, t)
}

func TestIntegration_AllMigrationsAreRolledBackSuccessfully(t *testing.T) {
	integration.SkipIfIntegrationTestsNotConfigured(t)
	integration.ClearDBOrFail(t)

	conn := integration.GetDatabaseMigrationsConnectionOrFail(t)
	defer conn.Close()

	l := zerologmocks.NewLogger()
	cm := &schemamocks.ConnectionManager{}
	cm.EXPECT().GetConnection().Return(conn, nil)
	cfg := &schemamocks.HasSchema{}
	cfg.EXPECT().GetRDBSchema().Return(integration.GetenvOrFail(t, "RDB_SCHEMA"))

	m, err := schema.NewMigrator(cm, cfg, l.Logger)

	require.Nil(t, err)

	err = m.Up(gomigrator.MigrateToLatest)

	require.Nil(t, err)

	integration.AssertTableExists("scm_integrations", conn, t)

	err = m.Down(gomigrator.MigrateToNothing)

	require.Nil(t, err)

	integration.AssertTableDoesNotExist("scm_integrations", conn, t)
}
