package main

import "flag"

type flags struct {
	fileName string
	debug    bool
}

func setupCommandLineFlags() flags {
	var flags flags

	flag.StringVar(&flags.fileName, "file", "", "(Required) Path to the file containing bank statement lines")
	flag.StringVar(&flags.fileName, "f", "", "Alias for -file")
	flag.BoolVar(&flags.debug, "debug", false, "(Optional) Enable debugging info")
	flag.BoolVar(&flags.debug, "d", false, "Alias for -debug")
	flag.Parse()

	return flags
}
