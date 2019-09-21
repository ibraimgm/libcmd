package libcfg_test

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/ibraimgm/libcfg"
)

func compareValue(t *testing.T, fixture int, a, b interface{}) {
	if a != b {
		t.Errorf("Case %d, wrong %T value: expected '%v', received '%v'", fixture, a, a, b)
	}
}

func compareArgs(t *testing.T, fixture int, a, b []string) bool {
	if len(a) != len(b) {
		t.Errorf("Case %d, wrong size of rest arguments: expected '%v', received '%v'", fixture, len(a), len(b))
		return false
	}

	for i := 0; i < len(a); i++ {
		if b[i] != a[i] {
			t.Errorf("Case %d, wrong args result at pos %d: expected '%v', received '%v'", fixture, i, a[i], b[i])
			return false
		}
	}

	return true
}

func TestOpt(t *testing.T) {
	tests := []struct {
		cmd      []string
		abool    bool
		aint     int
		auint    uint
		astring  string
		afloat32 float32
		afloat64 float64
		args     []string
	}{
		{cmd: []string{}},
		{cmd: []string{"-b", "-i", "5", "-u", "9", "-s", "foo"}, abool: true, aint: 5, auint: 9, astring: "foo"},
		{cmd: []string{"--abool", "--aint", "5", "--auint", "9", "--astring", "foo"}, abool: true, aint: 5, auint: 9, astring: "foo"},
		{cmd: []string{"--aint=5", "--astring=foo"}, aint: 5, astring: "foo"},
		{cmd: []string{"-b", "--abool=false"}},
		{cmd: []string{"-b", "--no-abool"}},
		{cmd: []string{"-i", "5", "--aint", "6", "-i", "7"}, aint: 7},
		{cmd: []string{"-u", "5", "--auint", "6", "-u", "7"}, auint: 7},
		{cmd: []string{"-b", "-i", "5", "foo", "bar"}, abool: true, aint: 5, args: []string{"foo", "bar"}},
		{cmd: []string{"foo", "bar"}, args: []string{"foo", "bar"}},
		{cmd: []string{"foo", "-i", "5"}, args: []string{"foo", "-i", "5"}},
		{cmd: []string{"-b", "--no-abool=true"}},
		{cmd: []string{"--no-abool=false"}, abool: true},
		{cmd: []string{"-s", "foo", "--astring="}},
		{cmd: []string{"--astring", "--aint", "5"}, astring: "--aint", args: []string{"5"}},
		{cmd: []string{"-i", "5", "-f32", "3.14", "-f64", "3.1415"}, aint: 5, afloat32: float32(3.14), afloat64: float64(3.1415)},
		{cmd: []string{"--afloat32", "3.14", "--afloat64", "3.1415"}, afloat32: float32(3.14), afloat64: float64(3.1415)},
		{cmd: []string{"--afloat32=3.14", "--afloat64=3.1415"}, afloat32: float32(3.14), afloat64: float64(3.1415)},
	}

	for i, test := range tests {
		app := libcfg.NewApp("", "")

		abool := app.Bool("abool", "b", false, "specifies a bool value")
		aint := app.Int("aint", "i", 0, "specifies an int value")
		auint := app.Uint("auint", "u", 0, "specifies an uint value")
		astring := app.String("astring", "s", "", "specifies a string value")
		afloat32 := app.Float32("afloat32", "f32", 0, "specifies a float32 value")
		afloat64 := app.Float64("afloat64", "f64", 0, "specifies a float64 value")

		if err := app.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error parsing args: %v", i, err)
			continue
		}

		compareValue(t, i, test.abool, *abool)
		compareValue(t, i, test.aint, *aint)
		compareValue(t, i, test.auint, *auint)
		compareValue(t, i, test.astring, *astring)
		compareValue(t, i, test.afloat32, *afloat32)
		compareValue(t, i, test.afloat64, *afloat64)
		compareArgs(t, i, test.args, app.Args())
	}
}

func TestOptDefault(t *testing.T) {
	tests := []struct {
		cmd      []string
		abool    bool
		aint     int
		auint    uint
		astring  string
		afloat32 float32
		afloat64 float64
		args     []string
	}{
		{cmd: []string{}, abool: true, aint: 8, auint: 16, afloat32: float32(3.14), afloat64: float64(3.1415), astring: "default"},
		{cmd: []string{"-b", "-i", "5", "-u", "9", "-s", "foo", "-f32", "5.5", "-f64", "5.555"}, abool: true, aint: 5, auint: 9, astring: "foo", afloat32: float32(5.5), afloat64: float64(5.555)},
		{cmd: []string{"--abool", "--aint", "5", "--auint", "9", "--astring", "foo", "--afloat32", "5.5", "--afloat64", "5.555"}, abool: true, aint: 5, auint: 9, astring: "foo", afloat32: float32(5.5), afloat64: float64(5.555)},
		{cmd: []string{"--aint=5", "--astring=foo"}, abool: true, aint: 5, auint: 16, astring: "foo", afloat32: float32(3.14), afloat64: float64(3.1415)},
		{cmd: []string{"-b", "--abool=false"}, aint: 8, auint: 16, astring: "default", afloat32: float32(3.14), afloat64: float64(3.1415)},
		{cmd: []string{"-b", "--no-abool"}, aint: 8, auint: 16, astring: "default", afloat32: float32(3.14), afloat64: float64(3.1415)},
		{cmd: []string{"-i", "5", "--aint", "6", "-i", "7", "-f32", "5.5", "--afloat32", "3.14"}, abool: true, aint: 7, auint: 16, astring: "default", afloat32: float32(3.14), afloat64: float64(3.1415)},
		{cmd: []string{"-u", "5", "--auint", "6", "-u", "7"}, abool: true, aint: 8, auint: 7, astring: "default", afloat32: float32(3.14), afloat64: float64(3.1415)},
		{cmd: []string{"-b", "-i", "5", "foo", "bar"}, abool: true, aint: 5, auint: 16, astring: "default", args: []string{"foo", "bar"}, afloat32: float32(3.14), afloat64: float64(3.1415)},
		{cmd: []string{"foo", "bar"}, abool: true, aint: 8, auint: 16, astring: "default", args: []string{"foo", "bar"}, afloat32: float32(3.14), afloat64: float64(3.1415)},
		{cmd: []string{"foo", "-i", "5"}, abool: true, aint: 8, auint: 16, astring: "default", args: []string{"foo", "-i", "5"}, afloat32: float32(3.14), afloat64: float64(3.1415)},
		{cmd: []string{"-b", "--no-abool=true"}, aint: 8, auint: 16, astring: "default", afloat32: float32(3.14), afloat64: float64(3.1415)},
		{cmd: []string{"--no-abool=false"}, abool: true, aint: 8, auint: 16, astring: "default", afloat32: float32(3.14), afloat64: float64(3.1415)},
		{cmd: []string{"-s", "foo", "--astring="}, abool: true, aint: 8, auint: 16, afloat32: float32(3.14), afloat64: float64(3.1415)},
		{cmd: []string{"--astring", "--aint", "5"}, abool: true, aint: 8, auint: 16, astring: "--aint", args: []string{"5"}, afloat32: float32(3.14), afloat64: float64(3.1415)},
	}

	for i, test := range tests {
		app := libcfg.NewApp("", "")

		abool := app.Bool("abool", "b", true, "specifies a bool value")
		aint := app.Int("aint", "i", 8, "specifies an int value")
		auint := app.Uint("auint", "u", 16, "specifies an uint value")
		astring := app.String("astring", "s", "default", "specifies a string value")
		afloat32 := app.Float32("afloat32", "f32", float32(3.14), "specifies a float32 value")
		afloat64 := app.Float64("afloat64", "f64", float64(3.1415), "specifies a float64 value")

		if err := app.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error parsing args: %v", i, err)
			continue
		}

		compareValue(t, i, test.abool, *abool)
		compareValue(t, i, test.aint, *aint)
		compareValue(t, i, test.auint, *auint)
		compareValue(t, i, test.astring, *astring)
		compareValue(t, i, test.afloat32, *afloat32)
		compareValue(t, i, test.afloat64, *afloat64)
		compareArgs(t, i, test.args, app.Args())
	}
}

func TestOptError(t *testing.T) {
	tests := []struct {
		cmd           []string
		expectedError string
	}{
		{cmd: []string{"-b", "-x"}, expectedError: "unknown argument: -x"},
		{cmd: []string{"-b", "--x"}, expectedError: "unknown argument: --x"},
		{cmd: []string{"--abool=X"}, expectedError: "is not a valid boolean value"},
		{cmd: []string{"-i", "a"}, expectedError: "is not a valid int value"},
		{cmd: []string{"-i"}, expectedError: "no value for argument: -i"},
		{cmd: []string{"-u", "a"}, expectedError: "is not a valid uint value"},
		{cmd: []string{"-u"}, expectedError: "no value for argument: -u"},
		{cmd: []string{"-s"}, expectedError: "no value for argument: -s"},
		{cmd: []string{"--aint"}, expectedError: "no value for argument: --aint"},
		{cmd: []string{"--auint"}, expectedError: "no value for argument: --auint"},
		{cmd: []string{"--astring"}, expectedError: "no value for argument: --astring"},
		{cmd: []string{"--aint="}, expectedError: "no value for argument: --aint"},
		{cmd: []string{"--aint=", "5"}, expectedError: "no value for argument: --aint"},
		{cmd: []string{"--auint="}, expectedError: "no value for argument: --auint"},
		{cmd: []string{"--auint=", "5"}, expectedError: "no value for argument: --auint"},
	}

	for i, test := range tests {
		app := libcfg.NewApp("", "")

		app.Bool("abool", "b", false, "specifies a bool value")
		app.Int("aint", "i", 0, "specifies an int value")
		app.Uint("auint", "u", 0, "specifies an uint value")
		app.String("astring", "s", "", "specifies a string value")

		err := app.RunArgs(test.cmd)

		if err == nil {
			t.Errorf("Case %d, argument parsing should return error", i)
			continue
		}

		if !strings.Contains(err.Error(), test.expectedError) {
			t.Errorf("Case %d, expected error '%s', but got '%s'", i, test.expectedError, err.Error())
		}
	}
}

func TestOptIntLimit(t *testing.T) {
	tests := []struct {
		cmd []string
		a   int8
		b   int16
		c   int32
		d   int64
	}{
		{cmd: []string{"-a", strconv.FormatInt(math.MaxInt8, 10)}, a: math.MaxInt8},
		{cmd: []string{"-b", strconv.FormatInt(math.MaxInt16, 10)}, b: math.MaxInt16},
		{cmd: []string{"-c", strconv.FormatInt(math.MaxInt32, 10)}, c: math.MaxInt32},
		{cmd: []string{"-d", strconv.FormatInt(math.MaxInt64, 10)}, d: math.MaxInt64},
	}

	for i, test := range tests {
		app := libcfg.NewApp("", "")

		a := app.Int8("", "a", 0, "specifies a int8 value")
		b := app.Int16("", "b", 0, "specifies a int16 value")
		c := app.Int32("", "c", 0, "specifies a int32 value")
		d := app.Int64("", "d", 0, "specifies a int64 value")

		if err := app.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error parsing args: %v", i, err)
			continue
		}

		compareValue(t, i, test.a, *a)
		compareValue(t, i, test.b, *b)
		compareValue(t, i, test.c, *c)
		compareValue(t, i, test.d, *d)
	}
}

func TestOptUintLimit(t *testing.T) {
	tests := []struct {
		cmd []string
		a   uint8
		b   uint16
		c   uint32
		d   uint64
	}{
		{cmd: []string{"-a", strconv.FormatUint(math.MaxUint8, 10)}, a: math.MaxUint8},
		{cmd: []string{"-b", strconv.FormatUint(math.MaxUint16, 10)}, b: math.MaxUint16},
		{cmd: []string{"-c", strconv.FormatUint(math.MaxUint32, 10)}, c: math.MaxUint32},
		{cmd: []string{"-d", strconv.FormatUint(math.MaxUint64, 10)}, d: math.MaxUint64},
	}

	for i, test := range tests {
		app := libcfg.NewApp("", "")

		a := app.Uint8("", "a", 0, "specifies a uint8 value")
		b := app.Uint16("", "b", 0, "specifies a uint16 value")
		c := app.Uint32("", "c", 0, "specifies a uint32 value")
		d := app.Uint64("", "d", 0, "specifies a uint64 value")

		if err := app.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error parsing args: %v", i, err)
			continue
		}

		compareValue(t, i, test.a, *a)
		compareValue(t, i, test.b, *b)
		compareValue(t, i, test.c, *c)
		compareValue(t, i, test.d, *d)
	}
}

func TestOptIntegerLimits(t *testing.T) {
	tests := []struct {
		cmd           []string
		expectedError string
	}{
		{cmd: []string{"--aint8", "1" + strconv.FormatInt(math.MaxInt8, 10)}, expectedError: "is not a valid int8 value"},
		{cmd: []string{"--aint16", "1" + strconv.FormatInt(math.MaxInt16, 10)}, expectedError: "is not a valid int16 value"},
		{cmd: []string{"--aint32", "1" + strconv.FormatInt(math.MaxInt32, 10)}, expectedError: "is not a valid int32 value"},
		{cmd: []string{"--aint64", "1" + strconv.FormatInt(math.MaxInt64, 10)}, expectedError: "is not a valid int64 value"},
		{cmd: []string{"--auint8", "1" + strconv.FormatUint(math.MaxUint8, 10)}, expectedError: "is not a valid uint8 value"},
		{cmd: []string{"--auint16", "1" + strconv.FormatUint(math.MaxUint16, 10)}, expectedError: "is not a valid uint16 value"},
		{cmd: []string{"--auint32", "1" + strconv.FormatUint(math.MaxUint32, 10)}, expectedError: "is not a valid uint32 value"},
		{cmd: []string{"--auint64", "1" + strconv.FormatUint(math.MaxUint64, 10)}, expectedError: "is not a valid uint64 value"},
		{cmd: []string{"--auint8", "-1"}, expectedError: "is not a valid uint8 value"},
		{cmd: []string{"--auint16", "-1"}, expectedError: "is not a valid uint16 value"},
		{cmd: []string{"--auint32", "-1"}, expectedError: "is not a valid uint32 value"},
		{cmd: []string{"--auint64", "-1"}, expectedError: "is not a valid uint64 value"},
	}

	for i, test := range tests {
		app := libcfg.NewApp("", "")

		app.Int8("aint8", "", 0, "specifies a int8 value")
		app.Int16("aint16", "", 0, "specifies a int16 value")
		app.Int32("aint32", "", 0, "specifies a int32 value")
		app.Int64("aint64", "", 0, "specifies a int64 value")
		app.Uint8("auint8", "", 0, "specifies a uint8 value")
		app.Uint16("auint16", "", 0, "specifies a uint16 value")
		app.Uint32("auint32", "", 0, "specifies a uint32 value")
		app.Uint64("auint64", "", 0, "specifies a uint64 value")

		err := app.RunArgs(test.cmd)

		if err == nil {
			t.Errorf("Case %d, argument parsing should return error", i)
			continue
		}

		if !strings.Contains(err.Error(), test.expectedError) {
			t.Errorf("Case %d, expected error '%s', but got '%s'", i, test.expectedError, err.Error())
		}
	}
}

func TestOptKeepValue(t *testing.T) {
	const keep = "keep"

	tests := []struct {
		cmd []string
		s1  string
		s2  string
	}{
		{cmd: []string{}, s1: keep, s2: "default"},
		{cmd: []string{"-s1", "a"}, s1: "a", s2: "default"},
		{cmd: []string{"--string1", "a"}, s1: "a", s2: "default"},
		{cmd: []string{"--string1=a"}, s1: "a", s2: "default"},
		{cmd: []string{"--string1="}, s1: "", s2: "default"},
		{cmd: []string{"-s2", "x"}, s1: keep, s2: "x"},
		{cmd: []string{"--string2", "x"}, s1: keep, s2: "x"},
		{cmd: []string{"--string2=x"}, s1: keep, s2: "x"},
		{cmd: []string{"--string2="}, s1: keep, s2: ""},
		{cmd: []string{"-s1", "a", "-s2", "x"}, s1: "a", s2: "x"},
		{cmd: []string{"--string1=", "-s2", "x"}, s1: "", s2: "x"},
		{cmd: []string{"-s1", "a", "-s2", "x", "--string1="}, s1: "", s2: "x"},
	}

	for i, test := range tests {
		app := libcfg.NewApp("", "")

		s1 := app.String("string1", "s1", "", "")
		s2 := app.String("string2", "s2", "default", "")

		*s1 = keep
		*s2 = ""

		if err := app.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error parsing args: %v", i, err)
			continue
		}

		compareValue(t, i, test.s1, *s1)
		compareValue(t, i, test.s2, *s2)
	}
}
