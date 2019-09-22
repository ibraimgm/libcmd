package libcmd_test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/ibraimgm/libcmd"
)

func compareHelpOutput(cmd *libcmd.Cmd, goldenfile string) error {
	bytes, err := ioutil.ReadFile(goldenfile)
	if err != nil {
		return err
	}
	expected := string(bytes)

	var b strings.Builder
	cmd.PrintHelp(&b)
	actual := b.String()

	if expected != actual {
		return fmt.Errorf("Wrong output. Expected:\n>>>\n%s\n<<<\nActual:\n>>>\n%s\n<<<", expected, actual)
	}

	return nil
}

func TestBasic(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")

	if err := compareHelpOutput(app.Cmd, "testdata/basic.golden"); err != nil {
		t.Error(err)
	}
}

func TestNoBrief(t *testing.T) {
	app := libcmd.NewApp("app", "")

	if err := compareHelpOutput(app.Cmd, "testdata/nobrief.golden"); err != nil {
		t.Error(err)
	}
}

func TestLong(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	if err := compareHelpOutput(app.Cmd, "testdata/long.golden"); err != nil {
		t.Error(err)
	}
}

func TestArgs(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", "s", "", "sets a string value")
	app.Int("aint", "i", 0, "sets a int value")

	if err := compareHelpOutput(app.Cmd, "testdata/args.golden"); err != nil {
		t.Error(err)
	}
}

func TestDefault(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", "s", "somevalue", "sets a string value")
	app.Int("aint", "i", 100, "sets a int value")

	if err := compareHelpOutput(app.Cmd, "testdata/default.golden"); err != nil {
		t.Error(err)
	}
}

func TestCommand(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", "s", "somevalue", "sets a string value")
	app.Int("aint", "i", 100, "sets a int value")

	app.Command("add", "Sums two numbers.", nil)
	app.Command("sub", "Subtract two numbers.", nil)

	if err := compareHelpOutput(app.Cmd, "testdata/command.golden"); err != nil {
		t.Error(err)
	}
}

func TestSubcommand(t *testing.T) {
	var ok bool

	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", "s", "somevalue", "sets a string value")
	app.Int("aint", "i", 100, "sets a int value")

	app.Command("add", "Sums two numbers.", func(cmd *libcmd.Cmd) {
		cmd.Long = "Runs a computation that returns the sum of two specified numbers."

		cmd.Match(func() {
			ok = true
			if err := compareHelpOutput(cmd, "testdata/subcommand.golden"); err != nil {
				t.Error(err)
			}
		})
	})
	app.Command("sub", "Subtract two numbers.", nil)

	if err := app.RunArgs([]string{"add"}); err != nil {
		t.Error(err)
	}

	if !ok {
		t.Error("The assertion did not run!")
	}
}
