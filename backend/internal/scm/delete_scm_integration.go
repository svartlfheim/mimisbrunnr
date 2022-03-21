package scm

import (
	"github.com/google/uuid"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/validation"
	"github.com/svartlfheim/mimisbrunnr/pkg/commands/result"
)

type deleteSCMIntegrationRepository interface {
	Delete(uuid.UUID) error
	Find(uuid.UUID) (*models.SCMIntegration, error)
}

type DeleteSCMIntegrationV1Response struct {
	deleted *scmIntegrationV1
	errors  []error
	status  result.Status
}

func (r *DeleteSCMIntegrationV1Response) Data() interface{} {
	return r.deleted
}

func (r *DeleteSCMIntegrationV1Response) Meta() interface{} {

	return map[string]interface{}{}
}

func (r *DeleteSCMIntegrationV1Response) Errors() []error {
	return r.errors
}

func (r *DeleteSCMIntegrationV1Response) ValidationErrors() []validation.ValidationError {
	return []validation.ValidationError{}
}

func (r *DeleteSCMIntegrationV1Response) Status() result.Status {
	return r.status
}

func (r *DeleteSCMIntegrationV1Response) IsListData() bool {
	return false
}

func handleDeleteSCMIntegration(repo deleteSCMIntegrationRepository, t scmIntegrationTransformerV1, id string) result.Result {
	uuid, err := uuid.Parse(id)

	if err != nil {
		return &DeleteSCMIntegrationV1Response{
			status: result.NotFound,
		}
	}

	existing, err := repo.Find(uuid)

	if err != nil {
		return &DeleteSCMIntegrationV1Response{
			errors: []error{
				err,
			},
			status: result.InternalError,
		}
	}

	if existing == nil {
		return &DeleteSCMIntegrationV1Response{
			status: result.NotFound,
		}
	}

	err = repo.Delete(uuid)

	if err != nil {
		return &DeleteSCMIntegrationV1Response{
			errors: []error{
				err,
			},
			status: result.InternalError,
		}
	}

	return &DeleteSCMIntegrationV1Response{
		status:  result.Okay,
		deleted: t.SCMIntegrationV1(existing),
	}
}
