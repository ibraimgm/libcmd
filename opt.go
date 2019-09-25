package libcmd

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

func (entry *optEntry) helpHeader() string {
	switch {
	case entry.short != "" && entry.long != "":
		return "-" + entry.short + ", --" + entry.long

	case entry.short != "":
		return "-" + entry.short

	default:
		return "--" + entry.long
	}
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
		b := entry.val.refValue.Bool()
		entry.val.refValue.SetBool(!b)
	}

	return nil
}

// add and optentry to the parser
func (cmd *Cmd) addOpt(entry *optEntry) {
	if entry.short == "" && entry.long == "" {
		return
	}

	cmd.optentries = append(cmd.optentries, entry)

	if entry.short != "" {
		cmd.shortopt["-"+entry.short] = entry
	}

	if entry.long != "" {
		cmd.longopt["--"+entry.long] = entry

		if entry.val.isBool {
			cmd.longopt["--no-"+entry.long] = entry
		}
	}
}

// find an entry (with '-' or '--')
func (cmd *Cmd) findOpt(entryName string) *optEntry {
	if entry, ok := cmd.shortopt[entryName]; ok {
		return entry
	}

	return cmd.longopt[entryName]
}

// parse all command-line arguments
func (cmd *Cmd) doParse(args []string) error {
	for i := 0; i < len(args); i++ {

		// parse the current argument
		arg := parseOptArg(args[i])

		// if it is an operand, bail out!
		if arg == nil {
			cmd.args = args[i:]
			return nil
		}

		// find the entry.
		// if no entry exists, this argument is unknown
		entry := cmd.findOpt(arg.name)
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

func (cmd *Cmd) doRun(args []string) error {
	cmd.configure()

	if err := cmd.doParse(args); err != nil {
		if cmd.errHandler != nil {
			err = cmd.errHandler(err)
		}

		if err != nil {
			return err
		}
	}

	for i := range cmd.optentries {
		cmd.optentries[i].val.useDefault()
	}

	if cmd.match != nil {
		cmd.match(cmd)
	}

	if len(cmd.args) >= 1 {
		name := cmd.args[0]

		if subCommand, ok := cmd.commands[name]; ok {
			if subCommand.callback != nil {
				subCommand.callback(subCommand)
			}

			err := subCommand.doRun(cmd.args[1:])
			cmd.args = subCommand.args

			return err
		}
	}

	// leaf command
	return cmd.runLeafCommand()
}

func (cmd *Cmd) runLeafCommand() error {
	var showHelp bool

	// check for the automatic handling of
	// the '-h' and '--help' flags
	if !cmd.options.SupressPrintHelpWhenSet {
		arg := cmd.shortopt["-h"]
		if arg == nil {
			arg = cmd.longopt["--help"]
		}

		if arg != nil {
			showHelp = arg.val.refValue.Bool()
		}
	}

	if showHelp {
		cmd.Help()
		return nil
	}

	// actual command, as defined by the user
	if cmd.run != nil {
		return cmd.run(cmd)
	}

	// if i'm the main app, do not show the help
	if cmd.parentCmd == nil {
		return nil
	}

	// the last resort is to run the help if we're a "partial"
	// subcommand
	if !cmd.options.SuppressPrintHelpPartialCommand {
		cmd.Help()
	}

	return nil
}

// Args returns the remaining non-parsed command line arguments.
func (cmd *Cmd) Args() []string {
	return cmd.args
}

// StringP defines a new string argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) StringP(target *string, long, short, defaultValue, help string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// BoolP defines a new bool argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) BoolP(target *bool, long, short string, defaultValue bool, help string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// IntP defines a new int argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) IntP(target *int, long, short string, defaultValue int, help string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Int8P defines a new int8 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Int8P(target *int8, long, short string, defaultValue int8, help string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Int16P defines a new int16 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Int16P(target *int16, long, short string, defaultValue int16, help string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Int32P defines a new int32 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Int32P(target *int32, long, short string, defaultValue int32, help string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Int64P defines a new int64 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Int64P(target *int64, long, short string, defaultValue int64, help string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// UintP defines a new uint argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) UintP(target *uint, long, short string, defaultValue uint, help string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Uint8P defines a new uint8 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Uint8P(target *uint8, long, short string, defaultValue uint8, help string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Uint16P defines a new uint16 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Uint16P(target *uint16, long, short string, defaultValue uint16, help string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Uint32P defines a new uint32 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Uint32P(target *uint32, long, short string, defaultValue uint32, help string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Uint64P defines a new uint64 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Uint64P(target *uint64, long, short string, defaultValue uint64, help string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Float32P defines a new float32 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Float32P(target *float32, long, short string, defaultValue float32, help string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Float64P defines a new float64 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Float64P(target *float64, long, short string, defaultValue float64, help string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// String defines a new string argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) String(long, short, defaultValue, help string) *string {
	target := new(string)
	cmd.StringP(target, long, short, defaultValue, help)
	return target
}

// Bool defines a new bool argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Bool(long, short string, defaultValue bool, help string) *bool {
	target := new(bool)
	cmd.BoolP(target, long, short, defaultValue, help)
	return target
}

// Int defines a new int argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Int(long, short string, defaultValue int, help string) *int {
	target := new(int)
	cmd.IntP(target, long, short, defaultValue, help)
	return target
}

// Int8 defines a new int8 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Int8(long, short string, defaultValue int8, help string) *int8 {
	target := new(int8)
	cmd.Int8P(target, long, short, defaultValue, help)
	return target
}

// Int16 defines a new int16 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Int16(long, short string, defaultValue int16, help string) *int16 {
	target := new(int16)
	cmd.Int16P(target, long, short, defaultValue, help)
	return target
}

// Int32 defines a new int32 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Int32(long, short string, defaultValue int32, help string) *int32 {
	target := new(int32)
	cmd.Int32P(target, long, short, defaultValue, help)
	return target
}

// Int64 defines a new int64 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Int64(long, short string, defaultValue int64, help string) *int64 {
	target := new(int64)
	cmd.Int64P(target, long, short, defaultValue, help)
	return target
}

// Uint defines a new uint argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Uint(long, short string, defaultValue uint, help string) *uint {
	target := new(uint)
	cmd.UintP(target, long, short, defaultValue, help)
	return target
}

// Uint8 defines a new uint8 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Uint8(long, short string, defaultValue uint8, help string) *uint8 {
	target := new(uint8)
	cmd.Uint8P(target, long, short, defaultValue, help)
	return target
}

// Uint16 defines a new uint16 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Uint16(long, short string, defaultValue uint16, help string) *uint16 {
	target := new(uint16)
	cmd.Uint16P(target, long, short, defaultValue, help)
	return target
}

// Uint32 defines a new uint32 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Uint32(long, short string, defaultValue uint32, help string) *uint32 {
	target := new(uint32)
	cmd.Uint32P(target, long, short, defaultValue, help)
	return target
}

// Uint64 defines a new uint64 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Uint64(long, short string, defaultValue uint64, help string) *uint64 {
	target := new(uint64)
	cmd.Uint64P(target, long, short, defaultValue, help)
	return target
}

// Float32 defines a new float32 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Float32(long, short string, defaultValue float32, help string) *float32 {
	target := new(float32)
	cmd.Float32P(target, long, short, defaultValue, help)
	return target
}

// Float64 defines a new float64 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Float64(long, short string, defaultValue float64, help string) *float64 {
	target := new(float64)
	cmd.Float64P(target, long, short, defaultValue, help)
	return target
}
