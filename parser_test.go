package libcfg_test

import (
	"testing"

	"github.com/ibraimgm/libcfg"
)

func TestParseCmdArgs(t *testing.T) {
	tests := []struct {
		args    []string
		abool   bool
		aint    int
		astring string
	}{
		{args: []string{}},
		{args: []string{}},
		{args: []string{"-b", "-i", "5", "-s", "foo"}, abool: true, aint: 5, astring: "foo"},
		{args: []string{"--abool", "--aint", "5", "--astring", "foo"}, abool: true, aint: 5, astring: "foo"},
		{args: []string{"--aint=5", "--astring=foo"}, aint: 5, astring: "foo"},
		{args: []string{"-b", "--abool=false"}},
		{args: []string{"-b", "--no-abool"}},
		{args: []string{"-i", "5", "--aint", "6", "-i", "7"}, aint: 7},
	}

	for i, test := range tests {
		cfg := libcfg.NewParser()

		abool := cfg.OptBool("abool", "b", false, "specifies a bool value")
		aint := cfg.OptInt("aint", "i", 0, "specifies an int value")
		astring := cfg.OptString("astring", "s", "", "specifies a string value")

		if err := cfg.ParseArgs(test.args); err != nil {
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
	}
}
