package libcmd

import (
	"io"
	"os"
)

// CmdCallback is a routine that runs when new command is configured.
// If the command don't need to define new arguments or subcommands, you
// can safely run your custom code directly in this callback.
type CmdCallback func(cmd *Cmd)

// MatchCallback is a callback that is run when a command matches and
// and the parsing happens without raising an error. You can safely assume that
// the argument values are correctly loaded; however, since this callback runs
// before the subcommand processing occurs, this may not be the 'final' command
// to be executed by the user.
//
// To run code when the final command is selected, take a look at the
// Run function and the RunCallback type.
type MatchCallback func()

// RunCallback is a callback that runs when a specified command is invoked. This
// callback only runs if the parsing did succeed.
//
// You can safely assume that all parsed args are correctly loaded and no other
// subcommands are pending. If you return an error value in this callback, it will be
// available as a return to the Run or RunArgs on the App instance.
type RunCallback func() error

// ErrCallback is a callback that executes when the parser encounters an error.
// If you wish to recover (or ignore) the error, return a nil value to force the parser
// to continue it's normal process.
type ErrCallback func(err error) error

// HelpCallback is a callback that allows the user to customize the help text.
type HelpCallback func(cmd *Cmd, out io.Writer)

// Options defines the configuration options of the main
// App instance. The zero value of this struct reflect the default
// set of options.
type Options struct {
	// When true, the automatic creation of help flags will
	// be suppressed.
	SuppressHelpFlag bool

	// When true, do not print the help automatically when a
	// help flag is set
	SupressPrintHelpWhenSet bool

	// when true, do not print the help automatically when a command with
	// subcommands and without a Run callback is executed
	SuppressPrintHelpPartialCommand bool

	// When set, redirect the help output to the specified writer.
	// When it is nil, the help text will be printed to Stdout
	HelpOutput io.Writer

	// function that overrides the auto-generated help text.
	OnHelp HelpCallback
}

// App defines the main application parser.
// An application can define one or more command-line arguments to parse, as well
// as define a chain of subcommands supported by the application.
//
// To get a new instance of App, use the NewApp function.
type App struct {
	*Cmd
}

// NewApp returns a new instance of an app parser.
func NewApp(name, brief string) *App {
	var app App

	app.Cmd = newCmd()
	app.Name = name
	app.Brief = brief
	app.options = Options{}

	return &app
}

// Run runs the parser, collecting the values from the command-line
// and running any needed subcommands.
func (app *App) Run() error {
	return app.RunArgs(os.Args[1:])
}

// RunArgs behave like run, but instead of looking to the command-line
// arguments, it takes an array of arguments as parameters.
func (app *App) RunArgs(args []string) error {
	return app.doRun(args)
}

// Configure change the App settings to the specified options object.
// Please see the Options struct documentation for details.
func (app *App) Configure(options Options) {
	app.options = options
}
