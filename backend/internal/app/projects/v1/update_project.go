package v1

import (
	"github.com/google/uuid"
	"github.com/svartlfheim/mimisbrunnr/internal/app/projects/rules"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/commandresult"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/validation"
)

type updateProjectIntegrationRepo interface {
	Find(uuid.UUID) (*models.SCMIntegration, error)
}

type updateProjectValidationRepository interface {
	FindByName(string) (*models.Project, error)
	FindByPathAndIntegrationID(string, uuid.UUID) (*models.Project, error)
}

type updateProjectRepository interface {
	updateProjectValidationRepository

	Patch(uuid.UUID, *models.ChangeSet) (*models.Project, error)
	Find(uuid.UUID) (*models.Project, error)
}

type UpdateProjectDTO struct {
	IntegrationID *string `json:"scm_integration_id" validate:"omitempty,uuid,exists"`
	Name          *string `json:"name" validate:"omitempty,gt=0,unique"`
	Path          *string `json:"path" validate:"omitempty,gt=0,uniqueperotherfield=scm_integration_id"`
}

func (dto UpdateProjectDTO) ToChangeSet(current *models.Project) *models.ChangeSet {
	cs := models.NewChangeSet()

	if dto.Name != nil && *dto.Name != current.GetName() {
		cs.RegisterChange("Name", *dto.Name)
	}

	if dto.Path != nil && *dto.Path != string(current.GetPath()) {
		cs.RegisterChange("Path", *dto.Path)
	}

	if dto.IntegrationID != nil && *dto.IntegrationID != current.GetSCMIntegration().GetID().String() {
		cs.RegisterChange("SCMIntegration", *dto.IntegrationID)
	}

	return cs
}

type updateIntegrationResponse struct {
	updated          *models.Project
	errors           []error
	status           commandresult.Status
	validationErrors []validation.ValidationError
	changeset        *models.ChangeSet
}

func (r *updateIntegrationResponse) Data() interface{} {
	return r.updated
}

func (r *updateIntegrationResponse) Meta() interface{} {
	modifiedFields := []string{}

	if r.changeset != nil {
		for k := range r.changeset.Changes {
			modifiedFields = append(modifiedFields, k)
		}
	}

	return map[string]interface{}{
		"modified": modifiedFields,
	}
}

func (r *updateIntegrationResponse) Errors() []error {
	return r.errors
}

func (r *updateIntegrationResponse) ValidationErrors() []validation.ValidationError {
	return r.validationErrors
}

func (r *updateIntegrationResponse) Status() commandresult.Status {
	return r.status
}

func Update(repo updateProjectRepository, iR updateProjectIntegrationRepo, v StructValidator, id string, dto UpdateProjectDTO) commandresult.Result {
	uuid, err := uuid.Parse(id)

	if err != nil {
		return &updateIntegrationResponse{
			status: commandresult.NotFound,
		}
	}

	existing, err := repo.Find(uuid)

	if err != nil {
		return &updateIntegrationResponse{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	if existing == nil {
		if err != nil {
			return &updateIntegrationResponse{
				status: commandresult.NotFound,
			}
		}
	}

	validationErrors, err := v.ValidateStruct(
		dto, 
		rules.Unique(repo, existing),
		rules.Exists(iR),
		rules.UniquePerIntegration(repo, existing),
	)

	if err != nil {
		return &updateIntegrationResponse{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	if len(validationErrors) > 0 {
		return &updateIntegrationResponse{
			status:           commandresult.Invalid,
			validationErrors: validationErrors,
		}
	}

	cs := dto.ToChangeSet(existing)

	if cs.IsEmpty() {
		return &updateIntegrationResponse{
			status:  commandresult.Okay,
			updated: existing,
		}
	}

	updated, err := repo.Patch(uuid, cs)


	if err != nil {
		return &updateIntegrationResponse{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	return &updateIntegrationResponse{
		status:    commandresult.Okay,
		updated:   updated,
		changeset: cs,
	}
}
