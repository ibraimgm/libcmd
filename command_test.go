package libcfg_test

import (
	"testing"

	"github.com/ibraimgm/libcfg"
)

func TestCommandSelecting(t *testing.T) {
	tests := []struct {
		cmd      []string
		c1       bool
		c11      bool
		c12      bool
		c2       bool
		unparsed []string
	}{
		{cmd: []string{}, unparsed: []string{}},
		{cmd: []string{"c1"}, unparsed: []string{}, c1: true},
		{cmd: []string{"c1", "a"}, unparsed: []string{"a"}, c1: true},
		{cmd: []string{"c1", "c11"}, unparsed: []string{}, c1: true, c11: true},
		{cmd: []string{"c1", "c12"}, unparsed: []string{}, c1: true, c12: true},
		{cmd: []string{"c1", "c12", "x", "y", "z"}, unparsed: []string{"x", "y", "z"}, c1: true, c12: true},
		{cmd: []string{"c2"}, unparsed: []string{}, c2: true},
		{cmd: []string{"a", "b", "c"}, unparsed: []string{"a", "b", "c"}},
		{cmd: []string{"c1", "a", "c12"}, unparsed: []string{"a", "c12"}, c1: true},
		{cmd: []string{"a", "c1", "b"}, unparsed: []string{"a", "c1", "b"}},
		{cmd: []string{"c2", "c1"}, unparsed: []string{"c1"}, c2: true},
		{cmd: []string{"c1", "c11", "c12"}, unparsed: []string{"c12"}, c1: true, c11: true},
	}

	for i, test := range tests {
		cfg := libcfg.NewParser()

		c1 := cfg.AddCommand("c1", "command 1")

		c11 := c1.AddCommand("c11", "command 1.1")
		c12 := c1.AddCommand("c12", "command 1.2")
		c2 := cfg.AddCommand("c2", "command 2")

		if err := cfg.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error running parser: %v", i, err)
			continue
		}

		if c1.Used() != test.c1 {
			t.Errorf("Case %d, 'used' flag for c1 expected to be '%v' but was '%v'", i, test.c1, c1.Used())
		}

		if c11.Used() != test.c11 {
			t.Errorf("Case %d, 'used' flag for c11 expected to be '%v' but was '%v'", i, test.c11, c11.Used())
		}

		if c12.Used() != test.c12 {
			t.Errorf("Case %d, 'used' flag for c12 expected to be '%v' but was '%v'", i, test.c12, c12.Used())
		}

		if c2.Used() != test.c2 {
			t.Errorf("Case %d, 'used' flag for c2 expected to be '%v' but was '%v'", i, test.c2, c2.Used())
		}

		args := cfg.Args()

		if len(args) != len(test.unparsed) {
			t.Errorf("Case %d, wrong size of unparsed arguments: expected '%v', received '%v'", i, len(test.unparsed), len(args))
			continue
		}

		for j, arg := range args {
			expected := test.unparsed[j]

			if arg != expected {
				t.Errorf("Case %d, wrong args result at pos %d: expected '%v', received '%v'", i, j, expected, arg)
			}
		}
	}
}
