package libcfg

import (
	"fmt"
	"strings"
)

// helper struct to determine what kind of argument was
// passesd (short, long, long with value, etc.)
type optArg struct {
	name    string
	value   string
	isShort bool
	isLong  bool
	isNeg   bool
	isEq    bool
}

func parseOptArg(argstr string) *optArg {
	var arg optArg

	arg.name = argstr
	arg.isLong = strings.HasPrefix(argstr, "--")
	arg.isNeg = strings.HasPrefix(argstr, "--no-")
	arg.isShort = !arg.isLong && strings.HasPrefix(argstr, "-")

	if !arg.isLong && !arg.isShort {
		return nil
	}

	if arg.isLong {
		splitted := strings.Split(arg.name, "=")

		if len(splitted) == 2 {
			arg.name = splitted[0]
			arg.value = splitted[1]
			arg.isEq = true
		}
	}

	return &arg
}

// inner struct to hold the values of each command line
// entry. Holds the definition provided by the user.
type optEntry struct {
	long  string
	short string
	help  string
	val   *variant
}

// try to 'fix' the options that have a 'natural' default
// value in command line. Examples:
// -b        : assumes 'true' in case of a boolean entry
// --bool    : same as above
// --no-bool : same as above, but assumes 'false'
func (entry *optEntry) fillAutoValue(arg *optArg) {
	if arg.value != "" {
		return
	}

	switch {
	case entry.val.isBool && arg.isNeg:
		arg.value = "false"
	case entry.val.isBool:
		arg.value = "true"
	}
}

// sets this entry value with the value from command-line
func (entry *optEntry) setValue(arg *optArg) error {
	// the option '--string=' is the only case where
	// an empty value should be accepted
	if arg.value == "" && !(entry.val.isStr && arg.isEq) {
		return fmt.Errorf("no value for argument: %s", arg.name)
	}

	if err := entry.val.setValueStr(arg.value); err != nil {
		return fmt.Errorf("error parsing argument '%s': %v", arg.name, err)
	}

	if entry.val.isBool && arg.isNeg && arg.isEq {
		b, _ := entry.val.ptrValue.(*bool)
		*b = !*b
	}

	return nil
}
