package scm

import (
	"github.com/rs/zerolog"
	"github.com/svartlfheim/mimisbrunnr/internal/validation"
	"github.com/svartlfheim/mimisbrunnr/pkg/commands/result"
)

type managerRepository interface {
	addSCMIntegrationRepository
}

type structValidator interface {
	ValidateStruct(s interface{}, opts ...validation.WithValidationExtension) ([]validation.ValidationError, error)
}

type Manager struct {
	logger    zerolog.Logger
	repo      managerRepository
	validator structValidator
}

func (m *Manager) AddV1(dto AddSCMIntegrationV1DTO) result.Result {

	return handleAddSCMIntegration(m.repo, m.validator, dto)
}

func NewManager(l zerolog.Logger, repo managerRepository, v structValidator) *Manager {
	return &Manager{
		logger:    l,
		repo:      repo,
		validator: v,
	}
}
