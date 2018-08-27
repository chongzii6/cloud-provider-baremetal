package main

import (
	goflag "flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	utilflag "k8s.io/apiserver/pkg/util/flag"
	"k8s.io/apiserver/pkg/util/logs"
	"k8s.io/kubernetes/cmd/cloud-controller-manager/app"
	_ "k8s.io/kubernetes/pkg/client/metrics/prometheus" // for client metric registration
	// NOTE: Importing all in-tree cloud-providers is not required when
	// implementing an out-of-tree cloud-provider.
	_ "k8s.io/kubernetes/pkg/cloudprovider/providers"
	_ "k8s.io/kubernetes/pkg/version/prometheus" // for version metric registration

	"github.com/spf13/pflag"
	"k8s.io/apiserver/pkg/server/healthz"
)

func init() {
	healthz.DefaultHealthz()
}

func main() {
	// s := options.NewCloudControllerManagerServer()
	// s.AddFlags(pflag.CommandLine)
	// addVersionFlag()

	// flag.InitFlags()
	// logs.InitLogs()
	// defer logs.FlushLogs()

	// printAndExitIfRequested()

	// cloud, err := cloudprovider.InitCloudProvider(baremetalcp.ProviderName, s.CloudConfigFile)
	// if err != nil {
	// 	glog.Fatalf("Cloud provider could not be initialized: %v", err)
	// }

	// glog.Info("Starting version ", version)
	// if err := app.Run(s, cloud); err != nil {
	// 	fmt.Fprintf(os.Stderr, "%v\n", err)
	// 	os.Exit(1)
	// }

	rand.Seed(time.Now().UTC().UnixNano())

	command := app.NewCloudControllerManagerCommand()

	// TODO: once we switch everything over to Cobra commands, we can go back to calling
	// utilflag.InitFlags() (by removing its pflag.Parse() call). For now, we have to set the
	// normalize func and add the go flag set by hand.
	pflag.CommandLine.SetNormalizeFunc(utilflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	// utilflag.InitFlags()
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

}
