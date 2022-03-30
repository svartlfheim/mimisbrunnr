package v1

import (
	"time"

	"github.com/google/uuid"
	"github.com/svartlfheim/mimisbrunnr/internal/app/projects/rules"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/commandresult"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/validation"
)

type addProjectValidationRepository interface {
	FindByName(string) (*models.Project, error)
	FindByPathAndIntegrationID(string, uuid.UUID) (*models.Project, error)
}

type addProjectIntegrationRepo interface {
	Find(uuid.UUID) (*models.SCMIntegration, error)
}

type addProjectRepo interface {
	addProjectValidationRepository

	Create(*models.Project) error
}

type AddProjectDTO struct {
	IntegrationID *string `json:"scm_integration_id" validate:"required,uuid,exists"`
	Name          *string `json:"name" validate:"required,gt=0,unique"`
	Path          *string `json:"path" validate:"required,gt=0,uniqueperotherfield=scm_integration_id"`
}

type addProjectResponse struct {
	created          *models.Project
	errors           []error
	status           commandresult.Status
	validationErrors []validation.ValidationError
}

func (r *addProjectResponse) Data() interface{} {
	return r.created
}

func (r *addProjectResponse) Meta() interface{} {
	return nil
}

func (r *addProjectResponse) Errors() []error {
	return r.errors
}

func (r *addProjectResponse) ValidationErrors() []validation.ValidationError {
	return r.validationErrors
}

func (r *addProjectResponse) Status() commandresult.Status {
	return r.status
}

func Add(repo addProjectRepo, iRepo addProjectIntegrationRepo, v StructValidator, dto AddProjectDTO) commandresult.Result {
	validationErrors, err := v.ValidateStruct(
		dto,
		rules.Unique(repo, nil),
		rules.Exists(iRepo),
		rules.UniquePerIntegration(repo, nil),
	)

	if err != nil {
		return &addProjectResponse{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	if len(validationErrors) > 0 {
		return &addProjectResponse{
			status:           commandresult.Invalid,
			validationErrors: validationErrors,
		}
	}

	integrationID, err := uuid.Parse(*dto.IntegrationID)

	if err != nil {
		// Shouild never get here thanks to validation
		// if we do, somethings not quite right
		return &addProjectResponse{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	i, err := iRepo.Find(integrationID)

	if err != nil {
		return &addProjectResponse{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	m := models.NewProject(
		uuid.New(),
		*dto.Name,
		*dto.Path,
		i,
		time.Now().UTC(),
		time.Now().UTC(),
	)

	if err := repo.Create(m); err != nil {
		return &addProjectResponse{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	return &addProjectResponse{
		status:  commandresult.Created,
		created: m,
	}
}
