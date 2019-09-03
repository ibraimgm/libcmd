package libcfg_test

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/ibraimgm/libcfg"
)

func TestParseOptArgs(t *testing.T) {
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
		cfg := libcfg.NewParser()

		abool := cfg.OptBool("abool", "b", false, "specifies a bool value")
		aint := cfg.OptInt("aint", "i", 0, "specifies an int value")
		auint := cfg.OptUint("auint", "u", 0, "specifies an uint value")
		astring := cfg.OptString("astring", "s", "", "specifies a string value")
		afloat32 := cfg.OptFloat32("afloat32", "f32", 0, "specifies a float32 value")
		afloat64 := cfg.OptFloat64("afloat64", "f64", 0, "specifies a float64 value")

		if err := cfg.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error parsing args: %v", i, err)
			continue
		}

		if *abool != test.abool {
			t.Errorf("Case %d, wrong boolean value: expected '%v', received '%v'", i, test.abool, *abool)
		}

		if *aint != test.aint {
			t.Errorf("Case %d, wrong int value: expected '%v', received '%v'", i, test.aint, *aint)
		}

		if *auint != test.auint {
			t.Errorf("Case %d, wrong uint value: expected '%v', received '%v'", i, test.auint, *auint)
		}

		if *astring != test.astring {
			t.Errorf("Case %d, wrong string value: expected '%v', received '%v'", i, test.astring, *astring)
		}

		if *afloat32 != test.afloat32 {
			t.Errorf("Case %d, wrong float32 value: expected '%v', received '%v'", i, test.afloat32, *afloat32)
		}

		if *afloat64 != test.afloat64 {
			t.Errorf("Case %d, wrong float64 value: expected '%v', received '%v'", i, test.afloat64, *afloat64)
		}

		args := cfg.Args()

		if len(test.args) != len(args) {
			t.Errorf("Case %d, wrong size of rest arguments: expected '%v', received '%v'", i, len(test.args), len(args))
			continue
		}

		for j := 0; j < len(test.args); j++ {
			if args[j] != test.args[j] {
				t.Errorf("Case %d, wrong args result at pos %d: expected '%v', received '%v'", i, j, test.args[j], args[j])
			}
		}
	}
}

func TestParseOptDefault(t *testing.T) {
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
		cfg := libcfg.NewParser()

		abool := cfg.OptBool("abool", "b", true, "specifies a bool value")
		aint := cfg.OptInt("aint", "i", 8, "specifies an int value")
		auint := cfg.OptUint("auint", "u", 16, "specifies an uint value")
		astring := cfg.OptString("astring", "s", "default", "specifies a string value")
		afloat32 := cfg.OptFloat32("afloat32", "f32", float32(3.14), "specifies a float32 value")
		afloat64 := cfg.OptFloat64("afloat64", "f64", float64(3.1415), "specifies a float64 value")

		if err := cfg.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error parsing args: %v", i, err)
			continue
		}

		if *abool != test.abool {
			t.Errorf("Case %d, wrong boolean value: expected '%v', received '%v'", i, test.abool, *abool)
		}

		if *aint != test.aint {
			t.Errorf("Case %d, wrong int value: expected '%v', received '%v'", i, test.aint, *aint)
		}

		if *auint != test.auint {
			t.Errorf("Case %d, wrong uint value: expected '%v', received '%v'", i, test.auint, *auint)
		}

		if *astring != test.astring {
			t.Errorf("Case %d, wrong string value: expected '%v', received '%v'", i, test.astring, *astring)
		}

		if *afloat32 != test.afloat32 {
			t.Errorf("Case %d, wrong float32 value: expected '%v', received '%v'", i, test.afloat32, *afloat32)
		}

		if *afloat64 != test.afloat64 {
			t.Errorf("Case %d, wrong float64 value: expected '%v', received '%v'", i, test.afloat64, *afloat64)
		}

		args := cfg.Args()

		if len(test.args) != len(args) {
			t.Errorf("Case %d, wrong size of rest arguments: expected '%v', received '%v'", i, len(test.args), len(args))
			continue
		}

		for j := 0; j < len(test.args); j++ {
			if args[j] != test.args[j] {
				t.Errorf("Case %d, wrong args result at pos %d: expected '%v', received '%v'", i, j, test.args[j], args[j])
			}
		}
	}
}

func TestParseOptError(t *testing.T) {
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
		cfg := libcfg.NewParser()

		cfg.OptBool("abool", "b", false, "specifies a bool value")
		cfg.OptInt("aint", "i", 0, "specifies an int value")
		cfg.OptUint("auint", "u", 0, "specifies an uint value")
		cfg.OptString("astring", "s", "", "specifies a string value")

		err := cfg.RunArgs(test.cmd)

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
		cfg := libcfg.NewParser()

		a := cfg.OptInt8("", "a", 0, "specifies a int8 value")
		b := cfg.OptInt16("", "b", 0, "specifies a int16 value")
		c := cfg.OptInt32("", "c", 0, "specifies a int32 value")
		d := cfg.OptInt64("", "d", 0, "specifies a int64 value")

		if err := cfg.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error parsing args: %v", i, err)
			continue
		}

		if *a != test.a {
			t.Errorf("Case %d, wrong value: expected '%v', received '%v'", i, test.a, *a)
		}

		if *b != test.b {
			t.Errorf("Case %d, wrong value: expected '%v', received '%v'", i, test.b, *b)
		}

		if *c != test.c {
			t.Errorf("Case %d, wrong value: expected '%v', received '%v'", i, test.c, *c)
		}

		if *d != test.d {
			t.Errorf("Case %d, wrong value: expected '%v', received '%v'", i, test.d, *d)
		}
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
		cfg := libcfg.NewParser()

		a := cfg.OptUint8("", "a", 0, "specifies a uint8 value")
		b := cfg.OptUint16("", "b", 0, "specifies a uint16 value")
		c := cfg.OptUint32("", "c", 0, "specifies a uint32 value")
		d := cfg.OptUint64("", "d", 0, "specifies a uint64 value")

		if err := cfg.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error parsing args: %v", i, err)
			continue
		}

		if *a != test.a {
			t.Errorf("Case %d, wrong value: expected '%v', received '%v'", i, test.a, *a)
		}

		if *b != test.b {
			t.Errorf("Case %d, wrong value: expected '%v', received '%v'", i, test.b, *b)
		}

		if *c != test.c {
			t.Errorf("Case %d, wrong value: expected '%v', received '%v'", i, test.c, *c)
		}

		if *d != test.d {
			t.Errorf("Case %d, wrong value: expected '%v', received '%v'", i, test.d, *d)
		}
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
		cfg := libcfg.NewParser()

		cfg.OptInt8("aint8", "", 0, "specifies a int8 value")
		cfg.OptInt16("aint16", "", 0, "specifies a int16 value")
		cfg.OptInt32("aint32", "", 0, "specifies a int32 value")
		cfg.OptInt64("aint64", "", 0, "specifies a int64 value")
		cfg.OptUint8("auint8", "", 0, "specifies a uint8 value")
		cfg.OptUint16("auint16", "", 0, "specifies a uint16 value")
		cfg.OptUint32("auint32", "", 0, "specifies a uint32 value")
		cfg.OptUint64("auint64", "", 0, "specifies a uint64 value")

		err := cfg.RunArgs(test.cmd)

		if err == nil {
			t.Errorf("Case %d, argument parsing should return error", i)
			continue
		}

		if !strings.Contains(err.Error(), test.expectedError) {
			t.Errorf("Case %d, expected error '%s', but got '%s'", i, test.expectedError, err.Error())
		}
	}
}
