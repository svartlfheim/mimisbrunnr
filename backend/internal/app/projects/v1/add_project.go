package v1

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

func (dto AddProjectDTO) Validate(v StructValidator, repo addProjectValidationRepository, integrationRepo addProjectIntegrationRepo) ([]validation.ValidationError, error) {
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

			m, err := repo.FindByPathAndIntegrationID(fl.Field().String(), integrationID)

			if err != nil {
				panic(err)
			}

			return m == nil
		})

		if err != nil {
			panic(err)
		}

		err = v.RegisterValidation("unique", func(fl validator.FieldLevel) bool {
			m, err := repo.FindByName(fl.Field().String())

			if err != nil {
				panic(err)
			}

			return m == nil
		})

		if err != nil {
			panic(err)
		}

		err = v.RegisterValidation("exists", func(fl validator.FieldLevel) bool {
			integrationID, err := uuid.Parse(fl.Field().String())

			if err != nil {
				return false
			}

			m, err := integrationRepo.Find(integrationID)

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

type addProjectV1Response struct {
	created          *TransformedProject
	errors           []error
	status           commandresult.Status
	validationErrors []validation.ValidationError
}

func (r *addProjectV1Response) Data() interface{} {
	return r.created
}

func (r *addProjectV1Response) Meta() interface{} {
	return map[string]interface{}{}
}

func (r *addProjectV1Response) Errors() []error {
	return r.errors
}

func (r *addProjectV1Response) ValidationErrors() []validation.ValidationError {
	return r.validationErrors
}

func (r *addProjectV1Response) Status() commandresult.Status {
	return r.status
}

func (r *addProjectV1Response) IsListData() bool {
	return false
}

func Add(repo addProjectRepo, iRepo addProjectIntegrationRepo, v StructValidator, t Transformer, dto AddProjectDTO) commandresult.Result {
	validationErrors, err := dto.Validate(v, repo, iRepo)

	if err != nil {
		return &addProjectV1Response{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	if len(validationErrors) > 0 {
		return &addProjectV1Response{
			status:           commandresult.Invalid,
			validationErrors: validationErrors,
		}
	}

	integrationID, err := uuid.Parse(*dto.IntegrationID)

	if err != nil {
		// Shouild never get here thanks to validation
		// if we do, somethings not quite right
		return &addProjectV1Response{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	i, err := iRepo.Find(integrationID)

	if err != nil {
		return &addProjectV1Response{
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
		return &addProjectV1Response{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	return &addProjectV1Response{
		status:  commandresult.Created,
		created: t.ProjectV1(m),
	}
}
