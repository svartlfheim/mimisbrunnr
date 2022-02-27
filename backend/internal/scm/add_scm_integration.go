package scm

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/validation"
	"github.com/svartlfheim/mimisbrunnr/pkg/commands/result"
)

type addSCMIntegrationRepository interface {
	Create(*models.SCMIntegration, *models.SCMAccessToken) error
}

type AddSCMIntegrationV1AccessToken struct {
	Name  string `json:"name" validate:"required"`
	Token string `json:"token" validate:"required"`
}

type AddSCMIntegrationV1DTO struct {
	Name        string                         `json:"name" validate:"required"`
	Type        string                         `json:"type" validate:"required"`
	Endpoint    string                         `json:"endpoint" validate:"required,gt=10"`
	AccessToken AddSCMIntegrationV1AccessToken `json:"access_token" validate:"required"`
}

type AddSCMIntegrationV1Response struct {
	created          *models.SCMIntegration
	errors           []error
	status           result.Status
	validationErrors []validation.ValidationError
}

func (r *AddSCMIntegrationV1Response) Data() interface{} {
	return *r.created
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

func handleAddSCMIntegration(repo addSCMIntegrationRepository, v structValidator, dto AddSCMIntegrationV1DTO) result.Result {
	validationErrors, err := v.ValidateStruct(dto, func(v *validator.Validate) *validator.Validate {
		// add struct level validation
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

	return &AddSCMIntegrationV1Response{
		errors: []error{
			errors.New("not implemented"),
		},
		status: result.InternalError,
	}
}
