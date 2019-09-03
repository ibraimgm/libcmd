package libcfg

// EnvOptStringVar creates a new parser setting to load a string value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the provided pointer.
func (cfg *CfgParser) EnvOptStringVar(target *string, long, short, defaultValue, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue, isStr: true}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.addEnv(&val, variables)
}

// EnvOptBoolVar creates a new parser setting to load a bool value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the provided pointer.
func (cfg *CfgParser) EnvOptBoolVar(target *bool, long, short string, defaultValue bool, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue, isBool: true}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.addEnv(&val, variables)
}

// EnvOptIntVar creates a new parser setting to load a int value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the provided pointer.
func (cfg *CfgParser) EnvOptIntVar(target *int, long, short string, defaultValue int, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.addEnv(&val, variables)
}

// EnvOptInt8Var creates a new parser setting to load a int8 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the provided pointer.
func (cfg *CfgParser) EnvOptInt8Var(target *int8, long, short string, defaultValue int8, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.addEnv(&val, variables)
}

// EnvOptInt16Var creates a new parser setting to load a int16 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the provided pointer.
func (cfg *CfgParser) EnvOptInt16Var(target *int16, long, short string, defaultValue int16, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.addEnv(&val, variables)
}

// EnvOptInt32Var creates a new parser setting to load a int32 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the provided pointer.
func (cfg *CfgParser) EnvOptInt32Var(target *int32, long, short string, defaultValue int32, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.addEnv(&val, variables)
}

// EnvOptInt64Var creates a new parser setting to load a int64 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the provided pointer.
func (cfg *CfgParser) EnvOptInt64Var(target *int64, long, short string, defaultValue int64, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.addEnv(&val, variables)
}

// EnvOptUintVar creates a new parser setting to load a uint value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the provided pointer.
func (cfg *CfgParser) EnvOptUintVar(target *uint, long, short string, defaultValue uint, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.addEnv(&val, variables)
}

// EnvOptUint8Var creates a new parser setting to load a uint8 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the provided pointer.
func (cfg *CfgParser) EnvOptUint8Var(target *uint8, long, short string, defaultValue uint8, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.addEnv(&val, variables)
}

// EnvOptUint16Var creates a new parser setting to load a uint16 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the provided pointer.
func (cfg *CfgParser) EnvOptUint16Var(target *uint16, long, short string, defaultValue uint16, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.addEnv(&val, variables)
}

// EnvOptUint32Var creates a new parser setting to load a uint32 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the provided pointer.
func (cfg *CfgParser) EnvOptUint32Var(target *uint32, long, short string, defaultValue uint32, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.addEnv(&val, variables)
}

// EnvOptUint64Var creates a new parser setting to load a uint64 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the provided pointer.
func (cfg *CfgParser) EnvOptUint64Var(target *uint64, long, short string, defaultValue uint64, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.addEnv(&val, variables)
}

// EnvOptFloat32Var creates a new parser setting to load a float32 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the provided pointer.
func (cfg *CfgParser) EnvOptFloat32Var(target *float32, long, short string, defaultValue float32, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.addEnv(&val, variables)
}

// EnvOptFloat64Var creates a new parser setting to load a float64 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the provided pointer.
func (cfg *CfgParser) EnvOptFloat64Var(target *float64, long, short string, defaultValue float64, help string, variables ...string) {
	val := variant{ptrValue: target, defaultValue: defaultValue}
	cfg.addOpt(&optEntry{long: long, short: short, help: help, val: &val})
	cfg.addEnv(&val, variables)
}

// EnvOptString creates a new parser setting to load a string value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the returned pointer.
func (cfg *CfgParser) EnvOptString(long, short, defaultValue, help string, variables ...string) *string {
	target := new(string)
	cfg.EnvOptStringVar(target, long, short, defaultValue, help, variables...)
	return target
}

// EnvOptBool creates a new parser setting to load a bool value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the returned pointer.
func (cfg *CfgParser) EnvOptBool(long, short string, defaultValue bool, help string, variables ...string) *bool {
	target := new(bool)
	cfg.EnvOptBoolVar(target, long, short, defaultValue, help, variables...)
	return target
}

// EnvOptInt creates a new parser setting to load a int value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the returned pointer.
func (cfg *CfgParser) EnvOptInt(long, short string, defaultValue int, help string, variables ...string) *int {
	target := new(int)
	cfg.EnvOptIntVar(target, long, short, defaultValue, help, variables...)
	return target
}

// EnvOptInt8 creates a new parser setting to load a int8 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the returned pointer.
func (cfg *CfgParser) EnvOptInt8(long, short string, defaultValue int8, help string, variables ...string) *int8 {
	target := new(int8)
	cfg.EnvOptInt8Var(target, long, short, defaultValue, help, variables...)
	return target
}

// EnvOptInt16 creates a new parser setting to load a int16 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the returned pointer.
func (cfg *CfgParser) EnvOptInt16(long, short string, defaultValue int16, help string, variables ...string) *int16 {
	target := new(int16)
	cfg.EnvOptInt16Var(target, long, short, defaultValue, help, variables...)
	return target
}

// EnvOptInt32 creates a new parser setting to load a int32 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the returned pointer.
func (cfg *CfgParser) EnvOptInt32(long, short string, defaultValue int32, help string, variables ...string) *int32 {
	target := new(int32)
	cfg.EnvOptInt32Var(target, long, short, defaultValue, help, variables...)
	return target
}

// EnvOptInt64 creates a new parser setting to load a int64 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the returned pointer.
func (cfg *CfgParser) EnvOptInt64(long, short string, defaultValue int64, help string, variables ...string) *int64 {
	target := new(int64)
	cfg.EnvOptInt64Var(target, long, short, defaultValue, help, variables...)
	return target
}

// EnvOptUint creates a new parser setting to load a uint value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the returned pointer.
func (cfg *CfgParser) EnvOptUint(long, short string, defaultValue uint, help string, variables ...string) *uint {
	target := new(uint)
	cfg.EnvOptUintVar(target, long, short, defaultValue, help, variables...)
	return target
}

// EnvOptUint8 creates a new parser setting to load a uint8 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the returned pointer.
func (cfg *CfgParser) EnvOptUint8(long, short string, defaultValue uint8, help string, variables ...string) *uint8 {
	target := new(uint8)
	cfg.EnvOptUint8Var(target, long, short, defaultValue, help, variables...)
	return target
}

// EnvOptUint16 creates a new parser setting to load a uint16 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the returned pointer.
func (cfg *CfgParser) EnvOptUint16(long, short string, defaultValue uint16, help string, variables ...string) *uint16 {
	target := new(uint16)
	cfg.EnvOptUint16Var(target, long, short, defaultValue, help, variables...)
	return target
}

// EnvOptUint32 creates a new parser setting to load a uint32 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the returned pointer.
func (cfg *CfgParser) EnvOptUint32(long, short string, defaultValue uint32, help string, variables ...string) *uint32 {
	target := new(uint32)
	cfg.EnvOptUint32Var(target, long, short, defaultValue, help, variables...)
	return target
}

// EnvOptUint64 creates a new parser setting to load a uint64 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the returned pointer.
func (cfg *CfgParser) EnvOptUint64(long, short string, defaultValue uint64, help string, variables ...string) *uint64 {
	target := new(uint64)
	cfg.EnvOptUint64Var(target, long, short, defaultValue, help, variables...)
	return target
}

// EnvOptFloat32 creates a new parser setting to load a float32 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the returned pointer.
func (cfg *CfgParser) EnvOptFloat32(long, short string, defaultValue float32, help string, variables ...string) *float32 {
	target := new(float32)
	cfg.EnvOptFloat32Var(target, long, short, defaultValue, help, variables...)
	return target
}

// EnvOptFloat64 creates a new parser setting to load a float64 value both
// from command-line and from the environment variables. The environment is checked first,
// and can be replaced with by the command-line options.
// The end value will be available on the returned pointer.
func (cfg *CfgParser) EnvOptFloat64(long, short string, defaultValue float64, help string, variables ...string) *float64 {
	target := new(float64)
	cfg.EnvOptFloat64Var(target, long, short, defaultValue, help, variables...)
	return target
}
