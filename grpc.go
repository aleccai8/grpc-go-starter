package grpc_go_starter

import (
	"fmt"

	"github.com/zhengheng7913/grpc-go-starter/filter"
	"github.com/zhengheng7913/grpc-go-starter/naming/registry"
	"github.com/zhengheng7913/grpc-go-starter/server"
	_ "go.uber.org/automaxprocs"
)

// NewServer 启动时序 初始化server->config->plugin->service
func NewServer(opts ...server.Option) *server.Server {

	path := ServerConfigPath()

	cfg, err := LoadConfig(path)

	if err != nil {
		panic("parse config failed: " + err.Error())
	}

	SetGlobalConfig(cfg)

	if err := Setup(cfg); err != nil {
		panic("setup plugin fail: " + err.Error())
	}

	return NewServerWithConfig(cfg, opts...)
}

//
func newServiceWithConfig(cfg *Config, conf *ServiceConfig, opt ...server.Option) server.Service {

	var (
		filters []filter.Filter
	)

	// 填充全局Port默认值
	if cfg.Server.Port > 0 && conf.Port == 0 {
		conf.Port = cfg.Server.Port
	}

	for _, name := range Deduplicate(cfg.Server.Filters, conf.Filters) { // 全局filter在前，且去重
		f := filter.GetServer(name)
		if f == nil {
			panic(fmt.Sprintf("filter %s no registered, do not configure", name))
		}
		filters = append(filters, f)
	}

	reg := registry.Get(conf.Name)
	if conf.Registry != "" && reg == nil {
		fmt.Printf("service:%s registry not exist\n", conf.Name)
	}
	opts := []server.Option{
		server.WithServiceName(conf.Name),
		server.WithPort(conf.Port),
		server.WithTarget(conf.Target),
		server.WithFilters(filters),
		server.WithRegistry(reg),
		server.WithAddress(conf.Address),
	}

	sc := server.Get(conf.Protocol)

	if sc == nil {
		panic("can not get service constructor:" + conf.Name)
	}

	opts = append(opts, opt...)

	return sc(opts...)
}

func NewServerWithConfig(cfg *Config, opts ...server.Option) *server.Server {
	s := server.NewServer()

	for _, conf := range cfg.Server.Services {
		s.AddService(conf.Name, newServiceWithConfig(cfg, conf, opts...))
	}
	return s
}
