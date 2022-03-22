package cmd

import (
	"github.com/svartlfheim/mimisbrunnr/internal/app/api"
	"github.com/svartlfheim/mimisbrunnr/internal/config"
)

func handleServe(srv *api.Server, cfg *config.AppConfig, args []string) error {
	return srv.Start(cfg)
}
