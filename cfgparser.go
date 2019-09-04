package libcfg

// CfgParser is a parser that can load configurations from the command line or the environment
// variables.
//
// CfgParser implements all the methods on EnvOptParser, and extra methods that allows the
// API user to run the parsing algorithm.
//
// You should never build an instance of CfgParser manually; to get a correctly configured and
// ready to use instance, use the NewParser function.
type CfgParser struct {
	commands   map[string]*commandImpl
	args       []string
	optentries []*optEntry
	shortopt   map[string]*optEntry
	longopt    map[string]*optEntry
	useEnv     bool
	enventries []*envEntry
	envvars    map[string]string
}

// NewParser returns a new CfgParser, ready to be used.
func NewParser() *CfgParser {
	return &CfgParser{
		commands:   make(map[string]*commandImpl),
		optentries: make([]*optEntry, 0),
		shortopt:   make(map[string]*optEntry),
		longopt:    make(map[string]*optEntry),
		useEnv:     true,
		enventries: make([]*envEntry, 0),
		envvars:    make(map[string]string),
	}
}

func newSubParser(original *CfgParser) *CfgParser {
	return NewParser()
}

// RunArgs loads all the configuration values according to
// the settings of the current parser, but assumes the values
// passed by args as command line arguments.
//
// Note that args must not include the program name
func (cfg *CfgParser) RunArgs(args []string) error {
	cfg.loadFromEnv()

	if err := cfg.doParse(args); err != nil {
		return err
	}

	for i := range cfg.optentries {
		cfg.optentries[i].val.useDefault()
	}

	if len(cfg.args) >= 1 {
		name := cfg.args[0]

		if cmd, ok := cfg.commands[name]; ok {
			newArgs, err := cmd.doRun(cfg.args[1:])

			cfg.args = newArgs
			return err
		}
	}

	return nil
}
