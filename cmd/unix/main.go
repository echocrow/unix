package main

import (
	"flag"
	"fmt"
	"os"
)

// version is the version of this app set at build-time.
var version = "0.0.0-dev"

func main() {
	var getVersion bool
	flag.BoolVar(&getVersion, "version", false, "Print the app version & exit.")

	flag.Parse()

	if getVersion {
		exitMessage(0, fmt.Sprint("unix ", version))
	}

	// @todo
}

func exitMessage(code int, msg string) {
	flag.CommandLine.SetOutput(nil)
	fmt.Fprintln(flag.CommandLine.Output(), msg)
	os.Exit(code)
}
