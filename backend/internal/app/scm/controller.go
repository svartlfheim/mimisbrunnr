package scm

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	v1 "github.com/svartlfheim/mimisbrunnr/internal/app/scm/v1"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/commandresult"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/validation"
)

type managerRepository interface {
	v1.RequiredRepository
}

type structValidator interface {
	v1.StructValidator
}

type Controller struct {
	logger    zerolog.Logger
	repo      managerRepository
	validator structValidator
}

func (m *Controller) AddV1(dto v1.AddIntegrationDTO) commandresult.Result {
	return v1.Add(m.repo, m.validator, dto)
}

func (m *Controller) GetV1(id string) commandresult.Result {
	return v1.Get(m.repo, id)
}

func (m *Controller) ListV1(dto v1.ListIntegrationsDTO) commandresult.Result {
	return v1.List(m.repo, m.validator, dto)
}

func (m *Controller) UpdateV1(id string, dto v1.UpdateIntegrationDTO) commandresult.Result {
	return v1.Update(m.repo, m.validator, id, dto)
}

func (m *Controller) DeleteV1(id string) commandresult.Result {
	return v1.Delete(m.repo, id)
}

func RegisterExtraValidations(v structValidator) {
	v.RegisterCustomValidation("scmintegrationtype", validation.CustomValidation{
		ValidatorFunc: func(fl validator.FieldLevel) bool {
			// It should fail type validation before here, so this should be safe...
			value := fl.Field().String()

			isValid := false

			for _, t := range models.AvailableSCMIntegrationTypes() {
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

			for _, o := range models.AvailableSCMIntegrationTypes() {
				opts = append(opts, string(o))
			}

			return map[string]string{
				"options": strings.Join(opts, ", "),
			}
		},
	})
}

func NewController(l zerolog.Logger, repo managerRepository, v structValidator) *Controller {
	RegisterExtraValidations(v)

	return &Controller{
		logger:    l,
		repo:      repo,
		validator: v,
	}
}
