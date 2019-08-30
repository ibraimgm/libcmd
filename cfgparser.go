package libcfg

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
