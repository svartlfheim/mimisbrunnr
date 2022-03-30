package generic

import "github.com/getkin/kin-openapi/openapi3"

func addValidationErrorSchema(doc *openapi3.T) {
	doc.Components.Schemas["validation_error_limit_params"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Title: "Type - (Validation Error Params) Numerically limited",
			Description: `Given for the following validation rules:
- lt
- lte
- gt
- gte`,
			Type: "object",
			Properties: openapi3.Schemas{
				"limit": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "integer",
					},
				},
			},
		},
	}
	doc.Components.Schemas["validation_error_uniqueperotherfield_params"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Title: "Type - (Validation Error Params) Linked to other field",
			Description: `Given for the following validation rules: 
- uniqueperotherfield`,
			Type: "object",
			Properties: openapi3.Schemas{
				"field": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: "string",
					},
				},
			},
		},
	}

	doc.Components.Schemas["validation_error_empty_params"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Title: "Type - (Validation Error Params) Empty",
			Description: `Given for the following validation rules:
- unique
- required
- exists
- uuid`,
			Type:       "object",
			Properties: openapi3.Schemas{},
		},
	}
	doc.Components.Schemas["validation_error"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Title: "Type - Validation Error",
			Type:  "array",
			Items: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					Type: "object",
					Properties: openapi3.Schemas{
						"path": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type:        "string",
								Description: "The name of the field for which this rule failed.",
							},
						},
						"message": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type:        "string",
								Description: "A human-readable message that explains the problem.",
							},
						},
						"params": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								OneOf: openapi3.SchemaRefs{
									{
										Ref: "#/components/schemas/validation_error_limit_params",
									},
									{
										Ref: "#/components/schemas/validation_error_uniqueperotherfield_params",
									},
									{
										Ref: "#/components/schemas/validation_error_empty_params",
									},
								},
								Discriminator: &openapi3.Discriminator{
									PropertyName: "rule",
									Mapping: map[string]string{
										"lt":                  "#/components/schemas/validation_error_limit_params",
										"lte":                 "#/components/schemas/validation_error_limit_params",
										"gt":                  "#/components/schemas/validation_error_limit_params",
										"gte":                 "#/components/schemas/validation_error_limit_params",
										"uniqueperotherfield": "#/components/schemas/validation_error_uniqueperotherfield_params",
										"required":            "#/components/schemas/validation_error_empty_params",
										"unique":              "#/components/schemas/validation_error_empty_params",
										"exists":              "#/components/schemas/validation_error_empty_params",
										"uuid":                "#/components/schemas/validation_error_empty_params",
									},
								},
							},
						},
						"rule": &openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "string",
								Description: `The name of the rule. Is machine-readable, so can be used to provide custom messages if the message does not suffice.

**required**: The field must be supplied in the request.

**lt**: For strings, the value must be shorter than the limit. For ints, the value must be smaller than the limit. 'limit' will be shown in the params field.

**lte**: For strings, the value must be shorter than or of equal length to the limit. For ints, the value must be smaller than or equal to the limit. 'limit' will be shown in the params field.

**gt**: For strings, the value must be longer than the limit. For ints, the value must be larger than the limit. 'limit' will be shown in the params field.

**gte**: For strings, the value must be longer than or of equal length to the limit. For ints, the value must be larger than or equal to the limit. 'limit' will be shown in the params field.

**uuid**: The value must be a valid uuid.

**exists**: The record being referenced must exist, typically used on a relationship definition.

**unique**: This value for this field must be unique across all of the records of this type.

**uniqueperotherfield**: This value for this field must be unique across all of the records of this type that share the same value for the field shown in the error message. 'field' will be shown in the params field, and shows the name of the field which this value must be unique across.
`,
								Enum: []interface{}{
									"required",
									"lt",
									"lte",
									"gt",
									"gte",
									"uuid",
									"exists",
									"unique",
									"uniqueperotherfield",
								},
							},
						},
					},
				},
			},
		},
	}
}