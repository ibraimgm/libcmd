package libcmd_test

import (
	"io/ioutil"
	"testing"

	"github.com/ibraimgm/libcmd"
)

func TestSubCommand(t *testing.T) {
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
		app := libcmd.NewApp("app", "")
		app.Configure(libcmd.Options{HelpOutput: ioutil.Discard})
		var c1, c11, c12, c2 bool

		i := i       //pin
		test := test //pin

		app.Command("c1", "", func(cmd *libcmd.Cmd) {
			cmd.Match(func(*libcmd.Cmd) {
				c1 = true
			})

			cmd.Command("c11", "", func(cmd *libcmd.Cmd) {
				cmd.Match(func(*libcmd.Cmd) {
					c11 = true
					compareArgs(t, i, test.unparsed, cmd.Args())
				})
			})

			cmd.CommandMatch("c12", "", func(*libcmd.Cmd) {
				c12 = true
			})
		})

		app.CommandMatch("c2", "", func(*libcmd.Cmd) {
			c2 = true
		})

		if err := app.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error running parser: %v", i, err)
			continue
		}

		compareValue(t, i, test.c1, c1)
		compareValue(t, i, test.c11, c11)
		compareValue(t, i, test.c12, c12)
		compareValue(t, i, test.c2, c2)

		compareArgs(t, i, test.unparsed, app.Args())
	}
}

func TestCommandArgReuse(t *testing.T) {
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
		app := libcmd.NewApp("", "")
		app.Configure(libcmd.Options{HelpOutput: ioutil.Discard})

		var c1, c2 bool

		str := app.String("str", "s", "default", "")
		b := app.Bool("bool", "b", false, "")

		app.Command("c1", "", func(cmd *libcmd.Cmd) {
			cmd.StringP(str, "str", "s", *str, "")
			cmd.BoolP(b, "bool", "b", *b, "")

			cmd.Match(func(*libcmd.Cmd) {
				c1 = true
			})
		})

		app.Command("c2", "", func(cmd *libcmd.Cmd) {
			cmd.Match(func(*libcmd.Cmd) {
				c2 = true
			})

			cmd.Err(func(err error) error {
				c2 = true
				return err
			})
		})

		err := app.RunArgs(test.cmd)
		if test.expectError && err == nil {
			t.Errorf("Case %d, should have returned error", i)
			continue
		} else if !test.expectError && err != nil {
			t.Errorf("Case %d, error running parser: %v", i, err)
			continue
		}

		compareValue(t, i, test.c1, c1)
		compareValue(t, i, test.c2, c2)
		compareValue(t, i, test.str, *str)
		compareValue(t, i, test.b, *b)
	}
}

func TestCommandArgSameName(t *testing.T) {
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
		app := libcmd.NewApp("", "")
		app.Configure(libcmd.Options{HelpOutput: ioutil.Discard})
		s1 := app.String("s1", "", "s1", "")
		s2 := app.String("s2", "", "s2", "")

		var c1, c2 bool

		// not really used - the commands will replace the pointer when run
		// the 'new' here only exists because the comparison at the end of the routine
		// always assume a valid pointer
		c1s1 := new(string)
		c2s2 := new(string)

		app.Command("c1", "", func(cmd *libcmd.Cmd) {
			c1s1 = cmd.String("s1", "", "", "")
			cmd.StringP(s2, "s2", "", *s2, "")

			cmd.Match(func(*libcmd.Cmd) {
				c1 = true
			})
		})

		app.Command("c2", "", func(cmd *libcmd.Cmd) {
			cmd.StringP(s1, "s1", "", *s1, "")
			c2s2 = cmd.String("s2", "", "", "")

			cmd.Match(func(*libcmd.Cmd) {
				c2 = true
			})
		})

		if err := app.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error running parser: %v", i, err)
			continue
		}

		compareValue(t, i, test.s1, *s1)
		compareValue(t, i, test.s2, *s2)
		compareValue(t, i, test.c1, c1)
		compareValue(t, i, test.c2, c2)
		compareValue(t, i, test.c1s1, *c1s1)
		compareValue(t, i, test.c2s2, *c2s2)
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		cmd []string
		c1  bool
		c11 bool
		c2  bool
		c21 bool
	}{
		{cmd: []string{}},
		{cmd: []string{"c1"}, c1: true},
		{cmd: []string{"c1", "c11"}, c11: true},
		{cmd: []string{"c2"}, c2: true},
		{cmd: []string{"c2", "c21"}, c21: true},
		{cmd: []string{"c11", "c1"}},
		{cmd: []string{"c21", "c2"}},
		{cmd: []string{"c2", "c1"}, c2: true},
		{cmd: []string{"c1", "c2"}, c1: true},
	}

	for i, test := range tests {
		var c1, c11, c2, c21 bool
		app := libcmd.NewApp("", "")

		app.Command("c1", "", func(cmd *libcmd.Cmd) {
			cmd.Run(func(*libcmd.Cmd) error {
				c1 = true
				return nil
			})

			cmd.CommandRun("c11", "", func(*libcmd.Cmd) error {
				c11 = true
				return nil
			})
		})

		app.Command("c2", "", func(cmd *libcmd.Cmd) {
			cmd.Run(func(*libcmd.Cmd) error {
				c2 = true
				return nil
			})

			cmd.CommandRun("c21", "", func(*libcmd.Cmd) error {
				c21 = true
				return nil
			})
		})

		if err := app.RunArgs(test.cmd); err != nil {
			t.Errorf("Case %d, error running parser: %v", i, err)
			continue
		}

		compareValue(t, i, test.c1, c1)
		compareValue(t, i, test.c11, c11)
		compareValue(t, i, test.c2, c2)
		compareValue(t, i, test.c21, c21)
	}
}

func TestNoMatch(t *testing.T) {
	var s string
	var c1 bool

	app := libcmd.NewApp("", "")
	app.StringP(&s, "", "", "xxxx", "")

	app.Command("", "", func(cmd *libcmd.Cmd) {
		c1 = true
	})

	if err := app.RunArgs([]string{"", ""}); err != nil {
		t.Errorf("Error running parser: %v", err)
		return
	}

	if s != "" {
		t.Errorf("Variable should not have a value (got: '%s')", s)
	}

	if c1 {
		t.Errorf("Command should not be called")
	}
}
