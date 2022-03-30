package v1

import (
	"github.com/google/uuid"
	"github.com/svartlfheim/mimisbrunnr/internal/app/scm/rules"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/commandresult"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/validation"
)

type updateIntegrationRepository interface {
	updateIntegrationValidationRepository

	Patch(uuid.UUID, *models.ChangeSet) (*models.SCMIntegration, error)
	Find(uuid.UUID) (*models.SCMIntegration, error)
}

type updateIntegrationValidationRepository interface {
	FindByName(string) (*models.SCMIntegration, error)
}

type UpdateIntegrationDTO struct {
	Name     *string `json:"name" validate:"omitempty,gt=0,unique"`
	Type     *string `json:"type" validate:"omitempty,gt=0,scmintegrationtype"`
	Endpoint *string `json:"endpoint" validate:"omitempty,gt=0"`
	Token    *string `json:"token" validate:"omitempty,gt=0"`
}

func (dto UpdateIntegrationDTO) ToChangeSet(current *models.SCMIntegration) *models.ChangeSet {
	cs := models.NewChangeSet()

	if dto.Name != nil && *dto.Name != current.GetName() {
		cs.RegisterChange("Name", *dto.Name)
	}

	if dto.Type != nil && *dto.Type != string(current.GetType()) {
		cs.RegisterChange("Type", *dto.Type)
	}

	if dto.Endpoint != nil && *dto.Endpoint != current.GetEndpoint() {
		cs.RegisterChange("Endpoint", *dto.Endpoint)
	}

	if dto.Token != nil && *dto.Token != current.GetToken() {
		cs.RegisterChange("Token", *dto.Token)
	}

	return cs
}

type updateIntegrationResponse struct {
	updated          *models.SCMIntegration
	errors           []error
	status           commandresult.Status
	validationErrors []validation.ValidationError
	changeset        *models.ChangeSet
}

func (r *updateIntegrationResponse) Data() interface{} {
	return r.updated
}

func (r *updateIntegrationResponse) Meta() interface{} {
	if r.changeset == nil {
		return nil
	}

	modifiedFields := []string{}

	for k := range r.changeset.Changes {
		modifiedFields = append(modifiedFields, k)
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

func Update(repo updateIntegrationRepository, v StructValidator, id string, dto UpdateIntegrationDTO) commandresult.Result {
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

	validationErrors, err := v.ValidateStruct(dto, rules.Unique(repo, existing))

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
