package main

import (
	"fmt"
	"os"

	"github.com/svartlfheim/mimisbrunnr/internal/cmdregistry"
	"github.com/svartlfheim/mimisbrunnr/internal/config"
	"github.com/svartlfheim/mimisbrunnr/internal/server"
)
var cfg *config.AppConfig = config.New()
var registry *cmdregistry.Registry = cmdregistry.NewRegistry()

func init() {
	err := registry.Register(
		server.NewServer(
			cfg,
			[]server.Controller{
				server.NewGitHostsController(),
				server.NewProjectsController(),
			},
		),
	)

	if err != nil {
		panic(err)
	}
}

func main() {
	args := os.Args

	if len(args) == 1 {
		fmt.Print(registry.GetHelp(nil))

		os.Exit(1)
	}

	cmd := args[1]
	h, err := registry.FindHandler(cmd)

	if err != nil {
		fmt.Print(registry.GetHelp(err))

		os.Exit(1)
	}

	if err := h.Handle(); err != nil {
		panic(err)
	}
}