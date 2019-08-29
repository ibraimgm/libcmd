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
	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fmt.Errorf("'%v' is not a valid int value", value)
	}
	*target = int(val)

	return nil
}
