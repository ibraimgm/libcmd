package libcfg_test

import (
	"testing"

	"github.com/ibraimgm/libcfg"
)

func TestEnvOpt(t *testing.T) {
	tests := []struct {
		env      map[string]string
		cmd      []string
		abool    bool
		astring  string
		aint     int
		auint    uint
		afloat32 float32
		afloat64 float64
	}{
		{env: map[string]string{}, cmd: []string{}},
		{
			env:   map[string]string{"B1": "true", "S1": "foo", "I1": "5", "U1": "7", "F1": "3.14", "D1": "3.1415"},
			cmd:   []string{},
			abool: true, astring: "foo", aint: 5, auint: 7, afloat32: float32(3.14), afloat64: float64(3.1415),
		},
		{
			env:   map[string]string{},
			cmd:   []string{"-b", "-s", "foo", "-i", "5", "-u", "7"},
			abool: true, astring: "foo", aint: 5, auint: 7,
		},
		{env: map[string]string{"B1": "true"}, cmd: []string{"--no-abool"}},
		{env: map[string]string{"S1": "foo"}, cmd: []string{"-s", "bar"}, astring: "bar"},
		{env: map[string]string{"S1": "foo", "S2": ""}, cmd: []string{}, astring: ""},
		{env: map[string]string{"S1": "foo"}, cmd: []string{"--astring="}},
		{env: map[string]string{"I1": "5", "I2": "aaaa"}, cmd: []string{}, aint: 5},
		{env: map[string]string{"I1": "5", "I2": ""}, cmd: []string{"-i", "10"}, aint: 10},
		{env: map[string]string{"I1": "5", "I2": ""}, cmd: []string{"-i", "10", "--aint", "15"}, aint: 15},
		{env: map[string]string{"U1": "5", "U2": "aaaa"}, cmd: []string{}, auint: 5},
		{env: map[string]string{"U1": "5", "U2": ""}, cmd: []string{"-u", "10"}, auint: 10},
		{env: map[string]string{"U1": "5", "U2": ""}, cmd: []string{"-u", "10", "--auint", "15"}, auint: 15},
	}

	for i, test := range tests {
		p := libcfg.NewParser()

		abool := p.Bool("abool", "b", false, "", "B1", "B2")
		astring := p.String("astring", "s", "", "", "S1", "S2")
		aint := p.Int("aint", "i", 0, "", "I1", "I2")
		auint := p.Uint("auint", "u", 0, "", "U1", "U2")
		afloat32 := p.Float32("afloat32", "f32", 0, "", "F1")
		afloat64 := p.Float64("afloat64", "f64", 0, "", "D1")

		i := i       //pin
		test := test //pin

		withEnv(test.env, func() {
			if err := p.RunArgs(test.cmd); err != nil {
				t.Errorf("Case %d, error running parser: %v", i, err)
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
		})
	}
}

func TestEnvOptFile(t *testing.T) {
	tests := []struct {
		env      map[string]string
		cmd      []string
		abool    bool
		astring  string
		aint     int
		auint    uint
		afloat32 float32
		afloat64 float64
	}{
		{env: map[string]string{}, cmd: []string{}},
		{
			env:   map[string]string{"B1": "true", "S1": "foo", "I1": "5", "U1": "7", "F1": "3.14", "D1": "3.1415"},
			cmd:   []string{},
			abool: true, astring: "foo", aint: 5, auint: 7, afloat32: float32(3.14), afloat64: float64(3.1415),
		},
		{
			env:   map[string]string{},
			cmd:   []string{"-b", "-s", "foo", "-i", "5", "-u", "7"},
			abool: true, astring: "foo", aint: 5, auint: 7,
		},
		{env: map[string]string{"B1": "true"}, cmd: []string{"--no-abool"}},
		{env: map[string]string{"S1": "foo"}, cmd: []string{"-s", "bar"}, astring: "bar"},
		{env: map[string]string{"S1": "foo", "S2": ""}, cmd: []string{}, astring: ""},
		{env: map[string]string{"S1": "foo"}, cmd: []string{"--astring="}},
		{env: map[string]string{"I1": "5", "I2": "aaaa"}, cmd: []string{}, aint: 5},
		{env: map[string]string{"I1": "5", "I2": ""}, cmd: []string{"-i", "10"}, aint: 10},
		{env: map[string]string{"I1": "5", "I2": ""}, cmd: []string{"-i", "10", "--aint", "15"}, aint: 15},
		{env: map[string]string{"U1": "5", "U2": "aaaa"}, cmd: []string{}, auint: 5},
		{env: map[string]string{"U1": "5", "U2": ""}, cmd: []string{"-u", "10"}, auint: 10},
		{env: map[string]string{"U1": "5", "U2": ""}, cmd: []string{"-u", "10", "--auint", "15"}, auint: 15},
	}

	for i, test := range tests {
		p := libcfg.NewParser()

		abool := p.Bool("abool", "b", false, "", "B1", "B2")
		astring := p.String("astring", "s", "", "", "S1", "S2")
		aint := p.Int("aint", "i", 0, "", "I1", "I2")
		auint := p.Uint("auint", "u", 0, "", "U1", "U2")
		afloat32 := p.Float32("afloat32", "f32", 0, "", "F1")
		afloat64 := p.Float64("afloat64", "f64", 0, "", "D1")

		i := i       //pin
		test := test //pin

		withFileEnv(test.env, func(filename string) {
			p.Configure(libcfg.Options{FilesOnly: true})

			if err := p.UseFile(filename); err != nil {
				t.Errorf("Case %d, error loading env file: %v", i, err)
			}

			if err := p.RunArgs(test.cmd); err != nil {
				t.Errorf("Case %d, error running parser: %v", i, err)
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
		})
	}
}

func TestEnvOptInt(t *testing.T) {
	tests := []struct {
		env map[string]string
		cmd []string
		a   int8
		b   int16
		c   int32
		d   int64
	}{
		{env: map[string]string{}, cmd: []string{}, a: 8, b: 16, c: 32, d: 64},
		{env: map[string]string{"A": "80"}, cmd: []string{}, a: 80, b: 16, c: 32, d: 64},
		{env: map[string]string{"A": "80"}, cmd: []string{"-a", "81"}, a: 81, b: 16, c: 32, d: 64},
		{env: map[string]string{"A": "wrong"}, cmd: []string{}, a: 8, b: 16, c: 32, d: 64},
		{env: map[string]string{"A": "wrong"}, cmd: []string{"-a", "81"}, a: 81, b: 16, c: 32, d: 64},
		{env: map[string]string{"A": "0"}, cmd: []string{}, a: 0, b: 16, c: 32, d: 64},
		{env: map[string]string{"A": ""}, cmd: []string{}, a: 0, b: 16, c: 32, d: 64},
		{env: map[string]string{"A": ""}, cmd: []string{"-a", "80"}, a: 80, b: 16, c: 32, d: 64},
		{env: map[string]string{"A": "80"}, cmd: []string{"-a", "0"}, a: 0, b: 16, c: 32, d: 64},
		{env: map[string]string{"A": "wrong"}, cmd: []string{"-a", "0"}, a: 0, b: 16, c: 32, d: 64},
		{env: map[string]string{"B": "160"}, cmd: []string{}, a: 8, b: 160, c: 32, d: 64},
		{env: map[string]string{"B": "160"}, cmd: []string{"-b", "161"}, a: 8, b: 161, c: 32, d: 64},
		{env: map[string]string{"B": "wrong"}, cmd: []string{}, a: 8, b: 16, c: 32, d: 64},
		{env: map[string]string{"B": "wrong"}, cmd: []string{"-b", "161"}, a: 8, b: 161, c: 32, d: 64},
		{env: map[string]string{"B": "0"}, cmd: []string{}, a: 8, b: 0, c: 32, d: 64},
		{env: map[string]string{"B": ""}, cmd: []string{}, a: 8, b: 0, c: 32, d: 64},
		{env: map[string]string{"B": ""}, cmd: []string{"-b", "161"}, a: 8, b: 161, c: 32, d: 64},
		{env: map[string]string{"B": "160"}, cmd: []string{"-b", "0"}, a: 8, b: 0, c: 32, d: 64},
		{env: map[string]string{"B": "wrong"}, cmd: []string{"-b", "0"}, a: 8, b: 0, c: 32, d: 64},
		{env: map[string]string{"C": "320"}, cmd: []string{}, a: 8, b: 16, c: 320, d: 64},
		{env: map[string]string{"C": "320"}, cmd: []string{"-c", "321"}, a: 8, b: 16, c: 321, d: 64},
		{env: map[string]string{"C": "wrong"}, cmd: []string{}, a: 8, b: 16, c: 32, d: 64},
		{env: map[string]string{"C": "wrong"}, cmd: []string{"-c", "321"}, a: 8, b: 16, c: 321, d: 64},
		{env: map[string]string{"C": "0"}, cmd: []string{}, a: 8, b: 16, c: 0, d: 64},
		{env: map[string]string{"C": ""}, cmd: []string{}, a: 8, b: 16, c: 0, d: 64},
		{env: map[string]string{"C": ""}, cmd: []string{"-c", "321"}, a: 8, b: 16, c: 321, d: 64},
		{env: map[string]string{"C": "320"}, cmd: []string{"-c", "0"}, a: 8, b: 16, c: 0, d: 64},
		{env: map[string]string{"C": "wrong"}, cmd: []string{"-c", "0"}, a: 8, b: 16, c: 0, d: 64},
		{env: map[string]string{"D": "640"}, cmd: []string{}, a: 8, b: 16, c: 32, d: 640},
		{env: map[string]string{"D": "640"}, cmd: []string{"-d", "641"}, a: 8, b: 16, c: 32, d: 641},
		{env: map[string]string{"D": "wrong"}, cmd: []string{}, a: 8, b: 16, c: 32, d: 64},
		{env: map[string]string{"D": "wrong"}, cmd: []string{"-d", "641"}, a: 8, b: 16, c: 32, d: 641},
		{env: map[string]string{"D": "0"}, cmd: []string{}, a: 8, b: 16, c: 32, d: 0},
		{env: map[string]string{"D": ""}, cmd: []string{}, a: 8, b: 16, c: 32, d: 0},
		{env: map[string]string{"D": ""}, cmd: []string{"-d", "641"}, a: 8, b: 16, c: 32, d: 641},
		{env: map[string]string{"D": "640"}, cmd: []string{"-d", "0"}, a: 8, b: 16, c: 32, d: 0},
		{env: map[string]string{"D": "wrong"}, cmd: []string{"-d", "0"}, a: 8, b: 16, c: 32, d: 0},
	}

	for i, test := range tests {
		p := libcfg.NewParser()

		a := p.Int8("int8", "a", 8, "", "A")
		b := p.Int16("int16", "b", 16, "", "B")
		c := p.Int32("int32", "c", 32, "", "C")
		d := p.Int64("int64", "d", 64, "", "D")

		i := i       //pin
		test := test //pin

		withEnv(test.env, func() {
			if err := p.RunArgs(test.cmd); err != nil {
				t.Errorf("Case %d, error running parser: %v", i, err)
			}

			if *a != test.a {
				t.Errorf("Case %d, wrong int8 value: expected '%v', received '%v'", i, test.a, *a)
			}

			if *b != test.b {
				t.Errorf("Case %d, wrong int16 value: expected '%v', received '%v'", i, test.b, *b)
			}

			if *c != test.c {
				t.Errorf("Case %d, wrong int32 value: expected '%v', received '%v'", i, test.c, *c)
			}

			if *d != test.d {
				t.Errorf("Case %d, wrong int64 value: expected '%v', received '%v'", i, test.d, *d)
			}
		})
	}
}

func TestEnvOptUint(t *testing.T) {
	tests := []struct {
		env map[string]string
		cmd []string
		a   uint8
		b   uint16
		c   uint32
		d   uint64
	}{
		{env: map[string]string{}, cmd: []string{}, a: 8, b: 16, c: 32, d: 64},
		{env: map[string]string{"A": "80"}, cmd: []string{}, a: 80, b: 16, c: 32, d: 64},
		{env: map[string]string{"A": "80"}, cmd: []string{"-a", "81"}, a: 81, b: 16, c: 32, d: 64},
		{env: map[string]string{"A": "wrong"}, cmd: []string{}, a: 8, b: 16, c: 32, d: 64},
		{env: map[string]string{"A": "wrong"}, cmd: []string{"-a", "81"}, a: 81, b: 16, c: 32, d: 64},
		{env: map[string]string{"A": "0"}, cmd: []string{}, a: 0, b: 16, c: 32, d: 64},
		{env: map[string]string{"A": ""}, cmd: []string{}, a: 0, b: 16, c: 32, d: 64},
		{env: map[string]string{"A": ""}, cmd: []string{"-a", "80"}, a: 80, b: 16, c: 32, d: 64},
		{env: map[string]string{"A": "80"}, cmd: []string{"-a", "0"}, a: 0, b: 16, c: 32, d: 64},
		{env: map[string]string{"A": "wrong"}, cmd: []string{"-a", "0"}, a: 0, b: 16, c: 32, d: 64},
		{env: map[string]string{"B": "160"}, cmd: []string{}, a: 8, b: 160, c: 32, d: 64},
		{env: map[string]string{"B": "160"}, cmd: []string{"-b", "161"}, a: 8, b: 161, c: 32, d: 64},
		{env: map[string]string{"B": "wrong"}, cmd: []string{}, a: 8, b: 16, c: 32, d: 64},
		{env: map[string]string{"B": "wrong"}, cmd: []string{"-b", "161"}, a: 8, b: 161, c: 32, d: 64},
		{env: map[string]string{"B": "0"}, cmd: []string{}, a: 8, b: 0, c: 32, d: 64},
		{env: map[string]string{"B": ""}, cmd: []string{}, a: 8, b: 0, c: 32, d: 64},
		{env: map[string]string{"B": ""}, cmd: []string{"-b", "161"}, a: 8, b: 161, c: 32, d: 64},
		{env: map[string]string{"B": "160"}, cmd: []string{"-b", "0"}, a: 8, b: 0, c: 32, d: 64},
		{env: map[string]string{"B": "wrong"}, cmd: []string{"-b", "0"}, a: 8, b: 0, c: 32, d: 64},
		{env: map[string]string{"C": "320"}, cmd: []string{}, a: 8, b: 16, c: 320, d: 64},
		{env: map[string]string{"C": "320"}, cmd: []string{"-c", "321"}, a: 8, b: 16, c: 321, d: 64},
		{env: map[string]string{"C": "wrong"}, cmd: []string{}, a: 8, b: 16, c: 32, d: 64},
		{env: map[string]string{"C": "wrong"}, cmd: []string{"-c", "321"}, a: 8, b: 16, c: 321, d: 64},
		{env: map[string]string{"C": "0"}, cmd: []string{}, a: 8, b: 16, c: 0, d: 64},
		{env: map[string]string{"C": ""}, cmd: []string{}, a: 8, b: 16, c: 0, d: 64},
		{env: map[string]string{"C": ""}, cmd: []string{"-c", "321"}, a: 8, b: 16, c: 321, d: 64},
		{env: map[string]string{"C": "320"}, cmd: []string{"-c", "0"}, a: 8, b: 16, c: 0, d: 64},
		{env: map[string]string{"C": "wrong"}, cmd: []string{"-c", "0"}, a: 8, b: 16, c: 0, d: 64},
		{env: map[string]string{"D": "640"}, cmd: []string{}, a: 8, b: 16, c: 32, d: 640},
		{env: map[string]string{"D": "640"}, cmd: []string{"-d", "641"}, a: 8, b: 16, c: 32, d: 641},
		{env: map[string]string{"D": "wrong"}, cmd: []string{}, a: 8, b: 16, c: 32, d: 64},
		{env: map[string]string{"D": "wrong"}, cmd: []string{"-d", "641"}, a: 8, b: 16, c: 32, d: 641},
		{env: map[string]string{"D": "0"}, cmd: []string{}, a: 8, b: 16, c: 32, d: 0},
		{env: map[string]string{"D": ""}, cmd: []string{}, a: 8, b: 16, c: 32, d: 0},
		{env: map[string]string{"D": ""}, cmd: []string{"-d", "641"}, a: 8, b: 16, c: 32, d: 641},
		{env: map[string]string{"D": "640"}, cmd: []string{"-d", "0"}, a: 8, b: 16, c: 32, d: 0},
		{env: map[string]string{"D": "wrong"}, cmd: []string{"-d", "0"}, a: 8, b: 16, c: 32, d: 0},
	}

	for i, test := range tests {
		p := libcfg.NewParser()

		a := p.Uint8("uint8", "a", 8, "", "A")
		b := p.Uint16("uint16", "b", 16, "", "B")
		c := p.Uint32("uint32", "c", 32, "", "C")
		d := p.Uint64("uint64", "d", 64, "", "D")

		i := i       //pin
		test := test //pin

		withEnv(test.env, func() {
			if err := p.RunArgs(test.cmd); err != nil {
				t.Errorf("Case %d, error running parser: %v", i, err)
			}

			if *a != test.a {
				t.Errorf("Case %d, wrong uint8 value: expected '%v', received '%v'", i, test.a, *a)
			}

			if *b != test.b {
				t.Errorf("Case %d, wrong uint16 value: expected '%v', received '%v'", i, test.b, *b)
			}

			if *c != test.c {
				t.Errorf("Case %d, wrong uint32 value: expected '%v', received '%v'", i, test.c, *c)
			}

			if *d != test.d {
				t.Errorf("Case %d, wrong uint64 value: expected '%v', received '%v'", i, test.d, *d)
			}
		})
	}
}
