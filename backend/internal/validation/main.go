package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type WithValidationExtension func(*validator.Validate) *validator.Validate

type ValidationError interface {
	Path() string
	Message() string
	Rule() string
	Parameters() map[string]interface{}
	ValueType() reflect.Type
}

type validationError struct {
	path      string
	rule      string
	param     string
	valueType reflect.Type
}

func (ve validationError) Path() string {
	return ve.path
}

func (ve validationError) Message() string {
	switch ve.Rule() {
	case "required":
		return "is required"
	case "gt":
		switch ve.ValueType().Name() {
		case "string":
			return fmt.Sprintf("must contain more than %s characters", ve.param)
		case "int":
			return fmt.Sprintf("must be larger than %s", ve.param)
		default:
			return "not large enough"
		}
	default:
		return "is invalid"
	}
}

func (ve validationError) Rule() string {
	return ve.rule
}

func (ve validationError) Parameters() map[string]interface{} {
	if ve.param == "" {
		return map[string]interface{}{}
	}

	switch ve.Rule() {
	case "gt", "lt":
		return map[string]interface{}{
			"limit": ve.param,
		}

	default:
		return map[string]interface{}{
			"param": ve.param,
		}
	}
}

func (ve validationError) ValueType() reflect.Type {
	return ve.valueType
}

type Validator struct {
	logger zerolog.Logger
}

func (v *Validator) make() *validator.Validate {
	validate := validator.New()

	// register function to get tag name from json tags.
	// See: https://github.com/go-playground/validator/blob/master/_examples/struct-level/main.go
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
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
		validationErrors = append(validationErrors, validationError{
			path:      strings.TrimPrefix(fieldErr.Namespace(), fmt.Sprintf("%s.", structName)),
			rule:      fieldErr.ActualTag(),
			param:     fieldErr.Param(),
			valueType: fieldErr.Type(),
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
