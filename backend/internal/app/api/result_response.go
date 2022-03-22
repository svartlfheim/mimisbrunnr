package api

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/commandresult"
)

type fieldError struct {
	Path       string            `json:"path"`
	Message    string            `json:"message"`
	Parameters map[string]string `json:"params"`
	Rule       string            `json:"rule"`
}

func buildValidationErrors(r commandresult.Result) []fieldError {
	fieldErrors := []fieldError{}
	for _, err := range r.ValidationErrors() {
		fieldErrors = append(fieldErrors, fieldError{
			Path:       err.Path(),
			Rule:       err.Rule(),
			Message:    err.Message(),
			Parameters: err.Parameters(),
		})
	}

	return fieldErrors
}

type response struct {
	Data             interface{}  `json:"data"`
	Meta             interface{}  `json:"meta"`
	ValidationErrors []fieldError `json:"validation_errors"`
}

func serveResponseForResult(res commandresult.Result, l zerolog.Logger, w http.ResponseWriter) {
	resp := response{
		Data:             res.Data(),
		Meta:             res.Meta(),
		ValidationErrors: buildValidationErrors(res),
	}

	if res.Status().Equals(commandresult.InternalError) {
		respondWithInternalError(w, l, res.Errors()...)
		return
	} else if res.Status().Equals(commandresult.BadInput) {
		respondWithBadRequest(w, l, res.Errors()...)
		return
	}

	body, err := json.Marshal(resp)

	if err != nil {
		respondWithInternalError(w, l, err)
		return
	}

	w.WriteHeader(res.Status().ToHTTP())
	_, err = w.Write(body)

	if err != nil {
		panic(err)
	}
}
