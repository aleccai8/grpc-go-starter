package config

import (
	"github.com/zhengheng7913/grpc-config/plugin"
	"sync/atomic"
)

const defaultConfigPath = "./grpc_go.yaml"

var Path = defaultConfigPath

// 架启动后解析yaml文件并赋值
var gm = atomic.Value{}

func init() {
	gm.Store(defaultConfig())
}

type Config struct {
	Global struct {
		Namespace     string `yaml:"namespace"`
		EnvName       string `yaml:"env_name"`
		ContainerName string `yaml:"container_name"`
		LocalIP       string `yaml:"local_ip"`
	}
	Server struct {
		App      string
		Server   string
		Protocol string   // 针对所有service的protocol 默认trpc
		Filter   []string // 针对所有service的拦截器
		Services []*ServiceConfig
	}
	Client  ClientConfig
	Plugins plugin.Config
}

type ServiceConfig struct {
	Name     string
	Protocol string
	Port     uint16
	Registry string
	Filter   []string
}

type ClientConfig struct {
}
