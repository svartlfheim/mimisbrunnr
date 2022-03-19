package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/svartlfheim/gomigrator"
	"github.com/svartlfheim/mimisbrunnr/internal/schema"
	schemamocks "github.com/svartlfheim/mimisbrunnr/test/mocks/schema"
	zerologmocks "github.com/svartlfheim/mimisbrunnr/test/mocks/zerolog"
)

func GetDatabaseMigrationsConnectionOrFail(t *testing.T) *sqlx.DB {
	driver := GetenvOrFail(t, "RDB_DRIVER");
	u := GetenvOrFail(t, "RDB_MIGRATIONS_USERNAME");
	pass := GetenvOrFail(t, "RDB_MIGRATIONS_PASSWORD");
	db := GetenvOrFail(t, "RDB_DATABASE");
	host := GetenvOrFail(t, "RDB_HOST");
	port := GetenvOrFail(t, "RDB_PORT");
	schema := GetenvOrFail(t, "RDB_SCHEMA"); 

	var err error
	var conn *sqlx.DB

	switch driver {
	case "postgres":
		conn, err = buildPostgresConnection(u, pass, db, host, port, schema)
	default:
		t.Errorf("unsupported db driver '%s'", driver)
		t.FailNow()
	}

	if err != nil {
		t.Errorf("error occurred connecting to database: %s", err.Error())
		t.FailNow()
	}

	return conn
}

func GetDatabaseConnectionOrFail(t *testing.T) *sqlx.DB {
	driver := GetenvOrFail(t, "RDB_DRIVER");
	u := GetenvOrFail(t, "RDB_USERNAME");
	pass := GetenvOrFail(t, "RDB_PASSWORD");
	db := GetenvOrFail(t, "RDB_DATABASE");
	host := GetenvOrFail(t, "RDB_HOST");
	port := GetenvOrFail(t, "RDB_PORT");
	schema := GetenvOrFail(t, "RDB_SCHEMA"); 

	var err error
	var conn *sqlx.DB

	switch driver {
	case "postgres":
		conn, err = buildPostgresConnection(u, pass, db, host, port, schema)
	default:
		t.Errorf("unsupported db driver '%s'", driver)
		t.FailNow()
	}

	if err != nil {
		t.Errorf("error occurred connecting to database: %s", err.Error())
		t.FailNow()
	}

	return conn
}

func buildPostgresConnection(u, pass, db, h, port, schema string) (*sqlx.DB, error) {
	connString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s search_path=%s sslmode=disable",
		u,
		pass,
		db,
		h,
		port,
		schema,
	)

	return sqlx.Connect("postgres", connString)
}

func clearPostgres(conn *sqlx.DB, schema string) error {
	rows, err := conn.Queryx(`
SELECT
	table_name
FROM
	information_schema.tables
WHERE 
	table_schema = $1
`, schema)

	if err != nil {
		return err
	}

	tables := []string{}

	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)

		if err != nil {
			return err
		}

		tables = append(tables, tableName)
	}


	if len(tables) == 0 {
		return nil
	}

	tableList := strings.Join(tables, ",")

	_, err = conn.Exec(
		fmt.Sprintf(`DROP TABLE %s CASCADE`, tableList),
	)
	
	return err
}

func ClearDBOrFail(t *testing.T) {
	conn := GetDatabaseMigrationsConnectionOrFail(t)
	defer conn.Close()
	driver := GetenvOrFail(t, "RDB_DRIVER");

	var err error

	switch driver {
	case "postgres":
		err = clearPostgres(conn, GetenvOrFail(t, "RDB_SCHEMA"))
	default:
		t.Errorf("unsupported db driver for ClearDBOrFail '%s'", driver)
		t.FailNow()
	}

	if err != nil {
		t.Errorf("error clearing database: %s", err.Error())
		t.FailNow()
	}
}

func ResetMigrationsOrFail(t *testing.T) {
	ClearDBOrFail(t)

	conn := GetDatabaseMigrationsConnectionOrFail(t)

	l := zerologmocks.NewLogger()

	cm := &schemamocks.ConnectionManager{}
	cm.EXPECT().GetConnection().Return(conn, nil)

	cfg := &schemamocks.HasSchema{}
	cfg.EXPECT().GetRDBSchema().Return(GetenvOrFail(t, "RDB_SCHEMA"))

	m, err := schema.NewMigrator(cm, cfg, l.Logger)

	require.Nil(t, err)

	err = m.Up(gomigrator.MigrateToLatest)

	require.Nil(t, err)
}


func AssertTableExists(table string, conn *sqlx.DB, t *testing.T) {
	
	if conn.DriverName() != "postgres" {
		t.Errorf("Unsupported DB driver %s", conn.DriverName())
		t.FailNow()
	}

	schema := GetenvOrFail(t, "RDB_SCHEMA")
	checkExistsQuery := `
	SELECT 
		COUNT(1)
	FROM 
		information_schema.tables 
	WHERE 
		table_schema=$1 AND table_name=$2;
	`
	rows, err := conn.Query(checkExistsQuery, schema, table)

	if err != nil {
		t.Errorf("Failed to run query in AssertTableExists: %s", err.Error())
	}

	rows.Next()
	var count int
	err = rows.Scan(&count)

	if err != nil {
		t.Errorf("Failed to scan result from AssertTableExists query: %s", err.Error())
	}

	if count == 0 {
		t.Errorf("table not found")
	}
}


func AssertTableDoesNotExist(table string, conn *sqlx.DB, t *testing.T) {
	
	if conn.DriverName() != "postgres" {
		t.Errorf("Unsupported DB driver %s", conn.DriverName())
		t.FailNow()
	}

	schema := GetenvOrFail(t, "RDB_SCHEMA")
	checkExistsQuery := `
	SELECT 
		COUNT(1)
	FROM 
		information_schema.tables 
	WHERE 
		table_schema=$1 AND table_name=$2;
	`
	rows, err := conn.Query(checkExistsQuery, schema, table)

	if err != nil {
		t.Errorf("Failed to run query in AssertTableDoesNotExist: %s", err.Error())
	}

	rows.Next()
	var count int
	err = rows.Scan(&count)

	if err != nil {
		t.Errorf("Failed to scan result from AssertTableDoesNotExist query: %s", err.Error())
	}

	if count > 0 {
		t.Errorf("table exists")
	}
}

func AssertRowCount(t *testing.T, q string, expect int, conn *sqlx.DB, message string) {
	rows, err := conn.Query(q)

	if err != nil {
		t.Errorf("Failed to run query in AssertRowCount: %s", err.Error())
	}

	count := 0

	for rows.Next() {
		count++
	}

	if err != nil {
		t.Errorf("Failed to scan result from AssertRowCount query: %s", err.Error())
	}

	assert.Equal(t, expect, count, message)
}