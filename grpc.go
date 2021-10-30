package grpc_config

import (
	"github.com/zhengheng7913/grpc-config/config"
	"github.com/zhengheng7913/grpc-config/filter"
	"github.com/zhengheng7913/grpc-config/server"
	_ "go.uber.org/automaxprocs"
)

// NewServer 启动时序 初始化server->config->plugin->service
func NewServer(opts ...server.Options) *server.Server {

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
func newServiceWithConfig(cfg *config.Config, serviceCfg *config.ServiceConfig, opts ...server.Options) server.Service {

	// 填充全局Port默认值
	if cfg.Server.Port > 0 && serviceCfg.Port == 0 {
		serviceCfg.Port = cfg.Server.Port
	}

	filterNum := len(cfg.Server.Filter)
	if filterNum > 0 {
		filters := make([]string, filterNum+len(serviceCfg.Filters))
		filters = append(filters, cfg.Server.Filter...)
		filters = append(filters, serviceCfg.Filters...)
		serviceCfg.Filters = filters
	}

	//注入所有filter
	filters := make([]server.Options, len(serviceCfg.Filters))

	for _, filterName := range serviceCfg.Filters {
		filters = append(filters, filter.Get(filterName)...)
	}

	sc := server.Get(serviceCfg.Protocol)

	if sc == nil {
		panic("can not get service constructor:" + serviceCfg.Name)
	}

	return sc(serviceCfg, opts...)
}

func NewServerWithConfig(cfg *config.Config, opts ...server.Options) *server.Server {
	s := server.NewServer()

	for _, srvCfg := range cfg.Server.Services {
		s.AddService(srvCfg.Name, newServiceWithConfig(cfg, srvCfg, opts...))
	}
	return s
}
