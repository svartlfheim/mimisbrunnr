package cmdregistry

import "fmt"

type ErrCommandNotFound struct {
	Command string
}

func (e ErrCommandNotFound) Error() string {
	return fmt.Sprintf("command %s does not exist", e.Command)
}

type ErrCommandAlreadyRegistered struct {
	Command string
	AttemptedType string
	RegisteredType string
}

func (e ErrCommandAlreadyRegistered) Error() string {
	return fmt.Sprintf("command %s (%s) is already registered with type: %s", e.Command, e.AttemptedType, e.RegisteredType)
}