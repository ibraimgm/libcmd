package libcfg_test

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"testing"

	"github.com/ibraimgm/libcfg"
)

func withEnv(env map[string]string, handler func()) {
	for k, v := range env {
		if err := os.Setenv(k, v); err != nil {
			panic(err)
		}
	}

	defer func() {
		for k := range env {
			if err := os.Unsetenv(k); err != nil {
				panic(err)
			}
		}
	}()

	handler()
}

func withFileEnv(env map[string]string, handler func(string)) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}

	file, err := ioutil.TempFile(dir, "")
	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(dir)

	for k, v := range env {
		line := fmt.Sprintf("%s=%s\n", k, v)

		if _, err := file.WriteString(line); err != nil {
			panic(err)
		}
	}

	handler(file.Name())
}

func TestEnv(t *testing.T) {
	tests := []struct {
		env      map[string]string
		abool    bool
		aint     int
		auint    uint
		afloat32 float32
		afloat64 float64
		astring  string
	}{
		{env: map[string]string{}},
		{env: map[string]string{
			"B1":  "true",
			"I1":  "5",
			"U1":  "9",
			"F1":  "3.14",
			"LF1": "3.1415",
			"S1":  "foo",
		}, abool: true, aint: 5, auint: 9, afloat32: float32(3.14), afloat64: float64(3.1415), astring: "foo"},
		{env: map[string]string{
			"B2":  "true",
			"I2":  "6",
			"U2":  "10",
			"F1":  "3.14",
			"LF1": "3.1415",
			"S2":  "bar",
		}, abool: true, aint: 6, auint: 10, afloat32: float32(3.14), afloat64: float64(3.1415), astring: "bar"},
		{env: map[string]string{
			"B3": "true",
			"I3": "7",
			"U3": "11",
			"S3": "baz",
		}, abool: true, aint: 7, auint: 11, astring: "baz"},
		{env: map[string]string{
			"F1": "3.14",
			"F2": "3.15",
		}, afloat32: float32(3.15)},
		{env: map[string]string{
			"B1": "false",
			"B2": "true",
		}, abool: true},
		{env: map[string]string{
			"B1": "true",
			"B2": "false",
		}},
		{env: map[string]string{
			"B1": "true",
			"B2": "false",
			"B3": "true",
		}, abool: true},
		{env: map[string]string{
			"I1": "1",
			"I2": "2",
		}, aint: 2},
		{env: map[string]string{
			"I1": "1",
			"I2": "0",
		}},
		{env: map[string]string{
			"I1": "1",
			"I2": "0",
			"I3": "3",
		}, aint: 3},
		{env: map[string]string{
			"I1": "1",
			"I2": "2",
			"I3": "",
		}},
		{env: map[string]string{
			"U1": "1",
			"U2": "2",
		}, auint: 2},
		{env: map[string]string{
			"U1": "1",
			"U2": "0",
		}},
		{env: map[string]string{
			"U1": "1",
			"U2": "0",
			"U3": "3",
		}, auint: 3},
		{env: map[string]string{
			"U1": "1",
			"U2": "2",
			"U3": "",
		}},
		{env: map[string]string{
			"S1": "foo",
			"S2": "bar",
		}, astring: "bar"},
		{env: map[string]string{
			"S1": "1",
			"S2": "",
		}},
		{env: map[string]string{
			"S1": "foo",
			"S2": "",
			"S3": "baz",
		}, astring: "baz"},
		{env: map[string]string{
			"S1": "foo",
			"S2": "bar",
			"S3": "",
		}},
	}

	for i, test := range tests {
		env := libcfg.NewEnvLoader()

		abool := env.Bool(false, "B1", "B2", "B3")
		aint := env.Int(0, "I1", "I2", "I3")
		auint := env.Uint(0, "U1", "U2", "U3")
		astring := env.String("", "S1", "S2", "S3")
		afloat32 := env.Float32(0, "F1", "F2")
		afloat64 := env.Float64(0, "LF1", "LF2")

		i := i       // pin scope
		test := test // pin scope

		withEnv(test.env, func() {
			env.LoadAll()

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
		})
	}
}

func TestEnvDefault(t *testing.T) {
	tests := []struct {
		env     map[string]string
		abool   bool
		aint    int
		auint   uint
		astring string
	}{
		{env: map[string]string{}, abool: true, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"B1": "true",
			"I1": "5",
			"U1": "9",
			"S1": "foo",
		}, abool: true, aint: 5, auint: 9, astring: "foo"},
		{env: map[string]string{
			"B2": "true",
			"I2": "6",
			"U2": "10",
			"S2": "bar",
		}, abool: true, aint: 6, auint: 10, astring: "bar"},
		{env: map[string]string{
			"B3": "true",
			"I3": "7",
			"U3": "11",
			"S3": "baz",
		}, abool: true, aint: 7, auint: 11, astring: "baz"},
		{env: map[string]string{
			"B1": "false",
			"B2": "true",
		}, abool: true, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"B1": "true",
			"B2": "false",
		}, abool: false, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"B1": "true",
			"B2": "false",
			"B3": "true",
		}, abool: true, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"I1": "1",
			"I2": "2",
		}, abool: true, aint: 2, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"I1": "1",
			"I2": "0",
		}, abool: true, aint: 0, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"I1": "1",
			"I2": "0",
			"I3": "3",
		}, abool: true, aint: 3, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"I1": "1",
			"I2": "2",
			"I3": "",
		}, abool: true, aint: 0, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"U1": "1",
			"U2": "2",
		}, abool: true, aint: 8, auint: 2, astring: "xyz"},
		{env: map[string]string{
			"U1": "1",
			"U2": "0",
		}, abool: true, aint: 8, auint: 0, astring: "xyz"},
		{env: map[string]string{
			"U1": "1",
			"U2": "0",
			"U3": "3",
		}, abool: true, aint: 8, auint: 3, astring: "xyz"},
		{env: map[string]string{
			"U1": "1",
			"U2": "2",
			"U3": "",
		}, abool: true, aint: 8, auint: 0, astring: "xyz"},
		{env: map[string]string{
			"S1": "foo",
			"S2": "bar",
		}, abool: true, aint: 8, auint: 16, astring: "bar"},
		{env: map[string]string{
			"S1": "1",
			"S2": "",
		}, abool: true, aint: 8, auint: 16, astring: ""},
		{env: map[string]string{
			"S1": "foo",
			"S2": "",
			"S3": "baz",
		}, abool: true, aint: 8, auint: 16, astring: "baz"},
		{env: map[string]string{
			"S1": "foo",
			"S2": "bar",
			"S3": "",
		}, abool: true, aint: 8, auint: 16, astring: ""},
	}

	for i, test := range tests {
		env := libcfg.NewEnvLoader()

		abool := env.Bool(true, "B1", "B2", "B3")
		aint := env.Int(8, "I1", "I2", "I3")
		auint := env.Uint(16, "U1", "U2", "U3")
		astring := env.String("xyz", "S1", "S2", "S3")

		i := i       // pin scope
		test := test // pin scope

		withEnv(test.env, func() {
			env.LoadAll()

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
		})
	}
}

func TestEnvIntLimit(t *testing.T) {
	tests := []struct {
		env map[string]string
		a   int8
		b   int16
		c   int32
		d   int64
	}{
		{env: map[string]string{"A": strconv.FormatInt(math.MaxInt8, 10)}, a: math.MaxInt8},
		{env: map[string]string{"B": strconv.FormatInt(math.MaxInt16, 10)}, b: math.MaxInt16},
		{env: map[string]string{"C": strconv.FormatInt(math.MaxInt32, 10)}, c: math.MaxInt32},
		{env: map[string]string{"D": strconv.FormatInt(math.MaxInt64, 10)}, d: math.MaxInt64},
	}

	for i, test := range tests {
		env := libcfg.NewEnvLoader()

		a := env.Int8(0, "A")
		b := env.Int16(0, "B")
		c := env.Int32(0, "C")
		d := env.Int64(0, "D")

		i := i       //pin
		test := test //pin

		withEnv(test.env, func() {
			env.LoadAll()

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
		})
	}
}

func TestEnvUintLimit(t *testing.T) {
	tests := []struct {
		env map[string]string
		a   uint8
		b   uint16
		c   uint32
		d   uint64
	}{
		{env: map[string]string{"A": strconv.FormatUint(math.MaxUint8, 10)}, a: math.MaxUint8},
		{env: map[string]string{"B": strconv.FormatUint(math.MaxUint16, 10)}, b: math.MaxUint16},
		{env: map[string]string{"C": strconv.FormatUint(math.MaxUint32, 10)}, c: math.MaxUint32},
		{env: map[string]string{"D": strconv.FormatUint(math.MaxUint64, 10)}, d: math.MaxUint64},
	}

	for i, test := range tests {
		env := libcfg.NewEnvLoader()

		a := env.Uint8(0, "A")
		b := env.Uint16(0, "B")
		c := env.Uint32(0, "C")
		d := env.Uint64(0, "D")

		i := i       //pin
		test := test //pin

		withEnv(test.env, func() {
			env.LoadAll()

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
		})
	}
}

func TestEnvFile(t *testing.T) {
	tests := []struct {
		env     map[string]string
		abool   bool
		aint    int
		auint   uint
		astring string
	}{
		{env: map[string]string{}},
		{env: map[string]string{
			"B1": "true",
			"I1": "5",
			"U1": "9",
			"S1": "foo",
		}, abool: true, aint: 5, auint: 9, astring: "foo"},
		{env: map[string]string{
			"B2": "true",
			"I2": "6",
			"U2": "10",
			"S2": "bar",
		}, abool: true, aint: 6, auint: 10, astring: "bar"},
		{env: map[string]string{
			"B3": "true",
			"I3": "7",
			"U3": "11",
			"S3": "baz",
		}, abool: true, aint: 7, auint: 11, astring: "baz"},
		{env: map[string]string{
			"B1": "false",
			"B2": "true",
		}, abool: true},
		{env: map[string]string{
			"B1": "true",
			"B2": "false",
		}},
		{env: map[string]string{
			"B1": "true",
			"B2": "false",
			"B3": "true",
		}, abool: true},
		{env: map[string]string{
			"I1": "1",
			"I2": "2",
		}, aint: 2},
		{env: map[string]string{
			"I1": "1",
			"I2": "0",
		}},
		{env: map[string]string{
			"I1": "1",
			"I2": "0",
			"I3": "3",
		}, aint: 3},
		{env: map[string]string{
			"I1": "1",
			"I2": "2",
			"I3": "",
		}},
		{env: map[string]string{
			"U1": "1",
			"U2": "2",
		}, auint: 2},
		{env: map[string]string{
			"U1": "1",
			"U2": "0",
		}},
		{env: map[string]string{
			"U1": "1",
			"U2": "0",
			"U3": "3",
		}, auint: 3},
		{env: map[string]string{
			"U1": "1",
			"U2": "2",
			"U3": "",
		}},
		{env: map[string]string{
			"S1": "foo",
			"S2": "bar",
		}, astring: "bar"},
		{env: map[string]string{
			"S1": "1",
			"S2": "",
		}},
		{env: map[string]string{
			"S1": "foo",
			"S2": "",
			"S3": "baz",
		}, astring: "baz"},
		{env: map[string]string{
			"S1": "foo",
			"S2": "bar",
			"S3": "",
		}},
	}

	for i, test := range tests {
		env := libcfg.NewEnvLoader()

		abool := env.Bool(false, "B1", "B2", "B3")
		aint := env.Int(0, "I1", "I2", "I3")
		auint := env.Uint(0, "U1", "U2", "U3")
		astring := env.String("", "S1", "S2", "S3")

		i := i       // pin scope
		test := test // pin scope

		withFileEnv(test.env, func(filename string) {
			if err := env.UseFile(filename); err != nil {
				t.Errorf("Case %d, error loading file: %v", i, err)
				return
			}

			env.LoadAll()

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
		})
	}
}

func TestEnvFileDefault(t *testing.T) {
	tests := []struct {
		env     map[string]string
		abool   bool
		aint    int
		auint   uint
		astring string
	}{
		{env: map[string]string{}, abool: true, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"B1": "true",
			"I1": "5",
			"U1": "9",
			"S1": "foo",
		}, abool: true, aint: 5, auint: 9, astring: "foo"},
		{env: map[string]string{
			"B2": "true",
			"I2": "6",
			"U2": "10",
			"S2": "bar",
		}, abool: true, aint: 6, auint: 10, astring: "bar"},
		{env: map[string]string{
			"B3": "true",
			"I3": "7",
			"U3": "11",
			"S3": "baz",
		}, abool: true, aint: 7, auint: 11, astring: "baz"},
		{env: map[string]string{
			"B1": "false",
			"B2": "true",
		}, abool: true, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"B1": "true",
			"B2": "false",
		}, abool: false, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"B1": "true",
			"B2": "false",
			"B3": "true",
		}, abool: true, aint: 8, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"I1": "1",
			"I2": "2",
		}, abool: true, aint: 2, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"I1": "1",
			"I2": "0",
		}, abool: true, aint: 0, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"I1": "1",
			"I2": "0",
			"I3": "3",
		}, abool: true, aint: 3, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"I1": "1",
			"I2": "2",
			"I3": "",
		}, abool: true, aint: 0, auint: 16, astring: "xyz"},
		{env: map[string]string{
			"U1": "1",
			"U2": "2",
		}, abool: true, aint: 8, auint: 2, astring: "xyz"},
		{env: map[string]string{
			"U1": "1",
			"U2": "0",
		}, abool: true, aint: 8, auint: 0, astring: "xyz"},
		{env: map[string]string{
			"U1": "1",
			"U2": "0",
			"U3": "3",
		}, abool: true, aint: 8, auint: 3, astring: "xyz"},
		{env: map[string]string{
			"U1": "1",
			"U2": "2",
			"U3": "",
		}, abool: true, aint: 8, auint: 0, astring: "xyz"},
		{env: map[string]string{
			"S1": "foo",
			"S2": "bar",
		}, abool: true, aint: 8, auint: 16, astring: "bar"},
		{env: map[string]string{
			"S1": "1",
			"S2": "",
		}, abool: true, aint: 8, auint: 16, astring: ""},
		{env: map[string]string{
			"S1": "foo",
			"S2": "",
			"S3": "baz",
		}, abool: true, aint: 8, auint: 16, astring: "baz"},
		{env: map[string]string{
			"S1": "foo",
			"S2": "bar",
			"S3": "",
		}, abool: true, aint: 8, auint: 16, astring: ""},
	}

	for i, test := range tests {
		env := libcfg.NewEnvLoader()

		abool := env.Bool(true, "B1", "B2", "B3")
		aint := env.Int(8, "I1", "I2", "I3")
		auint := env.Uint(16, "U1", "U2", "U3")
		astring := env.String("xyz", "S1", "S2", "S3")

		i := i       // pin scope
		test := test // pin scope

		withFileEnv(test.env, func(filename string) {
			if err := env.UseFile(filename); err != nil {
				t.Errorf("Case %d, error loading file: %v", i, err)
				return
			}

			env.LoadAll()

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
		})
	}
}

func TestEnvSpecialCases(t *testing.T) {
	tests := []struct {
		env    map[string]string
		result string
	}{
		{env: map[string]string{"A": "foo"}, result: "foo"},
		{env: map[string]string{"A": ""}, result: ""},
		{env: map[string]string{"A": "   foo"}, result: "   foo"},
		{env: map[string]string{"A": "foo   "}, result: "foo"},
		{env: map[string]string{"a": "foo"}, result: ""},
		{env: map[string]string{"A=": "foo"}, result: "=foo"},
		{env: map[string]string{"A": "=foo"}, result: "=foo"},
		{env: map[string]string{"A": "=fo=o="}, result: "=fo=o="},
		{env: map[string]string{"A": "#foo"}, result: ""},
		{env: map[string]string{"A#": "foo"}, result: ""},
		{env: map[string]string{"#A": "foo"}, result: ""},
		{env: map[string]string{"A": "foo#bar"}, result: "foo"},
		{env: map[string]string{"A": "foo #bar"}, result: "foo"},
	}

	for i, test := range tests {
		env := libcfg.NewEnvLoader()

		a := env.String("", "A")

		i := i       //pin
		test := test //pin

		withFileEnv(test.env, func(filename string) {
			if err := env.UseFile(filename); err != nil {
				t.Errorf("Case %d, error loading file: %v", i, err)
				return
			}

			env.LoadAll()

			if *a != test.result {
				t.Errorf("Case %d, wrong string value: expected '%v', received '%v'", i, test.result, *a)
			}
		})
	}
}

func TestEnvMultipleFiles(t *testing.T) {
	tests := []struct {
		file1  map[string]string
		file2  map[string]string
		result string
	}{
		{file1: map[string]string{"A1": "foo", "A2": "bar"}, file2: map[string]string{"A3": "baz"}, result: "baz"},
		{file1: map[string]string{"A1": "foo", "A3": "baz"}, file2: map[string]string{"A2": "bar"}, result: "baz"},
		{file1: map[string]string{"A1": "foo", "A2": "bar"}, file2: map[string]string{"A1": "baz"}, result: "bar"},
		{file1: map[string]string{"A1": "foo", "A3": "bar"}, file2: map[string]string{"A3": "baz"}, result: "baz"},
		{file1: map[string]string{"A1": "foo"}, file2: map[string]string{"A1": "bar"}, result: "bar"},
		{file1: map[string]string{"x": "y"}, file2: map[string]string{"A3": ""}, result: ""},
	}

	for i, test := range tests {
		env := libcfg.NewEnvLoader()

		a := env.String("default", "A1", "A2", "A3")

		i := i       //pin
		test := test //pin

		withFileEnv(test.file1, func(file1 string) {
			withFileEnv(test.file2, func(file2 string) {
				env.UseFiles(file1, "does-not-exist", file2)

				env.LoadAll()

				if *a != test.result {
					t.Errorf("Case %d, wrong string value: expected '%s', received '%s'", i, test.result, *a)
				}
			})
		})
	}
}

func TestEnvFileOnly(t *testing.T) {
	tests := []struct {
		env    map[string]string
		file   map[string]string
		result string
	}{
		{env: map[string]string{"A1": "foo", "A2": "bar"}, file: map[string]string{"A3": "baz"}, result: "baz"},
		{env: map[string]string{"A1": "foo", "A3": "baz"}, file: map[string]string{"A2": "bar"}, result: "bar"},
		{env: map[string]string{"A1": "foo", "A2": "bar"}, file: map[string]string{"A1": "baz"}, result: "baz"},
		{env: map[string]string{"A3": "bar"}, file: map[string]string{"A3": "baz", "A1": "foo"}, result: "baz"},
		{env: map[string]string{"A1": "foo"}, file: map[string]string{"A1": "bar"}, result: "bar"},
		{env: map[string]string{"A3": "foo"}, file: map[string]string{}, result: "default"},
		{env: map[string]string{"A3": "foo"}, file: map[string]string{"A3": ""}, result: ""},
	}

	for i, test := range tests {
		env := libcfg.NewEnvLoader()

		a := env.String("default", "A1", "A2", "A3")

		i := i       //pin
		test := test //pin

		withEnv(test.env, func() {
			withFileEnv(test.file, func(file string) {
				env.UseEnv(false)

				if err := env.UseFile(file); err != nil {
					t.Errorf("Case %d, error loading env from file: %v", i, err)
				}

				env.LoadAll()

				if *a != test.result {
					t.Errorf("Case %d, wrong string value: expected '%s', received '%s'", i, test.result, *a)
				}
			})
		})
	}
}

func TestEnvKeepValue(t *testing.T) {
	const keep = "keep"

	tests := []struct {
		env  map[string]string
		file map[string]string
		s1   string
		s2   string
	}{
		{env: map[string]string{}, file: map[string]string{}, s1: keep, s2: "default"},
		{env: map[string]string{"A": "a"}, file: map[string]string{}, s1: "a", s2: "default"},
		{env: map[string]string{}, file: map[string]string{"A": "a"}, s1: "a", s2: "default"},
		{env: map[string]string{"B": "b"}, file: map[string]string{"A": "a"}, s1: "b", s2: "default"},
		{env: map[string]string{"C": "", "B": "b"}, file: map[string]string{"A": "a"}, s1: "", s2: "default"},
		{env: map[string]string{}, file: map[string]string{}, s1: keep, s2: "default"},
		{env: map[string]string{"X": "x"}, file: map[string]string{}, s1: keep, s2: "x"},
		{env: map[string]string{}, file: map[string]string{"X": "x"}, s1: keep, s2: "x"},
		{env: map[string]string{"Y": "y"}, file: map[string]string{"X": "x"}, s1: keep, s2: "y"},
		{env: map[string]string{"Z": "", "Y": "Y"}, file: map[string]string{"X": "x"}, s1: keep, s2: ""},
	}

	for i, test := range tests {
		env := libcfg.NewEnvLoader()
		s1 := env.String("", "A", "B", "C")
		s2 := env.String("default", "X", "Y", "Z")

		*s1 = keep
		*s2 = ""

		i := i       //pin
		test := test //pin

		withEnv(test.env, func() {
			withFileEnv(test.file, func(file string) {
				if err := env.UseFile(file); err != nil {
					t.Errorf("Case %d, error loading file: %v", i, err)
					return
				}

				env.LoadAll()

				if *s1 != test.s1 {
					t.Errorf("Case %d, error on string1 value: expected '%s', received '%s'", i, test.s1, *s1)
				}

				if *s2 != test.s2 {
					t.Errorf("Case %d, error on string2 value: expected '%s', received '%s'", i, test.s2, *s2)
				}
			})
		})
	}
}
