package v1

import (
	"net/http"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"
	"github.com/svartlfheim/mimisbrunnr/internal/app/openapi/generic"
	v1 "github.com/svartlfheim/mimisbrunnr/internal/app/projects/v1"
	"github.com/svartlfheim/mimisbrunnr/internal/app/web"
)


func addCreateRequestBody(doc *openapi3.T) {
	exampleID := uuid.NewString()
	exampleName := "My awesome project"
	examplePath := "myorg/myrepo"
	

	dto := &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Title: "Request Body - Create Project V1",
			Type: "object",
			Description: `The name must be unique across all projects.

The path has a few quirks worth noting.

Firstly, it must be unique within the integration it is assigned to. An integration may have many projects, but no project within the integration may share the same path.

Secondly, it should represent the path in the SCM tool, that the repository project can be accessed on. For example if you're integration is for github, and the project repository is found at 'https://github.com/myorg/myrepo', the path here should be 'myorg/myrepo'.`,
			Required: []string{"name", "path", "scm_integration_id"},
			Properties: openapi3.Schemas{
				"scm_integration_id": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Description: "**MUST** already exist in the system.",
						Type: "string",
						Pattern: generic.UUIDFormat(),
						Format: "uuid",
					},
				},
				"name": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Description: "**MUST** be unique across all projects.",
						Type: "string",
						MinLength: 1,
					},
				},
				"path": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Description: "**MUST** be unique per the `scm_integration_id` given.",
						Type: "string",
						MinLength: 1,
					},
				},
			},
			Example: v1.AddProjectDTO{
				IntegrationID: &exampleID,
				Name: &exampleName,
				Path: &examplePath,
			},
		},
	}

	doc.Components.Schemas["add_project_v1_request_body"] = dto
}

func addCreateResponses(doc *openapi3.T) {
	successDesc := "The created project."
	doc.Components.Responses["create_project_v1_ok"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &successDesc,
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Properties: openapi3.Schemas{
								"data": &openapi3.SchemaRef{
									Ref: "#/components/schemas/project_v1",
								},
							},
						},
					},
				},
			},
		},
	}

	errorDesc := "The data submitted in the request was invalid. The example shows all possible validation errors for this request. Each response may contain one or more of these."
	doc.Components.Responses["create_project_v1_invalid"] = &openapi3.ResponseRef{
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
										Path:       "scm_integration_id",
										Message:    "must be a valid uuid",
										Parameters: map[string]string{},
										Rule:       "uuid",
									},
									{
										Path:       "scm_integration_id",
										Message:    "must already exist",
										Parameters: map[string]string{},
										Rule:       "exists",
									},
									{
										Path:       "scm_integration_id",
										Message:    "is required",
										Parameters: map[string]string{},
										Rule:       "required",
									},
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
										Path:       "name",
										Message:    "is required",
										Parameters: map[string]string{},
										Rule:       "required",
									},
									{
										Path:    "path",
										Message: "must contain more than 0 characters",
										Parameters: map[string]string{
											"limit": "0",
										},
										Rule: "gt",
									},
									{
										Path:    "path",
										Message: "value must be unique across all records of this type with the same value for: scm_integration_id",
										Parameters: map[string]string{
											"field": "scm_integration_id",
										},
										Rule: "uniqueperotherfield",
									},
									{
										Path:       "path",
										Message:    "is required",
										Parameters: map[string]string{},
										Rule:       "required",
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

func addCreateOperation(doc *openapi3.T) {
	op := openapi3.NewOperation()

	op.Description = `Creates a new project defined by the request body.

The chosen SCM integration ID must match an existing record, see: ` + "`GET v{apiVersion}/scm-integrations`."
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
		Ref: "#/components/responses/create_project_v1_ok",
	}
	op.Responses[strconv.Itoa(http.StatusUnprocessableEntity)] = &openapi3.ResponseRef{
		Ref: "#/components/responses/create_project_v1_invalid",
	}
	op.RequestBody = &openapi3.RequestBodyRef{
		Value: &openapi3.RequestBody{
			Required: true,
			Content: openapi3.Content{
				"application/json": {
					Schema: &openapi3.SchemaRef{
						Ref: "#/components/schemas/add_project_v1_request_body",
					},
				},
			},
		},
	}
	doc.AddOperation("/v1/projects", "POST", op)
}