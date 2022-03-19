package scm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/svartlfheim/mimisbrunnr/internal/scm"
	"github.com/svartlfheim/mimisbrunnr/internal/validation"
	"github.com/svartlfheim/mimisbrunnr/test/expectations"
	zerologmocks "github.com/svartlfheim/mimisbrunnr/test/mocks/zerolog"
	"github.com/svartlfheim/mimisbrunnr/test/pointto"
)

func Test_validation_for_add_scm_integration_v1(t *testing.T) {
	tests := []struct {
		name      string
		in        scm.AddSCMIntegrationV1DTO
		expect    []expectations.ValidationError
		expectErr error
	}{
		{
			name: "all required fields not present",
			in:   scm.AddSCMIntegrationV1DTO{},
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
					Path:  "access_token.name",
					Rule:  "required",
					Param: map[string]string{},
				},
				{
					Path:  "access_token.token",
					Rule:  "required",
					Param: map[string]string{},
				},
			},
			expectErr: nil,
		},
		{
			name: "all required fields empty but invalid values",
			in: scm.AddSCMIntegrationV1DTO{
				Name:     pointto.String(""),
				Type:     pointto.String(""),
				Endpoint: pointto.String(""),
				AccessToken: scm.AddSCMIntegrationV1AccessToken{
					Name:  pointto.String(""),
					Token: pointto.String(""),
				},
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
					Path: "access_token.name",
					Rule: "gt",
					Param: map[string]string{
						"limit": "0",
					},
				},
				{
					Path: "access_token.token",
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
			in: scm.AddSCMIntegrationV1DTO{
				Name:     pointto.String(""),
				Type:     pointto.String("somevalue"),
				Endpoint: pointto.String("not empty"),
				AccessToken: scm.AddSCMIntegrationV1AccessToken{
					Name:  nil,
					Token: pointto.String(""),
				},
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
					Path:  "access_token.name",
					Rule:  "required",
					Param: map[string]string{},
				},
				{
					Path: "access_token.token",
					Rule: "gt",
					Param: map[string]string{
						"limit": "0",
					},
				},
			},
			expectErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			l := zerologmocks.NewLogger()
			v := validation.NewValidator(l.Logger)

			res, err := v.ValidateStruct(test.in)

			require.Len(tt, res, len(test.expect))

			for i, e := range res {
				test.expect[i].Assert(tt, e)
			}

			assert.Equal(tt, test.expectErr, err)
		})
	}
}
