package openapi

import (
	"context"
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/svartlfheim/mimisbrunnr/internal/app/openapi/generic"
	"github.com/svartlfheim/mimisbrunnr/internal/app/openapi/projects"
	"github.com/svartlfheim/mimisbrunnr/internal/app/openapi/scm"
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

func buildDoc() (*openapi3.T, error) {
	doc := &openapi3.T{
		Info: &openapi3.Info{
			Title:   "Mimisbrunnr",
			Version: "v1",
		},
		Servers: openapi3.Servers{
			&openapi3.Server{
				URL:         "https://mimisbrunnr.local/api",
				Description: "Local development environment",
			},
		},
		Components: openapi3.Components{
			Parameters: openapi3.ParametersMap{},
			Responses:  openapi3.Responses{},
			Schemas:    openapi3.Schemas{},
		},
	}
	doc.OpenAPI = "3.0.0"

	if err := generic.Add(doc); err != nil {
		return doc, err
	}

	if err := scm.Add(doc); err != nil {
		return doc, err
	}

	if err := projects.Add(doc); err != nil {
		return doc, err
	}

	return doc, nil
}


type Generator struct {
	cached []byte
}

func NewGenerator() *Generator {
	return &Generator{
		cached: nil,
	}
}

func (g *Generator) Generate(pretty bool) ([]byte, error) {
	if g.cached != nil {
		return g.cached, nil
	}

	doc, err := buildDoc()

	if err != nil {
		return nil, err
	}

	var out []byte
	if pretty {
		out, err = json.MarshalIndent(doc, "", "  ")
	} else {
		out, err = json.Marshal(doc)
	}

	if err != nil {
		return nil, err
	}

	l := openapi3.NewLoader()

	/*
	  If we try and validate the doc as we created it in buildDoc the validation fails to handle refs correctly
	  Presumably there is some behaviour happening within the library that is required to load refs properly which only happens
	  upon loading the openapi spec.
	*/
	doc, err = l.LoadFromData(out)

	if err != nil {
		return nil, err
	}

	if err := doc.Validate(context.Background()); err != nil {
		return nil, err
	}

	g.cached = out
	return out, err
}
