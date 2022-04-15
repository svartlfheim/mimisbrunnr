package v1

import (
	"time"

	"github.com/google/uuid"
	"github.com/svartlfheim/mimisbrunnr/internal/app/scm/rules"
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
		time.Now().UTC(),
		time.Now().UTC(),
	)
}

type addIntegrationResponse struct {
	created          *models.SCMIntegration
	errors           []error
	status           commandresult.Status
	validationErrors []validation.ValidationError
}

func (r *addIntegrationResponse) Data() interface{} {
	return r.created
}

func (r *addIntegrationResponse) Meta() interface{} {
	return nil
}

func (r *addIntegrationResponse) Errors() []error {
	return r.errors
}

func (r *addIntegrationResponse) ValidationErrors() []validation.ValidationError {
	return r.validationErrors
}

func (r *addIntegrationResponse) Status() commandresult.Status {
	return r.status
}

func Add(repo addIntegrationRepository, v StructValidator, dto AddIntegrationDTO) commandresult.Result {
	validationErrors, err := v.ValidateStruct(dto, rules.Unique(repo, nil))

	if err != nil {
		return &addIntegrationResponse{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	if len(validationErrors) > 0 {
		return &addIntegrationResponse{
			status:           commandresult.Invalid,
			validationErrors: validationErrors,
		}
	}

	m := dto.ToModel()

	if err := repo.Create(m); err != nil {
		return &addIntegrationResponse{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	return &addIntegrationResponse{
		status:  commandresult.Created,
		created: m,
	}
}
