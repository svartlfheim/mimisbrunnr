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
	ProjectV1(m *models.Project) *TransformedProject
	ProjectListV1(m []*models.Project) []*TransformedProject
}

type RequiredRepository interface {
	addProjectRepo
	listProjectsRepository
	getProjectRepository
	deleteProjectRepository
	updateProjectRepository
}

type RequiredIntegrationRepository interface {
	addProjectIntegrationRepo
	updateProjectIntegrationRepo
}
