package grpc_go_starter

import (
	"fmt"

	"github.com/zhengheng7913/grpc-go-starter/config"
	"github.com/zhengheng7913/grpc-go-starter/filter"
	"github.com/zhengheng7913/grpc-go-starter/naming/registry"
	"github.com/zhengheng7913/grpc-go-starter/server"
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
func newServiceWithConfig(starter *server.ServiceStarter, opts ...server.Option) server.Service {

	var (
		filters []filter.Filter
	)

	global := starter.Global
	current := starter.Current

	// 填充全局Port默认值
	if global.Server.Port > 0 && current.Port == 0 {
		current.Port = global.Server.Port
	}

	for _, name := range Deduplicate(global.Server.Filters, current.Filters) { // 全局filter在前，且去重
		f := filter.GetServer(name)
		if f == nil {
			panic(fmt.Sprintf("filter %s no registered, do not configure", name))
		}
		filters = append(filters, f)
	}

	reg := registry.Get(current.Name)
	if current.Registry != "" && reg == nil {
		fmt.Printf("service:%s registry not exist\n", current.Name)
	}

	opts = append(opts,
		server.WithRegistry(reg),
		server.WithFilters(filters),
	)

	sc := server.Get(current.Protocol)

	if sc == nil {
		panic("can not get service constructor:" + current.Name)
	}

	return sc(starter, opts...)
}

func NewServerWithConfig(cfg *config.Config, opts ...server.Option) *server.Server {
	s := server.NewServer()

	for _, node := range cfg.Server.Services {
		decoder := &config.YamlNodeDecoder{Node: &node}
		starter := &server.ServiceStarter{
			Global:         cfg,
			Current:        nil,
			CurrentDecoder: decoder,
		}
		if err := decoder.Decode(starter.Current); err != nil {
			panic("decode service config failed: " + err.Error())
		}
		s.AddService(starter.Current.Name, newServiceWithConfig(starter, opts...))
	}
	return s
}
