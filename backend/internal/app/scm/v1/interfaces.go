package v1

import (
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/validation"
)

type StructValidator interface {
	ValidateStruct(s interface{}, opts ...validation.WithValidationExtension) ([]validation.ValidationError, error)
	RegisterCustomValidation(t string, cv validation.CustomValidation)
}

type Transformer interface {
	IntegrationV1(m *models.SCMIntegration) *TransformedSCMIntegration
	IntegrationListV1(m []*models.SCMIntegration) []*TransformedSCMIntegration
}

type RequiredRepository interface {
	addIntegrationRepository
	getIntegrationRepository
	listIntegrationsRepository
	updateIntegrationRepository
	deleteIntegrationRepository
}