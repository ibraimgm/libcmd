package libcfg

// OptParser is a parser that can load values from command-line options.
//
// Opt[type] functions create a setting to load a value of the associated
// type from the command-line parameters only.
// The functions with a 'Var' suffix expect the user to provide a existing
// pointer, while the functions without the suffix return a new pointer to
// the associated type.
// In either case, the value will be available to the pointer after the
// parsing.
//
// IMPORTANT: this interface is used for documentation purposes only. The actual
// parser does implements the full spectrum of the EnvOpt* methods. Client code
// is not expected to implement this interface.
type OptParser interface {
	OptStringVar(target *string, long, short, defaultValue, help string)
	OptBoolVar(target *bool, long, short string, defaultValue bool, help string)
	OptIntVar(target *int, long, short string, defaultValue int, help string)
	OptInt8Var(target *int8, long, short string, defaultValue int8, help string)
	OptInt16Var(target *int16, long, short string, defaultValue int16, help string)
	OptInt32Var(target *int32, long, short string, defaultValue int32, help string)
	OptInt64Var(target *int64, long, short string, defaultValue int64, help string)
	OptUintVar(target *uint, long, short string, defaultValue uint, help string)
	OptUint8Var(target *uint8, long, short string, defaultValue uint8, help string)
	OptUint16Var(target *uint16, long, short string, defaultValue uint16, help string)
	OptUint32Var(target *uint32, long, short string, defaultValue uint32, help string)
	OptUint64Var(target *uint64, long, short string, defaultValue uint64, help string)
	OptFloat32Var(target *float32, long, short string, defaultValue float32, help string)
	OptFloat64Var(target *float64, long, short string, defaultValue float64, help string)
	OptString(long, short, defaultValue, help string) *string
	OptBool(long, short string, defaultValue bool, help string) *bool
	OptInt(long, short string, defaultValue int, help string) *int
	OptInt8(long, short string, defaultValue int8, help string) *int8
	OptInt16(long, short string, defaultValue int16, help string) *int16
	OptInt32(long, short string, defaultValue int32, help string) *int32
	OptInt64(long, short string, defaultValue int64, help string) *int64
	OptUint(long, short string, defaultValue uint, help string) *uint
	OptUint8(long, short string, defaultValue uint8, help string) *uint8
	OptUint16(long, short string, defaultValue uint16, help string) *uint16
	OptUint32(long, short string, defaultValue uint32, help string) *uint32
	OptUint64(long, short string, defaultValue uint64, help string) *uint64
	OptFloat32(long, short string, defaultValue float32, help string) *float32
	OptFloat64(long, short string, defaultValue float64, help string) *float64
}

// EnvParser is a parser that can load values from environment variables.
//
// Env[type] functions create a setting to load a value of the associated
// type from the environment variables only.
// The functions with a 'Var' suffix expect the user to provide a existing
// pointer, while the functions without the suffix return a new pointer to
// the associated type. In either case, the value will be available to
// the pointer after the parsing.
//
// You can pass more than one environment variable name to these functions;
// they will be processed in order, and the last value will override the
// others.
//
// To better control the loading of values from environment, please take a
// look at the functions UseEnv, UseFile and UseFiles.
//
// IMPORTANT: this interface is used for documentation purposes only. The actual
// parser does implements the full spectrum of the EnvOpt* methods. Client code
// is not expected to implement this interface.
type EnvParser interface {
	UseEnv(shouldUse bool)
	UseFile(envfile string) error
	UseFiles(envfiles ...string)

	EnvStringVar(target *string, defaultValue string, variables ...string)
	EnvBoolVar(target *bool, defaultValue bool, variables ...string)
	EnvIntVar(target *int, defaultValue int, variables ...string)
	EnvInt8Var(target *int8, defaultValue int8, variables ...string)
	EnvInt16Var(target *int16, defaultValue int16, variables ...string)
	EnvInt32Var(target *int32, defaultValue int32, variables ...string)
	EnvInt64Var(target *int64, defaultValue int64, variables ...string)
	EnvUintVar(target *uint, defaultValue uint, variables ...string)
	EnvUint8Var(target *uint8, defaultValue uint8, variables ...string)
	EnvUint16Var(target *uint16, defaultValue uint16, variables ...string)
	EnvUint32Var(target *uint32, defaultValue uint32, variables ...string)
	EnvUint64Var(target *uint64, defaultValue uint64, variables ...string)
	EnvFloat32Var(target *float32, defaultValue float32, variables ...string)
	EnvFloat64Var(target *float64, defaultValue float64, variables ...string)
	EnvString(defaultValue string, variables ...string) *string
	EnvBool(defaultValue bool, variables ...string) *bool
	EnvInt(defaultValue int, variables ...string) *int
	EnvInt8(defaultValue int8, variables ...string) *int8
	EnvInt16(defaultValue int16, variables ...string) *int16
	EnvInt32(defaultValue int32, variables ...string) *int32
	EnvInt64(defaultValue int64, variables ...string) *int64
	EnvUint(defaultValue uint, variables ...string) *uint
	EnvUint8(defaultValue uint8, variables ...string) *uint8
	EnvUint16(defaultValue uint16, variables ...string) *uint16
	EnvUint32(defaultValue uint32, variables ...string) *uint32
	EnvUint64(defaultValue uint64, variables ...string) *uint64
	EnvFloat32(defaultValue float32, variables ...string) *float32
	EnvFloat64(defaultValue float64, variables ...string) *float64
}

// EnvOptParser combines the funcionality of a OptParser and a EnvParser,
// but also add extra helper functions to easily define settings to load
// both environment and command-line parser
//
// EnvOpt[type] functions create a setting to load a value of the associated
// type from both the environment variables and the command-line, in this order.
//
// This is equivalent to using an Env* and Opt* function in sequence, with the
// same pointer. The same logic of the 'Var' suffix, and the uses of environment
// variables is applyed here.
//
// Environment values are loaded first, and (possibly) overridden by command-line
// flags. This allows the application developer to neatly define a configuration that
// can be fixed by the user with environment variables, but can be easily override
// on the fly by command-line options.
//
// IMPORTANT: this interface is used for documentation purposes only. The actual
// parser does implements the full spectrum of the EnvOpt* methods. Client code
// is not expected to implement this interface.
type EnvOptParser interface {
	OptParser
	EnvParser

	EnvOptStringVar(target *string, long, short, defaultValue, help string, variables ...string)
	EnvOptBoolVar(target *bool, long, short string, defaultValue bool, help string, variables ...string)
	EnvOptIntVar(target *int, long, short string, defaultValue int, help string, variables ...string)
	EnvOptInt8Var(target *int8, long, short string, defaultValue int8, help string, variables ...string)
	EnvOptInt16Var(target *int16, long, short string, defaultValue int16, help string, variables ...string)
	EnvOptInt32Var(target *int32, long, short string, defaultValue int32, help string, variables ...string)
	EnvOptInt64Var(target *int64, long, short string, defaultValue int64, help string, variables ...string)
	EnvOptUintVar(target *uint, long, short string, defaultValue uint, help string, variables ...string)
	EnvOptUint8Var(target *uint8, long, short string, defaultValue uint8, help string, variables ...string)
	EnvOptUint16Var(target *uint16, long, short string, defaultValue uint16, help string, variables ...string)
	EnvOptUint32Var(target *uint32, long, short string, defaultValue uint32, help string, variables ...string)
	EnvOptUint64Var(target *uint64, long, short string, defaultValue uint64, help string, variables ...string)
	EnvOptFloat32Var(target *float32, long, short string, defaultValue float32, help string, variables ...string)
	EnvOptFloat64Var(target *float64, long, short string, defaultValue float64, help string, variables ...string)
	EnvOptString(long, short, defaultValue, help string, variables ...string) *string
	EnvOptBool(long, short string, defaultValue bool, help string, variables ...string) *bool
	EnvOptInt(long, short string, defaultValue int, help string, variables ...string) *int
	EnvOptInt8(long, short string, defaultValue int8, help string, variables ...string) *int8
	EnvOptInt16(long, short string, defaultValue int16, help string, variables ...string) *int16
	EnvOptInt32(long, short string, defaultValue int32, help string, variables ...string) *int32
	EnvOptInt64(long, short string, defaultValue int64, help string, variables ...string) *int64
	EnvOptUint(long, short string, defaultValue uint, help string, variables ...string) *uint
	EnvOptUint8(long, short string, defaultValue uint8, help string, variables ...string) *uint8
	EnvOptUint16(long, short string, defaultValue uint16, help string, variables ...string) *uint16
	EnvOptUint32(long, short string, defaultValue uint32, help string, variables ...string) *uint32
	EnvOptUint64(long, short string, defaultValue uint64, help string, variables ...string) *uint64
	EnvOptFloat32(long, short string, defaultValue float32, help string, variables ...string) *float32
	EnvOptFloat64(long, short string, defaultValue float64, help string, variables ...string) *float64
}
