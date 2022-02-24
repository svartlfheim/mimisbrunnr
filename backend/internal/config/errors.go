package config

import "fmt"

type ErrConfigFileDoesNotExist struct {
	Path string
}

func (e ErrConfigFileDoesNotExist) Error() string {
	return fmt.Sprintf("config file not found at path '%s'", e.Path)
}

type ErrCannotUnmarshalConfig struct {
	Message string
}

func (e ErrCannotUnmarshalConfig) Error() string {
	return fmt.Sprintf("config file contents could not be unmarshalled: %s", e.Message)
}

type ErrCouldNotProcessEnv struct {
	Message string
}

func (e ErrCouldNotProcessEnv) Error() string {
	return fmt.Sprintf("could not process env overrides for config: %s", e.Message)
}

type ErrFsUnusable struct {
	Message string
}

func (e ErrFsUnusable) Error() string {
	return fmt.Sprintf("FS not usable: %s", e.Message)
}
