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

type deleteProjectV1Response struct {
	deleted *TransformedProject
	errors  []error
	status  commandresult.Status
}

func (r *deleteProjectV1Response) Data() interface{} {
	return r.deleted
}

func (r *deleteProjectV1Response) Meta() interface{} {

	return map[string]interface{}{}
}

func (r *deleteProjectV1Response) Errors() []error {
	return r.errors
}

func (r *deleteProjectV1Response) ValidationErrors() []validation.ValidationError {
	return []validation.ValidationError{}
}

func (r *deleteProjectV1Response) Status() commandresult.Status {
	return r.status
}

func (r *deleteProjectV1Response) IsListData() bool {
	return false
}

func Delete(repo deleteProjectRepository, t Transformer, id string) commandresult.Result {
	uuid, err := uuid.Parse(id)

	if err != nil {
		return &deleteProjectV1Response{
			status: commandresult.NotFound,
		}
	}

	existing, err := repo.Find(uuid)

	if err != nil {
		return &deleteProjectV1Response{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	if existing == nil {
		return &deleteProjectV1Response{
			status: commandresult.NotFound,
		}
	}

	err = repo.Delete(uuid)

	if err != nil {
		return &deleteProjectV1Response{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	return &deleteProjectV1Response{
		status:  commandresult.Okay,
		deleted: t.ProjectV1(existing),
	}
}
