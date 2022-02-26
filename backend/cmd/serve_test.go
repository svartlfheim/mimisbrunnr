package cmd_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/svartlfheim/mimisbrunnr/cmd"
	"github.com/svartlfheim/mimisbrunnr/internal/config"
	cmdmocks "github.com/svartlfheim/mimisbrunnr/test/mocks/cmd"
)

func Test_ServeHandler_Handle_returns_error_from_server(t *testing.T) {
	cfg := &config.AppConfig{}

	serverMock := &cmdmocks.ServerForHandler{}
	defer serverMock.AssertExpectations(t)
	defer serverMock.AssertNumberOfCalls(t, "Start", 1)

	serverErr := errors.New("server returned error")
	serverMock.EXPECT().Start(cfg).Return(serverErr)

	h := cmd.NewServeHandler(serverMock)
	err := h.Handle(cfg, []string{})

	assert.Equal(t, err, serverErr)
}

func Test_ServeHandler_Handle_returns_nil_from_server(t *testing.T) {
	cfg := &config.AppConfig{}

	serverMock := &cmdmocks.ServerForHandler{}
	defer serverMock.AssertExpectations(t)
	defer serverMock.AssertNumberOfCalls(t, "Start", 1)

	serverMock.EXPECT().Start(cfg).Return(nil)

	h := cmd.NewServeHandler(serverMock)
	err := h.Handle(cfg, []string{})

	assert.Nil(t, err)
}

func Test_ServeHandler_Getters(t *testing.T) {
	serverMock := &cmdmocks.ServerForHandler{}
	defer serverMock.AssertExpectations(t)
	defer serverMock.AssertNumberOfCalls(t, "Start", 0)

	h := cmd.NewServeHandler(serverMock)

	assert.Equal(t, "serve", h.GetName())
	// Maybe remove this later...
	assert.Equal(t, "help for serve", h.GetHelp())
}
