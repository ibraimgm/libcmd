package libcfg

import (
	"fmt"
	"strconv"
)

func wrapOptErr(name string, err error) error {
	if err != nil {
		return fmt.Errorf("%v for argument: %s", err, name)
	}

	return nil
}

func setBool(target *bool, value string, reverse bool) error {
	val, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("'%v' is not a valid boolean value", value)
	}
	*target = val

	if reverse {
		*target = !*target
	}

	return nil
}

func setInt(target *int, value string) error {
	val, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return fmt.Errorf("'%v' is not a valid int value", value)
	}
	*target = int(val)

	return nil
}

func setInt8(target *int8, value string) error {
	val, err := strconv.ParseInt(value, 10, 8)
	if err != nil {
		return fmt.Errorf("'%v' is not a valid int8 value", value)
	}
	*target = int8(val)

	return nil
}

func setInt16(target *int16, value string) error {
	val, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		return fmt.Errorf("'%v' is not a valid int16 value", value)
	}
	*target = int16(val)

	return nil
}

func setInt32(target *int32, value string) error {
	val, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return fmt.Errorf("'%v' is not a valid int32 value", value)
	}
	*target = int32(val)

	return nil
}

func setInt64(target *int64, value string) error {
	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fmt.Errorf("'%v' is not a valid int64 value", value)
	}
	*target = val

	return nil
}

func setUint(target *uint, value string) error {
	val, err := strconv.ParseUint(value, 10, 0)
	if err != nil {
		return fmt.Errorf("'%v' is not a valid uint value", value)
	}
	*target = uint(val)

	return nil
}

func setUint8(target *uint8, value string) error {
	val, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return fmt.Errorf("'%v' is not a valid uint8 value", value)
	}
	*target = uint8(val)

	return nil
}

func setUint16(target *uint16, value string) error {
	val, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return fmt.Errorf("'%v' is not a valid uint16 value", value)
	}
	*target = uint16(val)

	return nil
}

func setUint32(target *uint32, value string) error {
	val, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return fmt.Errorf("'%v' is not a valid uint32 value", value)
	}
	*target = uint32(val)

	return nil
}

func setUint64(target *uint64, value string) error {
	val, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return fmt.Errorf("'%v' is not a valid uint64 value", value)
	}
	*target = val

	return nil
}
