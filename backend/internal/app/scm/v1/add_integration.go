package v1

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/commandresult"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/validation"
)


type addIntegrationValidationRepository interface {
	FindByName(string) (*models.SCMIntegration, error)
}

type addIntegrationRepository interface {
	addIntegrationValidationRepository

	Create(*models.SCMIntegration) error
}

type AddIntegrationDTO struct {
	Name     *string `json:"name" validate:"required,gt=0,unique"`
	Type     *string `json:"type" validate:"required,gt=0,scmintegrationtype"`
	Endpoint *string `json:"endpoint" validate:"required,gt=0"`
	Token    *string `json:"token" validate:"required,gt=0"`
}

func (dto AddIntegrationDTO) ToModel() *models.SCMIntegration {
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

func (dto AddIntegrationDTO) Validate(v StructValidator, repo addIntegrationValidationRepository) ([]validation.ValidationError, error) {
	return v.ValidateStruct(dto, func(v *validator.Validate) *validator.Validate {
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
}

type addIntegrationV1Response struct {
	created          *TransformedSCMIntegration
	errors           []error
	status           commandresult.Status
	validationErrors []validation.ValidationError
}

func (r *addIntegrationV1Response) Data() interface{} {
	return r.created
}

func (r *addIntegrationV1Response) Meta() interface{} {
	return map[string]interface{}{}
}

func (r *addIntegrationV1Response) Errors() []error {
	return r.errors
}

func (r *addIntegrationV1Response) ValidationErrors() []validation.ValidationError {
	return r.validationErrors
}

func (r *addIntegrationV1Response) Status() commandresult.Status {
	return r.status
}

func (r *addIntegrationV1Response) IsListData() bool {
	return false
}

func Add(repo addIntegrationRepository, v StructValidator, t Transformer, dto AddIntegrationDTO) commandresult.Result {
	validationErrors, err := dto.Validate(v, repo)

	if err != nil {
		return &addIntegrationV1Response{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	if len(validationErrors) > 0 {
		return &addIntegrationV1Response{
			status:           commandresult.Invalid,
			validationErrors: validationErrors,
		}
	}

	m := dto.ToModel()

	if err := repo.Create(m); err != nil {
		return &addIntegrationV1Response{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	return &addIntegrationV1Response{
		status:  commandresult.Created,
		created: t.IntegrationV1(m),
	}
}
