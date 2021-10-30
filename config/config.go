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
		App      string           `yaml:"app"'`
		Server   string           `yaml:"server"`
		Protocol string           `yaml:"protocol"` // 针对所有service的protocol 默认grpc
		Port     uint16           `yaml:"port"`
		Filter   []string         `yaml:"filter"` // 针对所有service的拦截器
		Services []*ServiceConfig `yaml:"services"`
	}
	Client  ClientConfig
	Plugins plugin.Config
}

type ServiceConfig struct {
	Name     string            `yaml:"name"`
	Protocol string            `yaml:"protocol"`
	Port     uint16            `yaml:"port"`
	Registry string            `yaml:"registry"`
	Filters  []string          `yaml:"filters"`
	Labels   map[string]string `yaml:"labels"`
}

type ClientConfig struct {
}
