package scm

type ResponseStatus string

func (rs ResponseStatus) Equals(other ResponseStatus) bool {
	return rs == other
}

const ResponseStatusCreated ResponseStatus = "created"
const ResponseStatusInternalError ResponseStatus = "internal_error"
const ResponseStatusInvalid ResponseStatus = "invalid"
