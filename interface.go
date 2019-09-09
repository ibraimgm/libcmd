package libcfg

// EnvLoader is a parser that can load values from environment variables.
//
// To load a value from an environment variable, use the method signature
// equivalent to the desired type (for example, use 'String' to load a string
// value, 'Int' to load an int value and so on) and store the returned pointer.
//
// After a call do LoadAll, the values will be available to use. If you want
// to provide the pointer yourself, use the methods with a 'P' suffix.
//
// You can pass more than one environment variable name to these functions;
// they will be processed in order, and the last value will override the
// others.
//
// To better control the loading of values from environment, please take a
// look at the functions UseEnv, UseFile and UseFiles.
type EnvLoader interface {

	// UseEnv sets whether the parser should consider environment variables
	// when loading env values. By default, when you define environment options
	// (for example, with String) the value to load is searched in the current
	// environment variables.
	//
	// If you call UseEnv(false) this behavior is disabled, and the only source of
	// values for environment variables will be calls to UseFile or UseFiles.
	//
	// This allows the usage of the the cfgParser with env values as a simple
	// file config loader.
	UseEnv(shouldUse bool)

	// UseFile loads the content of envfile into a in-memory environment value cache.
	// Each line of envfile is in the format 'key=value'. Comments start with '#' and
	// blank lines are ignored.
	//
	// You can call UseFile more than once, and newer values will override older values.
	// Values loaded with UseFile will always have higher priority than actual environment
	// values.
	UseFile(envfile string) error

	// UseFiles is a shortcut to multiple calls to UseFile.
	// Errors are ignored, so this funcion is useful to load a list of files when is acceptable
	// that some of them might not exist.
	UseFiles(envfiles ...string)

	// LoadAll loads all environment values in the associated pointers
	LoadAll()

	StringP(target *string, defaultValue string, variables ...string)
	BoolP(target *bool, defaultValue bool, variables ...string)
	IntP(target *int, defaultValue int, variables ...string)
	Int8P(target *int8, defaultValue int8, variables ...string)
	Int16P(target *int16, defaultValue int16, variables ...string)
	Int32P(target *int32, defaultValue int32, variables ...string)
	Int64P(target *int64, defaultValue int64, variables ...string)
	UintP(target *uint, defaultValue uint, variables ...string)
	Uint8P(target *uint8, defaultValue uint8, variables ...string)
	Uint16P(target *uint16, defaultValue uint16, variables ...string)
	Uint32P(target *uint32, defaultValue uint32, variables ...string)
	Uint64P(target *uint64, defaultValue uint64, variables ...string)
	Float32P(target *float32, defaultValue float32, variables ...string)
	Float64P(target *float64, defaultValue float64, variables ...string)
	String(defaultValue string, variables ...string) *string
	Bool(defaultValue bool, variables ...string) *bool
	Int(defaultValue int, variables ...string) *int
	Int8(defaultValue int8, variables ...string) *int8
	Int16(defaultValue int16, variables ...string) *int16
	Int32(defaultValue int32, variables ...string) *int32
	Int64(defaultValue int64, variables ...string) *int64
	Uint(defaultValue uint, variables ...string) *uint
	Uint8(defaultValue uint8, variables ...string) *uint8
	Uint16(defaultValue uint16, variables ...string) *uint16
	Uint32(defaultValue uint32, variables ...string) *uint32
	Uint64(defaultValue uint64, variables ...string) *uint64
	Float32(defaultValue float32, variables ...string) *float32
	Float64(defaultValue float64, variables ...string) *float64
}

// OptParser allows the loading of values from command-line arguments
// and optionally from environment variables
//
// To load a value from the command-line, use the method signature
// equivalent to the desired type (for example, use 'String' to load a string
// value, 'Int' to load an int value and so on) and store the returned pointer.
// If you also want to consider environment variables, specify one or more values
// to the variables parameter. If you want to provide the pointer yourself,
// use the methods with a 'P' suffix.
//
// Note that this interface does not allow the parser to run manually; to be able
// to configure and run a parser, you need an instance of RootParser.
//
// Environment values are loaded first, and (possibly) overridden by command-line
// flags. This allows the application developer to neatly define a configuration that
// can be fixed by the user with environment variables, but can be easily override
// on the fly by command-line options.
type OptParser interface {
	// Configure change the current parser configuration to the provided options value.
	// See Options for details.
	Configure(options Options)

	UseFile(envfile string) error
	UseFiles(envfiles ...string)

	StringP(target *string, long, short, defaultValue, help string, variables ...string)
	BoolP(target *bool, long, short string, defaultValue bool, help string, variables ...string)
	IntP(target *int, long, short string, defaultValue int, help string, variables ...string)
	Int8P(target *int8, long, short string, defaultValue int8, help string, variables ...string)
	Int16P(target *int16, long, short string, defaultValue int16, help string, variables ...string)
	Int32P(target *int32, long, short string, defaultValue int32, help string, variables ...string)
	Int64P(target *int64, long, short string, defaultValue int64, help string, variables ...string)
	UintP(target *uint, long, short string, defaultValue uint, help string, variables ...string)
	Uint8P(target *uint8, long, short string, defaultValue uint8, help string, variables ...string)
	Uint16P(target *uint16, long, short string, defaultValue uint16, help string, variables ...string)
	Uint32P(target *uint32, long, short string, defaultValue uint32, help string, variables ...string)
	Uint64P(target *uint64, long, short string, defaultValue uint64, help string, variables ...string)
	Float32P(target *float32, long, short string, defaultValue float32, help string, variables ...string)
	Float64P(target *float64, long, short string, defaultValue float64, help string, variables ...string)
	String(long, short, defaultValue, help string, variables ...string) *string
	Bool(long, short string, defaultValue bool, help string, variables ...string) *bool
	Int(long, short string, defaultValue int, help string, variables ...string) *int
	Int8(long, short string, defaultValue int8, help string, variables ...string) *int8
	Int16(long, short string, defaultValue int16, help string, variables ...string) *int16
	Int32(long, short string, defaultValue int32, help string, variables ...string) *int32
	Int64(long, short string, defaultValue int64, help string, variables ...string) *int64
	Uint(long, short string, defaultValue uint, help string, variables ...string) *uint
	Uint8(long, short string, defaultValue uint8, help string, variables ...string) *uint8
	Uint16(long, short string, defaultValue uint16, help string, variables ...string) *uint16
	Uint32(long, short string, defaultValue uint32, help string, variables ...string) *uint32
	Uint64(long, short string, defaultValue uint64, help string, variables ...string) *uint64
	Float32(long, short string, defaultValue float32, help string, variables ...string) *float32
	Float64(long, short string, defaultValue float64, help string, variables ...string) *float64
}

// RootParser is a parser that can be manually run.
//
// You can define values (like OptParser) and even load environment files (like and
// EnvLoader), and after a call do Run or RunArgs, the values will be available to use.
// Unparsed values are available by calling Args.
//
// It is also possible to add commands attached to this parser by
// using AddCommand.
type RootParser interface {
	OptParser

	// Run executes the parser and load all associated values
	Run() error

	// RunArgs behaves exactly like Run, but instead of pulling argument
	// options from the command line, you can pass the arguments to parse.
	// Note that the 'program name' (first argument when reading from command
	// line) should NOT be present.
	RunArgs(args []string) error

	// Returns the remaining unparsed args
	Args() []string

	// AddCommand creates a subcommand associated with this parser
	AddCommand(name string, description string) Command
}

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
type Command interface {
	OptParser

	Used() bool
	AddCommand(name string, description string) Command
	Inherit(optname string)
	InheritAll()
}
