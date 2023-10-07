package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type Registry struct {
	Host string
	Port int
}

type RegistryClient interface {
	Register(address string, port int, name string, tags []string, id string) error
	DeRegister(serviceId string) error
}

func NewRegistryClient(host string, port int) RegistryClient {
	return &Registry{
		Host: host,
		Port: port,
	}
}

// Register 将gRPC服务注册到consul
func (r *Registry) Register(address string, port int, name string, tags []string, id string) error {
	// 创建连接consul服务配置
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	// 健康检查
	check := &api.AgentServiceCheck{
		GRPC:     fmt.Sprintf("%s:%d", address, port), // 这里一定是外部可以访问的地址
		Timeout:  "30s",                               // 超时时间
		Interval: "6s",                                // 运行检查的频率
		// 指定时间后自动注销不健康的服务节点
		// 最小超时时间为15s，收获不健康服务的进程每30秒运行一次，因此触发注销的时间可能略长于配置的超时时间。
		DeregisterCriticalServiceAfter: "15s",
	}
	srv := &api.AgentServiceRegistration{
		ID:      id,      // 服务唯一ID
		Name:    name,    // 服务名称
		Tags:    tags,    // 为服务打标签
		Address: address, // 服务地址
		Port:    port,    // 服务端口
		Check:   check,   // 健康检查
	}
	return client.Agent().ServiceRegister(srv)
}

// DeRegister 注销方法
func (r *Registry) DeRegister(serviceId string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	err = client.Agent().ServiceDeregister(serviceId)
	return err
}
