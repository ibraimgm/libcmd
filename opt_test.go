package libcfg_test

import (
	"math"
	"strconv"
	"strings"
	"testing"

	"github.com/ibraimgm/libcfg"
)

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
		p := libcfg.NewParser()

		abool := p.Bool("abool", "b", false, "specifies a bool value")
		aint := p.Int("aint", "i", 0, "specifies an int value")
		auint := p.Uint("auint", "u", 0, "specifies an uint value")
		astring := p.String("astring", "s", "", "specifies a string value")
		afloat32 := p.Float32("afloat32", "f32", 0, "specifies a float32 value")
		afloat64 := p.Float64("afloat64", "f64", 0, "specifies a float64 value")

		if err := p.RunArgs(test.cmd); err != nil {
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

		args := p.Args()

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
		p := libcfg.NewParser()

		abool := p.Bool("abool", "b", true, "specifies a bool value")
		aint := p.Int("aint", "i", 8, "specifies an int value")
		auint := p.Uint("auint", "u", 16, "specifies an uint value")
		astring := p.String("astring", "s", "default", "specifies a string value")
		afloat32 := p.Float32("afloat32", "f32", float32(3.14), "specifies a float32 value")
		afloat64 := p.Float64("afloat64", "f64", float64(3.1415), "specifies a float64 value")

		if err := p.RunArgs(test.cmd); err != nil {
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

		args := p.Args()

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
		p := libcfg.NewParser()

		p.Bool("abool", "b", false, "specifies a bool value")
		p.Int("aint", "i", 0, "specifies an int value")
		p.Uint("auint", "u", 0, "specifies an uint value")
		p.String("astring", "s", "", "specifies a string value")

		err := p.RunArgs(test.cmd)

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
		p := libcfg.NewParser()

		a := p.Int8("", "a", 0, "specifies a int8 value")
		b := p.Int16("", "b", 0, "specifies a int16 value")
		c := p.Int32("", "c", 0, "specifies a int32 value")
		d := p.Int64("", "d", 0, "specifies a int64 value")

		if err := p.RunArgs(test.cmd); err != nil {
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
		p := libcfg.NewParser()

		a := p.Uint8("", "a", 0, "specifies a uint8 value")
		b := p.Uint16("", "b", 0, "specifies a uint16 value")
		c := p.Uint32("", "c", 0, "specifies a uint32 value")
		d := p.Uint64("", "d", 0, "specifies a uint64 value")

		if err := p.RunArgs(test.cmd); err != nil {
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
		p := libcfg.NewParser()

		p.Int8("aint8", "", 0, "specifies a int8 value")
		p.Int16("aint16", "", 0, "specifies a int16 value")
		p.Int32("aint32", "", 0, "specifies a int32 value")
		p.Int64("aint64", "", 0, "specifies a int64 value")
		p.Uint8("auint8", "", 0, "specifies a uint8 value")
		p.Uint16("auint16", "", 0, "specifies a uint16 value")
		p.Uint32("auint32", "", 0, "specifies a uint32 value")
		p.Uint64("auint64", "", 0, "specifies a uint64 value")

		err := p.RunArgs(test.cmd)

		if err == nil {
			t.Errorf("Case %d, argument parsing should return error", i)
			continue
		}

		if !strings.Contains(err.Error(), test.expectedError) {
			t.Errorf("Case %d, expected error '%s', but got '%s'", i, test.expectedError, err.Error())
		}
	}
}

func TestOptStrict(t *testing.T) {
	tests := []struct {
		cmd         []string
		args        []string
		expectError bool
	}{
		{cmd: []string{}, args: []string{}, expectError: false},
		{cmd: []string{"-s", "foo"}, args: []string{}, expectError: false},
		{cmd: []string{"-s", "foo", "bar"}, args: []string{"bar"}, expectError: true},
		{cmd: []string{"bar", "-s", "foo"}, args: []string{"bar", "-s", "foo"}, expectError: true},
		{cmd: []string{"bar", "-s", "foo", "baz"}, args: []string{"bar", "-s", "foo", "baz"}, expectError: true},
	}

	for i, test := range tests {
		p := libcfg.NewParser()
		p.Configure(libcfg.Options{StrictParsing: true})
		p.String("", "s", "", "")

		err := p.RunArgs(test.cmd)

		if test.expectError && err == nil {
			t.Errorf("Case %d should have returned an error", i)
		}

		if !test.expectError && err != nil {
			t.Errorf("Case %d, error running parser: %v", i, err)
		}

		args := p.Args()

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

func TestOptGreedy(t *testing.T) {
	tests := []struct {
		cmd     []string
		args    []string
		abool   bool
		aint    int
		astring string
		c1      bool
	}{
		{cmd: []string{}, args: []string{}},
		{cmd: []string{"-b"}, args: []string{}, abool: true},
		{cmd: []string{"-b", "foo"}, args: []string{"foo"}, abool: true},
		{cmd: []string{"-b", "foo", "-i", "5"}, args: []string{"foo"}, abool: true, aint: 5},
		{cmd: []string{"foo", "-b", "-i", "5"}, args: []string{"foo"}, abool: true, aint: 5},
		{cmd: []string{"-b", "-i", "5", "foo"}, args: []string{"foo"}, abool: true, aint: 5},
		{cmd: []string{"-b", "foo", "bar", "-i", "5"}, args: []string{"foo", "bar"}, abool: true, aint: 5},
		{cmd: []string{"foo", "bar", "-b", "-i", "5"}, args: []string{"foo", "bar"}, abool: true, aint: 5},
		{cmd: []string{"-b", "-i", "5", "foo", "bar"}, args: []string{"foo", "bar"}, abool: true, aint: 5},
		{cmd: []string{"-i", "5"}, args: []string{}, aint: 5},
		{cmd: []string{"-i", "5", "foo"}, args: []string{"foo"}, aint: 5},
		{cmd: []string{"foo", "-i", "5"}, args: []string{"foo"}, aint: 5},
		{cmd: []string{"-s", "foo"}, args: []string{}, astring: "foo"},
		{cmd: []string{"-s", "foo", "bar"}, args: []string{"bar"}, astring: "foo"},
		{cmd: []string{"--str=foo", "bar"}, args: []string{"bar"}, astring: "foo"},
		{cmd: []string{"--str=", "bar"}, args: []string{"bar"}},
		{cmd: []string{"baz", "-s", "foo", "bar"}, args: []string{"baz", "bar"}, astring: "foo"},
		{cmd: []string{"baz", "--str=foo", "bar"}, args: []string{"baz", "bar"}, astring: "foo"},
		{cmd: []string{"-b", "foo", "c1"}, args: []string{"foo"}, abool: true, c1: true},
		{cmd: []string{"-b", "c1", "foo"}, args: []string{"foo"}, abool: true, c1: true},
		{cmd: []string{"-s", "foo", "bar", "c1"}, args: []string{"bar"}, astring: "foo", c1: true},
		{cmd: []string{"-s", "c1", "bar", "foo"}, args: []string{"bar", "foo"}, astring: "c1", c1: false},
		{cmd: []string{"baz", "--str=foo", "bar", "c1", "baq"}, args: []string{"baz", "bar", "baq"}, astring: "foo", c1: true},
	}

	for i, test := range tests {
		p := libcfg.NewParser()
		p.Configure(libcfg.Options{Greedy: true})

		abool := p.Bool("", "b", false, "")
		aint := p.Int("", "i", 0, "")
		astring := p.String("str", "s", "", "")
		c1 := p.AddCommand("c1", "command 1")

		if err := p.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error parsing args: %v", i, err)
			continue
		}

		if *abool != test.abool {
			t.Errorf("Case %d, wrong boolean value: expected '%v', received '%v'", i, test.abool, *abool)
		}

		if *aint != test.aint {
			t.Errorf("Case %d, wrong int value: expected '%v', received '%v'", i, test.aint, *aint)
		}

		if *astring != test.astring {
			t.Errorf("Case %d, wrong string value: expected '%v', received '%v'", i, test.astring, *astring)
		}

		if c1.Used() != test.c1 {
			t.Errorf("Case %d, wrong used value for command: expected '%v', received '%v'", i, test.c1, c1.Used())
		}

		args := p.Args()

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
		p := libcfg.NewParser()

		s1 := p.String("string1", "s1", "", "")
		s2 := p.String("string2", "s2", "default", "")

		*s1 = keep
		*s2 = ""

		if err := p.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error parsing args: %v", i, err)
			continue
		}

		if *s1 != test.s1 {
			t.Errorf("Case %d, wrong string1 value: expected '%s', received '%s'", i, test.s1, *s1)
		}

		if *s2 != test.s2 {
			t.Errorf("Case %d, wrong string2 value: expected '%s', received '%s'", i, test.s2, *s2)
		}
	}
}
