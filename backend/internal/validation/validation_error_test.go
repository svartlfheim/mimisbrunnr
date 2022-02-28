package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validationError_Getters_for_rule_and_path(t *testing.T) {
	ve := validationError{
		path: "some.path",
		rule: "somerulename",
	}

	assert.Equal(t, "some.path", ve.Path())
	assert.Equal(t, "somerulename", ve.Rule())
}

func Test_valdiationError_message_generation(t *testing.T) {
	buildErr := func(r string, t string, p string) validationError {
		return validationError{
			rule:      r,
			valueType: t,
			param:     p,
		}
	}

	tests := []struct {
		name   string
		in     validationError
		expect string
	}{
		{
			name:   "default message",
			in:     buildErr("unknownrule", "", ""),
			expect: "is invalid",
		},

		// Required
		{
			name:   "for required rule",
			in:     buildErr("required", "", ""),
			expect: "is required",
		},

		// Greater than
		{
			name:   "greater than - string (no param)",
			in:     buildErr("gt", "string", ""),
			expect: "must contain more characters",
		},
		{
			name:   "greater than - int (no param)",
			in:     buildErr("gt", "int", ""),
			expect: "must be a larger number",
		},
		{
			name:   "greater than - default (no param)",
			in:     buildErr("gt", "bool", ""),
			expect: "must be larger",
		},
		{
			name:   "greater than - string (with param)",
			in:     buildErr("gt", "string", "8"),
			expect: "must contain more than 8 characters",
		},
		{
			name:   "greater than - int (with param)",
			in:     buildErr("gt", "int", "4"),
			expect: "must be larger than 4",
		},
		{
			name:   "greater than - default (with param)",
			in:     buildErr("gt", "bool", "10"),
			expect: "must be larger",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			assert.Equal(tt, test.expect, test.in.Message())
		})
	}
}

func Test_valdiationError_parameter_parsing(t *testing.T) {
	buildErr := func(r string, p string) validationError {
		return validationError{
			rule:  r,
			param: p,
		}
	}

	tests := []struct {
		name   string
		in     validationError
		expect map[string]string
	}{
		{
			name:   "empty if no param is set",
			in:     buildErr("gt", ""),
			expect: map[string]string{},
		},

		{
			name: "required params are generic",
			in:   buildErr("required", "somevalue"),
			expect: map[string]string{
				"param": "somevalue",
			},
		},

		{
			name: "gt params added as limit",
			in:   buildErr("gt", "8"),
			expect: map[string]string{
				"limit": "8",
			},
		},
		{
			name:   "gt empty param returns empty map",
			in:     buildErr("gt", ""),
			expect: map[string]string{},
		},

		{
			name: "lt params added as limit",
			in:   buildErr("lt", "5"),
			expect: map[string]string{
				"limit": "5",
			},
		},
		{
			name:   "lt empty param returns empty map",
			in:     buildErr("lt", ""),
			expect: map[string]string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			assert.Equal(tt, test.expect, test.in.Parameters())
		})
	}
}
