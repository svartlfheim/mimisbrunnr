package v1

import (
	"github.com/google/uuid"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/commandresult"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/validation"
)

type getIntegrationRepository interface {
	Find(uuid.UUID) (*models.SCMIntegration, error)
}

type getIntegrationV1Response struct {
	found  *TransformedSCMIntegration
	errors []error
	status commandresult.Status
}

func (r *getIntegrationV1Response) Data() interface{} {
	return r.found
}

func (r *getIntegrationV1Response) Meta() interface{} {
	return map[string]interface{}{}
}

func (r *getIntegrationV1Response) Errors() []error {
	return r.errors
}

func (r *getIntegrationV1Response) ValidationErrors() []validation.ValidationError {
	// There will never be validation errors here
	// Only a bad request or not found
	return []validation.ValidationError{}
}

func (r *getIntegrationV1Response) Status() commandresult.Status {
	return r.status
}

func (r *getIntegrationV1Response) IsListData() bool {
	return false
}

func Get(repo getIntegrationRepository, t Transformer, id string) commandresult.Result {
	uuid, err := uuid.Parse(id)

	if err != nil {
		return &getIntegrationV1Response{
			status: commandresult.NotFound,
		}
	}

	m, err := repo.Find(uuid)

	if err != nil {
		return &getIntegrationV1Response{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	if m == nil {
		return &getIntegrationV1Response{
			status: commandresult.NotFound,
		}
	}

	return &getIntegrationV1Response{
		status: commandresult.Okay,
		found:  t.IntegrationV1(m),
	}
}
