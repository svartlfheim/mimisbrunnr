package v1

import (
	"net/http"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/svartlfheim/mimisbrunnr/internal/app/openapi/generic"
)

func addDeleteResponses(doc *openapi3.T) {
	successDesc := "Returns the deleted SCM Integration."
	doc.Components.Responses["delete_scm_integration_v1_ok"] = &openapi3.ResponseRef{
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
}

func addDeleteOperation(doc *openapi3.T) {
	op := openapi3.NewOperation()

	op.Description = "Deletes the SCM Integration with the specified id."
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
		Ref: "#/components/responses/delete_scm_integration_v1_ok",
	}
	op.Parameters = openapi3.Parameters{
		{
			Value: &openapi3.Parameter{
				Name: "id",
				In: "path",
				Description: "The id of the SCM Integration to delete.",
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
	doc.AddOperation("/v1/scm-integrations/{id}", "DELETE", op)
}