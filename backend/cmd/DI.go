package cmd

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/svartlfheim/gomigrator"
	"github.com/svartlfheim/mimisbrunnr/internal/config"
	"github.com/svartlfheim/mimisbrunnr/internal/rdbconn"
	"github.com/svartlfheim/mimisbrunnr/internal/schema"
	"github.com/svartlfheim/mimisbrunnr/internal/scm"
	"github.com/svartlfheim/mimisbrunnr/internal/server"
	"github.com/svartlfheim/mimisbrunnr/internal/validation"
)

type DIContainer struct {
	Cfg *config.AppConfig

	// rdbconn.*
	RdbConnManager              *rdbconn.ConnectionManager
	RdbConnManagerForMigrations *rdbconn.ConnectionManager
	RdbConnOpener               *rdbconn.ConnectionOpener

	// scm.*
	PostgresSCMIntegrationsRepository *scm.PostgresRepository
	SCMIntegrationsManager            *scm.Manager

	// schema.*
	Migrator *gomigrator.Migrator

	// server.*
	Server                        *server.Server
	SCMIntegrationsController     *server.SCMIntegrationsController
	ProjectsController            *server.ProjectsController
	ErrorHandlingJsonUnmarshaller *server.ErrorHandlingJsonUnmarshaller

	// validation.*
	Validator *validation.Validator

	// External dependencies
	Logger zerolog.Logger
	Fs     afero.Fs
}

type commandHandler func(args []string) error

func (di *DIContainer) loadConfig(h commandHandler) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cfgPath, err := cmd.Flags().GetString("config")

		if err != nil {
			return errors.New("config path opt not found")
		}

		cfg, err := config.Load(cfgPath, di.Fs, "mimisbrunnr")

		if err != nil {
			return err
		}

		if logLevel, err := cmd.Flags().GetInt("log-level"); err != nil {
			di.Logger = di.Logger.Level(zerolog.Level(logLevel))
		}

		di.Cfg = cfg
		return h(args)
	}
}

func (di *DIContainer) GetCommands() []*cobra.Command {
	migrate := &cobra.Command{
		Use:   "migrate",
		Short: "Run migrations for app",
	}
	migrate.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "Create/update schema",
		RunE: di.loadConfig(func(args []string) error {
			return handleMigrationsUp(di.GetMigrator(), di.Cfg, args)
		}),
	})
	migrate.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "Rollback db schema",
		RunE: di.loadConfig(func(args []string) error {
			return handleMigrationsDown(di.GetMigrator(), di.Cfg, args)
		}),
	})
	migrate.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List the migrations",
		RunE: di.loadConfig(func(args []string) error {
			return handleMigrationsList(di.GetMigrator(), di.Cfg, args)
		}),
	})

	return []*cobra.Command{
		{
			Use:   "serve",
			Short: "Start the HTTP server",
			RunE: di.loadConfig(func(args []string) error {
				return handleServe(di.GetServer(), di.Cfg, args)
			}),
		},
		migrate,
	}
}

func (di *DIContainer) GetMigrator() *gomigrator.Migrator {
	if di.Migrator == nil {
		m, err := schema.NewMigrator(
			di.GetRDBConnManagerForMigrations(),
			di.Cfg,
			di.Logger,
		)

		if err != nil {
			di.Logger.Panic().Err(err).Msg(fmt.Sprintf("failed to build migrator: %s", err.Error()))
		}

		di.Migrator = m
	}

	return di.Migrator
}

func (di *DIContainer) GetRDBConnOpener() *rdbconn.ConnectionOpener {
	if di.RdbConnOpener == nil {
		o := rdbconn.NewConnectionOpener()

		di.RdbConnOpener = o
	}

	return di.RdbConnOpener
}

func (di *DIContainer) GetRDBConnManager() *rdbconn.ConnectionManager {
	if di.RdbConnManager == nil {
		built, err := rdbconn.NewConnectionManager(
			di.Logger,
			di.GetRDBConnOpener(),
			rdbconn.WithDriver(di.Cfg.GetRDBDriver()),
			rdbconn.WithHost(di.Cfg.GetRDBHost()),
			rdbconn.WithPort(di.Cfg.GetRDBPort()),
			rdbconn.WithUsername(di.Cfg.GetRDBUsername()),
			rdbconn.WithPassword(di.Cfg.GetRDBPassword()),
			rdbconn.WithSchema(di.Cfg.GetRDBSchema()),
			rdbconn.WithDatabase(di.Cfg.GetRDBDatabase()),
		)

		if err != nil {
			di.Logger.Panic().Err(err).Msg(fmt.Sprintf("failed to build RDB connection manager: %s", err.Error()))
		}

		di.RdbConnManager = built
	}

	return di.RdbConnManager
}

func (di *DIContainer) GetRDBConnManagerForMigrations() *rdbconn.ConnectionManager {
	if di.RdbConnManagerForMigrations == nil {
		built, err := rdbconn.NewConnectionManager(
			di.Logger,
			di.GetRDBConnOpener(),
			rdbconn.WithDriver(di.Cfg.GetRDBDriver()),
			rdbconn.WithHost(di.Cfg.GetRDBHost()),
			rdbconn.WithPort(di.Cfg.GetRDBPort()),
			rdbconn.WithUsername(di.Cfg.GetRDBMigrationsUsername()),
			rdbconn.WithPassword(di.Cfg.GetRDBMigrationsPassword()),
			rdbconn.WithSchema(di.Cfg.GetRDBSchema()),
			rdbconn.WithDatabase(di.Cfg.GetRDBDatabase()),
		)

		if err != nil {
			di.Logger.Panic().Err(err).Msg(fmt.Sprintf("failed to build RDB connection manager for migrations: %s", err.Error()))
		}

		di.RdbConnManagerForMigrations = built
	}

	return di.RdbConnManagerForMigrations
}

func (di *DIContainer) GetPostgresSCMIntegrationsRepository() *scm.PostgresRepository {
	if di.PostgresSCMIntegrationsRepository == nil {
		connManager := di.GetRDBConnManager()

		di.PostgresSCMIntegrationsRepository = scm.NewPostgresRepository(connManager, di.Logger)
	}

	return di.PostgresSCMIntegrationsRepository
}

func (di *DIContainer) GetSCMIntegrationsManager() *scm.Manager {
	if di.SCMIntegrationsManager == nil {
		di.SCMIntegrationsManager = scm.NewManager(
			di.Logger,
			di.GetPostgresSCMIntegrationsRepository(),
			di.GetValidator(),
		)
	}

	return di.SCMIntegrationsManager
}

func (di *DIContainer) GetErrorHandlingJsonUnmarshaller() *server.ErrorHandlingJsonUnmarshaller {
	if di.ErrorHandlingJsonUnmarshaller == nil {
		di.ErrorHandlingJsonUnmarshaller = server.NewErrorHandlingJsonUnmarshaller()
	}

	return di.ErrorHandlingJsonUnmarshaller
}

func (di *DIContainer) GetSCMIntegrationsController() *server.SCMIntegrationsController {
	if di.SCMIntegrationsController == nil {
		di.SCMIntegrationsController = server.NewSCMIntegrationsController(
			di.Logger,
			di.GetSCMIntegrationsManager(),
			di.GetErrorHandlingJsonUnmarshaller(),
		)
	}

	return di.SCMIntegrationsController
}

func (di *DIContainer) GetProjectsController() *server.ProjectsController {
	if di.ProjectsController == nil {
		di.ProjectsController = server.NewProjectsController()
	}

	return di.ProjectsController
}

func (di *DIContainer) GetServer() *server.Server {
	if di.Server == nil {
		s := server.NewServer(
			di.Logger,
			[]server.Controller{
				di.GetSCMIntegrationsController(),
				di.GetProjectsController(),
			},
		)

		di.Server = s
	}

	return di.Server
}

func (di *DIContainer) GetValidator() *validation.Validator {
	if di.Validator == nil {
		di.Validator = validation.NewValidator(di.Logger)
	}

	return di.Validator
}

func NewDIContainer(l zerolog.Logger, fs afero.Fs) *DIContainer {
	return &DIContainer{
		Logger: l,
		Fs:     fs,
	}
}
