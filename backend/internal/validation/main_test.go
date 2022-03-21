package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	zerologmocks "github.com/svartlfheim/mimisbrunnr/test/mocks/zerolog"
)

func pointToBool(b bool) *bool {
	return &b
}

func pointToString(s string) *string {
	return &s
}

func pointToInt(i int) *int {
	return &i
}

type subject1 struct {
	Field                             *string `json:"my_field" validate:"required"`
	SomeOtherField                    *int    `json:"some_other_field" validate:"required,gt=5"`
	NonRequiredField                  *string `json:"non_required_field"`
	FieldWithNoJsonTag                *bool
	FieldWithNoJsonTagButWithValidate *string `validate:"required"`
}

func Test_Validator_validateStruct(t *testing.T) {
	expectDefaultLogs := []map[string]interface{}{
		{
			"field":   "FieldWithNoJsonTag",
			"level":   "warn",
			"message": "field in struct for validation does not have a json tag",
		},
		{
			"field":   "FieldWithNoJsonTagButWithValidate",
			"level":   "warn",
			"message": "field in struct for validation does not have a json tag",
		},
	}

	tests := []struct {
		name   string
		subj   subject1
		opts   []WithValidationExtension
		expect []ValidationError
		// See above, for logs included every time
		expectLogs []map[string]interface{}
		expectErr  error
	}{
		{
			name: "no validation problems",
			subj: subject1{
				Field:                             pointToString("somevalue"),
				SomeOtherField:                    pointToInt(8),
				NonRequiredField:                  pointToString("im here"),
				FieldWithNoJsonTag:                pointToBool(true),
				FieldWithNoJsonTagButWithValidate: pointToString("somevalue"),
			},
			opts:      []WithValidationExtension{},
			expect:    []ValidationError{},
			expectErr: nil,
		},

		{
			name: "field is required",
			subj: subject1{
				Field:                             nil,
				SomeOtherField:                    pointToInt(8),
				NonRequiredField:                  pointToString("im here"),
				FieldWithNoJsonTag:                pointToBool(true),
				FieldWithNoJsonTagButWithValidate: pointToString("somevalue"),
			},
			opts: []WithValidationExtension{},
			expect: []ValidationError{
				Error{
					path:                   "my_field",
					rule:                   "required",
					param:                  "",
					valueType:              "nil",
					extraMessageGenerators: map[Rule]MessageGenerator{},
					extraParameterParsers:  map[Rule]ParameterParser{},
				},
			},
			expectErr: nil,
		},

		{
			name: "multiple fields are required",
			subj: subject1{
				Field:                             nil,
				SomeOtherField:                    nil,
				NonRequiredField:                  pointToString("im here"),
				FieldWithNoJsonTag:                pointToBool(true),
				FieldWithNoJsonTagButWithValidate: pointToString("somevalue"),
			},
			opts: []WithValidationExtension{},
			expect: []ValidationError{
				Error{
					path:                   "my_field",
					rule:                   "required",
					param:                  "",
					valueType:              "nil",
					extraMessageGenerators: map[Rule]MessageGenerator{},
					extraParameterParsers:  map[Rule]ParameterParser{},
				},
				Error{
					path:                   "some_other_field",
					rule:                   "required",
					param:                  "",
					valueType:              "nil",
					extraMessageGenerators: map[Rule]MessageGenerator{},
					extraParameterParsers:  map[Rule]ParameterParser{},
				},
			},
			expectErr: nil,
		},

		{
			name: "field with zero value does not trigger required rule",
			subj: subject1{
				Field:                             pointToString(""),
				SomeOtherField:                    pointToInt(0),
				NonRequiredField:                  pointToString("im here"),
				FieldWithNoJsonTag:                pointToBool(true),
				FieldWithNoJsonTagButWithValidate: pointToString("somevalue"),
			},
			opts: []WithValidationExtension{},
			expect: []ValidationError{
				Error{
					path:                   "some_other_field",
					rule:                   "gt",
					param:                  "5",
					valueType:              "int",
					extraMessageGenerators: map[Rule]MessageGenerator{},
					extraParameterParsers:  map[Rule]ParameterParser{},
				},
			},
			expectErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			l := zerologmocks.NewLogger()
			v := NewValidator(l.Logger)
			res, err := v.ValidateStruct(test.subj, test.opts...)

			assert.Equal(tt, test.expectErr, err)
			assert.Equal(tt, test.expect, res)
			l.AssertLogs(tt, append(expectDefaultLogs, test.expectLogs...))
		})
	}
}

/*
This test demonstrates the behavious of zero values during validation.
When not using pointers in the struct, the required validation will error for zero values.
*/
func Test_Validator_validateStruct_without_pointers_with_empty_values(t *testing.T) {
	subj := struct {
		MyField     string `json:"myfield" validate:"required"`
		MyIntField  int    `json:"my_int_field" validate:"required"`
		MyBoolField bool   `json:"my_bool_field" validate:"required"`
	}{
		MyField:     "",
		MyIntField:  0,
		MyBoolField: false,
	}

	l := zerologmocks.NewLogger()
	v := NewValidator(l.Logger)
	res, err := v.ValidateStruct(subj)

	assert.Equal(t, nil, err)
	assert.Equal(t, []ValidationError{
		Error{
			path:                   "myfield",
			rule:                   "required",
			param:                  "",
			valueType:              "string",
			extraMessageGenerators: map[Rule]MessageGenerator{},
			extraParameterParsers:  map[Rule]ParameterParser{},
		},
		Error{
			path:                   "my_int_field",
			rule:                   "required",
			param:                  "",
			valueType:              "int",
			extraMessageGenerators: map[Rule]MessageGenerator{},
			extraParameterParsers:  map[Rule]ParameterParser{},
		},
		Error{
			path:                   "my_bool_field",
			rule:                   "required",
			param:                  "",
			valueType:              "bool",
			extraMessageGenerators: map[Rule]MessageGenerator{},
			extraParameterParsers:  map[Rule]ParameterParser{},
		},
	}, res)

}

func Test_Validator_validateStruct_pointer_to_struct(t *testing.T) {
	subj := struct {
		MyField     string `json:"myfield" validate:"required"`
		MyIntField  *int   `json:"my_int_field" validate:"required,gt=5"`
		MyBoolField bool   `json:"my_bool_field" validate:"required"`
	}{
		MyField:     "",
		MyIntField:  pointToInt(3),
		MyBoolField: false,
	}

	l := zerologmocks.NewLogger()
	v := NewValidator(l.Logger)
	res, err := v.ValidateStruct(&subj)

	assert.Equal(t, nil, err)
	assert.Equal(t, []ValidationError{
		Error{
			path:                   "myfield",
			rule:                   "required",
			param:                  "",
			valueType:              "string",
			extraMessageGenerators: map[Rule]MessageGenerator{},
			extraParameterParsers:  map[Rule]ParameterParser{},
		},
		Error{
			path:                   "my_int_field",
			rule:                   "gt",
			param:                  "5",
			valueType:              "int",
			extraMessageGenerators: map[Rule]MessageGenerator{},
			extraParameterParsers:  map[Rule]ParameterParser{},
		},
		Error{
			path:                   "my_bool_field",
			rule:                   "required",
			param:                  "",
			valueType:              "bool",
			extraMessageGenerators: map[Rule]MessageGenerator{},
			extraParameterParsers:  map[Rule]ParameterParser{},
		},
	}, res)

}

func Test_Validator_validateStruct_bad_value(t *testing.T) {
	somestring := "somestring"
	l := zerologmocks.NewLogger()
	v := NewValidator(l.Logger)
	res, err := v.ValidateStruct(somestring)

	assert.IsType(t, &validator.InvalidValidationError{}, err)
	assert.Equal(t, []ValidationError{}, res)

	// check pointer is managed correctly
	res, err = v.ValidateStruct(&somestring)

	assert.IsType(t, &validator.InvalidValidationError{}, err)
	assert.Equal(t, []ValidationError{}, res)

	l.AssertLogs(t, []map[string]interface{}{
		{
			"level":   "error",
			"type":    "string",
			"message": "unknown validation error",
		},
		{
			"level":   "error",
			"type":    "string",
			"message": "unknown validation error",
		},
	}, zerologmocks.IgnoreFieldFilter("error"))
}
