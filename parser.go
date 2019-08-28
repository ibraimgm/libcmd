package libcfg

import (
	"fmt"
	"strconv"
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

func (arg *optArg) setBool(value *bool) error {
	bval, err := strconv.ParseBool(arg.value)
	if err != nil {
		return fmt.Errorf("'%v' is not a valid boolean value for argument '%s'", arg.value, arg.name)
	}
	*value = bval

	if arg.isNeg && arg.isEq {
		*value = !*value
	}

	return nil
}

func (arg *optArg) setInt(value *int) error {
	ival, err := strconv.ParseInt(arg.value, 10, 64)
	if err != nil {
		return fmt.Errorf("'%v' is not a valid int value for argument '%s'", arg.value, arg.name)
	}
	*value = int(ival)

	return nil
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
		return arg.setBool(v)

	case *int:
		return arg.setInt(v)

	case *string:
		*v = arg.value

	default:
		return fmt.Errorf("unrecognized entry type: %T", entry.valuePtr)
	}

	return nil
}

// CfgParser is a parser that can load configurations from the command line or the environment
// variables.
type CfgParser struct {
	args       []string
	optentries []*optEntry
	shortopt   map[string]*optEntry
	longopt    map[string]*optEntry
}

// NewParser returns a new CfgParser, ready to be used.
func NewParser() *CfgParser {
	return &CfgParser{
		optentries: make([]*optEntry, 0),
		shortopt:   make(map[string]*optEntry),
		longopt:    make(map[string]*optEntry),
	}
}

func (cfg *CfgParser) addOpt(entry *optEntry) {
	cfg.optentries = append(cfg.optentries, entry)

	if entry.short != "" {
		cfg.shortopt["-"+entry.short] = entry
	}

	if entry.long != "" {
		cfg.longopt["--"+entry.long] = entry

		if _, ok := entry.valuePtr.(*bool); ok {
			cfg.longopt["--no-"+entry.long] = entry
		}
	}
}

func (cfg *CfgParser) findEntry(entryName string) *optEntry {
	if entry, ok := cfg.shortopt[entryName]; ok {
		return entry
	}

	return cfg.longopt[entryName]
}

// ParseArgs parsers the arguments in args and load the configuration
// values according to the settings of the current parser.
// Note that args must not include the program name
func (cfg *CfgParser) ParseArgs(args []string) error {
	for i := 0; i < len(args); i++ {

		// parse the current argument
		arg := parseOptArg(args[i])

		// if it is not a param, break the parsing and collect
		// the rest of the list
		if arg == nil {
			cfg.args = args[i:]
			return nil
		}

		// find the entry.
		// if no entry exists, this argument is unknown
		entry := cfg.findEntry(arg.name)
		if entry == nil {
			return fmt.Errorf("unknown argument: %s", arg.name)
		}

		// some argument types have autmatic values in certain cases
		// fill them in, if necessary
		entry.fillAutoValue(arg)

		// check if we need to look at the next argument
		// long params with '=' should not be considered
		if !arg.isEq && arg.value == "" {
			if i+1 == len(args) {
				return fmt.Errorf("no value for argument: %s", arg.name)
			}

			arg.value = args[i+1]
			i++
		}

		if err := entry.setValue(arg); err != nil {
			return err
		}
	}

	return nil
}

// Args returns the remaining non-parsed command line arguments
func (cfg *CfgParser) Args() []string {
	return cfg.args
}

// OptBool creates a new parser setting to load a bool value from the command line only.
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptBool(long, short string, defaultValue bool, help string) *bool {
	val := new(bool)
	cfg.addOpt(&optEntry{long: long, short: short, help: help, defaultValue: defaultValue, valuePtr: val, isBool: true})
	return val
}

// OptInt creates a new parser setting to load a int value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptInt(long, short string, defaultValue int, help string) *int {
	val := new(int)
	cfg.addOpt(&optEntry{long: long, short: short, help: help, defaultValue: defaultValue, valuePtr: val})
	return val
}

// OptString creates a new parser setting to load a string value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptString(long, short, defaultValue, help string) *string {
	val := new(string)
	cfg.addOpt(&optEntry{long: long, short: short, help: help, defaultValue: defaultValue, valuePtr: val, isStr: true})
	return val
}
