package libcmd

type operand struct {
	name     string
	modifier string
}

// Cmd defines a (sub)command of the application.
// Since commands cannot do much by themselves, you should create
// your commands by calling the Command method in the App instance.
//
// Subcommands can be created by calling Command of an existing command.
type Cmd struct {
	// The name to be used to invoke the command
	Name string

	// The brief description of the command
	Brief string

	// The long description of the command
	Long string

	// The text to be shown at the 'Usage' line of the help.
	// Set this value to "-" to omit the usage line
	Usage string

	// Options for this command
	Options Options

	args        []string
	optentries  []*optEntry
	shortopt    map[string]*optEntry
	longopt     map[string]*optEntry
	callback    CmdCallback
	match       MatchCallback
	run         RunCallback
	errHandler  ErrCallback
	breadcrumbs string
	commands    map[string]*Cmd
	parentCmd   *Cmd
	operands    []operand
}

func newCmd() *Cmd {
	return &Cmd{
		args:       make([]string, 0),
		optentries: make([]*optEntry, 0),
		shortopt:   make(map[string]*optEntry),
		longopt:    make(map[string]*optEntry),
		commands:   make(map[string]*Cmd),
	}
}

// Command defines a new subcommand. The name is mandatory
// (otherwise, the command will never run), and the callback function gives
// the caller the opportunity to configure the newly created command, by
// defining new flags, subcommands, or specifying what to run when
// the command is activated.
func (cmd *Cmd) Command(name, brief string, callback CmdCallback) {
	if name == "" {
		return
	}

	c := newCmd()
	c.Name = name
	c.Brief = brief
	c.callback = callback
	c.breadcrumbs = cmd.breadcrumbs + " " + cmd.Name
	c.parentCmd = cmd

	cmd.commands[c.Name] = c
}

// Match registers a callback to run when this command matches.
// A command matches if it is invoked and the parsing is successful.
func (cmd *Cmd) Match(callback MatchCallback) {
	cmd.match = callback
}

// Run registers a callback to run when this command is matched
// (see Match) and no more subcommands were invoked.
func (cmd *Cmd) Run(callback RunCallback) {
	cmd.run = callback
}

// Err registers a handler to be run when the parsing fails.
func (cmd *Cmd) Err(handler ErrCallback) {
	cmd.errHandler = handler
}

// CommandMatch is a shortcut to Command() followed by Match() on the
// provided command.
func (cmd *Cmd) CommandMatch(name, brief string, callback MatchCallback) {
	cmd.Command(name, brief, nil)
	if c, ok := cmd.commands[name]; ok {
		c.Match(callback)
	}
}

// CommandRun is a shortcut to Command() followed by Run() on the
// provided command.
func (cmd *Cmd) CommandRun(name, brief string, callback RunCallback) {
	cmd.Command(name, brief, nil)
	if c, ok := cmd.commands[name]; ok {
		c.Run(callback)
	}
}

// AddOperand documents an expected operand.
// The modifier parameter can be either '?' for optional operands or '*'
// for repeating ones. The documentation is printed in the order that
// was used to add the operands, so it is advisable to put them in an order
// that makes sense for the user (required, optional and repeating, in this order).
//
// Note that this modifier is used only for documentation purposes; no special validation
// is done, except by the one documented in Options.StrictOperands.
func (cmd *Cmd) AddOperand(name string, modifer string) {
	cmd.operands = append(cmd.operands, operand{name: name, modifier: modifer})
}

// Operand returns the value of the named operand, if any.
// When specified using AddOperand, each unparsed arg is considered an operand and it's
// value is fetched - but not consumed -  from the Args() method.
//
// The behavior of this function is only guaranteed when used in a 'leaf' command or
// and Run() callback.
func (cmd *Cmd) Operand(name string) string {
	for i, op := range cmd.operands {
		if op.name == name && i < len(cmd.args) {
			return cmd.args[i]
		}
	}

	return ""
}

func (cmd *Cmd) setupHelp() {
	// no automatic '-h' flag
	if cmd.Options.SuppressHelpFlag {
		return
	}

	if (len(cmd.optentries) > 0 || len(cmd.commands) > 0 || len(cmd.operands) > 0) && cmd.shortopt["-h"] == nil {
		cmd.Bool("help", 'h', false, "Show this help message.")
	}
}
