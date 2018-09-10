package main

import (
	"os"

	"github.com/spf13/pflag"
)

var (
	version     = "unknwon"
	versionFlag bool
)

func addVersionFlag() {
	pflag.BoolVar(&versionFlag, "version", false, "Print version information and quit")
}

func printAndExitIfRequested() {
	if versionFlag {
		// fmt.Printf("%s %s\n", baremetalcp.ProviderName, version)
		os.Exit(0)
	}
}
