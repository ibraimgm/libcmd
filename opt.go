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
	long         string
	short        string
	help         string
	defaultValue interface{}
	isBool       bool
	isStr        bool
	valuePtr     interface{}
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
	case entry.isBool && arg.isNeg:
		arg.value = "false"
	case entry.isBool:
		arg.value = "true"
	}
}

// sets this entry value with the value from command-line
func (entry *optEntry) setValue(arg *optArg) error {
	// the option '--string=' is the only case where
	// an empty value should be accepted
	if arg.value == "" && !(entry.isStr && arg.isEq) {
		return fmt.Errorf("no value for argument: %s", arg.name)
	}

	switch v := entry.valuePtr.(type) {
	case *bool:
		return wrapOptErr(arg.name, setBool(v, arg.value, arg.isNeg && arg.isEq))

	case *int:
		return wrapOptErr(arg.name, setInt(v, arg.value))

	case *string:
		*v = arg.value

	default:
		return fmt.Errorf("unrecognized entry type: %T", entry.valuePtr)
	}

	return nil
}
