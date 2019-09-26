package libcmd

import (
	"fmt"
	"reflect"
	"strconv"
)

type variant struct {
	raw          interface{}
	refValue     reflect.Value
	defaultValue reflect.Value
	isBool       bool
	isStr        bool
	isSet        bool
}

func varFromInterface(target, defaultValue interface{}) *variant {
	ref := reflect.ValueOf(target)

	if ref.Kind() == reflect.Ptr {
		if ref.IsNil() {
			panic("nil pointer on variant creation")
		}

		ref = reflect.Indirect(ref)
	}

	def := reflect.ValueOf(defaultValue)

	return varFromReflect(target, ref, def)
}

func varFromReflect(raw interface{}, target, defaultValue reflect.Value) *variant {
	return &variant{
		raw:          raw,
		refValue:     target,
		defaultValue: defaultValue,
		isBool:       target.Kind() == reflect.Bool,
		isStr:        target.Kind() == reflect.String,
	}
}

func varFromCustom(target CustomArg, defaultValue string) *variant {
	return &variant{
		raw:          target,
		refValue:     reflect.ValueOf(target),
		defaultValue: reflect.ValueOf(defaultValue),
		isStr:        true,
	}
}

func (v *variant) setValue(value string) error {
	if v.refValue.Type().Implements(customArgType) {
		ca, _ := v.refValue.Interface().(CustomArg)
		return ca.Set(value)
	}

	converted, err := valueAsKind(value, v.refValue.Kind(), v.refValue.Type())
	if err != nil {
		return err
	}

	v.refValue.Set(converted)

	v.isSet = true
	return nil
}

func (v *variant) useDefault() {
	if v.isSet {
		return
	}

	if v.refValue.Type().Implements(customArgType) {
		ca, _ := v.refValue.Interface().(CustomArg)

		if ca.Get() == "" {
			return
		}

		str := v.defaultValue.String()
		ca.Set(str) //nolint: errcheck
		return
	}

	zero := reflect.Zero(v.refValue.Type())
	defaultIsZero := zero.Interface() == v.defaultValue.Interface()
	valueIsZero := zero.Interface() == v.refValue.Interface()

	// if we have a value, and our default is zero, we keep it
	if defaultIsZero && !valueIsZero {
		return
	}

	v.refValue.Set(v.defaultValue)
}

func (v *variant) defaultAsString() string {
	zero := reflect.Zero(v.refValue.Type())

	if zero.Interface() == v.defaultValue.Interface() {
		return ""
	}
	return fmt.Sprintf("%v", v.defaultValue.Interface())
}

// return the value converted to a *COMPATIBLE* kind, optionally casted to
// an *EXACT* type
func valueAsKind(value string, kind reflect.Kind, exactType reflect.Type) (reflect.Value, error) {
	// used only for numeric types; 'simpler' types are returned directly
	name := kind.String()
	size := bitSizeOf(kind)
	var parsed interface{}
	var err error

	switch kind {
	case reflect.String:
		return reflect.ValueOf(value), nil

	case reflect.Bool:
		vv, err := strconv.ParseBool(value)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("'%v' is not a valid boolean value", value)
		}
		return reflect.ValueOf(vv), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err = strconv.ParseInt(value, 10, size)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsed, err = strconv.ParseUint(value, 10, size)

	case reflect.Float32, reflect.Float64:
		parsed, err = strconv.ParseFloat(value, size)
	}

	if err != nil {
		return reflect.Value{}, fmt.Errorf("'%v' is not a valid %s value", value, name)
	}

	if parsed == nil {
		return reflect.Value{}, fmt.Errorf("unsupported type '%s' for value '%s'", name, value)
	}

	converted := reflect.ValueOf(parsed)

	if exactType != nil {
		converted = converted.Convert(exactType)
	}

	return converted, nil
}

// bit size of numeric types
func bitSizeOf(kind reflect.Kind) int {
	switch kind {
	case reflect.Int8, reflect.Uint8:
		return 8

	case reflect.Int16, reflect.Uint16:
		return 16

	case reflect.Int32, reflect.Uint32, reflect.Float32:
		return 32

	case reflect.Int64, reflect.Uint64, reflect.Float64:
		return 64

	default:
		return 0
	}
}
