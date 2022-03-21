package scm

import (
	"errors"
	"math"

	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/validation"
	"github.com/svartlfheim/mimisbrunnr/pkg/commands/result"
	"github.com/svartlfheim/mimisbrunnr/pkg/response/meta"
)

const defaultListLimit int = 20
const defaultListPage int = 1

type listSCMIntegrationsRepository interface {
	Paginate(int, int) ([]*models.SCMIntegration, error)
	Count() (int, error)
}

type ListSCMIntegrationsV1DTO struct {
	Page  *int `validate:"omitempty,gt=0" json:"page"`
	Limit *int `validate:"omitempty,gt=0,lte=100" json:"limit"`
}

type ListSCMIntegrationsV1Response struct {
	found            []*scmIntegrationV1
	errors           []error
	status           result.Status
	validationErrors []validation.ValidationError
	page             int
	limit            int
	total            int
}

func (r *ListSCMIntegrationsV1Response) Data() interface{} {
	return r.found
}

func (r *ListSCMIntegrationsV1Response) Meta() interface{} {
	if !r.status.Equals(result.Okay) {
		// Return no meta if the request was no okay
		// The values will likely be empty
		return map[string]interface{}{}
	}

	return map[string]interface{}{
		"pagination": meta.Pagination{
			Page:  r.page,
			Limit: r.limit,
			Count: len(r.found),
			Total: r.total,
		},
	}
}

func (r *ListSCMIntegrationsV1Response) Errors() []error {
	return r.errors
}

func (r *ListSCMIntegrationsV1Response) ValidationErrors() []validation.ValidationError {
	// There will never be validation errors here
	// Only a bad request or not found
	return r.validationErrors
}

func (r *ListSCMIntegrationsV1Response) Status() result.Status {
	return r.status
}

func (r *ListSCMIntegrationsV1Response) IsListData() bool {
	return true
}

func handleListSCMIntegrations(repo listSCMIntegrationsRepository, v structValidator, t scmIntegrationTransformerV1, dto ListSCMIntegrationsV1DTO) result.Result {
	validationErrors, err := v.ValidateStruct(dto)

	if err != nil {
		return &ListSCMIntegrationsV1Response{
			errors: []error{
				err,
			},
			status: result.InternalError,
		}
	}

	if len(validationErrors) > 0 {
		return &ListSCMIntegrationsV1Response{
			status:           result.Invalid,
			validationErrors: validationErrors,
		}
	}

	total, err := repo.Count()

	if err != nil {
		return &ListSCMIntegrationsV1Response{
			errors: []error{
				err,
			},
			status: result.InternalError,
		}
	}

	// Default to 20 per page
	resultLimit := defaultListLimit

	if dto.Limit != nil {
		resultLimit = *dto.Limit
	}

	page := defaultListPage

	if dto.Page != nil {
		page = *dto.Page
	}

	numPages := int(math.Ceil(float64(total) / float64(resultLimit)))

	if page > numPages {
		return &ListSCMIntegrationsV1Response{
			errors: []error{
				errors.New("page out of bounds"),
			},
			status: result.BadInput,
		}
	}

	list, err := repo.Paginate(page, resultLimit)

	if err != nil {
		return &ListSCMIntegrationsV1Response{
			errors: []error{
				err,
			},
			status: result.InternalError,
		}
	}

	return &ListSCMIntegrationsV1Response{
		status: result.Okay,
		found:  t.SCMIntegrationListV1(list),
		page:   page,
		limit:  resultLimit,
		total:  total,
	}
}
