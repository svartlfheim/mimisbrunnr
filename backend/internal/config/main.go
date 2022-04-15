package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type App interface {
	GetHTTPPort() string
	GetListenHost() string
	HTTPAPIEnabled() bool
	HTTPStaticServerEnabled() bool
	GetHTTPStaticContentPath() string
	HTTPFrontendEnabled() bool
	
	HTTPOpenAPIEnabled()
	GetOpenAPIPath() string

	GetRDBDriver() string
	GetRDBHost() string
	GetRDBPort() string
	GetRDBUsername() string
	GetRDBPassword() string
	GetRDBSchema() string
	GetRDBDatabase() string

	GetRDBMigrationsUsername() string
	GetRDBMigrationsPassword() string
}

type HTTPFrontendConfig struct {
	Enabled bool `yaml:"enabled"`
}

type HTTPAPIConfig struct {
	Enabled bool `yaml:"enabled"`
}

type HTTPStaticConfig struct {
	Path    string `yaml:"path"`
	Enabled bool   `yaml:"enabled"`
}

type HTTPEndpointsConfig struct {
	OpenAPIEnabled bool `yaml:"openapi_enabled"`
}

type HTTPConfig struct {
	Port       string             `yaml:"port"`
	ListenHost string             `yaml:"listen_host" split_words:"true"`
	API        HTTPAPIConfig      `yaml:"api"`
	Static     HTTPStaticConfig   `yaml:"static"`
	Frontend   HTTPFrontendConfig `yaml:"frontend"`
	Endpoints  HTTPEndpointsConfig `yaml:"endpoints"`
}

type RDBMigrationsConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type RDBConfig struct {
	Driver     string              `yaml:"driver"`
	Host       string              `yaml:"host"`
	Port       string              `yaml:"port"`
	Username   string              `yaml:"username"`
	Password   string              `yaml:"password"`
	Schema     string              `yaml:"schema"`
	Database   string              `yaml:"database"`
	Migrations RDBMigrationsConfig `yaml:"migrations"`
}

type OpenAPIConfig struct {
	Path string `yaml:"path"`
}

type AppConfig struct {
	HTTP    HTTPConfig    `yaml:"http"`
	OpenAPI OpenAPIConfig `yaml:"openapi"`
	RDB     RDBConfig     `yaml:"rdb"`
}

func (c *AppConfig) GetHTTPPort() string {
	return c.HTTP.Port
}

func (c *AppConfig) GetListenHost() string {
	return c.HTTP.ListenHost
}

func (c *AppConfig) HTTPAPIEnabled() bool {
	return c.HTTP.API.Enabled
}

func (c *AppConfig) HTTPStaticServerEnabled() bool {
	return c.HTTP.Static.Enabled
}

func (c *AppConfig) GetHTTPStaticContentPath() string {
	return c.HTTP.Static.Path
}

func (c *AppConfig) HTTPFrontendEnabled() bool {
	return c.HTTP.Frontend.Enabled
}

func (c *AppConfig) HTTPOpenAPIEnabled() bool {
	return c.HTTP.Endpoints.OpenAPIEnabled
}

func (c *AppConfig) GetOpenAPIPath() string {
	if c.OpenAPI.Path == "" {
		return "/static/openapi.json"
	}

	return c.OpenAPI.Path
}

func (c *AppConfig) GetRDBDriver() string {
	return c.RDB.Driver
}

func (c *AppConfig) GetRDBHost() string {
	return c.RDB.Host
}

func (c *AppConfig) GetRDBPort() string {
	return c.RDB.Port
}

func (c *AppConfig) GetRDBUsername() string {
	return c.RDB.Username
}

func (c *AppConfig) GetRDBPassword() string {
	return c.RDB.Password
}

func (c *AppConfig) GetRDBSchema() string {
	return c.RDB.Schema
}

func (c *AppConfig) GetRDBDatabase() string {
	return c.RDB.Database
}

func (c *AppConfig) GetRDBMigrationsUsername() string {
	if c.RDB.Migrations.Username == "" {
		return c.RDB.Username
	}

	return c.RDB.Migrations.Username
}

func (c *AppConfig) GetRDBMigrationsPassword() string {
	if c.RDB.Migrations.Password == "" {
		return c.RDB.Password
	}

	return c.RDB.Migrations.Password
}

func ensureFileExists(path string, fs afero.Fs) error {

	doesExist, err := afero.Exists(fs, path)

	if !doesExist {
		return ErrConfigFileDoesNotExist{
			Path: path,
		}
	}

	if err != nil {
		return ErrFsUnusable{
			Message: err.Error(),
		}
	}

	return nil
}

func unmarshalConfigToStruct(path string, fs afero.Fs, cfg *AppConfig) error {

	fileBytes, err := afero.ReadFile(fs, path)

	if err != nil {
		return ErrFsUnusable{
			Message: err.Error(),
		}
	}

	if err := yaml.Unmarshal(fileBytes, cfg); err != nil {
		return ErrCannotUnmarshalConfig{
			Message: err.Error(),
		}
	}

	return nil
}

func processEnvVars(appName string, cfg *AppConfig) error {
	if err := envconfig.Process(appName, cfg); err != nil {
		return ErrCouldNotProcessEnv{
			Message: err.Error(),
		}
	}

	return nil
}

func Load(path string, fs afero.Fs, appName string) (*AppConfig, error) {
	cfg := &AppConfig{
		HTTP: HTTPConfig{
			Port:       "8080",
			ListenHost: "0.0.0.0",
		},
		RDB: RDBConfig{
			Driver: "postgres",
			Host:   "localhost",
			Port:   "5432",
		},
	}

	if err := ensureFileExists(path, fs); err != nil {
		return cfg, err
	}

	if err := unmarshalConfigToStruct(path, fs, cfg); err != nil {
		return cfg, err
	}

	if err := processEnvVars(appName, cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
