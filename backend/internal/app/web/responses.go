package web

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/svartlfheim/mimisbrunnr/internal/pkg/commandresult"
)

type WritableResponse interface {
	StatusCode() int
	Body() ([]byte)
	Write(w http.ResponseWriter)
}

type responseBuilderTransformer interface {
	Transform(v int, item interface{}) (interface{}, error)
}

type OkayResponse struct {
	Status           int          `json:"-"`
	Data             interface{}  `json:"data"`
	Meta             interface{}  `json:"meta"`
	ValidationErrors []fieldError `json:"validation_errors"`
}

func (r OkayResponse) Body() ([]byte) {
	b, err := json.Marshal(r)

	if err != nil {
		panic(err)
	}
	return b
}

func (r OkayResponse) StatusCode() int {
	return r.Status
}

func (r OkayResponse) Write(w http.ResponseWriter) {
	writeResponse(w, r)
}

type ErrorResponse struct {
	Status  int     `json:"-"`
	Message string  `json:"message"`
	Errors  []string `json:"errors"`
}

func (r ErrorResponse) Body() ([]byte) {
	b, err := json.Marshal(r)

	if err != nil {
		panic(err)
	}
	return b
}

func (r ErrorResponse) StatusCode() int {
	return r.Status
}

func (r ErrorResponse) Write(w http.ResponseWriter) {
	writeResponse(w, r)
}

type EmptyResponse struct {
	Status  int `json:"-"`
}

func (r EmptyResponse) Body() ([]byte) {
	return []byte{}
}

func (r EmptyResponse) StatusCode() int {
	return r.Status
}

func (r EmptyResponse) Write(w http.ResponseWriter) {
	writeResponse(w, r)
}

type fieldError struct {
	Path       string            `json:"path"`
	Message    string            `json:"message"`
	Parameters map[string]string `json:"params"`
	Rule       string            `json:"rule"`
}

func writeResponse(w http.ResponseWriter, resp WritableResponse) {
	w.WriteHeader(resp.StatusCode())

	_, err := w.Write(resp.Body())

	if err != nil {
		panic(err)
	}
}

func buildValidationErrorsFromCommandResult(r commandresult.Result) []fieldError {
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

type ResponseBuilder struct {
	trans  responseBuilderTransformer
	logger zerolog.Logger
}

func (rb *ResponseBuilder) buildInternalErrorResponse(errs ...error) ErrorResponse {
	rb.logger.Error().Errs("errors", errs).Msg("internal error occurred")

	errsAsStrings := []string{}

	for _, err := range(errs) {
		errsAsStrings = append(errsAsStrings, err.Error())
	}
	return ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: "Oopsie, an unknown error occurred!",
		Errors:  errsAsStrings,
	}
}

func (rb *ResponseBuilder) buildBadRequestResponse(errs ...error) ErrorResponse {
	rb.logger.Warn().Errs("errors", errs).Msg("bad request data received")

	errsAsStrings := []string{}

	for _, err := range(errs) {
		errsAsStrings = append(errsAsStrings, err.Error())
	}

	return ErrorResponse{
		Status:  http.StatusBadRequest,
		Message: "Receieved invalid request input!",
		Errors:  errsAsStrings,
	}
}

func (rb *ResponseBuilder) FromCommandResult(apiVersion int, res commandresult.Result) WritableResponse {
	switch res.Status() {
	case commandresult.NotFound:
		return EmptyResponse{
			Status: res.Status().ToHTTP(),
		}
	case commandresult.InternalError:
		return rb.buildInternalErrorResponse(res.Errors()...)
	case commandresult.BadInput:
		return rb.buildBadRequestResponse(res.Errors()...)
	}

	transformedData, err := rb.trans.Transform(apiVersion, res.Data())

	if err != nil {
		rb.logger.Error().Err(err).Msg("api transformation error")
		return rb.buildInternalErrorResponse()
	}

	return OkayResponse{
		Status:           res.Status().ToHTTP(),
		Data:             transformedData,
		Meta:             res.Meta(),
		ValidationErrors: buildValidationErrorsFromCommandResult(res),
	}

}

func (rb *ResponseBuilder) FromUnmarshalError(err error) WritableResponse {
	switch err.(type) {
	case ErrBadRequestInputData, ErrEmptyRequestBodyNotAllowed:
		return rb.buildBadRequestResponse(err)
	default:
		return rb.buildInternalErrorResponse(err)
	}
}

func NewResponseBuilder(l zerolog.Logger, trans responseBuilderTransformer) *ResponseBuilder {
	return &ResponseBuilder{
		logger: l,
		trans: trans,
	}
}