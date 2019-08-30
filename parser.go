package libcfg

import (
	"fmt"
)

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

		if entry.val.isBool {
			cfg.longopt["--no-"+entry.long] = entry
		}
	}
}

func (cfg *CfgParser) findOpt(entryName string) *optEntry {
	if entry, ok := cfg.shortopt[entryName]; ok {
		return entry
	}

	return cfg.longopt[entryName]
}

// RunArgs loads all the configuration values according to
// the settingsof the current parser, but assumes the values
// passed by args as command line arguments.
// Note that args must not include the program name
func (cfg *CfgParser) RunArgs(args []string) error {
	if err := cfg.doParse(args); err != nil {
		return err
	}

	for i := range cfg.optentries {
		cfg.optentries[i].val.useDefault()
	}

	return nil
}

// parse all command-line arguments
func (cfg *CfgParser) doParse(args []string) error {
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
	}

	return nil
}

// Args returns the remaining non-parsed command line arguments
func (cfg *CfgParser) Args() []string {
	return cfg.args
}

// OptString creates a new parser setting to load a string value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptString(long, short, defaultValue, help string) *string {
	ptr := new(string)
	val := variant{ptrValue: ptr, defaultValue: defaultValue, isStr: true}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	return ptr
}

// OptBool creates a new parser setting to load a bool value from the command line only.
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptBool(long, short string, defaultValue bool, help string) *bool {
	ptr := new(bool)
	val := variant{ptrValue: ptr, defaultValue: defaultValue, isBool: true}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	return ptr
}

// OptInt creates a new parser setting to load a int value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptInt(long, short string, defaultValue int, help string) *int {
	ptr := new(int)
	val := variant{ptrValue: ptr, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	return ptr
}

// OptInt8 creates a new parser setting to load a int8 value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptInt8(long, short string, defaultValue int8, help string) *int8 {
	ptr := new(int8)
	val := variant{ptrValue: ptr, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	return ptr
}

// OptInt16 creates a new parser setting to load a int value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptInt16(long, short string, defaultValue int16, help string) *int16 {
	ptr := new(int16)
	val := variant{ptrValue: ptr, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	return ptr
}

// OptInt32 creates a new parser setting to load a int32 value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptInt32(long, short string, defaultValue int32, help string) *int32 {
	ptr := new(int32)
	val := variant{ptrValue: ptr, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	return ptr
}

// OptInt64 creates a new parser setting to load a int64 value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptInt64(long, short string, defaultValue int64, help string) *int64 {
	ptr := new(int64)
	val := variant{ptrValue: ptr, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	return ptr
}

// OptUint creates a new parser setting to load a uint value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptUint(long, short string, defaultValue uint, help string) *uint {
	ptr := new(uint)
	val := variant{ptrValue: ptr, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	return ptr
}

// OptUint8 creates a new parser setting to load a uint8 value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptUint8(long, short string, defaultValue uint8, help string) *uint8 {
	ptr := new(uint8)
	val := variant{ptrValue: ptr, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	return ptr
}

// OptUint16 creates a new parser setting to load a uint16 value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptUint16(long, short string, defaultValue uint16, help string) *uint16 {
	ptr := new(uint16)
	val := variant{ptrValue: ptr, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	return ptr
}

// OptUint32 creates a new parser setting to load a uint32 value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptUint32(long, short string, defaultValue uint32, help string) *uint32 {
	ptr := new(uint32)
	val := variant{ptrValue: ptr, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	return ptr
}

// OptUint64 creates a new parser setting to load a uint64 value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptUint64(long, short string, defaultValue uint64, help string) *uint64 {
	ptr := new(uint64)
	val := variant{ptrValue: ptr, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	return ptr
}
