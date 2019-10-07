package libcmd

import (
	"fmt"
	"reflect"
	"strings"
)

// CustomArg is the interface used to create customized argument types.
// As long as you can read and write to a string, you can use this.
//
// Note that a empty string ("") is assumed to be the zero value
// of your custom type
type CustomArg interface {
	Get() string
	Set(value string) error
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
