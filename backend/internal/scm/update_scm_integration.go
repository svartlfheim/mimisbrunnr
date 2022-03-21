package scm

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/validation"
	"github.com/svartlfheim/mimisbrunnr/pkg/commands/result"
)

type updateSCMIntegrationRepository interface {
	Patch(uuid.UUID, *models.ChangeSet) (*models.SCMIntegration, error)
	Find(uuid.UUID) (*models.SCMIntegration, error)
	FindByName(string) (*models.SCMIntegration, error)
}

type UpdateSCMIntegrationV1DTO struct {
	Name     *string `json:"name" validate:"omitempty,gt=0,unique"`
	Type     *string `json:"type" validate:"omitempty,gt=0,scmintegrationtype"`
	Endpoint *string `json:"endpoint" validate:"omitempty,gt=0"`
	Token    *string `json:"token" validate:"omitempty,gt=0"`
}

func (dto UpdateSCMIntegrationV1DTO) ToChangeSet(current *models.SCMIntegration) *models.ChangeSet {
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

type UpdateSCMIntegrationV1Response struct {
	updated          *scmIntegrationV1
	errors           []error
	status           result.Status
	validationErrors []validation.ValidationError
	changeset        *models.ChangeSet
}

func (r *UpdateSCMIntegrationV1Response) Data() interface{} {
	return r.updated
}

func (r *UpdateSCMIntegrationV1Response) Meta() interface{} {
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

func (r *UpdateSCMIntegrationV1Response) Errors() []error {
	return r.errors
}

func (r *UpdateSCMIntegrationV1Response) ValidationErrors() []validation.ValidationError {
	return r.validationErrors
}

func (r *UpdateSCMIntegrationV1Response) Status() result.Status {
	return r.status
}

func (r *UpdateSCMIntegrationV1Response) IsListData() bool {
	return false
}

func handleUpdateSCMIntegration(repo updateSCMIntegrationRepository, v structValidator, t scmIntegrationTransformerV1, id string, dto UpdateSCMIntegrationV1DTO) result.Result {
	uuid, err := uuid.Parse(id)

	if err != nil {
		return &UpdateSCMIntegrationV1Response{
			status: result.NotFound,
		}
	}

	existing, err := repo.Find(uuid)

	if err != nil {
		return &UpdateSCMIntegrationV1Response{
			errors: []error{
				err,
			},
			status: result.InternalError,
		}
	}

	if existing == nil {
		if err != nil {
			return &UpdateSCMIntegrationV1Response{
				status: result.NotFound,
			}
		}
	}

	validationErrors, err := v.ValidateStruct(dto, func(v *validator.Validate) *validator.Validate {
		err := v.RegisterValidation("unique", func(fl validator.FieldLevel) bool {
			m, err := repo.FindByName(fl.Field().String())

			if err != nil {
				panic(err)
			}

			return m == nil ||
				// Can reuse it's own name again, if it was included in dto
				m.ID.String() == existing.ID.String()
		})

		if err != nil {
			panic(err)
		}
		return v
	})

	if err != nil {
		return &UpdateSCMIntegrationV1Response{
			errors: []error{
				err,
			},
			status: result.InternalError,
		}
	}

	if len(validationErrors) > 0 {
		return &UpdateSCMIntegrationV1Response{
			status:           result.Invalid,
			validationErrors: validationErrors,
		}
	}

	cs := dto.ToChangeSet(existing)

	if cs.IsEmpty() {
		return &UpdateSCMIntegrationV1Response{
			status:  result.Okay,
			updated: t.SCMIntegrationV1(existing),
		}
	}

	updated, err := repo.Patch(uuid, cs)

	if err != nil {
		return &UpdateSCMIntegrationV1Response{
			errors: []error{
				err,
			},
			status: result.InternalError,
		}
	}

	return &UpdateSCMIntegrationV1Response{
		status:    result.Okay,
		updated:   t.SCMIntegrationV1(updated),
		changeset: cs,
	}
}
