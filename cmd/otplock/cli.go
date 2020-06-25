package main

import (
	"os"
	"strings"

	"gitlab.com/mjwhitta/cli"
	hl "gitlab.com/mjwhitta/hilighter"
	"gitlab.com/mjwhitta/otplock"
)

// Flags
type cliFlags struct {
	nocolor bool
	port    int
	unsafe  bool
	verbose bool
	version bool
}

var flags cliFlags

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
			"1: Invalid option\n",
			"2: Invalid argument\n",
			"3: Extra arguments\n",
			"4: Exception",
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
		"Show show stacktrace if error.",
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
		cli.Usage(ExtraArguments)
	}
}
