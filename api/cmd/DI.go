package cmd

import (
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/svartlfheim/mimisbrunnr/internal/cmdregistry"
	"github.com/svartlfheim/mimisbrunnr/internal/config"
	"github.com/svartlfheim/mimisbrunnr/internal/rdbconn"
	"github.com/svartlfheim/mimisbrunnr/internal/scm"
	"github.com/svartlfheim/mimisbrunnr/internal/server"
)

type DIContainer struct {
	Cfg *config.AppConfig
	RootCmdRegistry *cmdregistry.Registry
	RdbConnManager *rdbconn.ConnectionManager
	PostgresSCMIntegrationsRepository *scm.PostgresRepository
	Server *server.Server
	SCMIntegrationsController *server.SCMIntegrationsController
	ProjectsController *server.ProjectsController
	SCMIntegrationsManager *scm.Manager
	Logger zerolog.Logger
	Fs afero.Fs
}

func (di *DIContainer) GetRootCommandRegistry() *cmdregistry.Registry {
	if di.RootCmdRegistry == nil {
		r := cmdregistry.NewRegistry(di.Logger)

		if err := r.Register(NewServeHandler(di)); err != nil {
			di.Logger.Fatal().Err(err).Msg("failed to register serve handler")
		}

		di.RootCmdRegistry = r
	}

	return di.RootCmdRegistry
}


func (di *DIContainer) GetRDBConnManager() *rdbconn.ConnectionManager {
	if di.RdbConnManager == nil {
		built, err := rdbconn.NewConnectionManager(
			di.Logger,
			rdbconn.WithDriver(di.Cfg.GetRDBDriver()),
			rdbconn.WithHost(di.Cfg.GetRDBHost()),
			rdbconn.WithPort(di.Cfg.GetRDBPort()),
			rdbconn.WithUsername(di.Cfg.GetRDBUsername()),
			rdbconn.WithPassword(di.Cfg.GetRDBPassword()),
			rdbconn.WithSchema(di.Cfg.GetRDBSchema()),
			rdbconn.WithDatabase(di.Cfg.GetRDBDatabase()),
		)

		if err != nil {
			di.Logger.Fatal().Err(err).Msg("failed to build RDB connection manager")
		}

		di.RdbConnManager = built
	}

	return di.RdbConnManager
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
		)
	}

	return di.SCMIntegrationsManager
}

func (di *DIContainer) GetSCMIntegrationsController() *server.SCMIntegrationsController {
	if di.SCMIntegrationsController == nil {
		di.SCMIntegrationsController = server.NewSCMIntegrationsController(
			di.Logger,
			di.GetSCMIntegrationsManager(),
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

func NewDIContainer(l zerolog.Logger, fs afero.Fs, cfg *config.AppConfig) *DIContainer {
	return &DIContainer{
		Logger: l,
		Fs: fs,
		Cfg: cfg,
	}
}