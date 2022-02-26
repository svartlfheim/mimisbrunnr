package config

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	testenv "github.com/svartlfheim/mimisbrunnr/test/env"
)

func clearAllConfigEnvVars(t *testing.T) {
	err := testenv.Clear(
		"MIMISBRUNNR_HTTP_PORT",
		"MIMISBRUNNR_HTTP_LISTEN_HOST",
		"MIMISBRUNNR_RDB_DRIVER",
		"MIMISBRUNNR_RDB_HOST",
		"MIMISBRUNNR_RDB_PORT",
		"MIMISBRUNNR_RDB_SCHEMA",
		"MIMISBRUNNR_RDB_DATABASE",
		"MIMISBRUNNR_RDB_USERNAME",
		"MIMISBRUNNR_RDB_PASSWORD",
		"MIMISBRUNNR_RDB_MIGRATIONS_USERNAME",
		"MIMISBRUNNR_RDB_MIGRATIONS_PASSWORD",
	)

	if err != nil {
		t.Errorf("failed to clear env vars for config: %s", err.Error())
		t.FailNow()
	}
}

func Test_ensureFileExists(t *testing.T) {
	fs := afero.NewMemMapFs()

	if _, err := fs.Create("/some/path/that/exists"); err != nil {
		t.Errorf("failed to create file '/some/path/that/exists' in mock fs")
		t.FailNow()
	}

	assert.Nil(t, ensureFileExists("/some/path/that/exists", fs))
	assert.IsType(t, ErrConfigFileDoesNotExist{}, ensureFileExists("/does/not/exist", fs))
}

func Test_unmarshalConfigToStruct_valid_contents(t *testing.T) {
	configPath := "/some/config.yaml"
	configContents := `
http:
  port: 1010
  listen_host: 1.2.3.4

rdb:
  driver: somedriver
  host: somehost
  port: someport
  username: someusername
  password: somepassword
  schema: someschema
  database: somedatabase

  migrations:
    username: somemigrationusername
    password: somemigrationpassword
`
	fs := afero.NewMemMapFs()

	f, err := fs.Create(configPath)
	if err != nil {
		t.Errorf("failed to create file '%s' in mock fs", configPath)
		t.FailNow()
	}

	if _, err := f.Write([]byte(configContents)); err != nil {
		t.Errorf("failed to write to config file '%s' in mock fs", configPath)
		t.FailNow()
	}

	cfg := &AppConfig{}
	err = unmarshalConfigToStruct(configPath, fs, cfg)

	assert.Nil(t, err)

	assert.Equal(t, "1010", cfg.HTTP.Port)
	assert.Equal(t, "1.2.3.4", cfg.HTTP.ListenHost)
	assert.Equal(t, "somedriver", cfg.RDB.Driver)
	assert.Equal(t, "somehost", cfg.RDB.Host)
	assert.Equal(t, "someport", cfg.RDB.Port)
	assert.Equal(t, "someusername", cfg.RDB.Username)
	assert.Equal(t, "somepassword", cfg.RDB.Password)
	assert.Equal(t, "someschema", cfg.RDB.Schema)
	assert.Equal(t, "somedatabase", cfg.RDB.Database)
	assert.Equal(t, "somemigrationusername", cfg.RDB.Migrations.Username)
	assert.Equal(t, "somemigrationpassword", cfg.RDB.Migrations.Password)

	assert.Equal(t, "1010", cfg.GetHTTPPort())
	assert.Equal(t, "1.2.3.4", cfg.GetListenHost())
	assert.Equal(t, "somedriver", cfg.GetRDBDriver())
	assert.Equal(t, "somehost", cfg.GetRDBHost())
	assert.Equal(t, "someport", cfg.GetRDBPort())
	assert.Equal(t, "someusername", cfg.GetRDBUsername())
	assert.Equal(t, "somepassword", cfg.GetRDBPassword())
	assert.Equal(t, "someschema", cfg.GetRDBSchema())
	assert.Equal(t, "somedatabase", cfg.GetRDBDatabase())
	assert.Equal(t, "somemigrationusername", cfg.GetRDBMigrationsUsername())
	assert.Equal(t, "somemigrationpassword", cfg.GetRDBMigrationsPassword())
}

func Test_unmarshalConfigToStruct_no_file(t *testing.T) {
	configPath := "/some/config.yaml"
	fs := afero.NewMemMapFs()

	cfg := &AppConfig{}
	err := unmarshalConfigToStruct(configPath, fs, cfg)

	assert.IsType(t, ErrFsUnusable{}, err)
}

func Test_unmarshalConfigToStruct_invalid_contents(t *testing.T) {
	configPath := "/some/config.yaml"
	configContents := `
{this}
	ain't no kinda
yaml I
		Ever did see
`
	fs := afero.NewMemMapFs()

	f, err := fs.Create(configPath)
	if err != nil {
		t.Errorf("failed to create file '%s' in mock fs", configPath)
		t.FailNow()
	}

	if _, err := f.Write([]byte(configContents)); err != nil {
		t.Errorf("failed to write to config file '%s' in mock fs", configPath)
		t.FailNow()
	}

	cfg := &AppConfig{}
	err = unmarshalConfigToStruct(configPath, fs, cfg)

	assert.IsType(t, ErrCannotUnmarshalConfig{}, err)
}

func Test_processEnvVars(t *testing.T) {
	cfg := &AppConfig{
		HTTP: HTTPConfig{
			Port:       "8080",
			ListenHost: "0.0.0.0",
		},
		RDB: RDBConfig{
			Driver:   "postgres",
			Host:     "localhost",
			Port:     "5432",
			Username: "dummy",
			Password: "dummy",
			Schema:   "dummy",
			Database: "dummy",
			Migrations: RDBMigrationsConfig{
				Username: "dummymiguser",
				Password: "dummymigpassword",
			},
		},
	}

	clearAllConfigEnvVars(t)

	resetHttpPort := testenv.Override("MIMISBRUNNR_HTTP_PORT", "9898")
	resetHttpListenHost := testenv.Override("MIMISBRUNNR_HTTP_LISTEN_HOST", "9.8.7.6")
	resetRDBDriver := testenv.Override("MIMISBRUNNR_RDB_DRIVER", "mysql")
	resetRDBHost := testenv.Override("MIMISBRUNNR_RDB_HOST", "otherhost")
	resetRDBPort := testenv.Override("MIMISBRUNNR_RDB_PORT", "7896")
	resetRDBSchema := testenv.Override("MIMISBRUNNR_RDB_SCHEMA", "myschema")
	resetRDBDatabase := testenv.Override("MIMISBRUNNR_RDB_DATABASE", "mydb")
	resetRDBUsername := testenv.Override("MIMISBRUNNR_RDB_USERNAME", "myuser")
	resetRDBPassword := testenv.Override("MIMISBRUNNR_RDB_PASSWORD", "mypass")
	resetRDBMigrationsUsername := testenv.Override("MIMISBRUNNR_RDB_MIGRATIONS_USERNAME", "mymigrationuser")
	resetRDBMigrationsPassword := testenv.Override("MIMISBRUNNR_RDB_MIGRATIONS_PASSWORD", "mymigrationpass")

	defer func() {
		//nolint:errcheck
		resetHttpPort()
		//nolint:errcheck
		resetHttpListenHost()
		//nolint:errcheck
		resetRDBDriver()
		//nolint:errcheck
		resetRDBHost()
		//nolint:errcheck
		resetRDBPort()
		//nolint:errcheck
		resetRDBSchema()
		//nolint:errcheck
		resetRDBDatabase()
		//nolint:errcheck
		resetRDBUsername()
		//nolint:errcheck
		resetRDBPassword()
		//nolint:errcheck
		resetRDBMigrationsUsername()
		//nolint:errcheck
		resetRDBMigrationsPassword()
	}()

	err := processEnvVars("mimisbrunnr", cfg)

	assert.Nil(t, err)

	assert.Equal(t, "9898", cfg.GetHTTPPort())
	assert.Equal(t, "9.8.7.6", cfg.GetListenHost())
	assert.Equal(t, "mysql", cfg.GetRDBDriver())
	assert.Equal(t, "otherhost", cfg.GetRDBHost())
	assert.Equal(t, "7896", cfg.GetRDBPort())
	assert.Equal(t, "myuser", cfg.GetRDBUsername())
	assert.Equal(t, "mypass", cfg.GetRDBPassword())
	assert.Equal(t, "myschema", cfg.GetRDBSchema())
	assert.Equal(t, "mydb", cfg.GetRDBDatabase())
	assert.Equal(t, "mymigrationuser", cfg.GetRDBMigrationsUsername())
	assert.Equal(t, "mymigrationpass", cfg.GetRDBMigrationsPassword())
}

func Test_Load(t *testing.T) {
	configPath := "/some/config.yaml"
	configContents := `
http:
  port: 1010
  listen_host: 1.2.3.4

rdb:
  driver: somedriver
  host: somehost
  port: someport
  username: someusername
  password: somepassword
  schema: someschema
  database: somedatabase

  migrations:
    username: somemigrationusername
    password: somemigrationpassword
`
	fs := afero.NewMemMapFs()

	f, err := fs.Create(configPath)
	if err != nil {
		t.Errorf("failed to create file '%s' in mock fs", configPath)
		t.FailNow()
	}

	if _, err := f.Write([]byte(configContents)); err != nil {
		t.Errorf("failed to write to config file '%s' in mock fs", configPath)
		t.FailNow()
	}

	clearAllConfigEnvVars(t)

	resetRDBHost := testenv.Override("MIMISBRUNNR_RDB_HOST", "otherhost")
	resetRDBPort := testenv.Override("MIMISBRUNNR_RDB_PORT", "7896")

	defer func() {
		//nolint:errcheck
		resetRDBHost()
		//nolint:errcheck
		resetRDBPort()
	}()

	cfg, err := Load(configPath, fs, "mimisbrunnr")

	assert.Nil(t, err)

	assert.Equal(t, "1010", cfg.GetHTTPPort())
	assert.Equal(t, "1.2.3.4", cfg.GetListenHost())
	assert.Equal(t, "somedriver", cfg.GetRDBDriver())
	assert.Equal(t, "otherhost", cfg.GetRDBHost())
	assert.Equal(t, "7896", cfg.GetRDBPort())
	assert.Equal(t, "someusername", cfg.GetRDBUsername())
	assert.Equal(t, "somepassword", cfg.GetRDBPassword())
	assert.Equal(t, "someschema", cfg.GetRDBSchema())
	assert.Equal(t, "somedatabase", cfg.GetRDBDatabase())
	assert.Equal(t, "somemigrationusername", cfg.GetRDBMigrationsUsername())
	assert.Equal(t, "somemigrationpassword", cfg.GetRDBMigrationsPassword())
}

func Test_RDB_MigrationsUsername_defaults_to_rdb_username(t *testing.T) {
	cfg := &AppConfig{
		HTTP: HTTPConfig{
			Port:       "8080",
			ListenHost: "0.0.0.0",
		},
		RDB: RDBConfig{
			Driver:     "postgres",
			Host:       "localhost",
			Port:       "5432",
			Username:   "dummyuser",
			Password:   "dummypass",
			Schema:     "dummy",
			Database:   "dummy",
			Migrations: RDBMigrationsConfig{},
		},
	}

	assert.Equal(t, "dummyuser", cfg.GetRDBMigrationsUsername())
	assert.Equal(t, "dummypass", cfg.GetRDBMigrationsPassword())
}
