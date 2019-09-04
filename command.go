package libcfg

// Command is a runnable (sub)command of the main parser.
//
// Commands have a similar interface to the parser, but cannot
// run independently from the associated parser.
type Command interface {
	EnvOptParser

	Used() bool
	AddCommand(name string, description string) Command
}

type commandImpl struct {
	*CfgParser

	used        bool
	name        string
	description string
}

func (cmd *commandImpl) Used() bool {
	return cmd.used
}

func (cmd *commandImpl) AddCommand(name string, description string) Command {
	return cmd.CfgParser.AddCommand(name, description)
}

func (cmd *commandImpl) doRun(args []string) ([]string, error) {
	if err := cmd.CfgParser.RunArgs(args); err != nil {
		return args, err
	}

	cmd.used = true

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
		CfgParser:   subparser,
		name:        name,
		description: description,
		used:        false,
	}

	cfg.commands[cmd.name] = cmd

	return cmd
}
