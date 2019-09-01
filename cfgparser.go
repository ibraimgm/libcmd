package libcfg

// CfgParser is a parser that can load configurations from the command line or the environment
// variables.
type CfgParser struct {
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
		optentries: make([]*optEntry, 0),
		shortopt:   make(map[string]*optEntry),
		longopt:    make(map[string]*optEntry),
		useEnv:     true,
		enventries: make([]*envEntry, 0),
		envvars:    make(map[string]string),
	}
}

// RunArgs loads all the configuration values according to
// the settings of the current parser, but assumes the values
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

// RunEnv loads all the configuration values of the env options,
// using the current environment variables.
func (cfg *CfgParser) RunEnv() error {
	cfg.loadFromEnv()
	return nil
}
