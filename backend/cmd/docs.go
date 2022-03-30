package cmd

import (
	"fmt"

	"github.com/svartlfheim/mimisbrunnr/internal/app/openapi"
)

func handleDocsOpenAPI(g *openapi.Generator, args []string) error {
	out, err := g.Generate(true)
	if err != nil {
		return err
	}

	fmt.Print(string(out))

	return nil
}
