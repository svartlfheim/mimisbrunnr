package openapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"github.com/google/uuid"
	"github.com/spf13/afero"
	"github.com/svartlfheim/mimisbrunnr/internal/app/web"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
)

/*
Goal
openapi: 3.0.0

info:
  version: "v1"
  title: "Mimisbrunnr"
  description: Manages resources that can be used to display documentation from markdown files within SCM repsoitories (e.g. github, gitlab).

servers:
- url: https://mimisbrunnr.local/api
  description: Local development server


components:
  schemas:
    InternalError:
      type: object
      properties:
        errors:
          type: array
          items:
            type: string
        message:
          type: string

paths:
  /projects:
    get:
      description: "List all projects"
      responses:
        '500':
          description: "Something went wrong, there is nothing you can do."
          content:
            application/json:
              $ref
*/

type appConfig interface {
  GetOpenAPIGeneratePath() string
}

func pointToFloat64(i float64) *float64 {
  return &i
}

func pointToString(s string) *string {
  return &s
}

func buildDoc() (*openapi3.T, error) {
  doc := &openapi3.T{
    Info: &openapi3.Info{
      Title: "Mimisbrunnr",
      Version: "v1",
    },
    Servers: openapi3.Servers{
      &openapi3.Server{
        URL: "https://mimisbrunnr.local/api",
        Description:"Local development environment",
      },
    },
    Components: openapi3.Components{
      Parameters: openapi3.ParametersMap{},
      Responses: openapi3.Responses{},
      Schemas: openapi3.Schemas{},
    },
  }
  doc.OpenAPI = "3.0.0"

  doc.Components.Schemas["error_body"] = &openapi3.SchemaRef{
    Value: &openapi3.Schema{
      Properties: openapi3.Schemas{
        "message": &openapi3.SchemaRef{
          Value: &openapi3.Schema{
            Type: "string",
            Example: "Whoopsie, something went wrong!",
          },
        },
        "errors": &openapi3.SchemaRef{
          Value: &openapi3.Schema{
            Type: "array",
            Example: []string{
              "page out of bounds",
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

  scmIntegrationV1Ref, err := openapi3gen.NewSchemaRefForValue(web.SCMIntegrationV1{}, doc.Components.Schemas)

  if err != nil { 
    return nil, err
  }

  doc.Components.Schemas["scm_integration_v1"] = scmIntegrationV1Ref
  doc.Components.Schemas["scm_integration_v1_list"] = &openapi3.SchemaRef{
    Value: &openapi3.Schema{
      Type: "array",
      Items: &openapi3.SchemaRef{
        Ref: "#/components/schemas/scm_integration_v1",
      },
    },
  }

  projectV1Ref, err := openapi3gen.NewSchemaRefForValue(web.ProjectV1{}, doc.Components.Schemas)
  projectV1Ref.Value.Example = web.ProjectV1{
    ID: uuid.NewString(),
    Name: "My project",
    Path: "myorg/myrepo",
    SCMIntegration: &web.SCMIntegrationV1{
      ID: uuid.NewString(),
      Name: "My github",
      Type: string(models.GithubType),
      Endpoint: "https://github.com",
      Token: "mysupersecrettoken",
    },
  }
  

  if err != nil { 
    return nil, err
  }

  doc.Components.Schemas["project_v1"] = projectV1Ref
  doc.Components.Schemas["project_v1_list"] = &openapi3.SchemaRef{
    Value: &openapi3.Schema{
      Type: "array",
      Items: &openapi3.SchemaRef{
        Ref: "#/components/schemas/project_v1",
      },
    },
  }

  doc.Components.Responses["internal_server_error"] = &openapi3.ResponseRef{
    Value: &openapi3.Response{
      Description: pointToString("An internal error occurred, and there is nothing you can do."),
      Content: openapi3.Content{
        "application/json": &openapi3.MediaType{
          Schema: &openapi3.SchemaRef{
            Ref: "#/components/schemas/error_body",
          },
        },
      },
    },
  }
  
  doc.Components.Responses["bad_input"] = &openapi3.ResponseRef{
    Value: &openapi3.Response{
      Description: pointToString("The request could not be parsed correctly by the server."),
      Content: openapi3.Content{
        "application/json": &openapi3.MediaType{
          Schema: &openapi3.SchemaRef{
            Ref: "#/components/schemas/error_body",
          },
        },
      },
    },
  }

  doc.Components.Responses["list_projects_v1"] = &openapi3.ResponseRef{
    Value: &openapi3.Response{
      Description: pointToString("A list of all available projects within the bounds of the current page and limit selections."),
      Content: openapi3.Content{
        "application/json": &openapi3.MediaType{
          Schema: &openapi3.SchemaRef{
            Value: &openapi3.Schema{
              Properties: openapi3.Schemas{
                "data": &openapi3.SchemaRef{
                  Ref: "#/components/schemas/project_v1_list",
                },
              },
            },
          },
        },
      },
    },
  }

  doc.Components.Parameters["limit"] = &openapi3.ParameterRef{
    Value: &openapi3.Parameter{
      Name: "limit",
      In: "query",
      Description: "The amount of results to show per page.",
      Required: false,
      Schema: &openapi3.SchemaRef{
        Value: &openapi3.Schema{
          Default: 20,
          Type: "integer",
          Min: pointToFloat64(float64(1)),
          Max: pointToFloat64(float64(100)),
        },
      },
    },
  }

  doc.Components.Parameters["page"] = &openapi3.ParameterRef{
    Value: &openapi3.Parameter{
      Name: "page",
      In: "query",
      Description: "The page of results to show.",
      Required: false,
      Schema: &openapi3.SchemaRef{
        Value: &openapi3.Schema{
          Default: 1,
          Type: "integer",
          Min: pointToFloat64(float64(1)),
        },
      },
    },
  }

  getProjectsOpt := openapi3.NewOperation()
  getProjectsOpt.Description = "List all projects"
  getProjectsOpt.Responses = openapi3.NewResponses()
  getProjectsOpt.Responses[strconv.Itoa(http.StatusInternalServerError)] = &openapi3.ResponseRef{
    Ref:  "#/components/responses/internal_server_error",
  }
  getProjectsOpt.Responses[strconv.Itoa(http.StatusBadRequest)] = &openapi3.ResponseRef{
    Ref:  "#/components/responses/bad_input",
  }
  getProjectsOpt.Responses[strconv.Itoa(http.StatusOK)] = &openapi3.ResponseRef{
    Ref:  "#/components/responses/list_projects_v1",
  }
  getProjectsOpt.Parameters = []*openapi3.ParameterRef{
    {
      Ref: "#/components/parameters/limit",
    },
    {
      Ref: "#/components/parameters/page",
    },
  }
  doc.AddOperation("/v1/projects", "GET", getProjectsOpt)

  return doc, nil
}

func Generate(c appConfig, fs afero.Fs) error {
  doc, err := buildDoc()

  if err != nil {
    return err
  }

  out, err := json.MarshalIndent(doc, "", "  ")

  if err != nil {
     return err
  }

  l := openapi3.NewLoader()

  /*
  If we try and validate the doc as we created it in buildDoc the validation fails to handle refs correctly
  Presumably there is some behaviour happening within the library that is required to load refs properly which only happens
  upon loading the openapi spec.
  */
  doc, err = l.LoadFromData(out)

  if err != nil {
    return err
 }

  if err := doc.Validate(context.Background()); err != nil {
    return err
  }

  
  exists, err := afero.Exists(fs, c.GetOpenAPIGeneratePath());
  if err != nil {
    return err
  }

  if exists {
    if err := fs.Remove(c.GetOpenAPIGeneratePath()); err != nil {
      return err
    }
  } else {
    if err := fs.MkdirAll(filepath.Dir(c.GetOpenAPIGeneratePath()), os.ModeDir); err != nil {
      return err
    }
  }

  f, err := fs.Create(c.GetOpenAPIGeneratePath())

  if err != nil {
    return err
  }

  if _, err := f.Write(out); err != nil {
    return err
  }

  fmt.Printf("Generated spec at: %s\n", c.GetOpenAPIGeneratePath())

  return nil
}