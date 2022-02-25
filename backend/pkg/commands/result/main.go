package result

import (
	"net/http"

	"github.com/svartlfheim/mimisbrunnr/internal/validation"
)

type Status string

func (rs Status) Equals(other Status) bool {
	return rs == other
}

func (rs Status) ToHTTP() int {
	switch rs {
	case Created:
		return http.StatusCreated
	case InternalError:
		return http.StatusInternalServerError
	case Invalid:
		return http.StatusUnprocessableEntity
	}

	// Not implemented as default??
	return http.StatusNotImplemented
}

const Created Status = "created"
const InternalError Status = "internal_error"
const Invalid Status = "invalid"

type Result interface {
	Data() interface{}
	Meta() interface{}
	Errors() []error
	ValidationErrors() []validation.ValidationError
	Status() Status
	IsListData() bool
}
