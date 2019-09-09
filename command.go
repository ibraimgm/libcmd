package libcfg

type commandImpl struct {
	*cfgParser

	used         bool
	name         string
	description  string
	parentParser *cfgParser
}

func (cmd *commandImpl) Used() bool {
	return cmd.used
}

func (cmd *commandImpl) AddCommand(name string, description string) Command {
	return cmd.cfgParser.AddCommand(name, description)
}

func (cmd *commandImpl) Inherit(name string) {
	entry := cmd.parentParser.findByName(name)

	// entry does not exist: ignore
	if entry == nil {
		return
	}

	cmd.cfgParser.addOpt(entry)

	// check if we need to add it to env too
	if env := cmd.parentParser.envLoader.findEnvEntryByPtr(entry.val); env != nil {
		cmd.cfgParser.envLoader.addEnv(env.val, env.names)
	}
}

func (cmd *commandImpl) InheritAll() {
	for _, entry := range cmd.parentParser.optentries {
		opt := entry.long

		if opt == "" {
			opt = entry.short
		}

		cmd.Inherit(opt)
	}
}

func (cmd *commandImpl) doRun(args []string) ([]string, error) {
	cmd.used = true

	if err := cmd.cfgParser.RunArgs(args); err != nil {
		return args, err
	}

	return cmd.cfgParser.Args(), nil
}

// AddCommand adds a new parser command.
//
// When parsing the command-line, the first non-arg, non-value name found
// terminates the parser. If the parser has commands, the remaining args
// are passed to the command. This cycle will repeat itself until the entire
// command-line is parsed or an error is returned.
//
// You can check if a given command was used by looking at the return value
// of the Used method. If this method returns true, you can assume that the
// corresponding values are loaded in the pointer values.
func (cfg *cfgParser) AddCommand(name string, description string) Command {
	subparser := newSubParser(cfg)

	cmd := &commandImpl{
		parentParser: cfg,
		cfgParser:    subparser,
		name:         name,
		description:  description,
		used:         false,
	}

	cfg.commands[cmd.name] = cmd

	return cmd
}
