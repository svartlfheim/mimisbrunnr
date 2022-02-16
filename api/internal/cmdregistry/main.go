package cmdregistry

import (
	"fmt"
	"reflect"
)

type CommandHandler interface{
	GetName() string
	GetHelp() string
	Handle() error
}


type Registry struct {
	commands []CommandHandler
}

func (r *Registry) Register(h CommandHandler) error {
	attemptedType := reflect.TypeOf(h).Name()
	attemptedName := h.GetName()

	for _, registered := range(r.commands) {
		registeredType := reflect.TypeOf(registered).Name()
		if registered.GetName() != attemptedName {
			break
		}

		if (registeredType == attemptedType) {
			// skipping
			fmt.Printf("skipping registration of command '%s', as it's already present", h.GetName())
			break
		}

		return &ErrCommandAlreadyRegistered{
			Command: attemptedName,
			AttemptedType: attemptedType,
			RegisteredType: registeredType,
		}
	}

	r.commands = append(r.commands, h)

	return nil
}


func (r *Registry) FindHandler(c string) (CommandHandler, error) {
	for _, h := range(r.commands) {
		if h.GetName() == c {
			return h, nil
		}
	}

	return nil, &ErrCommandNotFound{
		Command: c,
	}
}

func (r *Registry) GetHelp(err error) string {
	msg := ""

	if err != nil {
		msg += fmt.Sprintf("%s\n\n", err.Error())
	}

	msg += "Available commands:\n\n"
	for _, i := range(r.commands) {
		msg += fmt.Sprintf("\t%s:\n%s\n\n\n", i.GetName(), i.GetHelp())
	}

	return msg
}

func NewRegistry() *Registry {
	return &Registry{}
}