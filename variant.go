package libcfg

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type variant struct {
	refValue     reflect.Value
	defaultValue reflect.Value
	isBool       bool
	isStr        bool
	isSet        bool
	isOpt        bool
}

func newVariant(target, defaultValue interface{}) *variant {
	ref := reflect.ValueOf(target)

	if ref.Kind() == reflect.Ptr {
		if ref.IsNil() {
			panic("nil pointer on variant creation")
		}

		ref = reflect.Indirect(ref)
	}

	def := reflect.ValueOf(defaultValue)

	return &variant{
		refValue:     ref,
		defaultValue: def,
		isBool:       ref.Kind() == reflect.Bool,
		isStr:        ref.Kind() == reflect.String,
	}
}

func (v *variant) setValue(value string) error {

	switch v.refValue.Kind() {
	case reflect.String:
		v.refValue.SetString(value)

	case reflect.Bool:
		vv, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid boolean value", value)
		}
		v.refValue.SetBool(vv)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if err := v.setAsInt(value); err != nil {
			return err
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if err := v.setAsUint(value); err != nil {
			return err
		}

	case reflect.Float32, reflect.Float64:
		if err := v.setAsFloat(value); err != nil {
			return err
		}
	}

	v.isSet = true
	return nil
}

func (v *variant) setAsInt(value string) error {
	name := v.refValue.Kind().String()
	var size int

	if name != "int" {
		bitSize, err := strconv.ParseInt(strings.TrimPrefix(name, "int"), 10, 0)
		if err != nil {
			return fmt.Errorf("invalid integer type '%s' for value '%s'", name, value)
		}

		size = int(bitSize)
	}

	vv, err := strconv.ParseInt(value, 10, int(size))
	if err != nil {
		return fmt.Errorf("'%v' is not a valid %s value", value, name)
	}

	v.refValue.SetInt(vv)
	return nil
}

func (v *variant) setAsUint(value string) error {
	name := v.refValue.Kind().String()
	var size int

	if name != "uint" {
		bitSize, err := strconv.ParseInt(strings.TrimPrefix(name, "uint"), 10, 0)
		if err != nil {
			return fmt.Errorf("invalid integer type '%s' for value '%s'", name, value)
		}

		size = int(bitSize)
	}

	vv, err := strconv.ParseUint(value, 10, size)
	if err != nil {
		return fmt.Errorf("'%v' is not a valid %s value", value, name)
	}

	v.refValue.SetUint(vv)
	return nil
}

func (v *variant) setAsFloat(value string) error {
	name := v.refValue.Kind().String()
	var size int

	switch name {
	case "float32":
		size = 32
	case "float64":
		size = 64
	default:
		return fmt.Errorf("invalid float type '%s' for value '%s'", name, value)
	}

	vv, err := strconv.ParseFloat(value, size)
	if err != nil {
		return fmt.Errorf("'%v' is not a valid %s value", value, name)
	}

	v.refValue.SetFloat(vv)
	return nil
}

func (v *variant) useDefault() {
	if v.isSet {
		return
	}

	v.refValue.Set(v.defaultValue)
}

func (v *variant) setToZero() {
	v.isSet = true

	zero := reflect.Zero(v.refValue.Type())
	v.refValue.Set(zero)
}
