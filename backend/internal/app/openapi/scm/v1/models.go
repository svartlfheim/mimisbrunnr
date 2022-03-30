package v1

import (
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"github.com/google/uuid"
	"github.com/svartlfheim/mimisbrunnr/internal/app/web"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
)

func addModels(doc *openapi3.T) error {
	doc.Components.Schemas["scm_integration_type"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Title: "Type - SCM Integration Type (enum)",
			Type:  "string",
			Enum: []interface{}{
				string(models.GithubType),
				string(models.GitlabType),
			},
		},
	}
	scmIntegrationV1Ref, err := openapi3gen.NewSchemaRefForValue(web.SCMIntegrationV1{}, doc.Components.Schemas)

	if err != nil {
		return err
	}

	scmIntegrationV1Ref.Value.Example = web.SCMIntegrationV1{
		ID:        uuid.NewString(),
		Name:      "Public Github",
		Type:      string(models.GithubType),
		Endpoint:  "https://github.com",
		Token:     "mysupersecrettoken",
		CreatedAt: time.RFC3339,
		UpdatedAt: time.RFC3339,
	}
	scmIntegrationV1Ref.Value.Properties["type"] = &openapi3.SchemaRef{
		Ref: "#/components/schemas/scm_integration_type",
	}
	scmIntegrationV1Ref.Value.Title = "Type - SCM Integration (V1)"
	doc.Components.Schemas["scm_integration_v1"] = scmIntegrationV1Ref
	doc.Components.Schemas["scm_integration_v1_list"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Title: "Type - SCM Integration List (V1)",
			Type:  "array",
			Items: &openapi3.SchemaRef{
				Ref: "#/components/schemas/scm_integration_v1",
			},
		},
	}

	return nil
}

