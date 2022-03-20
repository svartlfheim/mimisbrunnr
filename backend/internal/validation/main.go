package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type WithValidationExtension func(*validator.Validate) *validator.Validate

type CustomValidation struct {
	ValidatorFunc func(fl validator.FieldLevel) bool
	MessageGenerator MessageGenerator
	ParameterParser ParameterParser
}

type Validator struct {
	customValidations map[string]CustomValidation
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

func (v *Validator) customMessageGenerators() map[Rule]MessageGenerator {
	generators := map[Rule]MessageGenerator{}
	for k, cv := range(v.customValidations) {
		generators[Rule(k)] = cv.MessageGenerator
	}

	return generators
}


func (v *Validator) customParameterParsers() map[Rule]ParameterParser {
	parsers := map[Rule]ParameterParser{}
	for k, cv := range(v.customValidations) {
		parsers[Rule(k)] = cv.ParameterParser
	}

	return parsers
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

		validationErrors = append(validationErrors, Error{
			path:      strings.TrimPrefix(fieldErr.Namespace(), fmt.Sprintf("%s.", structName)),
			rule:      fieldErr.ActualTag(),
			param:     fieldErr.Param(),
			valueType: valueType,
			extraMessageGenerators: v.customMessageGenerators(),
			extraParameterParsers: v.customParameterParsers(),
		})
	}

	return validationErrors
}

func (v *Validator) ValidateStruct(s interface{}, opts ...WithValidationExtension) ([]ValidationError, error) {
	baseValidator := v.make()

	for tag, cv := range(v.customValidations) {
		baseValidator.RegisterValidation(tag, cv.ValidatorFunc)
	}

	for _, opt := range(opts) {
		baseValidator = opt(baseValidator)
	}

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

func (v *Validator) RegisterCustomValidation(t string, cv CustomValidation) {
	v.customValidations[t] = cv
}

func NewValidator(l zerolog.Logger) *Validator {
	return &Validator{
		logger: l,
		customValidations: map[string]CustomValidation{},
	}
}
