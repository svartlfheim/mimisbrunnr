package result

import "github.com/svartlfheim/mimisbrunnr/internal/validation"

type Status string

func (rs Status) Equals(other Status) bool {
	return rs == other
}

func (rs Status) ToHTTP() int {
	switch rs {
	case Created:
		return 201
	case InternalError:
		return 500
	case Invalid:
		return 422
	}

	// Not implemented as default??
	return 501
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
