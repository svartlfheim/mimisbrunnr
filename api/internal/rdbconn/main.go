package rdbconn

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

type ConnectionManager struct {
	logger zerolog.Logger
	activeConnection *sqlx.DB

	retries int

	driver string
	username string
	password string
	host string
	port string
	database string
	schema string
}

func (cm *ConnectionManager) openNewConnection() (*sqlx.DB, error) {
	if cm.driver != "postgres" {
		return nil, ErrUnsupportedDriver{
			Driver: cm.driver,
		}
	}

	connString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s search_path=%s sslmode=disable",
		cm.username,
		cm.password,
		cm.database,
		cm.host,
		cm.port,
		cm.schema,
	)

	return sqlx.Connect("postgres", connString)
}

func (cm *ConnectionManager) IsPostgres() bool {
	return cm.driver == "postgres"
}

func (cm *ConnectionManager) GetConnection() (*sqlx.DB, error) {
	err := cm.activeConnection.Ping();

	if err == nil {
		return cm.activeConnection, nil
	}

	for i := 0; i < cm.retries; i++ {
		attemptNumber := i+1
		cm.logger.Warn().Err(err).Int("attempt", attemptNumber).Msg("database ping failed, reconnecting")

		db, err := cm.openNewConnection();

		if err != nil  {
			cm.logger.Warn().Err(err).Int("attempt", attemptNumber).Msg("connection attempt failed")
			continue
		}

		err = db.Ping()

		if err != nil {
			cm.logger.Warn().Err(err).Int("attempt", attemptNumber).Msg("ping attempt failed")
			continue
		}

		cm.activeConnection = db
		return cm.activeConnection, nil
	}

	return nil, ErrConnectionGoneAway{
		Retries: cm.retries,
	}
}

func (cm *ConnectionManager) validateState() error {
	if cm.driver != "postgres" {
		return ErrUnsupportedDriver{
			Driver: cm.driver,
		}
	}

	unsetFields := []string{}

	if cm.host == "" {
		unsetFields = append(unsetFields, "host")
	}

	if cm.port == "" {
		unsetFields = append(unsetFields, "port")
	}

	if cm.database == "" {
		unsetFields = append(unsetFields, "database")
	}

	if cm.username == "" {
		unsetFields = append(unsetFields, "username")
	}

	if cm.password == "" {
		unsetFields = append(unsetFields, "password")
	}

	if cm.schema == "" {
		unsetFields = append(unsetFields, "schema")
	}

	if len(unsetFields) > 0 {
		return ErrConfigurationMissingFields{
			Fields: unsetFields,
		}
	}

	return nil
}

type WithParam func(*ConnectionManager)

func WithDriver(driver string) WithParam {
	return func(cm *ConnectionManager) {
		cm.driver = driver
	}
}

func WithUsername(username string) WithParam {
	return func(cm *ConnectionManager) {
		cm.username = username
	}
}

func WithPassword(password string) WithParam {
	return func(cm *ConnectionManager) {
		cm.password = password
	}
}

func WithHost(host string) WithParam {
	return func(cm *ConnectionManager) {
		cm.host = host
	}
}

func WithPort(port string) WithParam {
	return func(cm *ConnectionManager) {
		cm.port = port
	}
}

func WithDatabase(database string) WithParam {
	return func(cm *ConnectionManager) {
		cm.database = database
	}
}

func WithSchema(schema string) WithParam {
	return func(cm *ConnectionManager) {
		cm.schema = schema
	}
}

func WithRetries(retries int) WithParam {
	return func(cm *ConnectionManager) {
		cm.retries = retries
	}
}

func NewConnectionManager(l zerolog.Logger, opts ...WithParam) (*ConnectionManager, error) {
	cm := &ConnectionManager{
		logger: l,
		retries: 3,
	}

	for _, opt := range(opts) {
		opt(cm)
	}

	return cm, cm.validateState()
}