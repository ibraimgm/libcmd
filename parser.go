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
