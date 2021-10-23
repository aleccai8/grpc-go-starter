package grpc_config

import (
	"github.com/zhengheng7913/grpc-config/config"
	"github.com/zhengheng7913/grpc-config/filter"
	"github.com/zhengheng7913/grpc-config/protocol"
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

	// 这里需要获取自定义grpc的全局组件 filter registry
	gFilter := filter.Get(filter.GlobalFilterName)

	opts = append(opts, gFilter...)

	constructor := protocol.Get(serviceCfg.Protocol)

	if constructor == nil {
		panic("can not get service constructor:" + serviceCfg.Name)
	}

	return constructor(opts...)
}

func NewServerWithConfig(cfg *config.Config, opts ...server.Options) *server.Server {
	s := &server.Server{}

	for _, srvCfg := range cfg.Server.Services {
		s.AddService(srvCfg.Name, newServiceWithConfig(cfg, srvCfg, opts...))
	}
	return s
}
