package cmd

import (
	"github.com/svartlfheim/mimisbrunnr/internal/config"
	"github.com/svartlfheim/mimisbrunnr/internal/server"
)

type _server interface {
	Start(cfg server.ServerConfig) error
}

type ServeHandler struct {
	Server _server
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

func NewServeHandler(srv _server) *ServeHandler {
	return &ServeHandler{
		Server: srv,
	}
}
