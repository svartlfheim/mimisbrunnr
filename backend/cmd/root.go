package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "mmrbrnr",
		Short: "View documentation from your SCM repositories.",
		Long:  `<coming soon>`,
	}
)

func init() {
	rootCmd.PersistentFlags().IntP("log-level", "l", int(zerolog.InfoLevel), "Log level: -1=trace, 0=debug, 1=info, 2=warn, 3=error, 4=fatal, 5=panic")
	rootCmd.PersistentFlags().StringP("config", "c", "./mimisbrunnr.yaml", "Path to a YAML config file to load.")
}

func buildLogger() zerolog.Logger {
	return zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
}

func buildFS() afero.Fs {
	return afero.NewOsFs()
}

func Execute() error {

	fs := buildFS()
	logger := buildLogger()

	dic := NewDIContainer(logger, fs)

	for _, c := range dic.GetCommands() {
		rootCmd.AddCommand(c)
	}

	return rootCmd.Execute()
}
