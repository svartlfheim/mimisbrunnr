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

type getProjectV1Response struct {
	found  *TransformedProject
	errors []error
	status commandresult.Status
}

func (r *getProjectV1Response) Data() interface{} {
	return r.found
}

func (r *getProjectV1Response) Meta() interface{} {
	return map[string]interface{}{}
}

func (r *getProjectV1Response) Errors() []error {
	return r.errors
}

func (r *getProjectV1Response) ValidationErrors() []validation.ValidationError {
	// There will never be validation errors here
	// Only a bad request or not found
	return []validation.ValidationError{}
}

func (r *getProjectV1Response) Status() commandresult.Status {
	return r.status
}

func (r *getProjectV1Response) IsListData() bool {
	return false
}

func Get(repo getProjectRepository, t Transformer, id string) commandresult.Result {
	uuid, err := uuid.Parse(id)

	if err != nil {
		return &getProjectV1Response{
			status: commandresult.NotFound,
		}
	}

	m, err := repo.Find(uuid)

	if err != nil {
		return &getProjectV1Response{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	if m == nil {
		return &getProjectV1Response{
			status: commandresult.NotFound,
		}
	}

	return &getProjectV1Response{
		status: commandresult.Okay,
		found:  t.ProjectV1(m),
	}
}
