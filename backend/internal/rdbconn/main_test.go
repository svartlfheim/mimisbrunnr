package rdbconn

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/svartlfheim/mimisbrunnr/test/integration"
	rdbconnmocks "github.com/svartlfheim/mimisbrunnr/test/mocks/rdbconn"
	zerologmocks "github.com/svartlfheim/mimisbrunnr/test/mocks/zerolog"
)

func Test_NewConnectionManager_sets_opts_correctly(t *testing.T) {
	l := zerologmocks.NewLogger()
	cm, err := NewConnectionManager(
		l.Logger,
		&rdbconnmocks.ConnectionOpener{},
		WithDriver("postgres"),
		WithHost("1.2.3.4"),
		WithPort("8743"),
		WithUsername("myuser"),
		WithPassword("mypassword"),
		WithSchema("myschema"),
		WithDatabase("mydb"),
		WithRetries(5),
	)

	assert.Nil(t, err)

	assert.Equal(t, "postgres", cm.driver)
	assert.Equal(t, "1.2.3.4", cm.host)
	assert.Equal(t, "8743", cm.port)
	assert.Equal(t, "myuser", cm.username)
	assert.Equal(t, "mypassword", cm.password)
	assert.Equal(t, "myschema", cm.schema)
	assert.Equal(t, "mydb", cm.database)
	assert.Equal(t, 5, cm.retries)
}

func Test_ConnectionManager_validateState_only_allows_postgres(t *testing.T) {
	cm := &ConnectionManager{}

	drivers := []string{
		"mysql",
		"blah",
		"mssql",
		"sqlite",
		"garbage",
		"pgsql",
	}

	for _, d := range drivers {
		cm.driver = d
		assert.IsType(t, ErrUnsupportedDriver{}, cm.validateState())
	}
}

func Test_ConnectionManager_validateState_required_fields(t *testing.T) {
	cm := &ConnectionManager{
		driver: "postgres",
	}

	err := cm.validateState()
	assert.IsType(t, ErrConfigurationMissingFields{}, err)

	typedErr := err.(ErrConfigurationMissingFields)

	assert.IsType(t, []string{"host", "port", "database", "username", "password", "schema"}, typedErr.Fields)
}

func Test_ConnectionManager_IsPostgres(t *testing.T) {
	cm := &ConnectionManager{
		driver: "postgres",
	}

	assert.True(t, cm.IsPostgres())
}

func Test_ConnectionManager_GetConnection(t *testing.T) {
	getMockDB := func(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
		mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))

		if err != nil {
			t.Error("could not create postgres mock")
			t.FailNow()
		}

		sqlxDB := sqlx.NewDb(mockDB, "postgres")

		return sqlxDB, mock, func() {
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("sqlmock expectations failed: %s", err.Error())
				t.Fail()
			}

			mockDB.Close()
		}
	}
	buildConnectionManager := func(o *rdbconnmocks.ConnectionOpener) *ConnectionManager {
		l := zerologmocks.NewLogger()
		return &ConnectionManager{
			logger:           l.Logger,
			opener:           o,
			activeConnection: nil,

			retries: 3,

			driver:   "postgres",
			username: "myuser",
			password: "mypass",
			host:     "myhost",
			port:     "1234",
			database: "mydb",
			schema:   "myschema",
		}
	}

	t.Run("no connection - fails for non postgres driver", func(tt *testing.T) {
		opener := &rdbconnmocks.ConnectionOpener{}
		_, _, checkAndClose := getMockDB(tt)
		defer checkAndClose()

		defer opener.AssertExpectations(tt)
		defer opener.AssertNumberOfCalls(tt, "ForPostgres", 0)

		cm := buildConnectionManager(opener)
		cm.driver = "garbage"
		conn, err := cm.GetConnection()

		assert.Nil(tt, conn)
		assert.IsType(tt, ErrUnsupportedDriver{}, err)
	})

	t.Run("no connection - works first time", func(tt *testing.T) {
		opener := &rdbconnmocks.ConnectionOpener{}
		sqlxDB, mock, checkAndClose := getMockDB(tt)
		defer checkAndClose()

		mock.ExpectPing()

		opener.EXPECT().ForPostgres("myuser", "mypass", "mydb", "myhost", "1234", "myschema").Return(sqlxDB, nil)
		defer opener.AssertExpectations(tt)
		defer opener.AssertNumberOfCalls(tt, "ForPostgres", 1)

		cm := buildConnectionManager(opener)
		conn, err := cm.GetConnection()

		assert.Same(tt, sqlxDB, conn)
		assert.Nil(tt, err)
	})

	t.Run("healthy connection - returned", func(tt *testing.T) {
		opener := &rdbconnmocks.ConnectionOpener{}
		sqlxDB, mock, checkAndClose := getMockDB(tt)
		defer checkAndClose()

		mock.ExpectPing()

		defer opener.AssertExpectations(tt)
		defer opener.AssertNumberOfCalls(tt, "ForPostgres", 0)
		cm := buildConnectionManager(opener)
		cm.activeConnection = sqlxDB

		conn, err := cm.GetConnection()

		assert.Same(tt, sqlxDB, conn)
		assert.Nil(tt, err)
	})

	t.Run("healthy connection fails - new one returned first time", func(tt *testing.T) {
		opener := &rdbconnmocks.ConnectionOpener{}
		sqlxDB, mock, checkAndClose := getMockDB(tt)
		mock.ExpectPing().WillReturnError(errors.New("first ping failed"))
		defer checkAndClose()

		sqlxDB2, mock2, checkAndClose2 := getMockDB(tt)
		mock2.ExpectPing()
		defer checkAndClose2()

		opener.EXPECT().ForPostgres("myuser", "mypass", "mydb", "myhost", "1234", "myschema").Return(sqlxDB2, nil)
		defer opener.AssertExpectations(tt)
		defer opener.AssertNumberOfCalls(tt, "ForPostgres", 1)
		cm := buildConnectionManager(opener)
		cm.activeConnection = sqlxDB

		conn, err := cm.GetConnection()

		assert.Same(tt, sqlxDB2, conn)
		assert.Nil(tt, err)
	})

	t.Run("healthy connection fails - never becomes healthy", func(tt *testing.T) {
		opener := &rdbconnmocks.ConnectionOpener{}
		sqlxDB, mock, checkAndClose := getMockDB(tt)
		mock.ExpectPing().WillReturnError(errors.New("first ping failed"))
		defer checkAndClose()

		sqlxDB2, mock2, checkAndClose2 := getMockDB(tt)
		mock2.ExpectPing().WillReturnError(errors.New("second ping failed"))
		defer checkAndClose2()

		sqlxDB3, mock3, checkAndClose3 := getMockDB(tt)
		mock3.ExpectPing().WillReturnError(errors.New("third ping failed"))
		defer checkAndClose3()

		sqlxDB4, mock4, checkAndClose4 := getMockDB(tt)
		mock4.ExpectPing().WillReturnError(errors.New("fourth ping failed"))
		defer checkAndClose4()

		opener.EXPECT().ForPostgres("myuser", "mypass", "mydb", "myhost", "1234", "myschema").Return(sqlxDB2, nil).Once()
		opener.EXPECT().ForPostgres("myuser", "mypass", "mydb", "myhost", "1234", "myschema").Return(sqlxDB3, nil).Once()
		opener.EXPECT().ForPostgres("myuser", "mypass", "mydb", "myhost", "1234", "myschema").Return(sqlxDB4, nil).Once()
		defer opener.AssertExpectations(tt)
		defer opener.AssertNumberOfCalls(tt, "ForPostgres", 3)
		cm := buildConnectionManager(opener)
		cm.activeConnection = sqlxDB

		conn, err := cm.GetConnection()

		assert.IsType(tt, ErrConnectionGoneAway{}, err)
		assert.Nil(tt, conn)
	})

	t.Run("healthy connection fails - new one returned second time", func(tt *testing.T) {
		opener := &rdbconnmocks.ConnectionOpener{}
		sqlxDB, mock, checkAndClose := getMockDB(tt)
		mock.ExpectPing().WillReturnError(errors.New("first ping failed"))
		defer checkAndClose()

		sqlxDB2, mock2, checkAndClose2 := getMockDB(tt)
		mock2.ExpectPing()
		defer checkAndClose2()

		opener.EXPECT().ForPostgres("myuser", "mypass", "mydb", "myhost", "1234", "myschema").Return(nil, errors.New("could not open first time")).Once()
		opener.EXPECT().ForPostgres("myuser", "mypass", "mydb", "myhost", "1234", "myschema").Return(sqlxDB2, nil).Once()
		defer opener.AssertExpectations(tt)
		defer opener.AssertNumberOfCalls(tt, "ForPostgres", 2)
		cm := buildConnectionManager(opener)
		cm.activeConnection = sqlxDB

		conn, err := cm.GetConnection()

		assert.Same(tt, sqlxDB2, conn)
		assert.Nil(tt, err)
	})
}

func TestIntegration_connectionOpener_ForPostgres(t *testing.T) {
	integration.SkipIfIntegrationTestsNotConfigured(t)

	u := integration.GetenvOrFail(t, "RDB_USERNAME");
	pass := integration.GetenvOrFail(t, "RDB_PASSWORD");
	db := integration.GetenvOrFail(t, "RDB_DATABASE");
	host := integration.GetenvOrFail(t, "RDB_HOST");
	port := integration.GetenvOrFail(t, "RDB_PORT");
	schema := integration.GetenvOrFail(t, "RDB_SCHEMA"); 

	opener := NewConnectionOpener()
	conn, err := opener.ForPostgres(u, pass, db, host, port, schema)
	require.Nil(t, err)

	defer func() {
		//nolint:errcheck
		conn.Close()
	}()

}
