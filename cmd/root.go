package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/echocrow/unix/pkg/unix"
	"github.com/spf13/cobra"
)

// Execute executes the root command.
func Execute(
	version string,
	exit func(int),
) {
	cmd := newCmd(version)
	if err := cmd.Execute(); err != nil {
		exit(1)
	}
}

type cmdOptions struct {
	fromTz string
	addDur time.Duration
	toTz   string
	toLyt  string
}

func newCmd(version string) *cobra.Command {
	opts := cmdOptions{}

	cmd := &cobra.Command{
		Use:   "unix [TIME]",
		Short: "A simple UNIX timestamp and date converter",
		Long: heredoc.Doc(`
			Unix is a CLI that allows easy conversion between formatted dates and UNIX
			timestamps across different timezones and various date formats.
		`),
		Example: indentHeredoc(`
		  unix
		  unix '2005-03-18 01:58:31'
		  unix '1983-01-01 13:37:11' -f long
		  unix 1580650631
		  unix 1580650631 -f '%Y-%m-%d %H:%M:%S'
		  unix 1580650631 -z vienna -Z toronto
		  unix -Z Europe/London -f '%Y-%m-%d %H:%m:%m'
		`),
		Version: version,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return run(c, &opts, args)
		},
	}

	cmd.Flags().StringVarP(&opts.fromTz, "from", "z", "", flushHeredoc(`
		parse the input in a fixed timezone;
		e.g. "utc", "vienna", "America/New_York"
	`))
	cmd.Flags().StringVarP(&opts.toTz, "to", "Z", "", "convert the time to a different timezone")
	cmd.Flags().DurationVarP(&opts.addDur, "add", "a", 0, "add a time offset, e.g. \"10s\", \"30m\", \"-168h\"")
	cmd.Flags().StringVarP(&opts.toLyt, "format", "f", "", flushHeredoc(`
		format the time in a fixed layout;
		e.g. "long", "unix", "%Y-%m-%d %H:%m:%m"
	`))

	cmd.SetOut(os.Stdout)

	return cmd
}

func run(cmd *cobra.Command, opts *cmdOptions, args []string) error {
	timeS := ""
	if len(args) >= 1 {
		timeS = args[0]
	}
	return unixCmd(cmd, opts, timeS)
}

func unixCmd(cmd *cobra.Command, opts *cmdOptions, input string) error {
	t, srcLyt, err := unix.Parse(input, opts.fromTz)
	if err != nil {
		return err
	}

	if opts.toTz != "" {
		toLoc, err := unix.ParseLocation(opts.toTz)
		if err != nil {
			return err
		}
		t = t.In(toLoc)
	}

	t = t.Add(opts.addDur)

	if opts.toLyt == "" {
		if srcLyt != unix.UnixLayout {
			opts.toLyt = unix.UnixFormat
		} else {
			opts.toLyt = unix.LongFormat
		}
	}

	if out, err := unix.Format(t, opts.toLyt); err != nil {
		return err
	} else {
		cmd.Println(out)
	}

	return nil
}

func indentHeredoc(raw string) string {
	d := heredoc.Doc(raw + ".")
	return d[:len(d)-2]
}

func flushHeredoc(raw string) string {
	return strings.TrimSuffix(heredoc.Doc(raw), "\n")
}
