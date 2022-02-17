package result

type Status string

func (rs Status) Equals(other Status) bool {
	return rs == other
}

const Created Status = "created"
const InternalError Status = "internal_error"
const Invalid Status = "invalid"


type Result interface {
	Data() interface {}
	Meta() interface {}
	Errors() []error
	Status() Status
	IsListData() bool
}
