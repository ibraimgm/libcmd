package libcfg

// Options defines the parametersof a parser and/or command
//
// An empty options value represents the default parser configuration
type Options struct {

	// When true, the environment values should only consider files informed
	// with UseFile or UseFiles, and should not look for the actual environment
	// vars.
	//
	// If setting it to false on a subcommand, when you inherit a configuration,
	// the value might have already been set by a parent that allows the use
	// of env vars.
	FilesOnly bool

	// When true, if any args remain unparsed, the parser will return an error.
	// By default, the parser only returns an error if an unknown argument is passed.
	StrictParsing bool

	// When true, continue the parsing process until an error, a recognized command or the
	// end of command-line options is reached.
	//
	// Normally, when the parser finds a non-flag value, it stops the parsing. When this option
	// is set to true, the parser will 'ignore' the value and continue to parse the other options;
	// the 'ignored' arguments are collected and (re)inserted at the start of Args().
	//
	// This is useful, for example, in this scenario:
	// $ myapp file.md -o output.html
	//
	// In the example above, the parser would finish as soon as 'file.md' is encountered,
	// leaving the (valid) arguments unparsed. By setting Greedy as true, the first non-flag
	// argument is collected and the parsing continues.
	Greedy bool
}
