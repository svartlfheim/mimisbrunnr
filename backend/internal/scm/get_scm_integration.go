package scm

import (
	"github.com/google/uuid"
	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/validation"
	"github.com/svartlfheim/mimisbrunnr/pkg/commands/result"
)

type getSCMIntegrationRepository interface {
	Find(uuid.UUID) (*models.SCMIntegration, error)
}

type GetSCMIntegrationV1Response struct {
	found  *scmIntegrationV1
	errors []error
	status result.Status
}

func (r *GetSCMIntegrationV1Response) Data() interface{} {
	return r.found
}

func (r *GetSCMIntegrationV1Response) Meta() interface{} {
	return map[string]interface{}{}
}

func (r *GetSCMIntegrationV1Response) Errors() []error {
	return r.errors
}

func (r *GetSCMIntegrationV1Response) ValidationErrors() []validation.ValidationError {
	// There will never be validation errors here
	// Only a bad request or not found
	return []validation.ValidationError{}
}

func (r *GetSCMIntegrationV1Response) Status() result.Status {
	return r.status
}

func (r *GetSCMIntegrationV1Response) IsListData() bool {
	return false
}

func handleGetSCMIntegration(repo getSCMIntegrationRepository, t scmIntegrationTransformerV1, id string) result.Result {
	uuid, err := uuid.Parse(id)

	if err != nil {
		return &GetSCMIntegrationV1Response{
			status: result.NotFound,
		}
	}

	m, err := repo.Find(uuid)

	if err != nil {
		return &GetSCMIntegrationV1Response{
			errors: []error{
				err,
			},
			status: result.InternalError,
		}
	}

	if m == nil {
		return &GetSCMIntegrationV1Response{
			status: result.NotFound,
		}
	}

	return &GetSCMIntegrationV1Response{
		status: result.Okay,
		found:  t.SCMIntegrationV1(m),
	}
}
