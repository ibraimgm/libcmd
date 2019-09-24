package libcmd

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
	configured  bool
	options     Options
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

func (cmd *Cmd) configure() {
	if cmd.configured {
		return
	}

	if cmd.parentCmd != nil {
		cmd.parentCmd.configure()
		cmd.options = cmd.parentCmd.options
	}

	cmd.setupHelp()
	cmd.configured = true
}

func (cmd *Cmd) setupHelp() {
	if cmd.options.OnHelp == nil {
		cmd.options.OnHelp = automaticHelp
	}

	// no automatic '-h' flag
	if cmd.options.SuppressHelpFlag {
		return
	}

	if (len(cmd.optentries) > 0 || len(cmd.commands) > 0) && cmd.shortopt["-h"] == nil {
		cmd.Bool("help", "h", false, "Show this help message.")
	}
}
