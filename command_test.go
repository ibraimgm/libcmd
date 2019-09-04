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

func TestCommandOneInheritAll(t *testing.T) {
	tests := []struct {
		cmd         []string
		c1          bool
		c2          bool
		str         string
		b           bool
		expectError bool
	}{
		{cmd: []string{}, str: "default"},
		{cmd: []string{"-s", "foo", "-b"}, str: "foo", b: true},
		{cmd: []string{"c1", "-s", "foo", "-b"}, c1: true, str: "foo", b: true},
		{cmd: []string{"-b", "c1", "-s", "foo"}, c1: true, str: "foo", b: true},
		{cmd: []string{"-s", "foo", "c1", "-b"}, c1: true, str: "foo", b: true},
		{cmd: []string{"-b", "c1", "-s", "foo", "--no-bool"}, c1: true, str: "foo"},
		{cmd: []string{"-s", "foo", "--no-bool", "c1", "-s", "bar", "-b"}, c1: true, str: "bar", b: true},
		{cmd: []string{"c1", "-s", "c2"}, c1: true, str: "c2"},
		{cmd: []string{"c1", "-s", "foo", "c2"}, c1: true, str: "foo"},
		{cmd: []string{"c2", "-s", "foo", "-b"}, c2: true, str: "default", expectError: true},
		{cmd: []string{"-b", "c2", "-s", "foo"}, c2: true, b: true, str: "default", expectError: true},
		{cmd: []string{"-s", "foo", "c2", "-b"}, c2: true, str: "foo", expectError: true},
		{cmd: []string{"-b", "c2", "-s", "foo", "--no-bool"}, c2: true, b: true, str: "default", expectError: true},
		{cmd: []string{"-s", "foo", "--no-bool", "c2", "-s", "bar", "-b"}, c2: true, str: "foo", expectError: true},
		{cmd: []string{"c2", "-s", "c1"}, c2: true, str: "default", expectError: true},
		{cmd: []string{"c2", "-s", "foo", "c1"}, c2: true, str: "default", expectError: true},
	}

	for i, test := range tests {
		cfg := libcfg.NewParser()

		str := cfg.OptString("str", "s", "default", "")
		b := cfg.OptBool("bool", "b", false, "")

		c1 := cfg.AddCommand("c1", "command 1")
		c1.InheritAll()

		c2 := cfg.AddCommand("c2", "command 2")

		err := cfg.RunArgs(test.cmd)
		if test.expectError && err == nil {
			t.Errorf("Case %d, should have returned error", i)
			continue
		} else if !test.expectError && err != nil {
			t.Errorf("Case %d, error running parser: %v", i, err)
			continue
		}

		if c1.Used() != test.c1 {
			t.Errorf("Case %d, 'used' flag for c1 expected to be '%v' but was '%v'", i, test.c1, c1.Used())
		}

		if c2.Used() != test.c2 {
			t.Errorf("Case %d, 'used' flag for c2 expected to be '%v' but was '%v'", i, test.c2, c2.Used())
		}

		if *str != test.str {
			t.Errorf("Case %d, wrong value on string flag: expected '%s', received '%s'", i, test.str, *str)
		}

		if *b != test.b {
			t.Errorf("Case %d, wrong value on bool flag: expected '%v', received '%v'", i, test.b, *b)
		}
	}
}

func TestCommandPartialInherit(t *testing.T) {
	tests := []struct {
		cmd         []string
		c1          bool
		c2          bool
		str         string
		b           bool
		expectError bool
	}{
		{cmd: []string{}, str: "default"},
		{cmd: []string{"-s", "foo", "-b"}, str: "foo", b: true},
		{cmd: []string{"c1", "-s", "foo", "-b"}, c1: true, str: "foo", expectError: true},
		{cmd: []string{"-b", "c1", "-s", "foo"}, c1: true, str: "foo", b: true},
		{cmd: []string{"-s", "foo", "c1", "-b"}, c1: true, str: "foo", expectError: true},
		{cmd: []string{"-b", "c1", "-s", "foo", "--no-bool"}, c1: true, str: "foo", b: true, expectError: true},
		{cmd: []string{"-s", "foo", "--no-bool", "c1", "-s", "bar", "-b"}, c1: true, str: "bar", expectError: true},
		{cmd: []string{"c1", "-s", "c2"}, c1: true, str: "c2"},
		{cmd: []string{"c1", "-s", "foo", "c2"}, c1: true, str: "foo"},
		{cmd: []string{"c2", "-s", "foo", "-b"}, c2: true, str: "default", expectError: true},
		{cmd: []string{"-b", "c2", "-s", "foo"}, c2: true, b: true, str: "default", expectError: true},
		{cmd: []string{"-s", "foo", "c2", "-b"}, c2: true, str: "foo", b: true},
		{cmd: []string{"-b", "c2", "-s", "foo", "--no-bool"}, c2: true, b: true, str: "default", expectError: true},
		{cmd: []string{"-s", "foo", "--no-bool", "c2", "-s", "bar", "-b"}, c2: true, str: "foo", expectError: true},
		{cmd: []string{"c2", "-s", "c1"}, c2: true, str: "default", expectError: true},
		{cmd: []string{"c2", "-s", "foo", "c1"}, c2: true, str: "default", expectError: true},
		{cmd: []string{"-b", "c2", "--no-bool"}, c2: true, str: "default"},
		{cmd: []string{"-s", "foo", "-b", "c2"}, c2: true, str: "foo", b: true},
		{cmd: []string{"-s", "foo", "-b", "c2", "--no-bool"}, c2: true, str: "foo"},
	}

	for i, test := range tests {
		cfg := libcfg.NewParser()

		str := cfg.OptString("str", "s", "default", "")
		b := cfg.OptBool("bool", "b", false, "")

		c1 := cfg.AddCommand("c1", "command 1")
		c1.Inherit("s")

		c2 := cfg.AddCommand("c2", "command 2")
		c2.Inherit("b")

		err := cfg.RunArgs(test.cmd)
		if test.expectError && err == nil {
			t.Errorf("Case %d, should have returned error", i)
			continue
		} else if !test.expectError && err != nil {
			t.Errorf("Case %d, error running parser: %v", i, err)
			continue
		}

		if c1.Used() != test.c1 {
			t.Errorf("Case %d, 'used' flag for c1 expected to be '%v' but was '%v'", i, test.c1, c1.Used())
		}

		if c2.Used() != test.c2 {
			t.Errorf("Case %d, 'used' flag for c2 expected to be '%v' but was '%v'", i, test.c2, c2.Used())
		}

		if *str != test.str {
			t.Errorf("Case %d, wrong value on string flag: expected '%s', received '%s'", i, test.str, *str)
		}

		if *b != test.b {
			t.Errorf("Case %d, wrong value on bool flag: expected '%v', received '%v'", i, test.b, *b)
		}
	}
}

func TestCommandBothInheritAll(t *testing.T) {
	tests := []struct {
		cmd []string
		c1  bool
		c2  bool
		str string
		b   bool
	}{
		{cmd: []string{}, str: "default"},
		{cmd: []string{"-s", "foo", "-b"}, str: "foo", b: true},
		{cmd: []string{"c1", "-s", "foo", "-b"}, c1: true, str: "foo", b: true},
		{cmd: []string{"-b", "c1", "-s", "foo"}, c1: true, str: "foo", b: true},
		{cmd: []string{"-s", "foo", "c1", "-b"}, c1: true, str: "foo", b: true},
		{cmd: []string{"-b", "c1", "-s", "foo", "--no-bool"}, c1: true, str: "foo"},
		{cmd: []string{"-s", "foo", "--no-bool", "c1", "-s", "bar", "-b"}, c1: true, str: "bar", b: true},
		{cmd: []string{"c1", "-s", "c2"}, c1: true, str: "c2"},
		{cmd: []string{"c1", "-s", "foo", "c2"}, c1: true, str: "foo"},
		{cmd: []string{"c2", "-s", "foo", "-b"}, c2: true, str: "foo", b: true},
		{cmd: []string{"-b", "c2", "-s", "foo"}, c2: true, str: "foo", b: true},
		{cmd: []string{"-s", "foo", "c2", "-b"}, c2: true, str: "foo", b: true},
		{cmd: []string{"-b", "c2", "-s", "foo", "--no-bool"}, c2: true, str: "foo"},
		{cmd: []string{"-s", "foo", "--no-bool", "c2", "-s", "bar", "-b"}, c2: true, str: "bar", b: true},
		{cmd: []string{"c2", "-s", "c1"}, c2: true, str: "c1"},
		{cmd: []string{"c2", "-s", "foo", "c1"}, c2: true, str: "foo"},
	}

	for i, test := range tests {
		cfg := libcfg.NewParser()

		str := cfg.OptString("", "s", "default", "")
		b := cfg.OptBool("bool", "b", false, "")

		c1 := cfg.AddCommand("c1", "command 1")
		c1.InheritAll()

		c2 := cfg.AddCommand("c2", "command 2")
		c2.InheritAll()

		if err := cfg.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error running parser: %v", i, err)
			continue
		}

		if c1.Used() != test.c1 {
			t.Errorf("Case %d, 'used' flag for c1 expected to be '%v' but was '%v'", i, test.c1, c1.Used())
		}

		if c2.Used() != test.c2 {
			t.Errorf("Case %d, 'used' flag for c2 expected to be '%v' but was '%v'", i, test.c2, c2.Used())
		}

		if *str != test.str {
			t.Errorf("Case %d, wrong value on string flag: expected '%s', received '%s'", i, test.str, *str)
		}

		if *b != test.b {
			t.Errorf("Case %d, wrong value on bool flag: expected '%v', received '%v'", i, test.b, *b)
		}
	}
}

func TestCommandInheritByName(t *testing.T) {
	tests := []struct {
		cmd         []string
		toInherit   []string
		c1          bool
		aint        int
		auint       uint
		expectError bool
	}{
		{cmd: []string{}, toInherit: []string{}, aint: 8, auint: 16},
		{cmd: []string{"c1", "--aint", "10"}, toInherit: []string{"aint"}, c1: true, aint: 10, auint: 16},
		{cmd: []string{"c1", "-u", "10"}, toInherit: []string{"u"}, c1: true, aint: 8, auint: 10},
		{cmd: []string{"c1", "--aint", "10", "-u", "10"}, toInherit: []string{"aint"}, c1: true, aint: 10, auint: 16, expectError: true},
		{cmd: []string{"c1", "-u", "10", "--aint", "10"}, toInherit: []string{"u"}, c1: true, aint: 8, auint: 10, expectError: true},
		{cmd: []string{"c1", "--aint", "10", "-u", "10"}, toInherit: []string{"aint", "u"}, c1: true, aint: 10, auint: 10},
		{cmd: []string{"c1", "--aint", "10", "-u", "10"}, toInherit: []string{"x", "aint", "u"}, c1: true, aint: 10, auint: 10},
		{cmd: []string{"c1", "--aint", "10", "-u", "10"}, toInherit: []string{"aint", "x", "u"}, c1: true, aint: 10, auint: 10},
		{cmd: []string{"c1", "--aint", "10", "-u", "10"}, toInherit: []string{"aint", "u", "x"}, c1: true, aint: 10, auint: 10},
		{cmd: []string{"c1", "--aint", "10"}, toInherit: []string{"aint", "aint"}, c1: true, aint: 10, auint: 16},
		{cmd: []string{"c1", "-u", "10"}, toInherit: []string{"u", "u"}, c1: true, aint: 8, auint: 10},
		{cmd: []string{"c1", "--aint", "10", "-u", "10"}, toInherit: []string{"aint", "u", "aint", "u"}, c1: true, aint: 10, auint: 10},
	}

	for i, test := range tests {
		cfg := libcfg.NewParser()
		aint := cfg.OptInt("aint", "", 8, "")
		auint := cfg.OptUint("", "u", 16, "")

		c1 := cfg.AddCommand("c1", "command 1")
		for _, opt := range test.toInherit {
			c1.Inherit(opt)
		}

		err := cfg.RunArgs(test.cmd)
		if test.expectError && err == nil {
			t.Errorf("Case %d, should have returned error", i)
			continue
		} else if !test.expectError && err != nil {
			t.Errorf("Case %d, error running parser: %v", i, err)
			continue
		}

		if c1.Used() != test.c1 {
			t.Errorf("Case %d, 'used' flag for c1 expected to be '%v' but was '%v'", i, test.c1, c1.Used())
		}

		if *aint != test.aint {
			t.Errorf("Case %d, wrong int value: expected '%v', received '%v'", i, test.aint, *aint)
		}

		if *auint != test.auint {
			t.Errorf("Case %d, wrong uint value: expected '%v', received '%v'", i, test.auint, *auint)
		}
	}
}

func TestCommandInheritOverride(t *testing.T) {
	tests := []struct {
		cmd  []string
		c1   bool
		c2   bool
		s1   string
		s2   string
		c1s1 string
		c2s2 string
	}{
		{cmd: []string{}, s1: "s1", s2: "s2"},
		{cmd: []string{"c1"}, s1: "s1", s2: "s2", c1: true},
		{cmd: []string{"c1", "--s1=a"}, s1: "s1", s2: "s2", c1: true, c1s1: "a"},
		{cmd: []string{"--s1=a", "c1"}, s1: "a", s2: "s2", c1: true},
		{cmd: []string{"--s2=b", "c1"}, s1: "s1", s2: "b", c1: true},
		{cmd: []string{"--s1=a", "--s2=b", "c1", "--s1=c"}, s1: "a", s2: "b", c1: true, c1s1: "c"},
		{cmd: []string{"c2"}, s1: "s1", s2: "s2", c2: true},
		{cmd: []string{"c2", "--s1=a"}, s1: "a", s2: "s2", c2: true},
		{cmd: []string{"c2", "--s2=b"}, s1: "s1", s2: "s2", c2: true, c2s2: "b"},
		{cmd: []string{"--s1=a", "c2"}, s1: "a", s2: "s2", c2: true},
		{cmd: []string{"--s2=b", "c2"}, s1: "s1", s2: "b", c2: true},
		{cmd: []string{"--s1=a", "--s2=b", "c2", "--s2=c"}, s1: "a", s2: "b", c2: true, c2s2: "c"},
	}

	for i, test := range tests {
		cfg := libcfg.NewParser()
		s1 := cfg.OptString("s1", "", "s1", "")
		s2 := cfg.OptString("s2", "", "s2", "")

		c1 := cfg.AddCommand("c1", "command 1")
		c1.InheritAll()
		c1s1 := c1.OptString("s1", "", "", "")

		c2 := cfg.AddCommand("c2", "command 2")
		c2.InheritAll()
		c2s2 := c2.OptString("s2", "", "", "")

		if err := cfg.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error running parser: %v", i, err)
			continue
		}

		if *s1 != test.s1 {
			t.Errorf("Case %d, wrong value on s1 flag: expected '%s', received '%s'", i, test.s1, *s1)
		}

		if *s2 != test.s2 {
			t.Errorf("Case %d, wrong value on s2 flag: expected '%s', received '%s'", i, test.s2, *s2)
		}

		if c1.Used() != test.c1 {
			t.Errorf("Case %d, 'used' flag for c1 expected to be '%v' but was '%v'", i, test.c1, c1.Used())
		}

		if *c1s1 != test.c1s1 {
			t.Errorf("Case %d, wrong value on s1 flag of c1: expected '%s', received '%s'", i, test.c1s1, *c1s1)
		}

		if c2.Used() != test.c2 {
			t.Errorf("Case %d, 'used' flag for c2 expected to be '%v' but was '%v'", i, test.c2, c2.Used())
		}

		if *c2s2 != test.c2s2 {
			t.Errorf("Case %d, wrong value on s2 flag of c2: expected '%s', received '%s'", i, test.c2s2, *c2s2)
		}
	}
}

func TestCommandEnvInherit(t *testing.T) {
	tests := []struct {
		cmd  []string
		env  map[string]string
		file map[string]string
		s1   string
		s2   string
	}{
		{cmd: []string{"c1"}, env: map[string]string{}, file: map[string]string{}, s1: "s1", s2: "s2"},
		{cmd: []string{"c1"}, env: map[string]string{"A": "a"}, file: map[string]string{"B": "b"}, s1: "b", s2: "s2"},
		{cmd: []string{"c1"}, env: map[string]string{"B": "b"}, file: map[string]string{"A": "a"}, s1: "b", s2: "s2"},
		{cmd: []string{"c1"}, env: map[string]string{"A": "a", "C": "c"}, file: map[string]string{"B": "b"}, s1: "c", s2: "s2"},
		{cmd: []string{"c1"}, env: map[string]string{"X": "x"}, file: map[string]string{}, s1: "s1", s2: "x"},
		{cmd: []string{"c1"}, env: map[string]string{}, file: map[string]string{"X": "x"}, s1: "s1", s2: "x"},
		{cmd: []string{"c1"}, env: map[string]string{"Y": "y"}, file: map[string]string{"X": "x"}, s1: "s1", s2: "y"},
		{cmd: []string{"c1"}, env: map[string]string{"X": "x", "Y": "y"}, file: map[string]string{"Z": "z"}, s1: "s1", s2: "z"},
		{cmd: []string{"c1"}, env: map[string]string{"X": "x", "Z": "z"}, file: map[string]string{"y": "y"}, s1: "s1", s2: "z"},
	}

	for i, test := range tests {
		cfg := libcfg.NewParser()
		s1 := cfg.EnvOptString("s1", "", "s1", "", "A", "B", "C")
		cfg.OptString("s2", "", "", "") //don't care - will be overridden

		c1 := cfg.AddCommand("c1", "command 1")
		c1.InheritAll()
		s2 := c1.EnvOptString("s2", "", "s2", "", "X", "Y", "Z")

		i := i       //pin
		test := test //pin

		withEnv(test.env, func() {
			withFileEnv(test.file, func(file string) {
				if err := c1.UseFile(file); err != nil {
					t.Errorf("Case %d, error loading config file: %v", i, err)
				}

				if err := cfg.RunArgs(test.cmd); err != nil {
					t.Errorf("Case %d, error running parser: %v", i, err)
					return
				}

				if !c1.Used() {
					t.Errorf("Case %d, c1 should be used", i)
					return
				}

				if *s1 != test.s1 {
					t.Errorf("Case %d, wrong s1 value: expected '%s', received '%s'", i, test.s1, *s1)
				}

				if *s2 != test.s2 {
					t.Errorf("Case %d, wrong s2 value: expected '%s', received '%s'", i, test.s2, *s2)
				}
			})
		})
	}
}
