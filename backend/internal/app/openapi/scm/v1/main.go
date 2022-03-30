package v1

import "github.com/getkin/kin-openapi/openapi3"

func Add(doc *openapi3.T) error {
	if err := addModels(doc); err != nil {
		return err
	}

	addListResponses(doc)
	addListOperation(doc)

	addCreateRequestBody(doc)
	addCreateResponses(doc)
	addCreateOperation(doc)


	addUpdateRequestBody(doc)
	addUpdateResponses(doc)
	addUpdateOperation(doc)

	addDeleteResponses(doc)
	addDeleteOperation(doc)

	addGetResponses(doc)
	addGetOperation(doc)

	return nil
}