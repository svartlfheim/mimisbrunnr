package cmd

import (
	"github.com/svartlfheim/mimisbrunnr/internal/config"
	"github.com/svartlfheim/mimisbrunnr/internal/server"
)

func handleServe(srv *server.Server, cfg *config.AppConfig, args []string) error {
	return srv.Start(cfg)
}
