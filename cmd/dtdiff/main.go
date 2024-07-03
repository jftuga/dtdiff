package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/jftuga/dtdiff"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

// this constant was generated by ChatGPT and then manually refined
const usageTemplate string = `Usage:{{if .Runnable}}
 {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
 {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}
Aliases:
 {{.NameAndAliases}}{{end}}{{if .HasExample}}
Examples:
 {{.Example}}{{end}}{{if .HasAvailableSubCommands}}
Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
 {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Globals:
{{FlagUsagesCustom .LocalFlags "nonewline" "help" "version" | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableLocalFlags}}

Flag Group 1 (mutually exclusive with Flag Group 2):
{{FlagUsagesCustom .LocalFlags "start" "end" "stdin" "brief" | trimTrailingWhitespaces}}

Flag Group 2:
{{FlagUsagesCustom .LocalFlags "from" "add" "sub" "recurrence" "until" | trimTrailingWhitespaces}}

Durations:
years months weeks days
hours minutes seconds milliseconds microseconds nanoseconds
example: "1 year 2 months 3 days 4 hours 1 minute 6 seconds"

Brief Durations: (dates are upper, times are lower)
Y    M    W    D
h    m    s    ms    us    ns
examples: 1Y2M3W4D5h6m7s8ms9us1ns, "1Y 2M 3W 4D 5h 6m 7s 8ms 9us 1ns"

Relative Dates:
for the -s, -e, -F, and -U flags, you can use these shortcuts:
now
today (returns same value as now)
yesterday
tomorrow
example: dtdiff -F today -A 7h10m -U tomorrow
{{end}}{{if .HasAvailableInheritedFlags}}
Global Flags:
 {{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}
Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
 {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

var (
	start         string
	end           string
	from          string
	add           string
	sub           string
	recurrence    int
	until         string
	noNewline     bool
	readFromStdin bool
	brief         bool
	usageMsg      string

	rootCmd = &cobra.Command{
		Use:     "dtdiff",
		Version: dtdiff.PgmVersion,
		Short:   "dtdiff: output the difference between date, time or duration",
		Run: func(cmd *cobra.Command, args []string) {
			if (len(start) > 0 && len(end) > 0) || readFromStdin {
				computeStartEnd(start, end, brief)
				return
			}

			if len(from) > 0 && len(add) > 0 {
				if recurrence > 0 {
					computeAddSubWithRecurrence(from, add, 0, recurrence)
				} else if len(until) > 0 {
					computeUntil(from, until, add, 0)
				} else {
					computeAddSub(from, add, 0)
				}
				return
			}
			if len(from) > 0 && len(sub) > 0 {
				if recurrence > 0 {
					computeAddSubWithRecurrence(from, sub, 1, recurrence)
				} else if len(until) > 0 {
					computeUntil(from, until, sub, 1)
				} else {
					computeAddSub(from, sub, 1)
				}
				return
			}
			fmt.Fprintln(os.Stderr, usageMsg)
			os.Exit(0)
		},
	}
)

// FlagUsagesCustom customized to filter and format flags with types
// this function was generated by ChatGPT and then manually refined
func FlagUsagesCustom(flags *pflag.FlagSet, names ...string) string {
	var buf bytes.Buffer
	flags.VisitAll(func(flag *pflag.Flag) {
		for _, name := range names {
			if flag.Name == name {
				shorthand := ""
				if flag.Shorthand != "" {
					shorthand = fmt.Sprintf("-%s, ", flag.Shorthand)
				}
				name := flag.Name
				if flag.Value.Type() == "string" {
					name += " string"
				}
				if flag.Value.Type() == "int" {
					name += " int"
				}
				tabs := "\t"
				if len(name) <= 7 {
					tabs = "\t\t"
				}
				fmt.Fprintf(&buf, "  %s--%s%s%s\n", shorthand, name, tabs, flag.Usage)
			}
		}
	})
	return buf.String()
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&start, "start", "s", "", "start date, time, or a datetime")
	rootCmd.PersistentFlags().StringVarP(&end, "end", "e", "", "end date, time, or a datetime")
	rootCmd.PersistentFlags().StringVarP(&from, "from", "F", "", "a base date, time or datetime to use with -A or -S")
	rootCmd.PersistentFlags().StringVarP(&add, "add", "A", "", "add: a duration to use with -F, such as '1 day 2 hours 3 seconds'")
	rootCmd.PersistentFlags().StringVarP(&sub, "sub", "S", "", "subtract: a duration to use with -F, such as '5 months 4 weeks 3 days'")
	rootCmd.PersistentFlags().IntVarP(&recurrence, "recurrence", "R", 0, "repeat period this number of times (mutually exclusive with -U)")
	rootCmd.PersistentFlags().StringVarP(&until, "until", "U", "", "repeat period until date/time is exceeded")
	rootCmd.PersistentFlags().BoolVarP(&noNewline, "nonewline", "n", false, "do not output a newline character")
	rootCmd.PersistentFlags().BoolVarP(&readFromStdin, "stdin", "i", false, "read from STDIN instead of using -s/-e")
	rootCmd.PersistentFlags().BoolVarP(&brief, "brief", "b", false, "output in brief format, such as: 1Y2M3D4h5m6s7ms8us9ns")

	rootCmd.MarkFlagsRequiredTogether("start", "end")
	rootCmd.MarkFlagsMutuallyExclusive("add", "sub")
	rootCmd.MarkFlagsMutuallyExclusive("stdin", "start")
	rootCmd.MarkFlagsMutuallyExclusive("stdin", "end")
	rootCmd.MarkFlagsMutuallyExclusive("stdin", "from")
	rootCmd.MarkFlagsMutuallyExclusive("stdin", "add")
	rootCmd.MarkFlagsMutuallyExclusive("stdin", "sub")
	rootCmd.MarkFlagsMutuallyExclusive("stdin", "recurrence")
	rootCmd.MarkFlagsMutuallyExclusive("stdin", "until")
	rootCmd.MarkFlagsMutuallyExclusive("from", "start")
	rootCmd.MarkFlagsMutuallyExclusive("from", "end")
	rootCmd.MarkFlagsMutuallyExclusive("recurrence", "start")
	rootCmd.MarkFlagsMutuallyExclusive("recurrence", "end")
	rootCmd.MarkFlagsMutuallyExclusive("until", "start")
	rootCmd.MarkFlagsMutuallyExclusive("until", "end")
	rootCmd.MarkFlagsMutuallyExclusive("add", "start")
	rootCmd.MarkFlagsMutuallyExclusive("add", "end")
	rootCmd.MarkFlagsMutuallyExclusive("sub", "start")
	rootCmd.MarkFlagsMutuallyExclusive("sub", "end")
	rootCmd.MarkFlagsMutuallyExclusive("brief", "from")
	rootCmd.MarkFlagsMutuallyExclusive("brief", "add")
	rootCmd.MarkFlagsMutuallyExclusive("brief", "sub")
	rootCmd.MarkFlagsMutuallyExclusive("brief", "recurrence")
	rootCmd.MarkFlagsMutuallyExclusive("brief", "until")
	rootCmd.MarkFlagsMutuallyExclusive("until", "recurrence")

	versionTemplate := fmt.Sprintf("%s v%s\n%s\n", dtdiff.PgmName, dtdiff.PgmVersion, dtdiff.PgmUrl)
	rootCmd.SetVersionTemplate(versionTemplate)
	// Register the custom template function
	cobra.AddTemplateFunc("FlagUsagesCustom", func(flags *pflag.FlagSet, names ...string) string {
		return FlagUsagesCustom(flags, names...)
	})

	// Set the custom usage template
	rootCmd.SetUsageTemplate(usageTemplate)
	usageMsg = rootCmd.UsageString()
}

// either read one line containing a comma, then split start and end on this
// or read two lines with start on line one and end on line two
func getInput() (string, string) {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	line := input.Text()
	if strings.Contains(line, ",") {
		split := strings.Split(line, ",")
		if len(split) != 2 {
			fmt.Fprintf(os.Stderr, "invalid stdin input: %s\n", line)
			os.Exit(1)
		}
		return split[0], split[1]
	}
	input.Scan()
	end := input.Text()
	return line, end
}

// computeStartEnd used when -s and -e is given
func computeStartEnd(start, end string, brief bool) {
	if readFromStdin {
		start, end = getInput()
	}

	dt := dtdiff.New(start, end)
	dt.SetBrief(brief)
	format, _, err := dt.DtDiff()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if noNewline {
		fmt.Print(format)
	} else {
		fmt.Println(format)
	}
}

// computeAddSub used when -F is given along with
// add or subtract a duration from "from"
// index 0 = add; index = 1 = sub
func computeAddSub(from, period string, index int) {
	format := ""
	var err error
	if index == 0 {
		format, err = dtdiff.Add(from, period)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		format, err = dtdiff.Sub(from, period)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if noNewline {
		fmt.Print(format)
	} else {
		fmt.Println(format)
	}
}

// computeAddSubWithRecurrence is similar to computeAddSub
// but returns a slice of date/time intervals
// when -n is invoked, a comma-delimited output is used
// index 0 = add; index = 1 = sub
func computeAddSubWithRecurrence(from, period string, index, recurrence int) {
	var format []string
	var err error
	if index == 0 {
		format, err = dtdiff.AddWithRecurrence(from, period, recurrence)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		format, err = dtdiff.SubWithRecurrence(from, period, recurrence)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if noNewline {
		fmt.Print(strings.Join(format, ","))
	} else {
		for _, f := range format {
			fmt.Println(f)
		}
	}
}

func computeUntil(from, until, period string, index int) {
	var format []string
	var err error
	if index == 0 {
		format, err = dtdiff.AddUntil(from, until, period)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		format, err = dtdiff.SubUntil(from, until, period)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if noNewline {
		fmt.Print(strings.Join(format, ","))
	} else {
		for _, f := range format {
			fmt.Println(f)
		}
	}
}
