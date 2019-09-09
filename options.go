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

	// Number of (expected) 'target' params.
	// Normally, when the parser finds a non-flag value, it stops the parsing. If
	// you put value in Targets, the parser instead 'skips' the 'N' values that would
	// finish the parser execution, collects theses values, and continue the parsing.
	//
	// This is useful, for example, in this scenario:
	// $ myapp file.md -o output.html
	//
	// In the example above, the parser would finish as soon as 'file.md' is encountered,
	// leaving the (valid) arguments unparsed. By setting Targets as 1, the first non-flag
	// argument is collected and the parsing continues.
	//
	// You can get the collected values using the Targets method.
	Targets uint
}
