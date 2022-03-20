package server

import (
	"net/http"

	"github.com/rs/zerolog"
)

func respondWithInternalError(w http.ResponseWriter, l zerolog.Logger, errs... error) {
	l.Error().Errs("errors", errs).Msg("internal error occurred")

	w.WriteHeader(500)
	message := "Oopsie, an unknown error occurred!"

	if _, err := w.Write([]byte(message)); err != nil {
		l.Fatal().Err(err).Msg("failed to write internal error response")
	}
}

func handleError(w http.ResponseWriter, l zerolog.Logger, err error) bool {
	if err == nil {
		return false
	}

	statusCode := 500
	message := "Oopsie, an unknown error occurred!"

	switch err.(type) {
	case ErrBadRequestInputData, ErrEmptyRequestBodyNotAllowed:
		statusCode = 400
		message = err.Error()
		l.Warn().Err(err).Int("status-code", statusCode).Msg("bad request data received")
	default:
		respondWithInternalError(w, l, err)
		return true
	}

	w.WriteHeader(statusCode)

	if _, err := w.Write([]byte(message)); err != nil {
		l.Fatal().Err(err).Msg("failed to write response")
	}

	return true
}
