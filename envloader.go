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

type envLoaderImpl struct {
	useEnv     bool
	enventries []*envEntry
	envvars    map[string]string
}

func makeEnvLoader() *envLoaderImpl {
	return &envLoaderImpl{
		useEnv:     true,
		enventries: make([]*envEntry, 0),
		envvars:    make(map[string]string),
	}
}

// NewEnvLoader returns a new instance of an EnvLoader
func NewEnvLoader() EnvLoader {
	return makeEnvLoader()
}

func (env *envLoaderImpl) addEnv(val *variant, variables []string) {
	env.enventries = append(env.enventries, &envEntry{names: variables, val: val})
}

func (env *envLoaderImpl) LoadAll() {
	for i := range env.enventries {
		entry := env.enventries[i]

		if entry.val.isOpt {
			continue
		}

		env.findEnvValue(entry)
		entry.val.useDefault()
	}
}

func (env *envLoaderImpl) findEnvValue(entry *envEntry) {
	for _, name := range entry.names {
		var value string
		var ok bool

		if value, ok = env.envvars[name]; !ok && env.useEnv {
			value, ok = os.LookupEnv(name)
		}

		switch {
		case ok && (value != "" || entry.val.isStr):
			entry.val.setValue(value) //nolint: errcheck
		case ok && value == "":
			entry.val.setToZero()
		}
	}
}

func (env *envLoaderImpl) findEnvEntryByPtr(val *variant) *envEntry {
	for i := range env.enventries {
		entry := env.enventries[i]

		if entry.val == val {
			return entry
		}
	}

	return nil
}

func (env *envLoaderImpl) UseEnv(shouldUse bool) {
	env.useEnv = shouldUse
}

func (env *envLoaderImpl) UseFile(envfile string) error {
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

		env.envvars[key] = strings.TrimRight(value, " ")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (env *envLoaderImpl) UseFiles(envfiles ...string) {
	for _, file := range envfiles {
		env.UseFile(file) //nolint: errcheck
	}
}

func (env *envLoaderImpl) Bind(i interface{}) error {
	data, err := collectBindings(i)
	if err != nil {
		return err
	}

	for _, d := range data {
		if len(d.variables) > 0 {
			env.addEnv(d.val, d.variables)
		}
	}

	return nil
}

func (env *envLoaderImpl) StringP(target *string, defaultValue string, variables ...string) {
	env.addEnv(varFromInterface(target, defaultValue), variables)
}

func (env *envLoaderImpl) BoolP(target *bool, defaultValue bool, variables ...string) {
	env.addEnv(varFromInterface(target, defaultValue), variables)
}

func (env *envLoaderImpl) IntP(target *int, defaultValue int, variables ...string) {
	env.addEnv(varFromInterface(target, defaultValue), variables)
}

func (env *envLoaderImpl) Int8P(target *int8, defaultValue int8, variables ...string) {
	env.addEnv(varFromInterface(target, defaultValue), variables)
}

func (env *envLoaderImpl) Int16P(target *int16, defaultValue int16, variables ...string) {
	env.addEnv(varFromInterface(target, defaultValue), variables)
}

func (env *envLoaderImpl) Int32P(target *int32, defaultValue int32, variables ...string) {
	env.addEnv(varFromInterface(target, defaultValue), variables)
}

func (env *envLoaderImpl) Int64P(target *int64, defaultValue int64, variables ...string) {
	env.addEnv(varFromInterface(target, defaultValue), variables)
}

func (env *envLoaderImpl) UintP(target *uint, defaultValue uint, variables ...string) {
	env.addEnv(varFromInterface(target, defaultValue), variables)
}

func (env *envLoaderImpl) Uint8P(target *uint8, defaultValue uint8, variables ...string) {
	env.addEnv(varFromInterface(target, defaultValue), variables)
}

func (env *envLoaderImpl) Uint16P(target *uint16, defaultValue uint16, variables ...string) {
	env.addEnv(varFromInterface(target, defaultValue), variables)
}

func (env *envLoaderImpl) Uint32P(target *uint32, defaultValue uint32, variables ...string) {
	env.addEnv(varFromInterface(target, defaultValue), variables)
}

func (env *envLoaderImpl) Uint64P(target *uint64, defaultValue uint64, variables ...string) {
	env.addEnv(varFromInterface(target, defaultValue), variables)
}

func (env *envLoaderImpl) Float32P(target *float32, defaultValue float32, variables ...string) {
	env.addEnv(varFromInterface(target, defaultValue), variables)
}

func (env *envLoaderImpl) Float64P(target *float64, defaultValue float64, variables ...string) {
	env.addEnv(varFromInterface(target, defaultValue), variables)
}

func (env *envLoaderImpl) String(defaultValue string, variables ...string) *string {
	target := new(string)
	env.StringP(target, defaultValue, variables...)
	return target
}

func (env *envLoaderImpl) Bool(defaultValue bool, variables ...string) *bool {
	target := new(bool)
	env.BoolP(target, defaultValue, variables...)
	return target
}

func (env *envLoaderImpl) Int(defaultValue int, variables ...string) *int {
	target := new(int)
	env.IntP(target, defaultValue, variables...)
	return target
}

func (env *envLoaderImpl) Int8(defaultValue int8, variables ...string) *int8 {
	target := new(int8)
	env.Int8P(target, defaultValue, variables...)
	return target
}

func (env *envLoaderImpl) Int16(defaultValue int16, variables ...string) *int16 {
	target := new(int16)
	env.Int16P(target, defaultValue, variables...)
	return target
}

func (env *envLoaderImpl) Int32(defaultValue int32, variables ...string) *int32 {
	target := new(int32)
	env.Int32P(target, defaultValue, variables...)
	return target
}

func (env *envLoaderImpl) Int64(defaultValue int64, variables ...string) *int64 {
	target := new(int64)
	env.Int64P(target, defaultValue, variables...)
	return target
}

func (env *envLoaderImpl) Uint(defaultValue uint, variables ...string) *uint {
	target := new(uint)
	env.UintP(target, defaultValue, variables...)
	return target
}

func (env *envLoaderImpl) Uint8(defaultValue uint8, variables ...string) *uint8 {
	target := new(uint8)
	env.Uint8P(target, defaultValue, variables...)
	return target
}

func (env *envLoaderImpl) Uint16(defaultValue uint16, variables ...string) *uint16 {
	target := new(uint16)
	env.Uint16P(target, defaultValue, variables...)
	return target
}

func (env *envLoaderImpl) Uint32(defaultValue uint32, variables ...string) *uint32 {
	target := new(uint32)
	env.Uint32P(target, defaultValue, variables...)
	return target
}

func (env *envLoaderImpl) Uint64(defaultValue uint64, variables ...string) *uint64 {
	target := new(uint64)
	env.Uint64P(target, defaultValue, variables...)
	return target
}

func (env *envLoaderImpl) Float32(defaultValue float32, variables ...string) *float32 {
	target := new(float32)
	env.Float32P(target, defaultValue, variables...)
	return target
}

func (env *envLoaderImpl) Float64(defaultValue float64, variables ...string) *float64 {
	target := new(float64)
	env.Float64P(target, defaultValue, variables...)
	return target
}
