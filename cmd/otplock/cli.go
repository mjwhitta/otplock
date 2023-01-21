package main

import (
	"os"
	"strings"

	"github.com/mjwhitta/cli"
	hl "github.com/mjwhitta/hilighter"
	"github.com/mjwhitta/otplock"
)

// Exit status
const (
	Good = iota
	InvalidOption
	MissingOption
	InvalidArgument
	MissingArgument
	ExtraArgument
	Exception
)

// Flags
var flags struct {
	nocolor bool
	port    int
	unsafe  bool
	verbose bool
	version bool
}

func init() {
	// Configure cli package
	cli.Align = true
	cli.Authors = []string{"Miles Whittaker <mj@whitta.dev>"}
	cli.Banner = hl.Sprintf("%s [OPTIONS]", os.Args[0])
	cli.BugEmail = "otplock.bugs@whitta.dev"
	cli.ExitStatus = strings.Join(
		[]string{
			"Normally the exit status is 0. In the event of an error",
			"the exit status will be one of the below:\n\n",
			hl.Sprintf("%d: Invalid option\n", InvalidOption),
			hl.Sprintf("%d: Missing option\n", MissingOption),
			hl.Sprintf("%d: Invalid argument\n", InvalidArgument),
			hl.Sprintf("%d: Missing argument\n", MissingArgument),
			hl.Sprintf("%d: Extra argument\n", ExtraArgument),
			hl.Sprintf("%d: Exception", Exception),
		},
		" ",
	)
	cli.Info = "Encode paylods with a one-time-pad (OTP)."
	cli.Title = "One-Time-Padlock"

	// Parse cli flags
	cli.Flag(
		&flags.nocolor,
		"no-color",
		false,
		"Disable colorized output.",
	)
	cli.Flag(
		&flags.port,
		"p",
		"port",
		8080,
		"The port to listen on (default: 8080).",
	)
	cli.Flag(
		&flags.unsafe,
		"u",
		"unsafe",
		false,
		"Allow advanced config and unsafe use of shell commands.",
	)
	cli.Flag(
		&flags.verbose,
		"v",
		"verbose",
		false,
		"Show stacktrace, if error.",
	)
	cli.Flag(&flags.version, "V", "version", false, "Show version.")
	cli.Parse()
}

// Process cli flags and ensure no issues
func validate() {
	hl.Disable(flags.nocolor)

	// Short circuit if version was requested
	if flags.version {
		hl.Printf("otplock version %s\n", otplock.Version)
		os.Exit(Good)
	}

	// Validate cli flags
	if cli.NArg() > 1 {
		cli.Usage(ExtraArgument)
	}
}
