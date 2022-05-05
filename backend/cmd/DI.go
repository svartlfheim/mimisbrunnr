package cmd

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/svartlfheim/gomigrator"
	"github.com/svartlfheim/mimisbrunnr/internal/app/projects"
	"github.com/svartlfheim/mimisbrunnr/internal/app/scm"
	"github.com/svartlfheim/mimisbrunnr/internal/app/web"
	"github.com/svartlfheim/mimisbrunnr/internal/config"
	"github.com/svartlfheim/mimisbrunnr/internal/infra/httpsrv"
	"github.com/svartlfheim/mimisbrunnr/internal/infra/rdb"
	"github.com/svartlfheim/mimisbrunnr/internal/infra/rdb/postgres"
	"github.com/svartlfheim/mimisbrunnr/internal/infra/rdb/schema"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/validation"
)

type DIContainer struct {
	Cfg *config.AppConfig

	// rdb.*
	RdbConnManager              *rdb.ConnectionManager
	RdbConnManagerForMigrations *rdb.ConnectionManager
	RdbConnOpener               *rdb.ConnectionOpener

	// scm.*
	PostgresSCMIntegrationsRepository *postgres.SCMIntegrationsRepository
	SCMIntegrationsController         *scm.Controller

	// projects.*
	ProjectsController  *projects.Controller
	ProjectsRepository  *postgres.ProjectsRepository

	// schema.*
	Migrator *gomigrator.Migrator

	// web.*
	Server                        *httpsrv.Server
	SCMIntegrationsHandler        *web.SCMHandler
	ProjectsAPIHandler            *web.ProjectsHandler
	APITransformer                *web.Transformer
	APIResponseBuilder *web.ResponseBuilder
	ErrorHandlingJsonUnmarshaller *web.ErrorHandlingJsonUnmarshaller

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

	serve := &cobra.Command{
		Use:   "serve",
		Short: "Start the HTTP server",
		RunE: di.commandWrap(func(args []string) error {
			return handleServe(di.GetServer(), di.Cfg, args)
		}),
	}

	docs := &cobra.Command{
		Use:   "docs",
		Short: "Various commands around documentation for the app",
	}

	docs.AddCommand(&cobra.Command{
		Use:   "openapi",
		Short: "Generates openapi docs",
		RunE: di.commandWrap(func(args []string) error {
			return handleDocsOpenAPI(di.Cfg, di.Fs, args)
		}),
	})

	return []*cobra.Command{
		serve,
		migrate,
		docs,
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

func (di *DIContainer) GetPostgresSCMIntegrationsRepository() *postgres.SCMIntegrationsRepository {
	if di.PostgresSCMIntegrationsRepository == nil {
		connManager := di.GetRDBConnManager()

		di.PostgresSCMIntegrationsRepository = postgres.NewSCMIntegrationsRepository(di.Logger, connManager)
	}

	return di.PostgresSCMIntegrationsRepository
}

func (di *DIContainer) GetSCMIntegrationsController() *scm.Controller {
	if di.SCMIntegrationsController == nil {
		di.SCMIntegrationsController = scm.NewController(
			di.Logger,
			di.GetPostgresSCMIntegrationsRepository(),
			di.GetValidator(),
		)
	}

	return di.SCMIntegrationsController
}

func (di *DIContainer) GetPostgresProjectsRepository() *postgres.ProjectsRepository {
	if di.ProjectsRepository == nil {
		connManager := di.GetRDBConnManager()

		di.ProjectsRepository = postgres.NewProjectsRepository(di.Logger, connManager)
	}

	return di.ProjectsRepository
}

func (di *DIContainer) GetProjectsController() *projects.Controller {
	if di.ProjectsController == nil {
		di.ProjectsController = projects.NewController(
			di.Logger,
			di.GetPostgresProjectsRepository(),
			di.GetPostgresSCMIntegrationsRepository(),
			di.GetValidator(),
		)
	}

	return di.ProjectsController
}

func (di *DIContainer) GetErrorHandlingJsonUnmarshaller() *web.ErrorHandlingJsonUnmarshaller {
	if di.ErrorHandlingJsonUnmarshaller == nil {
		di.ErrorHandlingJsonUnmarshaller = web.NewErrorHandlingJsonUnmarshaller()
	}

	return di.ErrorHandlingJsonUnmarshaller
}

func (di *DIContainer) GetSCMIntegrationsHandler() *web.SCMHandler {
	if di.SCMIntegrationsHandler == nil {
		di.SCMIntegrationsHandler = web.NewSCMIntegrationsHandler(
			di.Logger,
			di.GetSCMIntegrationsController(),
			di.GetErrorHandlingJsonUnmarshaller(),
			di.GetAPIResponseBuilder(),
		)
	}

	return di.SCMIntegrationsHandler
}

func (di *DIContainer) GetProjectsAPIHandler() *web.ProjectsHandler {
	if di.ProjectsAPIHandler == nil {
		di.ProjectsAPIHandler = web.NewProjectsHandler(
			di.Logger,
			di.GetProjectsController(),
			di.GetErrorHandlingJsonUnmarshaller(),
			di.GetAPIResponseBuilder(),
		)
	}

	return di.ProjectsAPIHandler
}

func (di *DIContainer) GetAPITransformer() *web.Transformer {
	if di.APITransformer == nil {
		di.APITransformer = web.NewTransformer()
	}

	return di.APITransformer
}


func (di *DIContainer) GetAPIResponseBuilder() *web.ResponseBuilder {
	if di.APIResponseBuilder == nil {
		di.APIResponseBuilder = web.NewResponseBuilder(
			di.Logger,
			di.GetAPITransformer(),
		)
	}

	return di.APIResponseBuilder
}

func (di *DIContainer) GetServer() *httpsrv.Server {
	if di.Server == nil {
		s := httpsrv.NewServer(
			di.Logger,
			[]httpsrv.ApiController{
				di.GetSCMIntegrationsHandler(),
				di.GetProjectsAPIHandler(),
			},
		)

		di.Server = s
	}

	return di.Server
}

func (di *DIContainer) GetValidator() *validation.Validator {
	return validation.NewValidator(di.Logger)
}

func NewDIContainer(l zerolog.Logger, fs afero.Fs) *DIContainer {
	return &DIContainer{
		Logger: l,
		Fs:     fs,
	}
}
