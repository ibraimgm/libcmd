package libcfg_test

import (
	"testing"

	"github.com/ibraimgm/libcfg"
)

func TestBindOpt(t *testing.T) {
	// inner type for testing
	type bindOptStruct struct {
		Abool    bool    `long:"abool" short:"b" default:"true"`
		Aint     int     `long:"aint" short:"i" default:"8"`
		Auint    uint    `long:"auint" short:"u" default:"16"`
		Astring  string  `long:"astring" short:"s" default:"default"`
		Afloat32 float32 `long:"afloat32" short:"f32" default:"3.14"`
		Afloat64 float64 `long:"afloat64" short:"f64" default:"3.1415"`
	}

	// actual test values
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
		var s bindOptStruct

		if err := p.Bind(&s); err != nil {
			t.Errorf("Case %d, error binding on struct: %v", i, err)
			continue
		}

		if err := p.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error parsing args: %v", i, err)
			continue
		}

		if s.Abool != test.abool {
			t.Errorf("Case %d, wrong boolean value: expected '%v', received '%v'", i, test.abool, s.Abool)
		}

		if s.Aint != test.aint {
			t.Errorf("Case %d, wrong int value: expected '%v', received '%v'", i, test.aint, s.Aint)
		}

		if s.Auint != test.auint {
			t.Errorf("Case %d, wrong uint value: expected '%v', received '%v'", i, test.auint, s.Auint)
		}

		if s.Astring != test.astring {
			t.Errorf("Case %d, wrong string value: expected '%v', received '%v'", i, test.astring, s.Astring)
		}

		if s.Afloat32 != test.afloat32 {
			t.Errorf("Case %d, wrong float32 value: expected '%v', received '%v'", i, test.afloat32, s.Afloat32)
		}

		if s.Afloat64 != test.afloat64 {
			t.Errorf("Case %d, wrong float64 value: expected '%v', received '%v'", i, test.afloat64, s.Afloat64)
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

func TestBindEnv(t *testing.T) {
	type bindEnvStruct struct {
		Abool   bool   `env:"B1,B2,B3" default:"true"`
		Aint    int    `env:"I1,I2,I3" default:"8"`
		Auint   uint   `env:"U1,U2,U3" default:"16"`
		Astring string `env:"S1,S2,S3" default:"xyz"`
	}

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

		var s bindEnvStruct
		if err := env.Bind(&s); err != nil {
			t.Errorf("Case %d, error binding on struct: %v", i, err)
			continue
		}

		i := i       // pin scope
		test := test // pin scope

		withEnv(test.env, func() {
			env.LoadAll()

			if s.Abool != test.abool {
				t.Errorf("Case %d, wrong boolean value: expected '%v', received '%v'", i, test.abool, s.Abool)
			}

			if s.Aint != test.aint {
				t.Errorf("Case %d, wrong int value: expected '%v', received '%v'", i, test.aint, s.Aint)
			}

			if s.Auint != test.auint {
				t.Errorf("Case %d, wrong uint value: expected '%v', received '%v'", i, test.auint, s.Auint)
			}

			if s.Astring != test.astring {
				t.Errorf("Case %d, wrong string value: expected '%v', received '%v'", i, test.astring, s.Astring)
			}
		})
	}
}

func TestBindError(t *testing.T) {
	type bindErrorStruct struct {
		x int //nolint
	}

	tests := []interface{}{
		nil,
		bindErrorStruct{},
	}

	for i, test := range tests {
		p := libcfg.NewParser()
		if err := p.Bind(test); err == nil {
			t.Errorf("Case %d, parser should return error", i)
		}

		env := libcfg.NewEnvLoader()
		if err := env.Bind(test); err == nil {
			t.Errorf("Case %d, env loader should return error", i)
		}
	}
}

func TestBindSpecialCases(t *testing.T) {
	//nolint
	type bindSpecialStruct struct {
		x int `env:""`
		y int `default:"a"`
	}

	var s bindSpecialStruct

	if err := libcfg.NewParser().Bind(&s); err != nil {
		t.Errorf("Parser should not err on this (malformed) struct. Received: %v", err)
	}

	if err := libcfg.NewEnvLoader().Bind(&s); err != nil {
		t.Errorf("EnvLoader should not err on this (malformed) struct. Received: %v", err)
	}
}
