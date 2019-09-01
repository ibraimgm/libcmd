package libcfg_test

import (
	"os"
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

func TestParseEnvArgs(t *testing.T) {
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
		cfg := libcfg.NewParser()

		abool := cfg.EnvBool(false, "B1", "B2", "B3")
		aint := cfg.EnvInt(0, "I1", "I2", "I3")
		auint := cfg.EnvUint(0, "U1", "U2", "U3")
		astring := cfg.EnvString("", "S1", "S2", "S3")

		i := i       // pin scope
		test := test // pin scope

		withEnv(test.env, func() {
			if err := cfg.RunEnv(); err != nil {
				t.Errorf("Case %d, error loading from environment: %v", i, err)
				return
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
		})
	}
}
