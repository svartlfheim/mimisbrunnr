package generic

import "github.com/getkin/kin-openapi/openapi3"

func addMetaPaginationResponseSchema(doc *openapi3.T) {
	doc.Components.Schemas["pagination_response_meta"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Title: "Type - Pagination Meta",
			Properties: openapi3.Schemas{
				"page": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Description: "The current page of results",
						Type:        "integer",
					},
				},
				"per_page": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Description: "The maximum amount of results being displayed per page.",
						Type:        "integer",
					},
				},
				"total": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Description: "The total amount of results in the system.",
						Type:        "integer",
					},
				},
				"count": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Description: "The amount of results actually included in this page.",
						Type:        "integer",
					},
				},
			},
		},
	}
}