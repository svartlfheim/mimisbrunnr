package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
	IntegrationID *string `json:"scm_integration_id" validate:"required,uuid,exists"`
	Name          *string `json:"name" validate:"required,gt=0,unique"`
	Path          *string `json:"path" validate:"required,gt=0,uniqueperotherfield=scm_integration_id"`
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

func (dto UpdateProjectDTO) Validate(v StructValidator, r updateProjectValidationRepository, iR updateProjectIntegrationRepo, e *models.Project) ([]validation.ValidationError, error) {
	return v.ValidateStruct(dto, func(v *validator.Validate) *validator.Validate {
		err := v.RegisterValidation("uniqueperotherfield", func(fl validator.FieldLevel) bool {
			integrationIDField, kind, _, found := fl.GetStructFieldOKAdvanced2(fl.Parent(), "IntegrationID")

			if !found || kind.String() != "string" || integrationIDField.String() == "" {
				return false
			}

			integrationID, err := uuid.Parse(integrationIDField.String())

			if err != nil {
				return false
			}

			m, err := r.FindByPathAndIntegrationID(fl.Field().String(), integrationID)

			if err != nil {
				panic(err)
			}

			return m == nil ||
				// Can reuse it's own path again, if it was included in dto
				m.GetID().String() == e.GetID().String()
		})

		if err != nil {
			panic(err)
		}

		err = v.RegisterValidation("unique", func(fl validator.FieldLevel) bool {
			m, err := r.FindByName(fl.Field().String())

			if err != nil {
				panic(err)
			}

			return m == nil ||
				// Can reuse it's own name again, if it was included in dto
				m.GetID().String() == e.GetID().String()
		})

		if err != nil {
			panic(err)
		}

		err = v.RegisterValidation("exists", func(fl validator.FieldLevel) bool {
			integrationID, err := uuid.Parse(fl.Field().String())

			if err != nil {
				return false
			}

			m, err := iR.Find(integrationID)

			if err != nil {
				panic(err)
			}

			return m != nil
		})

		if err != nil {
			panic(err)
		}

		return v
	})
}

type updateIntegrationV1Response struct {
	updated          *TransformedProject
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

func Update(repo updateProjectRepository, iR updateProjectIntegrationRepo, v StructValidator, t Transformer, id string, dto UpdateProjectDTO) commandresult.Result {
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

	validationErrors, err := dto.Validate(v, repo, iR, existing)

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
			updated: t.ProjectV1(existing),
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
		updated:   t.ProjectV1(updated),
		changeset: cs,
	}
}
