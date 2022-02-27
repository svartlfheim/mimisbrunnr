package validation

import (
	"fmt"
)

type ValidationError interface {
	Path() string
	Message() string
	Rule() string
	Parameters() map[string]string
}

type Rule string
const requiredRule Rule = "required"
const greaterThanRule Rule = "gt"
const lessThanRule Rule = "lt"

type messageGenerator func(validationError) string
type parameterParser func(validationError) map[string]string

var messagesForRule map[Rule]messageGenerator = map[Rule]messageGenerator{
	requiredRule: func(validationError) string {
		return "is required"
	},
	greaterThanRule: func(ve validationError) string {
		limit, found := ve.Parameters()["limit"]

		switch ve.valueType {
		case "string":
			if ! found {
				return "must contain more characters"
			}

			return fmt.Sprintf("must contain more than %s characters", limit)
		case "int":
			if ! found {
				return "must be a larger number"
			}

			return fmt.Sprintf("must be larger than %s", limit)
		default:
			return "must be larger"
		}
	},
}

var parameterParsers map[Rule]parameterParser = map[Rule]parameterParser{
	greaterThanRule: func (ve validationError) map[string]string {
		return map[string]string{
			"limit": ve.param,
		}
	},
	lessThanRule: func (ve validationError) map[string]string {
		return map[string]string{
			"limit": ve.param,
		}
	},
}

type validationError struct {
	path      string
	rule      string
	param     string
	valueType string
}

func (ve validationError) Path() string {
	return ve.path
}

func (ve validationError) Message() string {
	rule := Rule(ve.Rule())

	if generator, found := messagesForRule[rule]; found {
		return generator(ve)
	}

	return "is invalid"
}

func (ve validationError) Rule() string {
	return ve.rule
}

func (ve validationError) Parameters() map[string]string {
	if ve.param == "" {
		return map[string]string{}
	}

	rule := Rule(ve.Rule())

	if parser, found := parameterParsers[rule]; found {
		return parser(ve)
	}

	return map[string]string{
		"param": ve.param,
	}
}
