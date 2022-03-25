package api

import (
	"fmt"
	"reflect"
)

type ErrBadRequestInputData struct {
	Message string
}

func (e ErrBadRequestInputData) Error() string {
	return e.Message
}

type ErrInternalError struct {
	Message string
}

func (e ErrInternalError) Error() string {
	return e.Message
}

type ErrStructFieldNotFoundForJsonFieldName struct {
	JSONField  string
	StructType reflect.Type
}

func (e ErrStructFieldNotFoundForJsonFieldName) Error() string {
	return fmt.Sprintf("no field in struct %s has tag 'json:\"%s\"'", e.StructType.Name(), e.JSONField)
}

type ErrEmptyRequestBodyNotAllowed struct{}

func (e ErrEmptyRequestBodyNotAllowed) Error() string {
	return "request body cannot be empty"
}

type ErrUnsupportedApiVersion struct {
	Version int
}

func (e ErrUnsupportedApiVersion) Error() string {
	return fmt.Sprintf("version %d is not supported", e.Version)
}

type ErrUnsupportedResourceType struct {
	Val interface{}
}

func (e ErrUnsupportedResourceType) Error() string {
	rtype := reflect.TypeOf(e.Val)
	return fmt.Sprintf("type %s is not supported", rtype)
}
