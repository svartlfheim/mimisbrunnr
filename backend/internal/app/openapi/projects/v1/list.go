package v1

import (
	"net/http"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/svartlfheim/mimisbrunnr/internal/app/openapi/generic"
	"github.com/svartlfheim/mimisbrunnr/internal/app/web"
)

func addListResponses(doc *openapi3.T) {
	successDesc := "A list of all available projects within the bounds of the current page and limit selections."

	doc.Components.Responses["list_projects_v1_ok"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &successDesc,
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Properties: openapi3.Schemas{
								"data": &openapi3.SchemaRef{
									Ref: "#/components/schemas/project_v1_list",
								},
								"meta": &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Properties: map[string]*openapi3.SchemaRef{
											"pagination": {
												Ref: "#/components/schemas/pagination_response_meta",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	invalidDesc := "The data submitted in the request was invalid. The example shows all possible validation errors for this request. Each response may contain one or more of these."
	doc.Components.Responses["list_projects_v1_invalid"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &invalidDesc,
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Properties: openapi3.Schemas{
								"validation_errors": &openapi3.SchemaRef{
									Ref: "#/components/schemas/validation_error",
								},
							},
							Example: web.OkayResponse{
								Data: nil,
								Meta: nil,
								ValidationErrors: []web.FieldError{
									{
										Path:    "page",
										Message: "must be less than or equal to {max page}",
										Parameters: map[string]string{
											"limit": "{max page}",
										},
										Rule: "lte",
									},
									{
										Path:    "page",
										Message: "must be larger than 0",
										Parameters: map[string]string{
											"limit": "0",
										},
										Rule: "gt",
									},
									{
										Path:    "limit",
										Message: "must be larger than 0",
										Parameters: map[string]string{
											"limit": "0",
										},
										Rule: "gt",
									},
									{
										Path:    "limit",
										Message: "must be less than or equal to 100",
										Parameters: map[string]string{
											"limit": "100",
										},
										Rule: "lte",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func addListOperation(doc *openapi3.T) {
	op := openapi3.NewOperation()

	op.Description = `Returns a list of all projects respecting the given limit and page query values.

Each project will include the serialised representation of the SCM Integration it is associated with.

The 'meta' section contains information about pagination to help with building UIs.

The {max page} value is calculated at runtime, and therefore cannot be included here. The maximum page value will be calculated based on the amount of records in the system and the defined limit (ceil(total/limit)). If you choose a page beyond this value you will see the validation error shown in the example 422 response.`
	op.Responses = openapi3.Responses{}
	op.Responses[strconv.Itoa(http.StatusInternalServerError)] = &openapi3.ResponseRef{
		Ref: "#/components/responses/internal_server_error",
	}
	op.Responses[strconv.Itoa(http.StatusBadRequest)] = &openapi3.ResponseRef{
		Ref: "#/components/responses/bad_input",
	}
	op.Responses[strconv.Itoa(http.StatusUnsupportedMediaType)] = &openapi3.ResponseRef{
		Ref: "#/components/responses/unsupported_media_type",
	}
	op.Responses[strconv.Itoa(http.StatusOK)] = &openapi3.ResponseRef{
		Ref: "#/components/responses/list_projects_v1_ok",
	}
	op.Responses[strconv.Itoa(http.StatusUnprocessableEntity)] = &openapi3.ResponseRef{
		Ref: "#/components/responses/list_projects_v1_invalid",
	}
	op.Parameters = []*openapi3.ParameterRef{
		generic.BuildPageParamRef(float64(1)),
		generic.BuildLimitParamRef(float64(1), float64(100), 20),
	}
	doc.AddOperation("/v1/projects", "GET", op)
}

