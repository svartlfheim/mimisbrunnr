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

type deleteIntegrationResponse struct {
	deleted *models.SCMIntegration
	errors  []error
	status  commandresult.Status
}

func (r *deleteIntegrationResponse) Data() interface{} {
	return r.deleted
}

func (r *deleteIntegrationResponse) Meta() interface{} {

	return nil
}

func (r *deleteIntegrationResponse) Errors() []error {
	return r.errors
}

func (r *deleteIntegrationResponse) ValidationErrors() []validation.ValidationError {
	return []validation.ValidationError{}
}

func (r *deleteIntegrationResponse) Status() commandresult.Status {
	return r.status
}

func Delete(repo deleteIntegrationRepository, id string) commandresult.Result {
	uuid, err := uuid.Parse(id)

	if err != nil {
		return &deleteIntegrationResponse{
			status: commandresult.NotFound,
		}
	}

	existing, err := repo.Find(uuid)

	if err != nil {
		return &deleteIntegrationResponse{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	if existing == nil {
		return &deleteIntegrationResponse{
			status: commandresult.NotFound,
		}
	}

	err = repo.Delete(uuid)

	if err != nil {
		return &deleteIntegrationResponse{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	return &deleteIntegrationResponse{
		status:  commandresult.Okay,
		deleted: existing,
	}
}
