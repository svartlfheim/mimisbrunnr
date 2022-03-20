package scm

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/validation"
	"github.com/svartlfheim/mimisbrunnr/pkg/commands/result"
)

type managerRepository interface {
	addSCMIntegrationRepository
}

type managerTransformer interface {
	scmIntegrationTransformerV1
}

type structValidator interface {
	ValidateStruct(s interface{}, opts ...validation.WithValidationExtension) ([]validation.ValidationError, error)
	RegisterCustomValidation(t string, cv validation.CustomValidation)
}

type Manager struct {
	logger    zerolog.Logger
	repo      managerRepository
	validator structValidator
	transformer managerTransformer
}

func (m *Manager) AddV1(dto AddSCMIntegrationV1DTO) result.Result {

	return handleAddSCMIntegration(m.repo, m.validator, m.transformer, dto)
}

func RegisterExtraValidations(v structValidator) {
	v.RegisterCustomValidation("scmintegrationtype", validation.CustomValidation{
		ValidatorFunc: func(fl validator.FieldLevel) bool {
			// It should fail type validation before here, so this should be safe...
			value := fl.Field().String()
	
			isValid := false
	
			for _, t := range(models.AvailableSCMIntegrationTypes()) {
				if string(t) == value {
					isValid = true
					break
				}
			}
	
			return isValid
		},
		MessageGenerator: func(validation.Error) string {
			return "invalid choice, must be one of the options"
		},
		ParameterParser: func(validation.Error) map[string]string {
			opts := []string{}

			for _, o := range(models.AvailableSCMIntegrationTypes()) {
				opts = append(opts, string(o))
			}

			return map[string]string{
				"options": strings.Join(opts, ", "),
			}
		},
	})
}

func NewManager(l zerolog.Logger, repo managerRepository, v structValidator, t managerTransformer) *Manager {
	RegisterExtraValidations(v)

	return &Manager{
		logger:    l,
		repo:      repo,
		validator: v,
		transformer: t,
	}
}
