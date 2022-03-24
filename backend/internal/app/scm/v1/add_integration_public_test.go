package v1_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/svartlfheim/mimisbrunnr/internal/app/scm"
	v1 "github.com/svartlfheim/mimisbrunnr/internal/app/scm/v1"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/validation"
	"github.com/svartlfheim/mimisbrunnr/test/expectations"
	zerologmocks "github.com/svartlfheim/mimisbrunnr/test/mocks/zerolog"
	"github.com/svartlfheim/mimisbrunnr/test/pointto"
)

func Test_validation_for_add_scm_integration_v1(t *testing.T) {
	tests := []struct {
		name      string
		in        v1.AddIntegrationDTO
		expect    []expectations.ValidationError
		expectErr error
	}{
		{
			name: "all required fields not present",
			in:   v1.AddIntegrationDTO{},
			expect: []expectations.ValidationError{
				{
					Path:  "name",
					Rule:  "required",
					Param: map[string]string{},
				},
				{
					Path:  "type",
					Rule:  "required",
					Param: map[string]string{},
				},
				{
					Path:  "endpoint",
					Rule:  "required",
					Param: map[string]string{},
				},
				{
					Path:  "token",
					Rule:  "required",
					Param: map[string]string{},
				},
			},
			expectErr: nil,
		},
		{
			name: "all required fields empty but invalid values",
			in: v1.AddIntegrationDTO{
				Name:     pointto.String(""),
				Type:     pointto.String(""),
				Endpoint: pointto.String(""),
				Token:    pointto.String(""),
			},
			expect: []expectations.ValidationError{
				{
					Path: "name",
					Rule: "gt",
					Param: map[string]string{
						"limit": "0",
					},
				},
				{
					Path: "type",
					Rule: "gt",
					Param: map[string]string{
						"limit": "0",
					},
				},
				{
					Path: "endpoint",
					Rule: "gt",
					Param: map[string]string{
						"limit": "0",
					},
				},
				{
					Path: "token",
					Rule: "gt",
					Param: map[string]string{
						"limit": "0",
					},
				},
			},
			expectErr: nil,
		},
		{
			name: "mixture of failures",
			in: v1.AddIntegrationDTO{
				Name:     pointto.String(""),
				Type:     pointto.String("somevalue"),
				Endpoint: pointto.String("not empty"),
				Token:    nil,
			},
			expect: []expectations.ValidationError{
				{
					Path: "name",
					Rule: "gt",
					Param: map[string]string{
						"limit": "0",
					},
				},
				{
					Path: "type",
					Rule: "scmintegrationtype",
					Param: map[string]string{
						"options": "github, gitlab",
					},
				},
				{
					Path:  "token",
					Rule:  "required",
					Param: map[string]string{},
				},
			},
			expectErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			l := zerologmocks.NewLogger()
			v := validation.NewValidator(l.Logger)
			scm.RegisterExtraValidations(v)

			res, err := v.ValidateStruct(test.in)

			require.Len(tt, res, len(test.expect))

			for i, e := range res {
				test.expect[i].Assert(tt, e)
			}

			assert.Equal(tt, test.expectErr, err)
		})
	}
}
