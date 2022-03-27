package cmd

import (
	"github.com/spf13/afero"
	"github.com/svartlfheim/mimisbrunnr/internal/app/openapi"
	"github.com/svartlfheim/mimisbrunnr/internal/config"
)

func handleDocsOpenAPI(cfg *config.AppConfig, fs afero.Fs, args []string) error {
	return openapi.Generate(cfg, fs)
}
