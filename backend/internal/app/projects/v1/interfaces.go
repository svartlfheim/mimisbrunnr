package v1

import (
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/validation"
)

type StructValidator interface {
	ValidateStruct(s interface{}, opts ...validation.WithValidationExtension) ([]validation.ValidationError, error)
	RegisterCustomValidation(t string, cv validation.CustomValidation)
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
