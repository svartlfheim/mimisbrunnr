package v1

import (
	"errors"
	"math"

	"github.com/svartlfheim/mimisbrunnr/internal/models"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/commandresult"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/validation"
	"github.com/svartlfheim/mimisbrunnr/pkg/response/meta"
)

const defaultListLimit int = 20
const defaultListPage int = 1

type listProjectsRepository interface {
	Paginate(int, int) ([]*models.Project, error)
	Count() (int, error)
}

type ListProjectsDTO struct {
	Page  *int `validate:"omitempty,gt=0" json:"page"`
	Limit *int `validate:"omitempty,gt=0,lte=100" json:"limit"`
}

type listProjectsV1Response struct {
	found            []*TransformedProject
	errors           []error
	status           commandresult.Status
	validationErrors []validation.ValidationError
	page             int
	limit            int
	total            int
}

func (r *listProjectsV1Response) Data() interface{} {
	return r.found
}

func (r *listProjectsV1Response) Meta() interface{} {
	if !r.status.Equals(commandresult.Okay) {
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

func (r *listProjectsV1Response) Errors() []error {
	return r.errors
}

func (r *listProjectsV1Response) ValidationErrors() []validation.ValidationError {
	// There will never be validation errors here
	// Only a bad request or not found
	return r.validationErrors
}

func (r *listProjectsV1Response) Status() commandresult.Status {
	return r.status
}

func (r *listProjectsV1Response) IsListData() bool {
	return true
}

func List(repo listProjectsRepository, v StructValidator, t Transformer, dto ListProjectsDTO) commandresult.Result {
	validationErrors, err := v.ValidateStruct(dto)

	if err != nil {
		return &listProjectsV1Response{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	if len(validationErrors) > 0 {
		return &listProjectsV1Response{
			status:           commandresult.Invalid,
			validationErrors: validationErrors,
		}
	}

	total, err := repo.Count()

	if err != nil {
		return &listProjectsV1Response{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
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

	if page != 1 && page > numPages {
		return &listProjectsV1Response{
			errors: []error{
				errors.New("page out of bounds"),
			},
			status: commandresult.BadInput,
		}
	}

	list, err := repo.Paginate(page, resultLimit)

	if err != nil {
		return &listProjectsV1Response{
			errors: []error{
				err,
			},
			status: commandresult.InternalError,
		}
	}

	return &listProjectsV1Response{
		status: commandresult.Okay,
		found:  t.ProjectListV1(list),
		page:   page,
		limit:  resultLimit,
		total:  total,
	}
}
