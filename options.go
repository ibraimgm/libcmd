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
}
