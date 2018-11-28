package consulsupport

import (
	"fmt"
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

// NewClient 通过指定的config和address创建consul client
func NewClient(config *consulapi.Config, consulAddress string) (*consulapi.Client, error) {
	config.Address = consulAddress
	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// DiscoveryService 用于从consul中获得指定服务名称的所有健康的服务
// 并返回所有服务对应的host:port
func DiscoveryService(client *consulapi.Client, serviceName string) ([]string, error) {
	var hosts []string

	serviceEntry, _, err := client.Health().Service(serviceName, "", true, &consulapi.QueryOptions{})
	if err != nil {
		return nil, err
	}

	for _, entry := range serviceEntry {
		for _, health := range entry.Checks {
			if health.ServiceName != serviceName {
				continue
			}
			ip := entry.Service.Address
			port := entry.Service.Port
			host := fmt.Sprintf("%s:%d", ip, port)
			hosts = append(hosts, host)
		}
	}

	if len(hosts) == 0 {
		return nil, fmt.Errorf("There are no %s service available", serviceName)
	}

	log.Printf("Discovery %s service: %v\n", serviceName, hosts)
	return hosts, nil
}

// GetConfig 从consul中获得对应key的value
func GetConfig(client *consulapi.Client, configName string) (string, error) {
	kv, _, err := client.KV().Get(configName, nil)
	if err != nil {
		return "", err
	}
	if kv == nil {
		return "", fmt.Errorf("There's no %s config in consul", configName)
	}

	log.Printf("Get config %s: %s\n", configName, string(kv.Value))

	return string(kv.Value), nil
}
