package generic

import "github.com/getkin/kin-openapi/openapi3"

func addErrorBodySchema(doc *openapi3.T) {
	doc.Components.Schemas["error_body"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Title: "Response Body - Error",
			Properties: openapi3.Schemas{
				"message": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type:    "string",
						Example: "Whoopsie, something went wrong!",
					},
				},
				"errors": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "array",
						Example: []string{
							"possible extra information, that may help with debugging",
						},
						Items: &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "string",
							},
						},
					},
				},
			},
		},
	}
}

func addErrorResponses(doc *openapi3.T) {
  serverErrorDesc := "An internal error occurred, and there is nothing you can do."
	doc.Components.Responses["internal_server_error"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &serverErrorDesc,
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{
						Ref: "#/components/schemas/error_body",
					},
				},
			},
		},
	}

  badInputDesc := "The request could not be parsed correctly by the server."
	doc.Components.Responses["bad_input"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &badInputDesc,
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{
						Ref: "#/components/schemas/error_body",
					},
				},
			},
		},
	}

  unsupportedMediaTypeDesc := `You requested an unsupported content type.

Ensure your request has the correct ` + "`Content-Type` header, `application/json` is the most supported."
	doc.Components.Responses["unsupported_media_type"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &unsupportedMediaTypeDesc,
		},
	}

	notFoundDesc := "You requested a resource that does not exist."
	doc.Components.Responses["not_found"] = &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &notFoundDesc,
		},
	}
}

