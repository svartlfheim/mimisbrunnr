package rules

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
)

type uniqueValidationRepo interface {
	FindByName(string) (*models.Project, error)
}

type uniquePerIntegrationRepo interface {
	FindByPathAndIntegrationID(string, uuid.UUID) (*models.Project, error)
}

type existsValidationRepo interface {
	Find(uuid.UUID) (*models.SCMIntegration, error)
}

func Unique(repo uniqueValidationRepo, existingRecord *models.Project) func(v *validator.Validate) (*validator.Validate) {
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

func UniquePerIntegration(repo uniquePerIntegrationRepo, existingRecord *models.Project) func(v *validator.Validate) (*validator.Validate) {
	return func(v *validator.Validate) *validator.Validate {
		err := v.RegisterValidation("uniqueperotherfield", func(fl validator.FieldLevel) bool {
			var integrationID uuid.UUID
			integrationIDField, _, _, _ := fl.GetStructFieldOKAdvanced2(fl.Parent(), "IntegrationID")


			if existingRecord == nil && (integrationIDField.Kind() == reflect.Ptr && integrationIDField.IsNil()) {
				return false
				// panic("IntegrationID not set in uniqueperotherfield rule when existing is nil")
			}

			if integrationIDField.Kind() == reflect.Ptr {
				integrationIDField = reflect.Indirect(integrationIDField)
			}

			if integrationIDField.Kind() == reflect.String {
				var err error
				integrationID, err = uuid.Parse(integrationIDField.String())

				if err != nil {
					return false
				}
			} else {
				integrationID = existingRecord.GetSCMIntegration().GetID()
			}

			m, err := repo.FindByPathAndIntegrationID(fl.Field().String(), integrationID)

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

func Exists(repo existsValidationRepo) func(v *validator.Validate) (*validator.Validate) {
	return func(v *validator.Validate) *validator.Validate {
		err := v.RegisterValidation("exists", func(fl validator.FieldLevel) bool {
			integrationID, err := uuid.Parse(fl.Field().String())

			if err != nil {
				return false
			}

			m, err := repo.Find(integrationID)

			if err != nil {
				panic(err)
			}

			return m != nil
		})

		if err != nil {
			panic(err)
		}

		return v
	}
}