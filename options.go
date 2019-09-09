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
}
