package cmd

import (
	"github.com/svartlfheim/mimisbrunnr/internal/config"
	"github.com/svartlfheim/mimisbrunnr/internal/server"
)

type serverForHandler interface {
	Start(cfg server.ServerConfig) error
}

type ServeHandler struct {
	Server serverForHandler
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

func NewServeHandler(srv serverForHandler) *ServeHandler {
	return &ServeHandler{
		Server: srv,
	}
}
