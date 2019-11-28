package libcmd

import (
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
	short rune
	help  []string
	val   *variant
}

func (entry *optEntry) helpHeader() string {
	var s string
	var kindSep string

	// 'short' name
	if entry.short != 0 {
		s = "-" + string(entry.short)
		kindSep = " "
	}

	// 'long' name, aligned
	if entry.long != "" {
		if entry.short != 0 {
			s += ", "
		}

		s += "--" + entry.long
		kindSep = "="
	}

	// value kind
	if len(entry.help) >= 2 {
		s += kindSep + entry.help[1]
	} else if !entry.val.isBool {
		if entry.val.refValue.Type().Implements(customArgType) {
			ca, _ := entry.val.refValue.Interface().(CustomArg)
			s += kindSep + ca.TypeName()
		} else {
			s += kindSep + entry.val.refValue.Kind().String()
		}
	}

	return s
}

func (entry *optEntry) helpExplain() string {
	var explain string

	// user supplyed text
	if len(entry.help) > 0 {
		explain = entry.help[0]
	}

	// expand template of custom types
	if entry.val.refValue.Type().Implements(customArgType) {
		ca, _ := entry.val.refValue.Interface().(CustomArg)
		explain = ca.Explain(explain)
	}

	// if it is still empty, put a generic message
	if explain == "" {
		explain = "Sets the argument value."
	}

	// adds the default value
	if def := entry.val.defaultAsString(); def != "" {
		explain += " (default: " + def + ")"
	}

	return explain
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
		return noValueErr{arg: arg.name}
	}

	if err := entry.val.setValue(arg.value); err != nil {
		return parserError{arg: arg.name, err: err}
	}

	if entry.val.isBool && arg.isNeg && arg.isEq {
		b := entry.val.refValue.Bool()
		entry.val.refValue.SetBool(!b)
	}

	return nil
}

// add and optentry to the parser
func (cmd *Cmd) addOpt(entry *optEntry) {
	if entry.short < 0 {
		entry.short = 0
	}

	if entry.short == 0 && entry.long == "" {
		return
	}

	cmd.optentries = append(cmd.optentries, entry)

	if entry.short != 0 {
		cmd.shortopt["-"+string(entry.short)] = entry
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

		// if is a bunch of flags in 'compressed' form,
		// process them all, except the last one
		if err := cmd.processMultiArgs(arg); err != nil {
			return err
		}

		// find the entry.
		// if no entry exists, this argument is unknown
		entry := cmd.findOpt(arg.name)
		if entry == nil {
			return unknownArgErr{arg: arg.name}
		}

		// some argument types have automatic values in certain cases
		// fill them in, if necessary
		entry.fillAutoValue(arg)

		// check if we need to look at the next argument
		// long params with '=' should not be considered
		if !arg.isEq && arg.value == "" {
			if i+1 == len(args) {
				return noValueErr{arg: arg.name}
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

// handles the compressed form, i. e. '-abc' instead of '-a -b -c'
// Every have to be a bool, except the last one. This routine sets
// the value of every flag, excet the last, and adjusts 'arg' so it points
// to the last flag only
func (cmd *Cmd) processMultiArgs(arg *optArg) error {
	if !arg.isShort || len(arg.name) <= 2 {
		return nil
	}

	names := strings.Split(arg.name, "")[1:]
	for i := 0; i < len(names)-1; i++ {
		entry := cmd.findOpt("-" + names[i])
		if entry == nil || !entry.val.isBool {
			return unknownArgErr{arg: arg.name}
		}

		if err := entry.val.setValue("true"); err != nil {
			return err
		}
	}

	arg.name = "-" + names[len(names)-1]
	return nil
}

func (cmd *Cmd) doRun(args []string) error {
	cmd.setupHelp()

	if err := cmd.doParse(args); err != nil {
		if cmd.errHandler != nil {
			err = cmd.errHandler(err)
		}

		if err != nil {
			return err
		}
	}

	for i := range cmd.optentries {
		if err := cmd.optentries[i].val.useDefault(); err != nil {
			return err
		}
	}

	if cmd.match != nil {
		cmd.match(cmd)
	}

	if len(cmd.args) >= 1 {
		name := cmd.args[0]

		if subCommand, ok := cmd.commands[name]; ok {
			subCommand.Options = cmd.Options
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
	if !cmd.Options.SupressPrintHelpWhenSet {
		arg := cmd.shortopt["-h"]
		if arg == nil {
			arg = cmd.longopt["--help"]
		}

		if arg != nil && arg.val.isBool {
			showHelp = arg.val.refValue.Bool()
		}
	}

	if showHelp {
		cmd.Help()
		return nil
	}

	// check for operands
	if err := cmd.checkOperands(); err != nil {
		return err
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
	if !cmd.Options.SuppressPrintHelpPartialCommand {
		cmd.Help()
	}

	return nil
}

func (cmd *Cmd) checkOperands() error {
	// if we're permissive, there's nothing to do
	if !cmd.Options.StrictOperands {
		return nil
	}

	// consider only the required ones
	var need int
	var hasOptionals bool
	for _, op := range cmd.operands {
		if op.modifier == "" {
			need++
		} else {
			hasOptionals = true
		}
	}

	// if at least one is optional, no need for an exact number of
	// arguments
	if hasOptionals && need > len(cmd.args) {
		return operandRequiredErr{required: need, got: len(cmd.args)}
	}

	// if e do not have optional arguments, we need an exact number
	if !hasOptionals && need != len(cmd.args) {
		return operandRequiredErr{required: need, got: len(cmd.args), exact: true}
	}

	// we should be good to go now
	return nil
}

// Args returns the remaining non-parsed command line arguments.
func (cmd *Cmd) Args() []string {
	return cmd.args
}

// StringP defines a new string argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) StringP(target *string, long string, short rune, defaultValue string, help ...string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// BoolP defines a new bool argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) BoolP(target *bool, long string, short rune, defaultValue bool, help ...string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// IntP defines a new int argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) IntP(target *int, long string, short rune, defaultValue int, help ...string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Int8P defines a new int8 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Int8P(target *int8, long string, short rune, defaultValue int8, help ...string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Int16P defines a new int16 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Int16P(target *int16, long string, short rune, defaultValue int16, help ...string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Int32P defines a new int32 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Int32P(target *int32, long string, short rune, defaultValue int32, help ...string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Int64P defines a new int64 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Int64P(target *int64, long string, short rune, defaultValue int64, help ...string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// UintP defines a new uint argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) UintP(target *uint, long string, short rune, defaultValue uint, help ...string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Uint8P defines a new uint8 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Uint8P(target *uint8, long string, short rune, defaultValue uint8, help ...string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Uint16P defines a new uint16 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Uint16P(target *uint16, long string, short rune, defaultValue uint16, help ...string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Uint32P defines a new uint32 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Uint32P(target *uint32, long string, short rune, defaultValue uint32, help ...string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Uint64P defines a new uint64 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Uint64P(target *uint64, long string, short rune, defaultValue uint64, help ...string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Float32P defines a new float32 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Float32P(target *float32, long string, short rune, defaultValue float32, help ...string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// Float64P defines a new float64 argument. After parsing, the argument value
// will be available in the specified pointer.
func (cmd *Cmd) Float64P(target *float64, long string, short rune, defaultValue float64, help ...string) {
	val := varFromInterface(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// CustomP defines a new argument with custom type. During parsing, the argument
// is manipulated via Get and Set methods of the CustomArg interface.
func (cmd *Cmd) CustomP(target CustomArg, long string, short rune, defaultValue string, help ...string) {
	val := varFromCustom(target, defaultValue)
	cmd.addOpt(&optEntry{long: long, short: short, help: help, val: val})
}

// ChoiceP defines a new string argument, ith tha values limited by choices.
// After parsing, the argument value will be available in the specified pointer.
//
// If the defaultValue is always considered 'valid', even when not listed on
// the choices parameter.
func (cmd *Cmd) ChoiceP(target *string, choices []string, long string, short rune, defaultValue string, help ...string) {
	// add the default value in the 'valid choices' list,
	// if it isn't present already
	valid := make([]string, len(choices), len(choices)+1)
	copy(valid, choices)

	var found bool
	for _, s := range valid {
		if s == defaultValue {
			found = true
			break
		}
	}

	if !found {
		valid = append(valid, defaultValue)
	}

	cmd.CustomP(newChoice(target, valid), long, short, defaultValue, help...)
}

// String defines a new string argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) String(long string, short rune, defaultValue string, help ...string) *string {
	target := new(string)
	cmd.StringP(target, long, short, defaultValue, help...)
	return target
}

// Bool defines a new bool argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Bool(long string, short rune, defaultValue bool, help ...string) *bool {
	target := new(bool)
	cmd.BoolP(target, long, short, defaultValue, help...)
	return target
}

// Int defines a new int argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Int(long string, short rune, defaultValue int, help ...string) *int {
	target := new(int)
	cmd.IntP(target, long, short, defaultValue, help...)
	return target
}

// Int8 defines a new int8 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Int8(long string, short rune, defaultValue int8, help ...string) *int8 {
	target := new(int8)
	cmd.Int8P(target, long, short, defaultValue, help...)
	return target
}

// Int16 defines a new int16 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Int16(long string, short rune, defaultValue int16, help ...string) *int16 {
	target := new(int16)
	cmd.Int16P(target, long, short, defaultValue, help...)
	return target
}

// Int32 defines a new int32 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Int32(long string, short rune, defaultValue int32, help ...string) *int32 {
	target := new(int32)
	cmd.Int32P(target, long, short, defaultValue, help...)
	return target
}

// Int64 defines a new int64 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Int64(long string, short rune, defaultValue int64, help ...string) *int64 {
	target := new(int64)
	cmd.Int64P(target, long, short, defaultValue, help...)
	return target
}

// Uint defines a new uint argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Uint(long string, short rune, defaultValue uint, help ...string) *uint {
	target := new(uint)
	cmd.UintP(target, long, short, defaultValue, help...)
	return target
}

// Uint8 defines a new uint8 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Uint8(long string, short rune, defaultValue uint8, help ...string) *uint8 {
	target := new(uint8)
	cmd.Uint8P(target, long, short, defaultValue, help...)
	return target
}

// Uint16 defines a new uint16 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Uint16(long string, short rune, defaultValue uint16, help ...string) *uint16 {
	target := new(uint16)
	cmd.Uint16P(target, long, short, defaultValue, help...)
	return target
}

// Uint32 defines a new uint32 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Uint32(long string, short rune, defaultValue uint32, help ...string) *uint32 {
	target := new(uint32)
	cmd.Uint32P(target, long, short, defaultValue, help...)
	return target
}

// Uint64 defines a new uint64 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Uint64(long string, short rune, defaultValue uint64, help ...string) *uint64 {
	target := new(uint64)
	cmd.Uint64P(target, long, short, defaultValue, help...)
	return target
}

// Float32 defines a new float32 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Float32(long string, short rune, defaultValue float32, help ...string) *float32 {
	target := new(float32)
	cmd.Float32P(target, long, short, defaultValue, help...)
	return target
}

// Float64 defines a new float64 argument. After parsing, the argument value
// will be available in the returned pointer.
func (cmd *Cmd) Float64(long string, short rune, defaultValue float64, help ...string) *float64 {
	target := new(float64)
	cmd.Float64P(target, long, short, defaultValue, help...)
	return target
}

// Choice defines a new string argument, ith tha values limited by choices.
// After parsing, the argument value will be available in the returned pointer.
//
// If the defaultValue is always considered 'valid', even when not listed on
// the choices parameter.
func (cmd *Cmd) Choice(choices []string, long string, short rune, defaultValue string, help ...string) *string {
	target := new(string)
	cmd.ChoiceP(target, choices, long, short, defaultValue, help...)
	return target
}
