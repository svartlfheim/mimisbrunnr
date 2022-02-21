package cmd

import (
	"github.com/svartlfheim/mimisbrunnr/internal/config"
	"github.com/svartlfheim/mimisbrunnr/internal/server"
)

type ServeHandler struct {
	Server *server.Server
}

func (s *ServeHandler) Handle(cfg *config.AppConfig, args []string) error {
	return s.Server.Start(cfg)
}

func (s *ServeHandler) GetName() string {
	return "serve"
}

func (s *ServeHandler) GetHelp() string {
	return "help for serve"
}

func NewServeHandler(di *DIContainer) *ServeHandler {
	return &ServeHandler{
		Server: di.GetServer(),
	}
}

