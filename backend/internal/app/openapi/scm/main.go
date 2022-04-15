package scm

import (
	"github.com/getkin/kin-openapi/openapi3"
	v1 "github.com/svartlfheim/mimisbrunnr/internal/app/openapi/scm/v1"
)

func Add(doc *openapi3.T) error {
	if err := v1.Add(doc); err != nil {
		return err
	}

	return nil
}