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

type getIntegrationResponse struct {
	found  *models.SCMIntegration
	errors []error
	status commandresult.Status
}

func (r *getIntegrationResponse) Data() interface{} {
	return r.found
}

func (r *getIntegrationResponse) Meta() interface{} {
	return nil
}

func (r *getIntegrationResponse) Errors() []error {
	return r.errors
}

func (r *getIntegrationResponse) ValidationErrors() []validation.ValidationError {
	// There will never be validation errors here
	// Only a bad request or not found
	return []validation.ValidationError{}
}

func (r *getIntegrationResponse) Status() commandresult.Status {
	return r.status
}

func Get(repo getIntegrationRepository, id string) commandresult.Result {
	uuid, err := uuid.Parse(id)

	if err != nil {
		return &getIntegrationResponse{
			status: commandresult.NotFound,
		}
	}

	m, err := repo.Find(uuid)

	if err != nil {
		return &getIntegrationResponse{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	if m == nil {
		return &getIntegrationResponse{
			status: commandresult.NotFound,
		}
	}

	return &getIntegrationResponse{
		status: commandresult.Okay,
		found:  m,
	}
}
