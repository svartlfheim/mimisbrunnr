package v1

import (
	"github.com/google/uuid"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/commandresult"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/validation"
)

type getProjectRepository interface {
	Find(uuid.UUID) (*models.Project, error)
}

type getProjectResponse struct {
	found  *models.Project
	errors []error
	status commandresult.Status
}

func (r *getProjectResponse) Data() interface{} {
	return r.found
}

func (r *getProjectResponse) Meta() interface{} {
	return nil
}

func (r *getProjectResponse) Errors() []error {
	return r.errors
}

func (r *getProjectResponse) ValidationErrors() []validation.ValidationError {
	// There will never be validation errors here
	// Only a bad request or not found
	return []validation.ValidationError{}
}

func (r *getProjectResponse) Status() commandresult.Status {
	return r.status
}

func Get(repo getProjectRepository, id string) commandresult.Result {
	uuid, err := uuid.Parse(id)

	if err != nil {
		return &getProjectResponse{
			status: commandresult.NotFound,
		}
	}

	m, err := repo.Find(uuid)

	if err != nil {
		return &getProjectResponse{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	if m == nil {
		return &getProjectResponse{
			status: commandresult.NotFound,
		}
	}

	return &getProjectResponse{
		status: commandresult.Okay,
		found:  m,
	}
}
