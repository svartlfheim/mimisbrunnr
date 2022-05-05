package rdb

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type connectionOpener interface {
	ForPostgres(u string, pass string, db string, h string, port string, s string) (*sqlx.DB, error)
}

type ConnectionOpener struct{}

func (c *ConnectionOpener) ForPostgres(u string, pass string, db string, h string, port string, s string) (*sqlx.DB, error) {
	connString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s search_path=%s sslmode=disable",
		u,
		pass,
		db,
		h,
		port,
		s,
	)

	return sqlx.Connect("postgres", connString)
}

type ConnectionManager struct {
	logger           zerolog.Logger
	opener           connectionOpener
	activeConnection *sqlx.DB

	retries int

	driver   string
	username string
	password string
	host     string
	port     string
	database string
	schema   string
}

func (cm *ConnectionManager) openNewConnection() (*sqlx.DB, error) {
	if cm.driver != "postgres" {
		return nil, ErrUnsupportedDriver{
			Driver: cm.driver,
		}
	}

	return cm.opener.ForPostgres(
		cm.username,
		cm.password,
		cm.database,
		cm.host,
		cm.port,
		cm.schema,
	)
}

func (cm *ConnectionManager) IsPostgres() bool {
	return cm.driver == "postgres"
}

func (cm *ConnectionManager) GetConnection() (*sqlx.DB, error) {
	var err error

	if cm.activeConnection != nil {
		err = cm.activeConnection.Ping()

		if err == nil {
			return cm.activeConnection, nil
		}
	}

	for i := 0; i < cm.retries; i++ {
		attemptNumber := i + 1
		if cm.activeConnection != nil {
			cm.logger.Warn().Err(err).Int("attempt", attemptNumber).Msg("database ping failed, reconnecting")
		} else {
			cm.logger.Debug().Int("attempt", attemptNumber).Msg("opening database connection for first time")
		}

		db, err := cm.openNewConnection()

		// fmt.Printf("\n\nopen err: %#v\n\n", err)
		if err != nil {
			if _, ok := err.(ErrUnsupportedDriver); ok {
				return nil, err
			}

			cm.logger.Warn().Err(err).Int("attempt", attemptNumber).Msg("connection attempt failed")
			continue
		}

		err = db.Ping()

		// fmt.Printf("\n\n%#v\n\n", err)
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

func NewConnectionManager(l zerolog.Logger, o connectionOpener, opts ...WithParam) (*ConnectionManager, error) {
	cm := &ConnectionManager{
		logger:  l,
		opener:  o,
		retries: 3,
	}

	for _, opt := range opts {
		opt(cm)
	}

	return cm, cm.validateState()
}

func NewConnectionOpener() *ConnectionOpener {
	return &ConnectionOpener{}
}
