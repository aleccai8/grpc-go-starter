package grpc_config

import (
	"fmt"
	"github.com/zhengheng7913/grpc-config/config"
	"github.com/zhengheng7913/grpc-config/filter"
	"github.com/zhengheng7913/grpc-config/naming/registry"
	"github.com/zhengheng7913/grpc-config/server"
	_ "go.uber.org/automaxprocs"
)

// NewServer 启动时序 初始化server->config->plugin->service
func NewServer(opts ...server.Option) *server.Server {

	path := config.ServerConfigPath()

	cfg, err := config.LoadConfig(path)

	if err != nil {
		panic("parse config failed: " + err.Error())
	}

	config.SetGlobalConfig(cfg)

	if err := config.Setup(cfg); err != nil {
		panic("setup plugin fail: " + err.Error())
	}

	return NewServerWithConfig(cfg, opts...)
}

//
func newServiceWithConfig(cfg *config.Config, serviceCfg *config.ServiceConfig, opts ...server.Option) server.Service {

	var (
		filters []filter.Filter
	)
	// 填充全局Port默认值
	if cfg.Server.Port > 0 && serviceCfg.Port == 0 {
		serviceCfg.Port = cfg.Server.Port
	}

	for _, name := range Deduplicate(cfg.Server.Filters, serviceCfg.Filters) { // 全局filter在前，且去重
		f := filter.GetServer(name)
		if f == nil {
			panic(fmt.Sprintf("filter %s no registered, do not configure", name))
		}
		filters = append(filters, f)
	}

	reg := registry.Get(serviceCfg.Name)
	if serviceCfg.Registry != "" && reg == nil {
		fmt.Printf("service:%s registry not exist\n", serviceCfg.Name)
	}

	opts = append(opts,
		server.WithRegistry(reg),
		server.WithFilters(filters),
	)

	sc := server.Get(serviceCfg.Protocol)

	if sc == nil {
		panic("can not get service constructor:" + serviceCfg.Name)
	}

	return sc(serviceCfg, opts...)
}

func NewServerWithConfig(cfg *config.Config, opts ...server.Option) *server.Server {
	s := server.NewServer()

	for _, srvCfg := range cfg.Server.Services {
		s.AddService(srvCfg.Name, newServiceWithConfig(cfg, srvCfg, opts...))
	}
	return s
}
