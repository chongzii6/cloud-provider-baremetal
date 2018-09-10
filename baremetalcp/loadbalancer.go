package baremetalcp

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chongzii6/haproxy-kube-agent/agent"

	"github.com/golang/glog"

	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

const configMapAnnotationKey = "k8s.co/cloud-provider-config"
const serviceForwardMethodAnnotationKey = "k8s.co/keepalived-forward-method"

//BmLoadBalancer loadbalancer
type BmLoadBalancer struct {
	kubeClient *kubernetes.Clientset
	config     HTConfig
}

var _ cloudprovider.LoadBalancer = &BmLoadBalancer{}

//ErrTimeout = timeout
var ErrTimeout = fmt.Errorf("ErrTimeout")

//NewBMLoadBalancer new BmLoadBalancer struct
func NewBMLoadBalancer(kubeClient *kubernetes.Clientset, config HTConfig) cloudprovider.LoadBalancer {
	return &BmLoadBalancer{kubeClient, config}
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

	loadBalancerName := cloudprovider.GetLoadBalancerName(service)
	ip, err := k.config.GetLoadBalancer(loadBalancerName)
	if err == nil {
		return &v1.LoadBalancerStatus{
			Ingress: []v1.LoadBalancerIngress{{IP: ip}},
		}, true, nil
	}

	return nil, false, nil
}

// EnsureLoadBalancer creates a new load balancer 'name', or updates the existing one. Returns the status of the balancer
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (k *BmLoadBalancer) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	// EnsureLoadBalancer(clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error)
	// return k.syncLoadBalancer(service)
	glog.V(4).Infof("EnsureLoadBalancer(%v, %v, %v, %v, %v, %v, %v)", clusterName, service.Namespace, service.Name, service.Spec.LoadBalancerIP, service.Spec.Ports, nodes, service.Annotations)
	// glog.Infof("EnsureLoadBalancer: %s/%s", clusterName, service.GetName())
	if len(nodes) == 0 {
		return nil, fmt.Errorf("there are no available nodes for LoadBalancer service %s/%s", service.Namespace, service.Name)
	}

	// loadBalancerIP := service.Spec.LoadBalancerIP
	loadBalancerIP, err := k.addLBReq(service, nodes, false)

	if err == nil {
		status := &v1.LoadBalancerStatus{}
		status.Ingress = []v1.LoadBalancerIngress{{IP: loadBalancerIP}}
		return status, nil
	}
	return nil, err
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
	if len(nodes) == 0 {
		return fmt.Errorf("there are no available nodes for LoadBalancer service %s/%s", service.Namespace, service.Name)
	}
	_, err := k.addLBReq(service, nodes, true)
	return err
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
	svcName := service.GetName()
	glog.Infof("EnsureLoadBalancerDeleted: %s/%s", clusterName, svcName)
	loadBalancerName := cloudprovider.GetLoadBalancerName(service)
	loadBalancerIP := service.Spec.LoadBalancerIP
	lbChannel := k.config.GetLBChannel(loadBalancerIP)
	req := &agent.Request{
		Action:  agent.DELETE,
		LbName:  loadBalancerName,
		SvcName: svcName,
	}

	done := make(chan error)
	go func() {
		_, err := k.config.WaitforLbReady(loadBalancerName)
		done <- err
	}()

	<-time.After(time.Second)
	err := k.config.SendReq(lbChannel, req)
	if err != nil {
		log.Println(err)
		return err
	}
	err = <-done
	if err == ErrTimeout {
		err = nil
	}
	return err
}

// The LB needs to be configured with instance addresses on the same
// subnet as the LB (aka opts.SubnetID).  Currently we're just
// guessing that the node's InternalIP is the right address - and that
// should be sufficient for all "normal" cases.
func (k *BmLoadBalancer) nodeAddressForLB(node *v1.Node) (string, error) {
	addrs := node.Status.Addresses
	if len(addrs) == 0 {
		return "", fmt.Errorf("ErrNoAddressFound")
	}

	for _, addr := range addrs {
		if addr.Type == v1.NodeInternalIP {
			return addr.Address, nil
		}
	}

	return addrs[0].Address, nil
}

func (k *BmLoadBalancer) addLBReq(service *v1.Service, nodes []*v1.Node, update bool) (string, error) {

	svcName := service.Name
	loadBalancerName := cloudprovider.GetLoadBalancerName(service)

	_ = service.Annotations
	loadBalancerIP := service.Spec.LoadBalancerIP
	lbChannel := k.config.GetLBChannel(loadBalancerIP)

	ports := service.Spec.Ports
	for _, port := range ports {

		endps := []agent.Endpoint{}
		for _, node := range nodes {
			addr, err := k.nodeAddressForLB(node)
			if err != nil {
				continue
			}

			ep := agent.Endpoint{
				Name: node.Name,
				IP:   addr,
				Port: int(port.NodePort),
			}
			endps = append(endps, ep)
		}

		var act agent.RequestType
		if update {
			act = agent.UPDATE
		} else {
			act = agent.ADD
		}

		req := &agent.Request{
			Action:     act,
			LbName:     loadBalancerName,
			TargetPort: int(port.Port),
			Endpoints:  endps,
			SvcName:    svcName,
		}

		var ip string
		done := make(chan error)
		go func() {
			var err error
			ip, err = k.config.WaitforLbReady(loadBalancerName)
			done <- err
		}()

		<-time.After(time.Second)

		err := k.config.SendReq(lbChannel, req)
		if err != nil {
			log.Println(err)
			return "", err
		}
		err = <-done
		return ip, err
	}
	return "", fmt.Errorf("ErrNoServicePort")
}
