package libcmd_test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/ibraimgm/libcmd"
)

func compareHelpOutput(app *libcmd.App, options libcmd.Options, args []string, goldenfile string) error {
	bytes, err := ioutil.ReadFile(goldenfile)
	if err != nil {
		return err
	}
	expected := string(bytes)

	var b strings.Builder
	options.HelpOutput = &b
	app.Configure(options)

	if len(args) == 0 {
		app.Help()
	} else if err := app.RunArgs(args); err != nil {
		return err
	}

	actual := b.String()

	if expected != actual {
		return fmt.Errorf("Wrong output. Expected:\n>>>\n%s\n<<<\nActual:\n>>>\n%s\n<<<", expected, actual)
	}

	return nil
}

func TestBasic(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")

	if err := compareHelpOutput(app, libcmd.Options{}, []string{}, "testdata/basic.golden"); err != nil {
		t.Error(err)
	}
}

func TestNoBrief(t *testing.T) {
	app := libcmd.NewApp("app", "")

	if err := compareHelpOutput(app, libcmd.Options{}, []string{}, "testdata/nobrief.golden"); err != nil {
		t.Error(err)
	}
}

func TestLong(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	if err := compareHelpOutput(app, libcmd.Options{}, []string{}, "testdata/long.golden"); err != nil {
		t.Error(err)
	}
}

func TestArgs(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", "s", "", "sets a string value")
	app.Int("aint", "i", 0, "sets a int value")

	if err := compareHelpOutput(app, libcmd.Options{}, []string{"-h"}, "testdata/args.golden"); err != nil {
		t.Error(err)
	}
}

func TestArgsPartial(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("", "s", "", "sets a string value")
	app.Int("aint", "", 0, "sets a int value")

	if err := compareHelpOutput(app, libcmd.Options{}, []string{"-h"}, "testdata/partial.golden"); err != nil {
		t.Error(err)
	}
}

func TestArgsNoHelp(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", "s", "", "sets a string value")
	app.Int("aint", "i", 0, "sets a int value")

	if err := compareHelpOutput(app, libcmd.Options{SuppressHelpFlag: true}, []string{}, "testdata/nohelp.golden"); err != nil {
		t.Error(err)
	}
}

func TestDefault(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", "s", "somevalue", "sets a string value")
	app.Int("aint", "i", 100, "sets a int value")

	if err := compareHelpOutput(app, libcmd.Options{}, []string{"-h"}, "testdata/default.golden"); err != nil {
		t.Error(err)
	}
}

func TestCommand(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.Command("add", "Sums two numbers.", nil)
	app.Command("sub", "Subtract two numbers.", nil)

	if err := compareHelpOutput(app, libcmd.Options{}, []string{"-h"}, "testdata/command.golden"); err != nil {
		t.Error(err)
	}
}

func TestCommandArgs(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", "s", "somevalue", "sets a string value")
	app.Int("aint", "i", 100, "sets a int value")

	app.Command("add", "Sums two numbers.", nil)
	app.Command("sub", "Subtract two numbers.", nil)

	if err := compareHelpOutput(app, libcmd.Options{}, []string{"-h"}, "testdata/commandargs.golden"); err != nil {
		t.Error(err)
	}
}

func TestSubcommand(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", "s", "somevalue", "sets a string value")
	app.Int("aint", "i", 100, "sets a int value")

	app.Command("add", "Sums two numbers.", func(cmd *libcmd.Cmd) {
		cmd.Long = "Runs a computation that returns the sum of two specified numbers."
	})
	app.Command("sub", "Subtract two numbers.", nil)

	if err := compareHelpOutput(app, libcmd.Options{}, []string{"add"}, "testdata/subcommand.golden"); err != nil {
		t.Error(err)
	}
}

func TestSubcommandShallow(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", "s", "somevalue", "sets a string value")
	app.Int("aint", "i", 100, "sets a int value")

	app.Command("add", "Sums two numbers.", func(cmd *libcmd.Cmd) {
		cmd.Long = "Runs a computation that returns the sum of two specified numbers."
		cmd.Command("deep", "A deep subcommand.", nil)
	})
	app.Command("sub", "Subtract two numbers.", nil)

	if err := compareHelpOutput(app, libcmd.Options{}, []string{"add"}, "testdata/shallow.golden"); err != nil {
		t.Error(err)
	}
}

func TestSubcommandDeep(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", "s", "somevalue", "sets a string value")
	app.Int("aint", "i", 100, "sets a int value")

	app.Command("add", "Sums two numbers.", func(cmd *libcmd.Cmd) {
		cmd.Long = "Runs a computation that returns the sum of two specified numbers."
		cmd.Command("deep", "A deep subcommand.", func(cmd *libcmd.Cmd) {
			cmd.Long = "This is a deep subcommand."
		})
	})
	app.Command("sub", "Subtract two numbers.", nil)

	if err := compareHelpOutput(app, libcmd.Options{}, []string{"add", "deep"}, "testdata/deep.golden"); err != nil {
		t.Error(err)
	}
}
