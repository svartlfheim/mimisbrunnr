package v1

import (
	"net/http"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/svartlfheim/mimisbrunnr/internal/app/openapi/generic"
)

func addGetResponses(doc *openapi3.T) {
	successDesc := "Returns the project."
	doc.Components.Responses["get_project_v1_ok"] = &openapi3.ResponseRef{
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
}

func addGetOperation(doc *openapi3.T) {
	op := openapi3.NewOperation()

	op.Description = "Gets the project with the specified id."
	op.Responses = openapi3.Responses{}
	op.Responses[strconv.Itoa(http.StatusInternalServerError)] = &openapi3.ResponseRef{
		Ref: "#/components/responses/internal_server_error",
	}
	op.Responses[strconv.Itoa(http.StatusUnsupportedMediaType)] = &openapi3.ResponseRef{
		Ref: "#/components/responses/unsupported_media_type",
	}
	op.Responses[strconv.Itoa(http.StatusNotFound)] = &openapi3.ResponseRef{
		Ref: "#/components/responses/not_found",
	}
	op.Responses[strconv.Itoa(http.StatusOK)] = &openapi3.ResponseRef{
		Ref: "#/components/responses/get_project_v1_ok",
	}
	op.Parameters = openapi3.Parameters{
		{
			Value: &openapi3.Parameter{
				Name: "id",
				In: "path",
				Description: "The id of the project to get.",
				Required: true,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Description: "**MUST** exist in the system.",
						Type: "string",
						Pattern: generic.UUIDFormat(),
						Format: "uuid",
					},
				},
			},
		},
	}
	doc.AddOperation("/v1/projects/{id}", "GET", op)
}