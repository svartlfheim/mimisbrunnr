package scm

import (
	"github.com/rs/zerolog"
	"github.com/svartlfheim/mimisbrunnr/pkg/commands/result"
)

type managerRepository interface {
	addSCMIntegrationRepository
}


type Manager struct {
	logger zerolog.Logger
	repo managerRepository
}

func (m *Manager) Add(dto AddSCMIntegrationDTO) (result.Result) {

	return handleAddSCMIntegration(m.repo, dto)
}


func NewManager(l zerolog.Logger, repo managerRepository) *Manager {
	return &Manager{
		logger: l,
		repo: repo,
	}
}