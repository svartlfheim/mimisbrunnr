package server

import "github.com/svartlfheim/mimisbrunnr/pkg/commands/result"

type fieldError struct {
	Path       string                 `json:"path"`
	Message    string                 `json:"message"`
	Parameters map[string]interface{} `json:"params"`
	Rule       string                 `json:"rule"`
}

type invalidDataResponse struct {
	Errors []fieldError `json:"validation_errors"`
}

func buildInvalidResponseBody(r result.Result) invalidDataResponse {
	fieldErrors := []fieldError{}
	for _, err := range r.ValidationErrors() {
		fieldErrors = append(fieldErrors, fieldError{
			Path:       err.Path(),
			Rule:       err.Rule(),
			Message:    err.Message(),
			Parameters: err.Parameters(),
		})
	}

	return invalidDataResponse{
		Errors: fieldErrors,
	}
}
