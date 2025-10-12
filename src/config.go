package main

import "flag"

type flags struct {
	bank           string
	filePath       string
	debug          bool
	openDefaultApp bool
}

func setupCommandLineFlags() flags {
	var flags flags

	flag.StringVar(&flags.filePath, "file", "", "(Required) Path to the file containing bank statement lines")
	flag.StringVar(&flags.filePath, "f", "", "Alias for -file")
	flag.StringVar(&flags.bank, "bank", "", "(Required) Name of the bank that generated the statement")
	flag.StringVar(&flags.bank, "b", "", "Alias for -bank")
	flag.BoolVar(&flags.debug, "debug", false, "(Optional) Enable debugging info")
	flag.BoolVar(&flags.debug, "d", false, "Alias for -debug")
	flag.BoolVar(&flags.openDefaultApp, "open-default-app", true, "(Optional) Open output file with the default application (by default true)")
	flag.BoolVar(&flags.openDefaultApp, "o", true, "Alias for -open-default-app")
	flag.Parse()

	checkRequiredFlags(flags)

	return flags
}

func checkRequiredFlags(flags flags) {
	if flags.filePath == "" {
		panic("File path was not provided.")
	}

	if flags.bank == "" {
		panic("Bank was not provided.")
	}
}
