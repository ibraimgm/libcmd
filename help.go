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
// You can override/disable the behavior descirbed in (1) by
// defining a flag with the same name; likewise, providing a Run callback
// in the cases described by (2) will prevent the automatic help in those
// scenarios.
//
// If you need to override the actual text of the help, set the OnHelp
// field of the associated App instance.
func (cmd *Cmd) Help() {
	cmd.PrintHelp(os.Stdout)
}

// PrintHelp prints the help text to the specified writer.
// It functions exactly like the Help method.
func (cmd *Cmd) PrintHelp(writer io.Writer) {
	if cmd.helpHandler != nil {
		cmd.helpHandler(cmd, writer)
	} else {
		automaticHelp(cmd, writer)
	}
}

func automaticHelp(cmd *Cmd, writer io.Writer) {
	fullname := strings.TrimSpace(cmd.breadcrumbs + " " + cmd.Name)

	if cmd.Brief != "" {
		fmt.Fprintf(writer, "%s - %s\n", fullname, cmd.Brief)
	} else {
		fmt.Fprintln(writer, fullname)
	}

	if cmd.Long != "" {
		fmt.Fprintf(writer, "\n%s\n", cmd.Long)
	}

	if len(cmd.optentries) > 0 {
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

	if len(cmd.commands) > 0 {
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
}
