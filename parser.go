package libcfg

import (
	"fmt"
	"strconv"
	"strings"
)

type cfgEntry struct {
	long         string
	short        string
	help         string
	defaultValue interface{}
	isBool       bool
	valuePtr     interface{}
}

type CfgParser struct {
	entries []*cfgEntry
}

func NewParser() *CfgParser {
	return &CfgParser{
		entries: make([]*cfgEntry, 0),
	}
}

func (cfg *CfgParser) ParseArgs(args []string) error {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		var val string

		isLong := strings.HasPrefix(arg, "--")
		isShort := !isLong && strings.HasPrefix(arg, "-")
		var isLongWithValue bool

		// is not a param
		if !isShort && !isLong {
			return fmt.Errorf("unexpected argument: %s", arg)
		}

		// long param might use '='
		if isLong {
			splitted := strings.Split(arg, "=")
			if len(splitted) == 2 {
				arg = splitted[0]
				val = splitted[1]
				isLongWithValue = true
			}
		}

		// find the entry
		var entry *cfgEntry
		for j := 0; j < len(cfg.entries); j++ {
			short := "-" + cfg.entries[j].short
			long := "--" + cfg.entries[j].long
			nolong := "--no-" + cfg.entries[j].long

			if short == arg || long == arg {
				entry = cfg.entries[j]
				break
			}

			if cfg.entries[j].isBool && nolong == arg {
				entry = cfg.entries[j]
				val = "false"
				break
			}
		}

		if entry == nil {
			return fmt.Errorf("unknown argument: %s", arg)
		}

		if entry.isBool && val == "" {
			val = "true"
		}

		// short or a long without value, must look the next argument
		if (isShort || !isLongWithValue) && val == "" {
			if i+1 == len(args) {
				return fmt.Errorf("no value for argument: %s", arg)
			}

			val = args[i+1]
			i++
		}

		switch v := entry.valuePtr.(type) {
		case *bool:
			if bval, err := strconv.ParseBool(val); err != nil {
				return fmt.Errorf("'%v' is not a valid boolean value for argument '%s'", val, arg)
			} else {
				*v = bval
			}

		case *int:
			if ival, err := strconv.ParseInt(val, 10, 64); err != nil {
				return fmt.Errorf("'%v' is not a valid int value for argument '%s'", val, arg)
			} else {
				*v = int(ival)
			}

		case *string:
			*v = val

		default:
			return fmt.Errorf("unrecognized entry type: %T", entry.valuePtr)
		}
	}

	return nil
}

func (cfg *CfgParser) OptBool(long, short string, defaultValue bool, help string) *bool {
	val := new(bool)
	cfg.entries = append(cfg.entries, &cfgEntry{long: long, short: short, help: help, defaultValue: defaultValue, valuePtr: val, isBool: true})
	return val
}

func (cfg *CfgParser) OptInt(long, short string, defaultValue int, help string) *int {
	val := new(int)
	cfg.entries = append(cfg.entries, &cfgEntry{long: long, short: short, help: help, defaultValue: defaultValue, valuePtr: val})
	return val
}

func (cfg *CfgParser) OptString(long, short, defaultValue, help string) *string {
	val := new(string)
	cfg.entries = append(cfg.entries, &cfgEntry{long: long, short: short, help: help, defaultValue: defaultValue, valuePtr: val})
	return val
}
