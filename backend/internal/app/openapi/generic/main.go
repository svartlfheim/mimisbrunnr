package generic

import (
	"github.com/getkin/kin-openapi/openapi3"
)

func BuildLimitParamRef(min float64, max float64, def int) *openapi3.ParameterRef {
	return &openapi3.ParameterRef{
		Value: &openapi3.Parameter{
			Name:        "limit",
			In:          "query",
			Description: "The amount of results to show per page.",
			Required:    false,
			Schema: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Default: def,
					Type:    "integer",
					Min:     &min,
					Max:     &max,
				},
			},
		},
	}
}

func BuildPageParamRef(minPage float64) *openapi3.ParameterRef {
	return &openapi3.ParameterRef{
		Value: &openapi3.Parameter{
			Name:        "page",
			In:          "query",
			Description: "The page of results to show.",
			Required:    false,
			Schema: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Default: 1,
					Type:    "integer",
					Min:     &minPage,
				},
			},
		},
	}
}

// See the following file in the validation library
// https://github.com/go-playground/validator/blob/v10.0.0/regexes.go
func UUIDFormat() string {
	return "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
}


func Add(doc *openapi3.T) error {
  addErrorBodySchema(doc)
  addErrorResponses(doc)
  addValidationErrorSchema(doc)
  addMetaPaginationResponseSchema(doc)

	return nil
}
