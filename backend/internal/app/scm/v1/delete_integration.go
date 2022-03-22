package v1

import (
	"github.com/google/uuid"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/commandresult"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/validation"
)

type deleteIntegrationRepository interface {
	Delete(uuid.UUID) error
	Find(uuid.UUID) (*models.SCMIntegration, error)
}

type deleteIntegrationV1Response struct {
	deleted *TransformedSCMIntegration
	errors  []error
	status  commandresult.Status
}

func (r *deleteIntegrationV1Response) Data() interface{} {
	return r.deleted
}

func (r *deleteIntegrationV1Response) Meta() interface{} {

	return map[string]interface{}{}
}

func (r *deleteIntegrationV1Response) Errors() []error {
	return r.errors
}

func (r *deleteIntegrationV1Response) ValidationErrors() []validation.ValidationError {
	return []validation.ValidationError{}
}

func (r *deleteIntegrationV1Response) Status() commandresult.Status {
	return r.status
}

func (r *deleteIntegrationV1Response) IsListData() bool {
	return false
}

func Delete(repo deleteIntegrationRepository, t Transformer, id string) commandresult.Result {
	uuid, err := uuid.Parse(id)

	if err != nil {
		return &deleteIntegrationV1Response{
			status: commandresult.NotFound,
		}
	}

	existing, err := repo.Find(uuid)

	if err != nil {
		return &deleteIntegrationV1Response{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	if existing == nil {
		return &deleteIntegrationV1Response{
			status: commandresult.NotFound,
		}
	}

	err = repo.Delete(uuid)

	if err != nil {
		return &deleteIntegrationV1Response{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	return &deleteIntegrationV1Response{
		status:  commandresult.Okay,
		deleted: t.IntegrationV1(existing),
	}
}
