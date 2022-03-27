package cmd

import (
	"github.com/svartlfheim/mimisbrunnr/internal/config"
	"github.com/svartlfheim/mimisbrunnr/internal/infra/httpsrv"
)

func handleServe(srv *httpsrv.Server, cfg *config.AppConfig, args []string) error {
	return srv.Start(cfg)
}
