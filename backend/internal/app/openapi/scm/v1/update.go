package v1

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	v1 "github.com/svartlfheim/mimisbrunnr/internal/app/scm/v1"
	"github.com/svartlfheim/mimisbrunnr/internal/app/web"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
)


func addUpdateRequestBody(doc *openapi3.T) {
	exampleName := "My private gitlab"
	exampleEndpoint := "private.gitlab.com"
	exampleType := string(models.GitlabType)
	exampleToken := "mysupersecrettoken"
	
	dto := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Title: "Request Body - Update SCM integration V1",
			Type: "object",
			Description: "The name must be unique across all SCM integrations.",
			Properties: openapi3.Schemas{
				"name": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Description: "**MUST** be unique across all SCM Integrations.",
						Type: "string",
						MinLength: 1,
					},
				},
				"type": &openapi3.SchemaRef{
					Ref: "#/components/schemas/scm_integration_type",
				},
				"endpoint": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Description: "The http endpoint to use for API requests.",
						Type: "string",
						MinLength: 1,
					},
				},
				"token": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Description: "The access token to use for authentication with the API.",
						Type: "string",
						MinLength: 1,
					},
				},
			},
			Example: v1.AddIntegrationDTO{
				Name: &exampleName,
				Type: &exampleType,
				Token: &exampleToken,
				Endpoint: &exampleEndpoint,
			},
		},
	}

	doc.Components.Schemas["update_scm_integration_v1_request_body"] = dto
}

func addUpdateResponses(doc *openapi3.T) {
	successDesc := "The updated SCM Integration."
	doc.Components.Responses["update_scm_integration_v1_ok"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &successDesc,
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Properties: openapi3.Schemas{
								"data": &openapi3.SchemaRef{
									Ref: "#/components/schemas/scm_integration_v1",
								},
							},
						},
					},
				},
			},
		},
	}

	typeOpts := []string{}
	for _, o := range models.AvailableSCMIntegrationTypes() {
		typeOpts = append(typeOpts, string(o))
	}

	errorDesc := "The data submitted in the request was invalid. The example shows all possible validation errors for this request. Each response may contain one or more of these."
	doc.Components.Responses["update_scm_integration_v1_invalid"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &errorDesc,
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
										Path:    "name",
										Message: "must contain more than 0 characters",
										Parameters: map[string]string{
											"limit": "0",
										},
										Rule: "gt",
									},
									{
										Path:       "name",
										Message:    "value must be unique across all records of this type",
										Parameters: map[string]string{},
										Rule:       "unique",
									},
									{
										Path:    "type",
										Message: "must contain more than 0 characters",
										Parameters: map[string]string{
											"limit": "0",
										},
										Rule: "gt",
									},
									{
										Path:       "type",
										Message:    "invalid choice, must be one of the options",
										Parameters: map[string]string{
											"options": strings.Join(typeOpts, ", "),
										},
										Rule:       "scmintegrationtype",
									},
									{
										Path:    "endpoint",
										Message: "must contain more than 0 characters",
										Parameters: map[string]string{
											"limit": "0",
										},
										Rule: "gt",
									},
									{
										Path:    "token",
										Message: "must contain more than 0 characters",
										Parameters: map[string]string{
											"limit": "0",
										},
										Rule: "gt",
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

func addUpdateOperation(doc *openapi3.T) {
	op := openapi3.NewOperation()

	op.Description = "Updates a new SCM Integration defined by the request body."
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
		Ref: "#/components/responses/update_scm_integration_v1_ok",
	}
	op.Responses[strconv.Itoa(http.StatusUnprocessableEntity)] = &openapi3.ResponseRef{
		Ref: "#/components/responses/update_scm_integration_v1_invalid",
	}
	op.RequestBody = &openapi3.RequestBodyRef{
		Value: &openapi3.RequestBody{
			Required: true,
			Content: openapi3.Content{
				"application/json": {
					Schema: &openapi3.SchemaRef{
						Ref: "#/components/schemas/update_scm_integration_v1_request_body",
					},
				},
			},
		},
	}
	doc.AddOperation("/v1/scm-integrations", "PATCH", op)
}