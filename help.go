package libcmd

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Help prints the help text in the stdout.
// Normally, this method is not manually called, since it will be
// automatically executed on the following scenarios:
//
//   1. When '-h' or '--help' is used;
//   2. In a multi-command scenario, when a given command has subcommands
//      but does not have a Run callback;
//
// In the case defined by (1), you can suppress this behavior by setting
// the option SupressPrintHelpWhenSet. This will not call the help function when
// '-h' or '--help' is passed, but the flags will still be defined if you have at
// least one more flag defined. To suppress even the '-h' and '--help' flags
// creation, set the SuppressHelpFlag option.
//
// In the sceenario defined by (2), you can either define a Run callback for the
// command or set the SuppressPrintHelpPartialCommand configuration option.
//
// Last, but not least, if you need to override the actual text of the help, set
// the OnHelp field of the App options instance.
func (cmd *Cmd) Help() {
	cmd.configure()
	output := cmd.options.HelpOutput

	if output == nil {
		output = os.Stdout
	}

	cmd.PrintHelp(output)
}

// PrintHelp prints the help text to the specified writer.
// It functions exactly like the Help method.
func (cmd *Cmd) PrintHelp(writer io.Writer) {
	cmd.options.OnHelp(cmd, writer)
}

func automaticHelp(cmd *Cmd, writer io.Writer) {
	printHelpBrief(cmd, writer)
	printHelpLong(cmd, writer)
	printHelpOptions(cmd, writer)
	printHelpCommands(cmd, writer)
}

func printHelpBrief(cmd *Cmd, writer io.Writer) {
	fullname := strings.TrimSpace(cmd.breadcrumbs + " " + cmd.Name)

	if cmd.Brief != "" {
		fmt.Fprintf(writer, "%s - %s\n", fullname, cmd.Brief)
	} else {
		fmt.Fprintln(writer, fullname)
	}
}

func printHelpLong(cmd *Cmd, writer io.Writer) {
	if cmd.Long != "" {
		fmt.Fprintf(writer, "\n%s\n", cmd.Long)
	}
}

func printHelpOptions(cmd *Cmd, writer io.Writer) {
	if len(cmd.optentries) == 0 {
		return
	}
	sort.Slice(cmd.optentries, func(i, j int) bool {
		return cmd.optentries[i].helpHeader() < cmd.optentries[j].helpHeader()
	})

	fmt.Fprintf(writer, "\nOptions:\n")

	for _, entry := range cmd.optentries {
		fmt.Fprintf(writer, "\n  %s", entry.helpHeader())

		if defHelp := entry.val.defaultAsString(); defHelp != "" {
			fmt.Fprintf(writer, " (default: %s)", defHelp)
		}

		fmt.Fprintln(writer)

		if entry.help != "" {
			fmt.Fprintf(writer, "      %s\n", entry.help)
		}
	}

}

func printHelpCommands(cmd *Cmd, writer io.Writer) {
	if len(cmd.commands) == 0 {
		return
	}

	largest := 0
	keys := make([]string, 0, len(cmd.commands))

	for k := range cmd.commands {
		keys = append(keys, k)

		if len(k) > largest {
			largest = len(k)
		}
	}

	sort.Strings(keys)
	fmt.Fprintf(writer, "\nCommands:\n")

	for _, k := range keys {
		c := cmd.commands[k]
		fmt.Fprintf(writer, "  %-*s   %s\n", largest, c.Name, c.Brief)
	}

}
