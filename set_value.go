package libcfg

import (
	"fmt"
	"strconv"
)

type variant struct {
	ptrValue     interface{}
	defaultValue interface{}
	isBool       bool
	isStr        bool
}

func (v *variant) setValueStr(value string) error { //nolint: gocyclo
	switch v := v.ptrValue.(type) {
	case *string:
		*v = value

	case *bool:
		vv, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid boolean value", value)
		}
		*v = vv

	case *int:
		vv, err := strconv.ParseInt(value, 10, 0)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid int value", value)
		}
		*v = int(vv)

	case *int8:
		vv, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid int8 value", value)
		}
		*v = int8(vv)

	case *int16:
		vv, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid int16 value", value)
		}
		*v = int16(vv)

	case *int32:
		vv, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid int32 value", value)
		}
		*v = int32(vv)

	case *int64:
		vv, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid int64 value", value)
		}
		*v = int64(vv)

	case *uint:
		vv, err := strconv.ParseUint(value, 10, 0)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid uint value", value)
		}
		*v = uint(vv)

	case *uint8:
		vv, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid uint8 value", value)
		}
		*v = uint8(vv)

	case *uint16:
		vv, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid uint16 value", value)
		}
		*v = uint16(vv)

	case *uint32:
		vv, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid uint32 value", value)
		}
		*v = uint32(vv)

	case *uint64:
		vv, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("'%v' is not a valid uint64 value", value)
		}
		*v = uint64(vv)
	}

	return nil
}
