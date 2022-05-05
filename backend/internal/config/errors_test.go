package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/svartlfheim/mimisbrunnr/internal/config"
)

func Test_ErrConfigFileDoesNotExist(t *testing.T) {
	err := config.ErrConfigFileDoesNotExist{
		Path: "/some/path",
	}

	assert.Equal(t, "config file not found at path '/some/path'", err.Error())
}

func Test_ErrCannotUnmarshalConfig(t *testing.T) {
	err := config.ErrCannotUnmarshalConfig{
		Message: "bad times",
	}

	assert.Equal(t, "config file contents could not be unmarshalled: bad times", err.Error())
}

func Test_ErrCouldNotProcessEnv(t *testing.T) {
	err := config.ErrCouldNotProcessEnv{
		Message: "bad times",
	}

	assert.Equal(t, "could not process env overrides for config: bad times", err.Error())
}

func Test_ErrFsUnusable(t *testing.T) {
	err := config.ErrFsUnusable{
		Message: "bad times",
	}

	assert.Equal(t, "FS not usable: bad times", err.Error())
}
