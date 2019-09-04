package libcfg

// Command is a runnable (sub)command of the main parser.
//
// Commands have a similar interface to the parser, but cannot
// run independently from the associated parser.
//
// By default, commands have an empty context, separate from the
// original parser. However, you can 'inherit' (with Inherit or InheritAll)
// selected command line options of the parent parser or command,
// and reuse the same pointer.
// This allows the end user to write:
//
// myApp --foo=1 commandname --foo=2
//
// While this is not that useful for humans, it gives some leeway to
// machine-generated command-line options. If you want the same name,
// with a different meaning, just define you options normally.
//
// Note that you cannot inherit a 'env only' variable, but the cache of
// environment values collected is always cloned to the child command.
type Command interface {
	EnvOptParser

	Used() bool
	AddCommand(name string, description string) Command
	Inherit(optname string)
	InheritAll()
}

type commandImpl struct {
	*CfgParser

	used         bool
	name         string
	description  string
	parentParser *CfgParser
}

func (cmd *commandImpl) Used() bool {
	return cmd.used
}

func (cmd *commandImpl) AddCommand(name string, description string) Command {
	return cmd.CfgParser.AddCommand(name, description)
}

func (cmd *commandImpl) Inherit(name string) {
	entry := cmd.parentParser.findByName(name)

	// entry does not exist: ignore
	if entry == nil {
		return
	}

	cmd.CfgParser.addOpt(entry)

	// check if we need to add it to env too
	if env := cmd.parentParser.findEnvEntryByPtr(entry.val); env != nil {
		cmd.CfgParser.addEnv(env.val, env.names)
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

	if err := cmd.CfgParser.RunArgs(args); err != nil {
		return args, err
	}

	return cmd.CfgParser.Args(), nil
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
func (cfg *CfgParser) AddCommand(name string, description string) Command {
	subparser := newSubParser(cfg)

	cmd := &commandImpl{
		parentParser: cfg,
		CfgParser:    subparser,
		name:         name,
		description:  description,
		used:         false,
	}

	cfg.commands[cmd.name] = cmd

	return cmd
}
