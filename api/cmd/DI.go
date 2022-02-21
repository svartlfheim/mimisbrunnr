package cmd

import (
	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/svartlfheim/mimisbrunnr/internal/cmdregistry"
	"github.com/svartlfheim/mimisbrunnr/internal/config"
	"github.com/svartlfheim/mimisbrunnr/internal/githosts"
	"github.com/svartlfheim/mimisbrunnr/internal/rdbconn"
	"github.com/svartlfheim/mimisbrunnr/internal/server"
)

type DIContainer struct {
	Cfg *config.AppConfig
	RootCmdRegistry *cmdregistry.Registry
	RdbConnManager *rdbconn.ConnectionManager
	PostgresGitHostsRepository *githosts.PostgresRepository
	Server *server.Server
	GitHostsController *server.GitHostsController
	ProjectsController *server.ProjectsController
	GitHostsManager *githosts.Manager
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

func (di *DIContainer) GetPostgresGitHostsRepository() *githosts.PostgresRepository {
	if di.PostgresGitHostsRepository == nil {
		connManager := di.GetRDBConnManager()

		di.PostgresGitHostsRepository = githosts.NewPostgresRepository(connManager, di.Logger)
	}

	return di.PostgresGitHostsRepository
}

func (di *DIContainer) GetGitHostsManager() *githosts.Manager {
	if di.GitHostsManager == nil {
		di.GitHostsManager = githosts.NewManager(
			di.Logger,
			di.GetPostgresGitHostsRepository(),
		)
	}

	return di.GitHostsManager
}

func (di *DIContainer) GetGitHostsController() *server.GitHostsController {
	if di.GitHostsController == nil {
		di.GitHostsController = server.NewGitHostsController(
			di.Logger,
			di.GetGitHostsManager(),
		)
	}

	return di.GitHostsController
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
				di.GetGitHostsController(),
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