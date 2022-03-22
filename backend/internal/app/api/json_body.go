package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

type ErrorHandlingJsonUnmarshaller struct{}

func (m *ErrorHandlingJsonUnmarshaller) findStructFieldByJsonName(t reflect.Type, name string) (reflect.StructField, error) {
	for i := 0; i < t.NumField(); i++ {
		sField := t.Field(i)
		if sField.Tag.Get("json") == name {
			return sField, nil
		}
	}

	return reflect.StructField{}, ErrStructFieldNotFoundForJsonFieldName{
		StructType: t,
		JSONField:  name,
	}
}

func (m *ErrorHandlingJsonUnmarshaller) buildUnmarshalTypeError(s interface{}, err *json.UnmarshalTypeError) error {
	rval := reflect.ValueOf(s)

	for rval.Kind() == reflect.Ptr {
		// maybe a bad idea...don't think you can have inifinitely recursive pointer?...
		rval = reflect.Indirect(rval)
	}

	t := rval.Type()

	sField, sFieldErr := m.findStructFieldByJsonName(t, err.Field)

	if sFieldErr != nil {
		return ErrInternalError{
			Message: fmt.Sprintf("error encountered building unmarshal error: %s", sFieldErr.Error()),
		}
	}

	requiredType := sField.Type

	for requiredType.Kind() == reflect.Ptr {
		requiredType = sField.Type.Elem()
	}

	return ErrBadRequestInputData{
		Message: fmt.Sprintf("could not parse JSON, field %s, got %s, expected %s", err.Field, err.Value, requiredType.Name()),
	}
}

func (m *ErrorHandlingJsonUnmarshaller) Unmarshal(r *http.Request, into interface{}) error {
	if r.ContentLength == 0 {
		return ErrEmptyRequestBodyNotAllowed{}
	}

	rval := reflect.ValueOf(into)

	if rval.Kind() != reflect.Ptr {
		return ErrInternalError{
			Message: "argument 'into' passed to parseJSONBody must be a pointer",
		}
	}

	err := json.NewDecoder(r.Body).Decode(into)

	if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
		return m.buildUnmarshalTypeError(into, typeErr)
	}

	if err != nil {
		return ErrBadRequestInputData{
			Message: fmt.Sprintf("invalid json: %s", err.Error()),
		}
	}

	return nil
}

func NewErrorHandlingJsonUnmarshaller() *ErrorHandlingJsonUnmarshaller {
	return &ErrorHandlingJsonUnmarshaller{}
}
