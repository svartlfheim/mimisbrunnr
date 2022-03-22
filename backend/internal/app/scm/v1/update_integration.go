package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

func (dto UpdateIntegrationDTO) Validate(v StructValidator, repo updateIntegrationValidationRepository, existing *models.SCMIntegration) ([]validation.ValidationError, error) {
	return v.ValidateStruct(dto, func(v *validator.Validate) *validator.Validate {
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
}

type updateIntegrationV1Response struct {
	updated          *TransformedSCMIntegration
	errors           []error
	status           commandresult.Status
	validationErrors []validation.ValidationError
	changeset        *models.ChangeSet
}

func (r *updateIntegrationV1Response) Data() interface{} {
	return r.updated
}

func (r *updateIntegrationV1Response) Meta() interface{} {
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

func (r *updateIntegrationV1Response) Errors() []error {
	return r.errors
}

func (r *updateIntegrationV1Response) ValidationErrors() []validation.ValidationError {
	return r.validationErrors
}

func (r *updateIntegrationV1Response) Status() commandresult.Status {
	return r.status
}

func (r *updateIntegrationV1Response) IsListData() bool {
	return false
}

func Update(repo updateIntegrationRepository, v StructValidator, t Transformer, id string, dto UpdateIntegrationDTO) commandresult.Result {
	uuid, err := uuid.Parse(id)

	if err != nil {
		return &updateIntegrationV1Response{
			status: commandresult.NotFound,
		}
	}

	existing, err := repo.Find(uuid)

	if err != nil {
		return &updateIntegrationV1Response{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	if existing == nil {
		if err != nil {
			return &updateIntegrationV1Response{
				status: commandresult.NotFound,
			}
		}
	}

	validationErrors, err := dto.Validate(v, repo, existing)

	if err != nil {
		return &updateIntegrationV1Response{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	if len(validationErrors) > 0 {
		return &updateIntegrationV1Response{
			status:           commandresult.Invalid,
			validationErrors: validationErrors,
		}
	}

	cs := dto.ToChangeSet(existing)

	if cs.IsEmpty() {
		return &updateIntegrationV1Response{
			status:  commandresult.Okay,
			updated: t.IntegrationV1(existing),
		}
	}

	updated, err := repo.Patch(uuid, cs)

	if err != nil {
		return &updateIntegrationV1Response{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	return &updateIntegrationV1Response{
		status:    commandresult.Okay,
		updated:   t.IntegrationV1(updated),
		changeset: cs,
	}
}
