package baremetalcp

import (
	"context"

	"github.com/golang/glog"

	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

const configMapAnnotationKey = "k8s.co/cloud-provider-config"
const serviceForwardMethodAnnotationKey = "k8s.co/keepalived-forward-method"

type BmLoadBalancer struct {
	kubeClient      *kubernetes.Clientset
	namespace, name string
}

var _ cloudprovider.LoadBalancer = &BmLoadBalancer{}

func NewBMLoadBalancer(kubeClient *kubernetes.Clientset, ns, name string) cloudprovider.LoadBalancer {
	return &BmLoadBalancer{kubeClient, ns, name}
}

// TODO: Break this up into different interfaces (LB, etc) when we have more than one type of service
// GetLoadBalancer returns whether the specified load balancer exists, and
// if so, what its status is.
// Implementations must treat the *v1.Service parameter as read-only and not modify it.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (k *BmLoadBalancer) GetLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error) {
	// GetLoadBalancer(clusterName string, service *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error)
	// cm, err := k.getConfigMap()

	// if err != nil {
	// 	return nil, false, err
	// }

	// cfg, err := configFrom(cm)

	// if err != nil {
	// 	return nil, false, err
	// }

	// for _, svc := range cfg.Services {
	// 	if svc.UID == string(service.UID) {
	// 		return &v1.LoadBalancerStatus{
	// 			Ingress: []v1.LoadBalancerIngress{{IP: svc.IP}},
	// 		}, true, nil
	// 	}
	// }

	glog.Infof("GetLoadBalancer: %s/%s", clusterName, service.GetName())

	return nil, false, nil
}

// EnsureLoadBalancer creates a new load balancer 'name', or updates the existing one. Returns the status of the balancer
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (k *BmLoadBalancer) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	// EnsureLoadBalancer(clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error)
	// return k.syncLoadBalancer(service)
	glog.Infof("EnsureLoadBalancer: %s/%s", clusterName, service.GetName())
	return nil, nil
}

// UpdateLoadBalancer updates hosts under the specified load balancer.
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (k *BmLoadBalancer) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	// UpdateLoadBalancer(clusterName string, service *v1.Service, nodes []*v1.Node) error
	// _, err := k.syncLoadBalancer(service)
	// return err
	glog.Infof("UpdateLoadBalancer: %s/%s", clusterName, service.GetName())
	return nil
}

// EnsureLoadBalancerDeleted deletes the specified load balancer if it
// exists, returning nil if the load balancer specified either didn't exist or
// was successfully deleted.
// This construction is useful because many cloud providers' load balancers
// have multiple underlying components, meaning a Get could say that the LB
// doesn't exist even if some part of it is still laying around.
// Implementations must treat the *v1.Service parameter as read-only and not modify it.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (k *BmLoadBalancer) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	// EnsureLoadBalancerDeleted(clusterName string, service *v1.Service) error
	// return k.deleteLoadBalancer(service)
	glog.Infof("EnsureLoadBalancerDeleted: %s/%s", clusterName, service.GetName())
	return nil
}
