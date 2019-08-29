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

		if _, ok := entry.valuePtr.(*bool); ok {
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
	val := new(string)
	cfg.addOpt(&optEntry{long: long, short: short, help: help, defaultValue: defaultValue, valuePtr: val, isStr: true})
	return val
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

// OptInt8 creates a new parser setting to load a int8 value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptInt8(long, short string, defaultValue int8, help string) *int8 {
	val := new(int8)
	cfg.addOpt(&optEntry{long: long, short: short, help: help, defaultValue: defaultValue, valuePtr: val})
	return val
}

// OptInt16 creates a new parser setting to load a int value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptInt16(long, short string, defaultValue int16, help string) *int16 {
	val := new(int16)
	cfg.addOpt(&optEntry{long: long, short: short, help: help, defaultValue: defaultValue, valuePtr: val})
	return val
}

// OptInt32 creates a new parser setting to load a int32 value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptInt32(long, short string, defaultValue int32, help string) *int32 {
	val := new(int32)
	cfg.addOpt(&optEntry{long: long, short: short, help: help, defaultValue: defaultValue, valuePtr: val})
	return val
}

// OptInt64 creates a new parser setting to load a int64 value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptInt64(long, short string, defaultValue int64, help string) *int64 {
	val := new(int64)
	cfg.addOpt(&optEntry{long: long, short: short, help: help, defaultValue: defaultValue, valuePtr: val})
	return val
}

// OptUint creates a new parser setting to load a uint value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptUint(long, short string, defaultValue uint, help string) *uint {
	val := new(uint)
	cfg.addOpt(&optEntry{long: long, short: short, help: help, defaultValue: defaultValue, valuePtr: val})
	return val
}

// OptUint8 creates a new parser setting to load a uint8 value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptUint8(long, short string, defaultValue uint8, help string) *uint8 {
	val := new(uint8)
	cfg.addOpt(&optEntry{long: long, short: short, help: help, defaultValue: defaultValue, valuePtr: val})
	return val
}

// OptUint16 creates a new parser setting to load a uint16 value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptUint16(long, short string, defaultValue uint16, help string) *uint16 {
	val := new(uint16)
	cfg.addOpt(&optEntry{long: long, short: short, help: help, defaultValue: defaultValue, valuePtr: val})
	return val
}

// OptUint32 creates a new parser setting to load a uint32 value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptUint32(long, short string, defaultValue uint32, help string) *uint32 {
	val := new(uint32)
	cfg.addOpt(&optEntry{long: long, short: short, help: help, defaultValue: defaultValue, valuePtr: val})
	return val
}

// OptUint64 creates a new parser setting to load a uint64 value from the command line only
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) OptUint64(long, short string, defaultValue uint64, help string) *uint64 {
	val := new(uint64)
	cfg.addOpt(&optEntry{long: long, short: short, help: help, defaultValue: defaultValue, valuePtr: val})
	return val
}
