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
const uniqueRule Rule = "unique"
const uniquePerOtherFieldRule Rule = "uniqueperotherfield"
const lessThanOrEqualToRule Rule = "lte"
const existsRule Rule = "exists"
const uuidRule Rule = "uuid"

type MessageGenerator func(Error) string
type ParameterParser func(Error) map[string]string

var messagesForRule map[Rule]MessageGenerator = map[Rule]MessageGenerator{
	requiredRule: func(Error) string {
		return "is required"
	},
	uniqueRule: func(ve Error) string {
		return "value must be unique across all records of this type"
	},

	uuidRule: func(ve Error) string {
		return "must be a valid uuid"
	},

	lessThanOrEqualToRule: func(ve Error) string {
		limit, found := ve.Parameters()["limit"]

		switch ve.valueType {
		case "string":
			if !found {
				return "must contain less characters"
			}

			return fmt.Sprintf("must contain less than or equal to %s characters", limit)
		case "int":
			if !found {
				return "must be a smaller number"
			}

			return fmt.Sprintf("must be less than or equal to %s", limit)
		default:
			return "must be smaller"
		}
	},

	greaterThanRule: func(ve Error) string {
		limit, found := ve.Parameters()["limit"]

		switch ve.valueType {
		case "string":
			if !found {
				return "must contain more characters"
			}

			return fmt.Sprintf("must contain more than %s characters", limit)
		case "int":
			if !found {
				return "must be a larger number"
			}

			return fmt.Sprintf("must be larger than %s", limit)
		default:
			return "must be larger"
		}
	},

	existsRule: func(ve Error) string {
		return "must already exist"
	},

	uniquePerOtherFieldRule: func(ve Error) string {
		// Parameters func will panic if it wasn't found
		field := ve.Parameters()["field"]

		return fmt.Sprintf("value must be unique across all records of this type with the same value for: %s", field)
	},
}

var parameterParsers map[Rule]ParameterParser = map[Rule]ParameterParser{
	greaterThanRule: func(ve Error) map[string]string {
		return map[string]string{
			"limit": ve.param,
		}
	},
	lessThanRule: func(ve Error) map[string]string {
		return map[string]string{
			"limit": ve.param,
		}
	},
	lessThanOrEqualToRule: func(ve Error) map[string]string {
		return map[string]string{
			"limit": ve.param,
		}
	},

	uniquePerOtherFieldRule: func(ve Error) map[string]string {
		if ve.param == "" {
			panic(fmt.Errorf("invalid %s rule config, no field defined", string(uniquePerOtherFieldRule)))
		}

		return map[string]string{
			"field": ve.param,
		}
	},
}

type Error struct {
	path                   string
	rule                   string
	param                  string
	valueType              string
	extraMessageGenerators map[Rule]MessageGenerator
	extraParameterParsers  map[Rule]ParameterParser
}

func (ve Error) Path() string {
	return ve.path
}

func (ve Error) Message() string {
	rule := Rule(ve.Rule())

	allRules := map[Rule]MessageGenerator{}

	for r, mg := range messagesForRule {
		allRules[r] = mg
	}

	for r, mg := range ve.extraMessageGenerators {
		allRules[r] = mg
	}

	if generator, found := allRules[rule]; found {
		return generator(ve)
	}

	return "is invalid"
}

func (ve Error) Rule() string {
	return ve.rule
}

func (ve Error) Parameters() map[string]string {
	rule := Rule(ve.Rule())

	// If there is a custom rule, go straight through to it
	// The built-in ones should only be run if the param is not empty
	// It makes sens for the custom ones to avoid this behaviour
	// They should account for an empty param if required
	if parser, found := ve.extraParameterParsers[rule]; found {
		return parser(ve)
	}

	if ve.param == "" {
		return map[string]string{}
	}

	if parser, found := parameterParsers[rule]; found {
		return parser(ve)
	}

	return map[string]string{
		"param": ve.param,
	}
}

func NewError(path string, rule string, param string, valueType string, extraMessageGenerators map[Rule]MessageGenerator, extraParameterParsers map[Rule]ParameterParser) Error {
	return Error{
		path:                   path,
		rule:                   rule,
		param:                  param,
		valueType:              valueType,
		extraMessageGenerators: extraMessageGenerators,
		extraParameterParsers:  extraParameterParsers,
	}
}
