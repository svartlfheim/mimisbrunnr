package scm

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/validation"
	"github.com/svartlfheim/mimisbrunnr/pkg/commands/result"
)

type addSCMIntegrationRepository interface {
	Create(*models.SCMIntegration) error
	FindByName(string) (*models.SCMIntegration, error)
}

type AddSCMIntegrationV1DTO struct {
	Name     *string `json:"name" validate:"required,gt=0,unique"`
	Type     *string `json:"type" validate:"required,gt=0,scmintegrationtype"`
	Endpoint *string `json:"endpoint" validate:"required,gt=0"`
	Token    *string `json:"token" validate:"required,gt=0"`
}

func (dto AddSCMIntegrationV1DTO) ToModel() *models.SCMIntegration {
	return models.NewSCMIntegration(
		uuid.New(),
		*dto.Name,
		models.SCMIntegrationType(*dto.Type),
		*dto.Endpoint,
		*dto.Token,
		time.Now(),
		time.Now(),
	)
}

type AddSCMIntegrationV1Response struct {
	created          *scmIntegrationV1
	errors           []error
	status           result.Status
	validationErrors []validation.ValidationError
}

func (r *AddSCMIntegrationV1Response) Data() interface{} {
	return r.created
}

func (r *AddSCMIntegrationV1Response) Meta() interface{} {
	return map[string]interface{}{}
}

func (r *AddSCMIntegrationV1Response) Errors() []error {
	return r.errors
}

func (r *AddSCMIntegrationV1Response) ValidationErrors() []validation.ValidationError {
	return r.validationErrors
}

func (r *AddSCMIntegrationV1Response) Status() result.Status {
	return r.status
}

func (r *AddSCMIntegrationV1Response) IsListData() bool {
	return false
}

func handleAddSCMIntegration(repo addSCMIntegrationRepository, v structValidator, t scmIntegrationTransformerV1, dto AddSCMIntegrationV1DTO) result.Result {
	validationErrors, err := v.ValidateStruct(dto, func(v *validator.Validate) *validator.Validate {
		err := v.RegisterValidation("unique", func(fl validator.FieldLevel) bool {
			m, err := repo.FindByName(fl.Field().String())

			if err != nil {
				panic(err)
			}

			return m == nil
		})

		if err != nil {
			panic(err)
		}

		return v
	})

	if err != nil {
		return &AddSCMIntegrationV1Response{
			errors: []error{
				err,
			},
			status: result.InternalError,
		}
	}

	if len(validationErrors) > 0 {
		return &AddSCMIntegrationV1Response{
			status:           result.Invalid,
			validationErrors: validationErrors,
		}
	}

	m := dto.ToModel()

	if err := repo.Create(m); err != nil {
		return &AddSCMIntegrationV1Response{
			errors: []error{
				err,
			},
			status: result.InternalError,
		}
	}

	return &AddSCMIntegrationV1Response{
		status:  result.Created,
		created: t.SCMIntegrationV1(m),
	}
}
