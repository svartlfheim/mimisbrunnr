package rules

import (
	"github.com/go-playground/validator/v10"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
)

type uniqueValidationRepo interface {
	FindByName(string) (*models.SCMIntegration, error)
}

func Unique(repo uniqueValidationRepo, existingRecord *models.SCMIntegration) func(v *validator.Validate) *validator.Validate {
	return func(v *validator.Validate) *validator.Validate {
		err := v.RegisterValidation("unique", func(fl validator.FieldLevel) bool {
			m, err := repo.FindByName(fl.Field().String())

			if err != nil {
				panic(err)
			}

			isOkay := m == nil
			if existingRecord == nil {
				return isOkay
			}

			return isOkay || m.GetID().String() == existingRecord.GetID().String()
		})

		if err != nil {
			panic(err)
		}

		return v
	}
}
