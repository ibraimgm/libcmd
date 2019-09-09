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

	if err := entry.val.setValue(arg.value); err != nil {
		return fmt.Errorf("error parsing argument '%s': %v", arg.name, err)
	}

	if entry.val.isBool && arg.isNeg && arg.isEq {
		b, _ := entry.val.ptrValue.(*bool)
		*b = !*b
	}

	return nil
}

// add and optentry to the parser
func (cfg *cfgParser) addOpt(entry *optEntry) {
	cfg.optentries = append(cfg.optentries, entry)

	if entry.short != "" {
		cfg.shortopt["-"+entry.short] = entry
	}

	if entry.long != "" {
		cfg.longopt["--"+entry.long] = entry

		if entry.val.isBool {
			cfg.longopt["--no-"+entry.long] = entry
		}
	}
}

// find an entry (with '-' or '--')
func (cfg *cfgParser) findOpt(entryName string) *optEntry {
	if entry, ok := cfg.shortopt[entryName]; ok {
		return entry
	}

	return cfg.longopt[entryName]
}

// try to find an entry by it's "canonical" name
func (cfg *cfgParser) findByName(optName string) *optEntry {
	if entry, ok := cfg.longopt["--"+optName]; ok {
		return entry
	}

	if entry, ok := cfg.shortopt["-"+optName]; ok {
		return entry
	}

	return nil
}

// parse all command-line arguments
func (cfg *cfgParser) doParse(args []string) error {
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
		entry := cfg.findOpt(arg.name)
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
		entry.val.isOpt = true
	}

	return nil
}

// Args returns the remaining non-parsed command line arguments.
func (cfg *cfgParser) Args() []string {
	return cfg.args
}

func (cfg *cfgParser) StringP(target *string, long, short, defaultValue, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue, isStr: true}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.envLoader.addEnv(&val, variables)
}

func (cfg *cfgParser) BoolP(target *bool, long, short string, defaultValue bool, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue, isBool: true}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.envLoader.addEnv(&val, variables)
}

func (cfg *cfgParser) IntP(target *int, long, short string, defaultValue int, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.envLoader.addEnv(&val, variables)
}

func (cfg *cfgParser) Int8P(target *int8, long, short string, defaultValue int8, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.envLoader.addEnv(&val, variables)
}

func (cfg *cfgParser) Int16P(target *int16, long, short string, defaultValue int16, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.envLoader.addEnv(&val, variables)
}

func (cfg *cfgParser) Int32P(target *int32, long, short string, defaultValue int32, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.envLoader.addEnv(&val, variables)
}

func (cfg *cfgParser) Int64P(target *int64, long, short string, defaultValue int64, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.envLoader.addEnv(&val, variables)
}

func (cfg *cfgParser) UintP(target *uint, long, short string, defaultValue uint, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.envLoader.addEnv(&val, variables)
}

func (cfg *cfgParser) Uint8P(target *uint8, long, short string, defaultValue uint8, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.envLoader.addEnv(&val, variables)
}

func (cfg *cfgParser) Uint16P(target *uint16, long, short string, defaultValue uint16, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.envLoader.addEnv(&val, variables)
}

func (cfg *cfgParser) Uint32P(target *uint32, long, short string, defaultValue uint32, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.envLoader.addEnv(&val, variables)
}

func (cfg *cfgParser) Uint64P(target *uint64, long, short string, defaultValue uint64, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.envLoader.addEnv(&val, variables)
}

func (cfg *cfgParser) Float32P(target *float32, long, short string, defaultValue float32, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.envLoader.addEnv(&val, variables)
}

func (cfg *cfgParser) Float64P(target *float64, long, short string, defaultValue float64, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.envLoader.addEnv(&val, variables)
}

func (cfg *cfgParser) String(long, short, defaultValue, help string, variables ...string) *string {
	target := new(string)
	cfg.StringP(target, long, short, defaultValue, help, variables...)
	return target
}

func (cfg *cfgParser) Bool(long, short string, defaultValue bool, help string, variables ...string) *bool {
	target := new(bool)
	cfg.BoolP(target, long, short, defaultValue, help, variables...)
	return target
}

func (cfg *cfgParser) Int(long, short string, defaultValue int, help string, variables ...string) *int {
	target := new(int)
	cfg.IntP(target, long, short, defaultValue, help, variables...)
	return target
}

func (cfg *cfgParser) Int8(long, short string, defaultValue int8, help string, variables ...string) *int8 {
	target := new(int8)
	cfg.Int8P(target, long, short, defaultValue, help, variables...)
	return target
}

func (cfg *cfgParser) Int16(long, short string, defaultValue int16, help string, variables ...string) *int16 {
	target := new(int16)
	cfg.Int16P(target, long, short, defaultValue, help, variables...)
	return target
}

func (cfg *cfgParser) Int32(long, short string, defaultValue int32, help string, variables ...string) *int32 {
	target := new(int32)
	cfg.Int32P(target, long, short, defaultValue, help, variables...)
	return target
}

func (cfg *cfgParser) Int64(long, short string, defaultValue int64, help string, variables ...string) *int64 {
	target := new(int64)
	cfg.Int64P(target, long, short, defaultValue, help, variables...)
	return target
}

func (cfg *cfgParser) Uint(long, short string, defaultValue uint, help string, variables ...string) *uint {
	target := new(uint)
	cfg.UintP(target, long, short, defaultValue, help, variables...)
	return target
}

func (cfg *cfgParser) Uint8(long, short string, defaultValue uint8, help string, variables ...string) *uint8 {
	target := new(uint8)
	cfg.Uint8P(target, long, short, defaultValue, help, variables...)
	return target
}

func (cfg *cfgParser) Uint16(long, short string, defaultValue uint16, help string, variables ...string) *uint16 {
	target := new(uint16)
	cfg.Uint16P(target, long, short, defaultValue, help, variables...)
	return target
}

func (cfg *cfgParser) Uint32(long, short string, defaultValue uint32, help string, variables ...string) *uint32 {
	target := new(uint32)
	cfg.Uint32P(target, long, short, defaultValue, help, variables...)
	return target
}

func (cfg *cfgParser) Uint64(long, short string, defaultValue uint64, help string, variables ...string) *uint64 {
	target := new(uint64)
	cfg.Uint64P(target, long, short, defaultValue, help, variables...)
	return target
}

func (cfg *cfgParser) Float32(long, short string, defaultValue float32, help string, variables ...string) *float32 {
	target := new(float32)
	cfg.Float32P(target, long, short, defaultValue, help, variables...)
	return target
}

func (cfg *cfgParser) Float64(long, short string, defaultValue float64, help string, variables ...string) *float64 {
	target := new(float64)
	cfg.Float64P(target, long, short, defaultValue, help, variables...)
	return target
}
