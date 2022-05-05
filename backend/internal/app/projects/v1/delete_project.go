package v1

import (
	"github.com/google/uuid"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/commandresult"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/validation"
)

type deleteProjectRepository interface {
	Delete(uuid.UUID) error
	Find(uuid.UUID) (*models.Project, error)
}

type deleteProjectResponse struct {
	deleted *models.Project
	errors  []error
	status  commandresult.Status
}

func (r *deleteProjectResponse) Data() interface{} {
	return r.deleted
}

func (r *deleteProjectResponse) Meta() interface{} {

	return map[string]interface{}{}
}

func (r *deleteProjectResponse) Errors() []error {
	return r.errors
}

func (r *deleteProjectResponse) ValidationErrors() []validation.ValidationError {
	return []validation.ValidationError{}
}

func (r *deleteProjectResponse) Status() commandresult.Status {
	return r.status
}

func Delete(repo deleteProjectRepository, id string) commandresult.Result {
	uuid, err := uuid.Parse(id)

	if err != nil {
		return &deleteProjectResponse{
			status: commandresult.NotFound,
		}
	}

	existing, err := repo.Find(uuid)

	if err != nil {
		return &deleteProjectResponse{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	if existing == nil {
		return &deleteProjectResponse{
			status: commandresult.NotFound,
		}
	}

	err = repo.Delete(uuid)

	if err != nil {
		return &deleteProjectResponse{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	return &deleteProjectResponse{
		status:  commandresult.Okay,
		deleted: existing,
	}
}
