package baremetalcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chongzii6/haproxy-kube-agent/agent"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/coreos/etcd/pkg/transport"
	uuid "github.com/satori/go.uuid"
)

//HTConfig for cloud htnm
type HTConfig struct {
	Global struct {
		Cafile    string   `gcfg:"cafile"`
		Keyfile   string   `gcfg:"keyfile"`
		Certfile  string   `gcfg:"certfile"`
		Agentkey  string   `gcfg:"agentkey"`
		Reqkey    string   `gcfg:"reqkey"`
		Endpoints []string `gcfg:"endpoints"`
		DefaultLB string   `gcfg:"defaultlb"`
	}
}

const (
	dialTimeout     = 5 * time.Second
	lbCreateTimeout = 30 * time.Second
)

func (c *HTConfig) newClientCfg() (*clientv3.Config, error) {
	var cfgtls *transport.TLSInfo
	tlsinfo := transport.TLSInfo{}

	cfg := &clientv3.Config{
		Endpoints:   c.Global.Endpoints,
		DialTimeout: dialTimeout,
	}

	tlsinfo.CertFile = c.Global.Certfile
	tlsinfo.KeyFile = c.Global.Keyfile
	tlsinfo.TrustedCAFile = c.Global.Cafile
	cfgtls = &tlsinfo

	clientTLS, err := cfgtls.ClientConfig()
	if err != nil {
		return nil, err
	}
	cfg.TLS = clientTLS

	return cfg, nil
}

func (c *HTConfig) newClient() (*clientv3.Client, error) {
	cfg, err := c.newClientCfg()
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	client, err := clientv3.New(*cfg)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return client, nil
}

//EtcdGet get
func (c *HTConfig) EtcdGet(key string) (string, error) {
	client, err := c.newClient()
	if err != nil {
		return "", err
	}
	defer client.Close()

	var resp *clientv3.GetResponse
	if resp, err = client.Get(context.Background(), key); err != nil {
		log.Fatalln(err)
		return "", err
	}
	log.Println("resp: ", resp)

	var ret string
	for _, kv := range resp.Kvs {
		ret += string(kv.Value)
	}
	return ret, nil
}

//EtcdPut put
func (c *HTConfig) EtcdPut(key string, val string) error {
	client, err := c.newClient()
	if err != nil {
		return err
	}
	defer client.Close()

	var resp *clientv3.PutResponse
	if resp, err = client.Put(context.Background(), key, val); err != nil {
		log.Fatalln(err)
		return err
	}

	log.Println("resp: ", resp)
	return nil
}

//EtcdWatch watch key
func (c *HTConfig) EtcdWatch(key string, timeout time.Duration) (string, error) {
	client, err := c.newClient()
	if err != nil {
		return "", err
	}
	defer client.Close()

	wc := client.Watch(context.Background(), key, clientv3.WithPrefix())

	log.Printf("watching: %s\n", key)
	for {
		select {
		case <-time.After(timeout):
			return "", ErrTimeout
		case resp := <-wc:
			for _, e := range resp.Events {
				log.Printf("%s key:%s, value:%s\n", e.Type, e.Kv.Key, e.Kv.Value)
				if e.Type == mvccpb.PUT {
					// err = HandleReq(e.Kv.Key, e.Kv.Value)
				}
				return string(e.Kv.Value), nil
			}
		}
	}
}

//GetLoadBalancer retrieve from etcd
func (c *HTConfig) GetLoadBalancer(name string) (string, error) {
	lbkey := fmt.Sprintf("%s/%s", c.Global.Agentkey, name)
	text, err := c.EtcdGet(lbkey)
	if err != nil {
		return "", err
	}

	st := &agent.LBState{}
	err = json.Unmarshal([]byte(text), st)
	if err != nil {
		return "", err
	}

	return st.HostIP, nil
}

//SendReq put request to etcd
func (c *HTConfig) SendReq(loadBalancerIP string, req *agent.Request) error {
	u1 := uuid.NewV4()
	key := fmt.Sprintf("%s/%s/%s", c.Global.Reqkey, loadBalancerIP, u1)
	by, err := json.Marshal(req)
	if err == nil {
		c.EtcdPut(key, string(by))
	}

	return err
}

//WaitforResp wait for response of request
func (c *HTConfig) WaitforResp(lbChannel string, lbName string) (string, error) {
	key := fmt.Sprintf("%s/%s/%s", c.Global.Agentkey, lbChannel, lbName)
	ip, err := c.EtcdWatch(key, lbCreateTimeout)
	return ip, err
}

//GetLBChannel get channel name by loadbalancerIP
func (c *HTConfig) GetLBChannel(LbIP string) string {
	if LbIP == "" {
		if c.Global.DefaultLB == "" {
			return "any"
		} else {
			return c.Global.DefaultLB
		}
	} else {
		return LbIP
	}
}

// type config struct {
// 	Services []serviceConfig `json:"services"`
// }

// func (c *config) allocateIP(cidr string) (string, error) {
// 	possible, err := Hosts(cidr)
// 	if err != nil {
// 		return "", err
// 	}

// Outer:
// 	for _, ip := range possible {
// 		for _, svc := range c.Services {
// 			// if this 'ip' candidate is already in use,
// 			// break the inner loop to move onto the next IP address
// 			if svc.IP == ip {
// 				continue Outer
// 			}
// 		}

// 		// if we get to this point, then 'ip' hasn't been allocated already
// 		return ip, nil
// 	}

// 	return "", fmt.Errorf("ip cidr pool exhausted. increase size of cidr or remove some loadbalancers")
// }

// func (c *config) encode() ([]byte, error) {
// 	return json.Marshal(c)
// }

// func (c *config) ensureService(cfg serviceConfig) {
// 	for i, s := range c.Services {
// 		if s.UID == cfg.UID {
// 			glog.Infof("updating service with uid '%s' in config: %s->%s(%s)", cfg.UID, s.IP, cfg.IP, cfg.ForwardMethod)
// 			c.Services[i] = cfg
// 			return
// 		}
// 	}
// 	glog.Infof("adding new service '%s': %s(%s)", cfg.UID, cfg.IP, cfg.ForwardMethod)
// 	c.Services = append(c.Services, cfg)
// 	glog.Infof("there are now %d services in config", len(c.Services))
// }

// func (c *config) deleteService(cfg serviceConfig) {
// 	for i, s := range c.Services {
// 		if s.UID == cfg.UID {
// 			glog.Infof("deleted service with uid %s, ip: %s from config", s.UID, s.IP)
// 			c.Services = append(c.Services[:i], c.Services[i+1:]...)
// 			return
// 		}
// 	}
// }

// type serviceConfig struct {
// 	UID              string `json:"uid"`
// 	IP               string `json:"ip"`
// 	ServiceNamespace string `json:"serviceNamespace"`
// 	ServiceName      string `json:"serviceName"`
// 	ForwardMethod    string `json:"forwardMethod,omitempty"`
// }

// func configFrom(cm *v1.ConfigMap) (*config, error) {
// 	cfg := config{}
// 	if c, ok := cm.Annotations[configMapAnnotationKey]; ok {
// 		err := json.Unmarshal([]byte(c), &cfg)

// 		if err != nil {
// 			return nil, fmt.Errorf("error getting cloud provider config from annotation: %s", err.Error())
// 		}
// 	}
// 	return &cfg, nil
// }

// func (c *config) toConfigMapData() map[string]string {
// 	d := make(map[string]string, len(c.Services))
// 	for _, s := range c.Services {

// 		if s.ForwardMethod != "" {
// 			d[s.IP] = s.ServiceNamespace + "/" + s.ServiceName + ":" + s.ForwardMethod
// 		} else {
// 			d[s.IP] = s.ServiceNamespace + "/" + s.ServiceName
// 		}
// 	}

// 	return d
// }

// // from: https://gist.github.com/kotakanbe/d3059af990252ba89a82
// func Hosts(cidr string) ([]string, error) {
// 	ip, ipnet, err := net.ParseCIDR(cidr)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var ips []string
// 	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
// 		ips = append(ips, ip.String())
// 	}
// 	// remove network address and broadcast address
// 	return ips[1 : len(ips)-1], nil
// }

// func inc(ip net.IP) {
// 	for j := len(ip) - 1; j >= 0; j-- {
// 		ip[j]++
// 		if ip[j] > 0 {
// 			break
// 		}
// 	}
// }
