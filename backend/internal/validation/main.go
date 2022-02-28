package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type WithValidationExtension func(*validator.Validate) *validator.Validate

type Validator struct {
	logger zerolog.Logger
}

func (v *Validator) make() *validator.Validate {
	validate := validator.New()

	// register function to get tag name from json tags.
	// See: https://github.com/go-playground/validator/blob/master/_examples/struct-level/main.go
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		jsonName := fld.Tag.Get("json")

		if jsonName == "" {
			v.logger.Warn().Str("field", fld.Name).Msg("field in struct for validation does not have a json tag")
		}

		return jsonName
	})

	return validate
}

func (v *Validator) transformErrors(s interface{}, errs validator.ValidationErrors) []ValidationError {
	rval := reflect.ValueOf(s)

	for rval.Kind() == reflect.Ptr {
		rval = reflect.Indirect(rval)
	}

	structName := rval.Type().Name()

	validationErrors := []ValidationError{}
	for _, fieldErr := range errs {

		var valueType string

		errval := reflect.ValueOf(fieldErr.Value())

		if errval.Kind() == reflect.Ptr && errval.IsNil() {
			valueType = "nil"
		} else {
			valueType = errval.Type().Name()
		}

		validationErrors = append(validationErrors, validationError{
			path:      strings.TrimPrefix(fieldErr.Namespace(), fmt.Sprintf("%s.", structName)),
			rule:      fieldErr.ActualTag(),
			param:     fieldErr.Param(),
			valueType: valueType,
		})
	}

	return validationErrors
}

func (v *Validator) ValidateStruct(s interface{}, opts ...WithValidationExtension) ([]ValidationError, error) {
	baseValidator := v.make()
	err := baseValidator.Struct(s)

	if err == nil {
		return []ValidationError{}, nil
	}

	// this check is only needed when your code could produce
	// an invalid value for validation such as interface with nil
	// value most including myself do not usually have code like this.
	if _, ok := err.(validator.ValidationErrors); !ok {
		rval := reflect.ValueOf(s)

		for rval.Kind() == reflect.Ptr {
			rval = reflect.Indirect(rval)
		}

		v.logger.Error().Str("type", rval.Type().Name()).Err(err).Msg("unknown validation error")
		return []ValidationError{}, err
	}

	// from here you can create your own error messages in whatever language you wish
	return v.transformErrors(s, err.(validator.ValidationErrors)), nil
}

func NewValidator(l zerolog.Logger) *Validator {
	return &Validator{
		logger: l,
	}
}
