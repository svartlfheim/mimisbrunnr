package server

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog"
)

func respondWithInternalError(w http.ResponseWriter, l zerolog.Logger, errs ...error) {
	l.Error().Errs("errors", errs).Msg("internal error occurred")

	w.WriteHeader(http.StatusInternalServerError)
	message := "Oopsie, an unknown error occurred!"

	if _, err := w.Write([]byte(message)); err != nil {
		l.Fatal().Err(err).Msg("failed to write internal error response")
	}
}

func respondWithBadRequest(w http.ResponseWriter, l zerolog.Logger, errs ...error) {
	l.Warn().Errs("errors", errs).Msg("bad request data received")

	w.WriteHeader(http.StatusBadRequest)
	errorMessages := []string{}

	for _, err := range errs {
		errorMessages = append(errorMessages, err.Error())
	}

	resp := struct {
		Errors []string `json:"errors"`
	}{
		Errors: errorMessages,
	}

	body, err := json.Marshal(resp)

	if err != nil {
		l.Fatal().Err(err).Msg("failed to marshal bad request error response")
	}

	if _, err := w.Write([]byte(body)); err != nil {
		l.Fatal().Err(err).Msg("failed to write bad request error response")
	}
}

func handleError(w http.ResponseWriter, l zerolog.Logger, err error) bool {
	if err == nil {
		return false
	}

	switch err.(type) {
	case ErrBadRequestInputData, ErrEmptyRequestBodyNotAllowed:
		respondWithBadRequest(w, l, err)
		return true
	default:
		respondWithInternalError(w, l, err)
		return true
	}
}
