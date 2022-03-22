package cmd

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/svartlfheim/gomigrator"
	"github.com/svartlfheim/mimisbrunnr/internal/app/api"
	"github.com/svartlfheim/mimisbrunnr/internal/app/scm"
	"github.com/svartlfheim/mimisbrunnr/internal/config"
	"github.com/svartlfheim/mimisbrunnr/internal/infra/rdb"
	"github.com/svartlfheim/mimisbrunnr/internal/infra/rdb/schema"
	scmpostgres "github.com/svartlfheim/mimisbrunnr/internal/infra/rdb/scm/postgres"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/validation"
)

type DIContainer struct {
	Cfg *config.AppConfig

	// rdb.*
	RdbConnManager              *rdb.ConnectionManager
	RdbConnManagerForMigrations *rdb.ConnectionManager
	RdbConnOpener               *rdb.ConnectionOpener

	// scm.*
	PostgresSCMIntegrationsRepository *scmpostgres.Repository
	SCMIntegrationsManager            *scm.Controller
	SCMIntegrationsTransformer        *scm.Transformer

	// schema.*
	Migrator *gomigrator.Migrator

	// api.*
	Server                        *api.Server
	SCMIntegrationsController     *api.SCMHandler
	ProjectsController            *api.ProjectsHandler
	ErrorHandlingJsonUnmarshaller *api.ErrorHandlingJsonUnmarshaller

	// validation.*
	Validator *validation.Validator

	// External dependencies
	Logger zerolog.Logger
	Fs     afero.Fs
}

type commandHandler func(args []string) error

func (di *DIContainer) commandWrap(h commandHandler) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cfgPath, err := cmd.Flags().GetString("config")

		if err != nil {
			return errors.New("config path opt not found")
		}

		cfg, err := config.Load(cfgPath, di.Fs, "mimisbrunnr")

		if err != nil {
			return err
		}

		if logLevel, err := cmd.Flags().GetInt("log-level"); err == nil {
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
		RunE: di.commandWrap(func(args []string) error {
			return handleMigrationsUp(di.GetMigrator(), di.Cfg, args)
		}),
	})
	migrate.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "Rollback db schema",
		RunE: di.commandWrap(func(args []string) error {
			return handleMigrationsDown(di.GetMigrator(), di.Cfg, args)
		}),
	})
	migrate.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List the migrations",
		RunE: di.commandWrap(func(args []string) error {
			return handleMigrationsList(di.GetMigrator(), di.Cfg, args)
		}),
	})

	return []*cobra.Command{
		{
			Use:   "serve",
			Short: "Start the HTTP server",
			RunE: di.commandWrap(func(args []string) error {
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

func (di *DIContainer) GetRDBConnOpener() *rdb.ConnectionOpener {
	if di.RdbConnOpener == nil {
		o := rdb.NewConnectionOpener()

		di.RdbConnOpener = o
	}

	return di.RdbConnOpener
}

func (di *DIContainer) GetRDBConnManager() *rdb.ConnectionManager {
	if di.RdbConnManager == nil {
		built, err := rdb.NewConnectionManager(
			di.Logger,
			di.GetRDBConnOpener(),
			rdb.WithDriver(di.Cfg.GetRDBDriver()),
			rdb.WithHost(di.Cfg.GetRDBHost()),
			rdb.WithPort(di.Cfg.GetRDBPort()),
			rdb.WithUsername(di.Cfg.GetRDBUsername()),
			rdb.WithPassword(di.Cfg.GetRDBPassword()),
			rdb.WithSchema(di.Cfg.GetRDBSchema()),
			rdb.WithDatabase(di.Cfg.GetRDBDatabase()),
		)

		if err != nil {
			di.Logger.Panic().Err(err).Msg(fmt.Sprintf("failed to build RDB connection manager: %s", err.Error()))
		}

		di.RdbConnManager = built
	}

	return di.RdbConnManager
}

func (di *DIContainer) GetRDBConnManagerForMigrations() *rdb.ConnectionManager {
	if di.RdbConnManagerForMigrations == nil {
		built, err := rdb.NewConnectionManager(
			di.Logger,
			di.GetRDBConnOpener(),
			rdb.WithDriver(di.Cfg.GetRDBDriver()),
			rdb.WithHost(di.Cfg.GetRDBHost()),
			rdb.WithPort(di.Cfg.GetRDBPort()),
			rdb.WithUsername(di.Cfg.GetRDBMigrationsUsername()),
			rdb.WithPassword(di.Cfg.GetRDBMigrationsPassword()),
			rdb.WithSchema(di.Cfg.GetRDBSchema()),
			rdb.WithDatabase(di.Cfg.GetRDBDatabase()),
		)

		if err != nil {
			di.Logger.Panic().Err(err).Msg(fmt.Sprintf("failed to build RDB connection manager for migrations: %s", err.Error()))
		}

		di.RdbConnManagerForMigrations = built
	}

	return di.RdbConnManagerForMigrations
}

func (di *DIContainer) GetPostgresSCMIntegrationsRepository() *scmpostgres.Repository {
	if di.PostgresSCMIntegrationsRepository == nil {
		connManager := di.GetRDBConnManager()

		di.PostgresSCMIntegrationsRepository = scmpostgres.NewRepository(di.Logger, connManager)
	}

	return di.PostgresSCMIntegrationsRepository
}

func (di *DIContainer) GetSCMIntegrationsTransformer() *scm.Transformer {
	if di.SCMIntegrationsTransformer == nil {
		di.SCMIntegrationsTransformer = scm.NewTransformer()
	}

	return di.SCMIntegrationsTransformer
}

func (di *DIContainer) GetSCMIntegrationsManager() *scm.Controller {
	if di.SCMIntegrationsManager == nil {
		di.SCMIntegrationsManager = scm.NewController(
			di.Logger,
			di.GetPostgresSCMIntegrationsRepository(),
			di.GetValidator(),
			di.GetSCMIntegrationsTransformer(),
		)
	}

	return di.SCMIntegrationsManager
}

func (di *DIContainer) GetErrorHandlingJsonUnmarshaller() *api.ErrorHandlingJsonUnmarshaller {
	if di.ErrorHandlingJsonUnmarshaller == nil {
		di.ErrorHandlingJsonUnmarshaller = api.NewErrorHandlingJsonUnmarshaller()
	}

	return di.ErrorHandlingJsonUnmarshaller
}

func (di *DIContainer) GetSCMIntegrationsController() *api.SCMHandler {
	if di.SCMIntegrationsController == nil {
		di.SCMIntegrationsController = api.NewSCMIntegrationsController(
			di.Logger,
			di.GetSCMIntegrationsManager(),
			di.GetErrorHandlingJsonUnmarshaller(),
		)
	}

	return di.SCMIntegrationsController
}

func (di *DIContainer) GetProjectsController() *api.ProjectsHandler {
	if di.ProjectsController == nil {
		di.ProjectsController = api.NewProjectsController()
	}

	return di.ProjectsController
}

func (di *DIContainer) GetServer() *api.Server {
	if di.Server == nil {
		s := api.NewServer(
			di.Logger,
			[]api.Controller{
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
