package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/svartlfheim/mimisbrunnr/cmd"
	"github.com/svartlfheim/mimisbrunnr/internal/cmdregistry"
	"github.com/svartlfheim/mimisbrunnr/internal/config"
)

func buildLogger() zerolog.Logger {
	return zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
}

func buildFS() afero.Fs {
	return afero.NewOsFs()
}

func showHelp(r *cmdregistry.Registry, err error) {
	flag.Usage()
	fmt.Println("")
	fmt.Print(r.GetHelp(err))
}

func main() {
	logger := buildLogger()
	fs := buildFS()

	logLevel := flag.Int("log", int(zerolog.InfoLevel), "Log level: -1=trace, 0=debug, 1=info, 2=warn, 3=error, 4=fatal, 5=panic")
	cfgPath := flag.String("f", "./mimisbrunnr.yaml", "Path to a YAML config file to load.")
	flag.Parse()

	cfg, err := config.Load(*cfgPath, fs)

	dic := cmd.NewDIContainer(logger, fs, cfg)

	fmt.Printf("%#v\n", cfg)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to load config")
	}

	logger = logger.Level(zerolog.Level(*logLevel))

	allArgs := flag.Args()

	r := dic.GetRootCommandRegistry()

	if len(allArgs) < 1 {
		showHelp(r, nil)

		os.Exit(1)
	}

	cmd := allArgs[0]

	args := allArgs[1:]

	h, err := r.FindHandler(cmd)

	if err != nil {
		showHelp(r, err)

		os.Exit(1)
	}

	if err := h.Handle(cfg, args); err != nil {
		panic(err)
	}
}
