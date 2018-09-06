package baremetalcp

import (
	"fmt"
	"io"
	"os"

	"github.com/golang/glog"
	gcfg "gopkg.in/gcfg.v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
)

const (
	ProviderName = "htnm"
)

func readConfig(config io.Reader) (HTConfig, error) {
	if config == nil {
		err := fmt.Errorf("cloud provider config file is missing. Please start with --cloud-provider=%s --cloud-config=[path_to_config_file]", ProviderName)
		return HTConfig{}, err
	}

	var cfg HTConfig
	err := gcfg.ReadInto(&cfg, config)
	return cfg, err
}

func init() {
	// cloudprovider.RegisterCloudProvider(ProviderName, newBmCloudProvider)
	cloudprovider.RegisterCloudProvider(ProviderName, func(config io.Reader) (cloudprovider.Interface, error) {
		cfg, err := readConfig(config)
		if err != nil {
			glog.Errorf("Photon Cloud Provider: failed to read in cloud provider config file. Error[%v]", err)
			return nil, err
		}

		return newBmCloudProvider(cfg)
	})
}

type BmCloudProvider struct {
	lb cloudprovider.LoadBalancer
}

var _ cloudprovider.Interface = &BmCloudProvider{}

func newBmCloudProvider(config HTConfig) (cloudprovider.Interface, error) {
	ns := os.Getenv("CLOUDPROVIDER_NAMESPACE")
	cm := os.Getenv("CLOUDPROVIDER_CONFIG_MAP")

	cfg, err := rest.InClusterConfig()

	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes client config: %s", err.Error())
	}

	clients, err := kubernetes.NewForConfig(cfg)

	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes client: %s", err.Error())
	}

	return &BmCloudProvider{NewBMLoadBalancer(clients, ns, cm)}, nil
}

// Initialize provides the cloud with a kubernetes client builder and may spawn goroutines
// to perform housekeeping activities within the cloud provider.
func (k *BmCloudProvider) Initialize(clientBuilder controller.ControllerClientBuilder) {
	return
}

// LoadBalancer returns a loadbalancer interface. Also returns true if the interface is supported, false otherwise.
func (k *BmCloudProvider) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	// lb, err := os.NewLoadBalancerV2()
	return k.lb, true
}

// Instances returns an instances interface. Also returns true if the interface is supported, false otherwise.
func (k *BmCloudProvider) Instances() (cloudprovider.Instances, bool) {
	return nil, false
}

// Zones returns a zones interface. Also returns true if the interface is supported, false otherwise.
func (k *BmCloudProvider) Zones() (cloudprovider.Zones, bool) {
	return nil, false
}

// Clusters returns a clusters interface.  Also returns true if the interface is supported, false otherwise.
func (k *BmCloudProvider) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

// Routes returns a routes interface along with whether the interface is supported.
func (k *BmCloudProvider) Routes() (cloudprovider.Routes, bool) {
	return nil, false
}

// ProviderName returns the cloud provider ID.
func (k *BmCloudProvider) ProviderName() string {
	return ProviderName
}

// ScrubDNS provides an opportunity for cloud-provider-specific code to process DNS settings for pods.
func (k *BmCloudProvider) ScrubDNS(nameservers, searches []string) (nsOut, srchOut []string) {
	return nil, nil
}

// HasClusterID returns true if a ClusterID is required and set
func (k *BmCloudProvider) HasClusterID() bool {
	return true
}

// type zones struct{}

// func (z zones) GetZone() (cloudprovider.Zone, error) {
// 	return cloudprovider.Zone{FailureDomain: "FailureDomain1", Region: "Region1"}, nil
// }
