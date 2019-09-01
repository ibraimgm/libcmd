package libcfg

import (
	"bufio"
	"os"
	"strings"
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
		entry.val.useDefault()
	}
}

func (cfg *CfgParser) findEnvValue(entry *envEntry) {
	for _, name := range entry.names {
		var value string
		var ok bool

		if value, ok = cfg.envvars[name]; !ok && cfg.useEnv {
			value, ok = os.LookupEnv(name)
		}

		switch {
		case ok && (value != "" || entry.val.isStr):
			entry.val.setValue(value) //nolint: errheck
		case ok && value == "":
			entry.val.unsetValue()
		}
	}
}

// UseEnv sets whether the parser should consider environment variables
// when loading env values. By default, when you define environment options
// (for example, with EnvString) the value to load is searched in the current
// environment variables.
//
// If you call UseEnv(false) this behavior is disabled, and the only source of
// values for environment variables will be calls to UseFile or UseFiles.
//
// This allows the usage of the the CfgParser with env values as a simple
// file config loader.
func (cfg *CfgParser) UseEnv(shouldUse bool) {
	cfg.useEnv = shouldUse
}

// UseFile loads the content of envfile into a in-memory environment value cache.
// Each line of envfile is in the format 'key=value'. Comments start with '#' and
// blank lines are ignored.
//
// You can call UseFile more than once, and newer values will override older values.
// Values loaded with UseFile will always have higher priority than actual environment
// values.
func (cfg *CfgParser) UseFile(envfile string) error {
	file, err := os.Open(envfile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if comment := strings.Index(line, "#"); comment >= 0 {
			line = line[0:comment]
		}

		index := strings.Index(line, "=")

		if index < 0 {
			continue
		}

		key := line[0:index]
		value := line[index+1:]

		cfg.envvars[key] = strings.TrimRight(value, " ")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// UseFiles is a shortcut to multiple calls to UseFile.
// Errors are ignored, so this funcion is useful to load a list of files when is acceptable
// that some of them might not exist
func (cfg *CfgParser) UseFiles(envfiles ...string) {
	for _, file := range envfiles {
		cfg.UseFile(file)
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
	cfg.EnvStringVar(target, defaultValue, variables...)
	return target
}

// EnvBool creates a new parser setting to load a boolean value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvBool(defaultValue bool, variables ...string) *bool {
	target := new(bool)
	cfg.EnvBoolVar(target, defaultValue, variables...)
	return target
}

// EnvInt creates a new parser setting to load a int value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvInt(defaultValue int, variables ...string) *int {
	target := new(int)
	cfg.EnvIntVar(target, defaultValue, variables...)
	return target
}

// EnvInt8 creates a new parser setting to load a int8 value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvInt8(defaultValue int8, variables ...string) *int8 {
	target := new(int8)
	cfg.EnvInt8Var(target, defaultValue, variables...)
	return target
}

// EnvInt16 creates a new parser setting to load a int16 value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvInt16(defaultValue int16, variables ...string) *int16 {
	target := new(int16)
	cfg.EnvInt16Var(target, defaultValue, variables...)
	return target
}

// EnvInt32 creates a new parser setting to load a int32 value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvInt32(defaultValue int32, variables ...string) *int32 {
	target := new(int32)
	cfg.EnvInt32Var(target, defaultValue, variables...)
	return target
}

// EnvInt64 creates a new parser setting to load a int64 value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvInt64(defaultValue int64, variables ...string) *int64 {
	target := new(int64)
	cfg.EnvInt64Var(target, defaultValue, variables...)
	return target
}

// EnvUint creates a new parser setting to load a uint value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvUint(defaultValue uint, variables ...string) *uint {
	target := new(uint)
	cfg.EnvUintVar(target, defaultValue, variables...)
	return target
}

// EnvUint8 creates a new parser setting to load a uint8 value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvUint8(defaultValue uint8, variables ...string) *uint8 {
	target := new(uint8)
	cfg.EnvUint8Var(target, defaultValue, variables...)
	return target
}

// EnvUint16 creates a new parser setting to load a uint16 value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvUint16(defaultValue uint16, variables ...string) *uint16 {
	target := new(uint16)
	cfg.EnvUint16Var(target, defaultValue, variables...)
	return target
}

// EnvUint32 creates a new parser setting to load a uint32 value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvUint32(defaultValue uint32, variables ...string) *uint32 {
	target := new(uint32)
	cfg.EnvUint32Var(target, defaultValue, variables...)
	return target
}

// EnvUint64 creates a new parser setting to load a uint64 value from the
// specified environment variables
// After parsing, the value will be available on the returned pointer.
func (cfg *CfgParser) EnvUint64(defaultValue uint64, variables ...string) *uint64 {
	target := new(uint64)
	cfg.EnvUint64Var(target, defaultValue, variables...)
	return target
}
