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
	projectV1Ref, err := openapi3gen.NewSchemaRefForValue(web.ProjectV1{}, doc.Components.Schemas)

	if err != nil {
		return err
	}

	projectV1Ref.Value.Example = web.ProjectV1{
		ID:   uuid.NewString(),
		Name: "My project",
		Path: "myorg/myrepo",
		SCMIntegration: &web.SCMIntegrationV1{
			ID:        uuid.NewString(),
			Name:      "Public Github",
			Type:      string(models.GithubType),
			Endpoint:  "https://github.com",
			Token:     "mysupersecrettoken",
			CreatedAt: time.RFC3339,
			UpdatedAt: time.RFC3339,
		},
		CreatedAt: time.RFC3339,
		UpdatedAt: time.RFC3339,
	}

	projectV1Ref.Value.Properties["scm_integration"] = &openapi3.SchemaRef{
		Ref: "#/components/schemas/scm_integration_v1",
	}
	projectV1Ref.Value.Title = "Type - Project (V1)"
	doc.Components.Schemas["project_v1"] = projectV1Ref
	doc.Components.Schemas["project_v1_list"] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Title: "Type - Project (V1) List",
			Type:  "array",
			Items: &openapi3.SchemaRef{
				Ref: "#/components/schemas/project_v1",
			},
		},
	}

	return nil
}