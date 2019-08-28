package libcfg_test

import (
	"strings"
	"testing"

	"github.com/ibraimgm/libcfg"
)

func TestParseCmdArgs(t *testing.T) {
	tests := []struct {
		cmd     []string
		abool   bool
		aint    int
		astring string
		args    []string
	}{
		{cmd: []string{}},
		{cmd: []string{"-b", "-i", "5", "-s", "foo"}, abool: true, aint: 5, astring: "foo"},
		{cmd: []string{"--abool", "--aint", "5", "--astring", "foo"}, abool: true, aint: 5, astring: "foo"},
		{cmd: []string{"--aint=5", "--astring=foo"}, aint: 5, astring: "foo"},
		{cmd: []string{"-b", "--abool=false"}},
		{cmd: []string{"-b", "--no-abool"}},
		{cmd: []string{"-i", "5", "--aint", "6", "-i", "7"}, aint: 7},
		{cmd: []string{"-b", "-i", "5", "foo", "bar"}, abool: true, aint: 5, args: []string{"foo", "bar"}},
		{cmd: []string{"foo", "bar"}, args: []string{"foo", "bar"}},
		{cmd: []string{"foo", "-i", "5"}, args: []string{"foo", "-i", "5"}},
		{cmd: []string{"-b", "--no-abool=true"}},
		{cmd: []string{"--no-abool=false"}, abool: true},
		{cmd: []string{"-s", "foo", "--astring="}},
		{cmd: []string{"--astring", "--aint", "5"}, astring: "--aint", args: []string{"5"}},
	}

	for i, test := range tests {
		cfg := libcfg.NewParser()

		abool := cfg.OptBool("abool", "b", false, "specifies a bool value")
		aint := cfg.OptInt("aint", "i", 0, "specifies an int value")
		astring := cfg.OptString("astring", "s", "", "specifies a string value")

		if err := cfg.ParseArgs(test.cmd); err != nil {
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

func TestParseCmdError(t *testing.T) {
	tests := []struct {
		cmd           []string
		expectedError string
	}{
		{cmd: []string{"-b", "-x"}, expectedError: "unknown argument: -x"},
		{cmd: []string{"-b", "--x"}, expectedError: "unknown argument: --x"},
		{cmd: []string{"--abool=X"}, expectedError: "is not a valid boolean value"},
		{cmd: []string{"-i", "a"}, expectedError: "is not a valid int value"},
		{cmd: []string{"-i"}, expectedError: "no value for argument: -i"},
		{cmd: []string{"-s"}, expectedError: "no value for argument: -s"},
		{cmd: []string{"--aint"}, expectedError: "no value for argument: --aint"},
		{cmd: []string{"--astring"}, expectedError: "no value for argument: --astring"},
		{cmd: []string{"--aint="}, expectedError: "no value for argument: --aint"},
		{cmd: []string{"--aint=", "5"}, expectedError: "no value for argument: --aint"},
	}

	for i, test := range tests {
		cfg := libcfg.NewParser()

		cfg.OptBool("abool", "b", false, "specifies a bool value")
		cfg.OptInt("aint", "i", 0, "specifies an int value")
		cfg.OptString("astring", "s", "", "specifies a string value")

		err := cfg.ParseArgs(test.cmd)

		if err == nil {
			t.Errorf("Case %d, argument parsing should return error", i)
			continue
		}

		if !strings.Contains(err.Error(), test.expectedError) {
			t.Errorf("Case %d, expected error '%s', but got '%s'", i, test.expectedError, err.Error())
		}
	}
}
