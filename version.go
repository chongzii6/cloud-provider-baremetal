package main

import (
	"fmt"
	"os"

	"github.com/chongzii6/cloud-provider-baremetal/baremetalcp"
	"github.com/spf13/pflag"
)

var (
	version     string = "unknwon"
	versionFlag bool
)

func addVersionFlag() {
	pflag.BoolVar(&versionFlag, "version", false, "Print version information and quit")
}

func printAndExitIfRequested() {
	if versionFlag {
		fmt.Printf("%s %s\n", baremetalcp.ProviderName, version)
		os.Exit(0)
	}
}
