package grpc_go_starter

import (
	"fmt"
	"github.com/zhengheng7913/grpc-go-starter/client"
	"github.com/zhengheng7913/grpc-go-starter/filter"
	"github.com/zhengheng7913/grpc-go-starter/naming/discovery"
)

func NewClients(opt ...client.Option) *client.Clients {
	cfg, err := LoadSetup()
	if err != nil {
		panic(fmt.Errorf("load setup error: %s", err))
	}
	return NewClientsWithConfig(cfg, opt...)
}

func newClientWithConfig(cfg *Config, conf *ClientConfig, opt ...client.Option) client.Client {
	var (
		filters []filter.Filter
	)
	for _, name := range Deduplicate(cfg.Server.Filters, conf.Filters) { // 全局filter在前，且去重
		f := filter.GetClient(name)
		if f == nil {
			panic(fmt.Sprintf("filter %s no registered, do not configure", name))
		}
		filters = append(filters, f)
	}
	var dis discovery.Discovery = nil
	if conf.Discovery != "" {
		dis = discovery.Get(conf.Discovery)
	}
	if conf.Discovery != "" && dis == nil {
		fmt.Printf("service:%s discovery not exist\n", conf.ServiceName)
	}
	opts := []client.Option{
		client.WithNamespace(conf.Namespace),
		client.WithServiceName(conf.ServiceName),
		//client.WithFilters(filters),
		client.WithDiscovery(dis),
	}

	opts = append(opts, opt...)

	cc := client.Get(conf.Protocol)
	if cc == nil {
		panic(fmt.Errorf("can not get client constructor: %s ", conf.ServiceName))
	}

	return cc(opts...)
}

func NewClientsWithConfig(cfg *Config, opts ...client.Option) *client.Clients {
	clients := client.NewClients()
	for _, conf := range cfg.Client.Clients {
		clients.AddClient(conf.Name, newClientWithConfig(cfg, conf, opts...))
	}
	return clients
}
