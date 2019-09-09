package libcfg

import (
	"os"
)

// parser that can load configurations from the command line or the environment variables.
type cfgParser struct {
	commands   map[string]*commandImpl
	args       []string
	optentries []*optEntry
	shortopt   map[string]*optEntry
	longopt    map[string]*optEntry
	envLoader  *envLoaderImpl
}

func makeCfgParser() *cfgParser {
	return &cfgParser{
		commands:   make(map[string]*commandImpl),
		optentries: make([]*optEntry, 0),
		shortopt:   make(map[string]*optEntry),
		longopt:    make(map[string]*optEntry),
		envLoader:  makeEnvLoader(),
	}
}

// NewParser returns a new cfgParser, ready to be used.
func NewParser() RootParser {
	p := makeCfgParser()
	p.Configure(Options{})

	return p
}

func newSubParser(original *cfgParser) *cfgParser {
	return makeCfgParser()
}

func (cfg *cfgParser) Run() error {
	return cfg.RunArgs(os.Args[1:])
}

func (cfg *cfgParser) RunArgs(args []string) error {
	cfg.envLoader.LoadAll()

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

func (cfg *cfgParser) Configure(options Options) {
	cfg.envLoader.UseEnv(!options.FilesOnly)
}

func (cfg *cfgParser) UseFile(envfile string) error {
	return cfg.envLoader.UseFile(envfile)
}

func (cfg *cfgParser) UseFiles(envfiles ...string) {
	cfg.envLoader.UseFiles(envfiles...)
}
