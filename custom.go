package libcmd

import (
	"fmt"
	"reflect"
	"strings"
)

// CustomArg is the interface used to create customized argument types.
// As long as you can read and write to a string, you can use this.
//
// TypeName is the name to be used to represent this type on help messages
// (e. g. int, string, value, foo). This will only be used if the user does
// not supply a sencond help message.
//
// Explain is a (optional) short string describing the custom type. This
// method receives a template string(that might be empty) in Printf style
// and 'injects' the type explanation on it. It may return different results
// based on the existence or not of a template string.
//
// Note that a empty string ("") is assumed to be the zero value
// of your custom type
type CustomArg interface {
	Get() string
	Set(value string) error
	TypeName() string
	Explain(template string) string
}

var customArgType = reflect.TypeOf(new(CustomArg)).Elem()

type choiceString struct {
	value   *string
	choices []string
}

func newChoice(target *string, choices []string) *choiceString {
	return &choiceString{
		value:   target,
		choices: choices,
	}
}

func (c *choiceString) Get() string {
	return *c.value
}

func (c *choiceString) Set(value string) error {
	if value == "" {
		*c.value = value
		return nil
	}

	for _, s := range c.choices {
		if s == value {
			*c.value = value
			return nil
		}
	}

	return parserError{err: fmt.Errorf("'%s' is not a valid value (possible values: %s)", value, strings.Join(c.choices, ","))}
}

func (c *choiceString) TypeName() string {
	return "value"
}

func (c *choiceString) Explain(template string) string {
	choices := strings.Join(c.choices, ",")
	choices = strings.Trim(choices, ",")

	if strings.Contains(template, "%s") {
		return fmt.Sprintf(template, choices)
	} else if template != "" {
		return template
	}

	return "Valid values: " + choices + "."
}
