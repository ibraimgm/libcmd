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
	output := cmd.Options.HelpOutput

	if output == nil {
		output = os.Stdout
	}

	cmd.PrintHelp(output)
}

// PrintHelp prints the help text to the specified writer.
// It functions exactly like the Help method.
func (cmd *Cmd) PrintHelp(writer io.Writer) {
	handler := automaticHelp

	if cmd.Options.OnHelp != nil {
		handler = cmd.Options.OnHelp
	}

	handler(cmd, writer)
}

func automaticHelp(cmd *Cmd, writer io.Writer) {
	printHelpBrief(cmd, writer)
	printHelpUsage(cmd, writer)
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

func printHelpUsage(cmd *Cmd, writer io.Writer) {
	// do not print usage line
	if cmd.Usage == "-" {
		return
	}

	// use a custom usage line
	if cmd.Usage != "" {
		fmt.Fprintf(writer, "\nUSAGE: %s\n", cmd.Usage)
		return
	}

	// compute a usage line
	usage := strings.TrimSpace(cmd.breadcrumbs + " " + cmd.Name)

	if len(cmd.optentries) > 0 {
		usage += " [OPTIONS...]"
	}

	var params string
	operands := getHelpOperands(cmd)
	commands := getHelpCommands(cmd)

	switch {
	case operands != "" && commands != "":
		params = "[" + operands + " | " + commands + "]"

	case commands != "":
		params = commands

	case operands == "OPERANDS...":
		params = "[" + operands + "]"

	default:
		params = operands
	}

	fmt.Fprintf(writer, "\nUSAGE: %s\n", strings.TrimSpace(usage+" "+params))
}

func getHelpOperands(cmd *Cmd) string {
	if len(cmd.operands) == 0 {
		if cmd.Options.StrictOperands || len(cmd.commands) > 0 {
			return ""
		}

		return "OPERANDS..."
	}

	var operands string
	for _, op := range cmd.operands {
		operand := op.name

		switch op.modifier {
		case "*":
			operand = "[" + operand + "...]"
		case "?":
			operand = "[" + operand + "]"
		default:
			operand += op.modifier
		}

		operands += " " + operand
	}

	return strings.TrimSpace(operands)
}

func getHelpCommands(cmd *Cmd) string {
	if len(cmd.commands) > 0 {
		return "COMMAND"
	}

	return ""
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
		fmt.Fprintf(writer, "  %-24s  %s\n", entry.helpHeader(), entry.helpExplain())
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
