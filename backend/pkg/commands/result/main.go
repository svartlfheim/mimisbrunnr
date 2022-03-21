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
	case NotFound:
		return http.StatusNotFound
	case BadInput:
		return http.StatusBadRequest
	case Okay:
		return http.StatusOK
	}

	// Not implemented as default??
	return http.StatusNotImplemented
}

const Created Status = "created"
const InternalError Status = "internal_error"
const Invalid Status = "invalid"
const NotFound Status = "not_found"
const Okay Status = "ok"
const BadInput Status = "bad_input"

type Result interface {
	Data() interface{}
	Meta() interface{}
	Errors() []error
	ValidationErrors() []validation.ValidationError
	Status() Status
	IsListData() bool
}
