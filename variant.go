package libcfg

import (
	"fmt"
	"reflect"
	"strconv"
)

type variant struct {
	ptrValue     interface{}
	defaultValue interface{}
	isBool       bool
	isStr        bool
	isSet        bool
}

func (v *variant) setValue(value string) error { //nolint: gocyclo
	switch val := v.ptrValue.(type) {
	case *string:
		*val = value

	case *bool:
		vv, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid boolean value", value)
		}
		*val = vv

	case *int:
		vv, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid int value", value)
		}
		*val = int(vv)

	case *int8:
		vv, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid int8 value", value)
		}
		*val = int8(vv)

	case *int16:
		vv, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid int16 value", value)
		}
		*val = int16(vv)

	case *int32:
		vv, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid int32 value", value)
		}
		*val = int32(vv)

	case *int64:
		vv, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid int64 value", value)
		}
		*val = int64(vv)

	case *uint:
		vv, err := strconv.ParseUint(value, 10, 0)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid uint value", value)
		}
		*val = uint(vv)

	case *uint8:
		vv, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid uint8 value", value)
		}
		*val = uint8(vv)

	case *uint16:
		vv, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid uint16 value", value)
		}
		*val = uint16(vv)

	case *uint32:
		vv, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid uint32 value", value)
		}
		*val = uint32(vv)

	case *uint64:
		vv, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid uint64 value", value)
		}
		*val = uint64(vv)
	}

	v.isSet = true
	return nil
}

func (v *variant) useDefault() {
	if v.isSet {
		return
	}

	// TODO: check performance of this versus a type switch
	ptr := reflect.Indirect(reflect.ValueOf(v.ptrValue))
	val := reflect.ValueOf(v.defaultValue)

	ptr.Set(val)
}

func (v *variant) unsetValue() {
	v.isSet = false

	ptr := reflect.Indirect(reflect.ValueOf(v.ptrValue))
	val := reflect.Zero(reflect.TypeOf(v.defaultValue))
	ptr.Set(val)
}
