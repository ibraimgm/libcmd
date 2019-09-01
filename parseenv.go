package libcfg

import (
	"os"
)

type envEntry struct {
	names []string
	val   *variant
}

func (cfg *CfgParser) addEnv(val *variant, variables []string) {
	cfg.enventries = append(cfg.enventries, &envEntry{names: variables, val: val})
}

func (cfg *CfgParser) loadFromEnv() {
	for i := range cfg.enventries {
		entry := cfg.enventries[i]

		cfg.findEnvValue(entry)
	}
}

func (cfg *CfgParser) findEnvValue(entry *envEntry) {
	for _, name := range entry.names {
		var value string
		var ok bool

		if value, ok = cfg.envvars[name]; !ok {
			value, ok = os.LookupEnv(name)
		}

		switch {
		case ok && value != "":
			entry.val.setValue(value) //nolint: errheck
		case ok && value == "":
			entry.val.unsetValue()
		}
	}
}

// EnvStringVar creates a new parser setting to load a string value from the
// specified environment variables
// After parsing, the value will be available on the provided pointer.
func (cfg *CfgParser) EnvStringVar(target *string, defaultValue string, variables ...string) {
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue, isStr: true}, variables)
}

// EnvBoolVar creates a new parser setting to load a boolean value from the
// specified environment variables
// After parsing, the value will be available on the provided pointer.
func (cfg *CfgParser) EnvBoolVar(target *bool, defaultValue bool, variables ...string) {
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue, isBool: true}, variables)
}

// EnvIntVar creates a new parser setting to load a int value from the
// specified environment variables
// After parsing, the value will be available on the provided pointer.
func (cfg *CfgParser) EnvIntVar(target *int, defaultValue int, variables ...string) {
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
}

// EnvInt8Var creates a new parser setting to load a int8 value from the
// specified environment variables
// After parsing, the value will be available on the provided pointer.
func (cfg *CfgParser) EnvInt8Var(target *int8, defaultValue int8, variables ...string) {
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
}

// EnvInt16Var creates a new parser setting to load a int16 value from the
// specified environment variables
// After parsing, the value will be available on the provided pointer.
func (cfg *CfgParser) EnvInt16Var(target *int16, defaultValue int16, variables ...string) {
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
}

// EnvInt32Var creates a new parser setting to load a int32 value from the
// specified environment variables
// After parsing, the value will be available on the provided pointer.
func (cfg *CfgParser) EnvInt32Var(target *int32, defaultValue int32, variables ...string) {
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
}

// EnvInt64Var creates a new parser setting to load a int64 value from the
// specified environment variables
// After parsing, the value will be available on the provided pointer.
func (cfg *CfgParser) EnvInt64Var(target *int64, defaultValue int64, variables ...string) {
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
}

// EnvUintVar creates a new parser setting to load a uint value from the
// specified environment variables
// After parsing, the value will be available on the provided pointer.
func (cfg *CfgParser) EnvUintVar(target *uint, defaultValue uint, variables ...string) {
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
}

// EnvUint8Var creates a new parser setting to load a uint8 value from the
// specified environment variables
// After parsing, the value will be available on the provided pointer.
func (cfg *CfgParser) EnvUint8Var(target *uint8, defaultValue uint8, variables ...string) {
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
}

// EnvUint16Var creates a new parser setting to load a uint16 value from the
// specified environment variables
// After parsing, the value will be available on the provided pointer.
func (cfg *CfgParser) EnvUint16Var(target *uint16, defaultValue uint16, variables ...string) {
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
}

// EnvUint32Var creates a new parser setting to load a uint32 value from the
// specified environment variables
// After parsing, the value will be available on the provided pointer.
func (cfg *CfgParser) EnvUint32Var(target *uint32, defaultValue uint32, variables ...string) {
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
}

// EnvUint64Var creates a new parser setting to load a uint64 value from the
// specified environment variables
// After parsing, the value will be available on the provided pointer.
func (cfg *CfgParser) EnvUint64Var(target *uint64, defaultValue uint64, variables ...string) {
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
}

// EnvString creates a new parser setting to load a string value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvString(defaultValue string, variables ...string) *string {
	target := new(string)
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue, isStr: true}, variables)
	return target
}

// EnvBool creates a new parser setting to load a boolean value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvBool(defaultValue bool, variables ...string) *bool {
	target := new(bool)
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue, isBool: true}, variables)
	return target
}

// EnvInt creates a new parser setting to load a int value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvInt(defaultValue int, variables ...string) *int {
	target := new(int)
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
	return target
}

// EnvInt8 creates a new parser setting to load a int8 value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvInt8(defaultValue int8, variables ...string) *int8 {
	target := new(int8)
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
	return target
}

// EnvInt16 creates a new parser setting to load a int16 value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvInt16(defaultValue int16, variables ...string) *int16 {
	target := new(int16)
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
	return target
}

// EnvInt32 creates a new parser setting to load a int32 value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvInt32(defaultValue int32, variables ...string) *int32 {
	target := new(int32)
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
	return target
}

// EnvInt64 creates a new parser setting to load a int64 value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvInt64(defaultValue int64, variables ...string) *int64 {
	target := new(int64)
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
	return target
}

// EnvUint creates a new parser setting to load a uint value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvUint(defaultValue uint, variables ...string) *uint {
	target := new(uint)
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
	return target
}

// EnvUint8 creates a new parser setting to load a uint8 value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvUint8(defaultValue uint8, variables ...string) *uint8 {
	target := new(uint8)
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
	return target
}

// EnvUint16 creates a new parser setting to load a uint16 value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvUint16(defaultValue uint16, variables ...string) *uint16 {
	target := new(uint16)
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
	return target
}

// EnvUint32 creates a new parser setting to load a uint32 value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvUint32(defaultValue uint32, variables ...string) *uint32 {
	target := new(uint32)
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
	return target
}

// EnvUint64 creates a new parser setting to load a uint64 value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvUint64(defaultValue uint64, variables ...string) *uint64 {
	target := new(uint64)
	cfg.addEnv(&variant{ptrValue: target, defaultValue: defaultValue}, variables)
	return target
}
