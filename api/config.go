package api

import (
	"flag"
	"github.com/zhengheng7913/grpc-go-starter/pkg/config"
	"github.com/zhengheng7913/grpc-go-starter/pkg/plugin"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"sync/atomic"
)

const defaultConfigPath = "./grpc_go.yaml"

// 架启动后解析yaml文件并赋值
var gm = atomic.Value{}

// SetGlobalConfig 设置全局配置对象
func SetGlobalConfig(cfg *Config) {
	gm.Store(cfg)
}

// ConfigPath 获取服务启动配置文件路径
//	最高优先级：服务主动修改ServerConfigPath变量
//	第二优先级：服务通过--conf或者-conf传入配置文件路径
//	第三优先级：默认路径./grpc_go.yaml
func ConfigPath() string {
	if Path == defaultConfigPath {
		flag.StringVar(&Path, "conf", defaultConfigPath, "server config path")
		flag.Parse()
	}
	return Path
}

func LoadSetup() (*Config, error) {
	cfg, ok := gm.Load().(*Config)
	if ok {
		return cfg, nil
	}
	path := ConfigPath()
	cfg, err := LoadConfig(path)
	if err != nil {
		return nil, err
	}
	err = Setup(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadConfig(path string) (*Config, error) {
	cfg, err := parseConfigFromFile(path)
	if err != nil {
		return nil, err
	}
	if err := repairServerConfig(cfg); err != nil {
		return nil, err
	}
	if err := repairClientConfig(cfg); err != nil {
		return nil, err
	}
	SetGlobalConfig(cfg)
	return cfg, nil
}

// Setup 加载全局配置，加载插件
func Setup(cfg *Config) error {

	// 装载插件
	if cfg.Plugins != nil {
		if err := cfg.Plugins.Setup(); err != nil {
			return err
		}
	}

	return nil
}

var Path = defaultConfigPath

type Config struct {
	Includes []string `yaml:"includes"`
	Global   struct {
		Namespace     string `yaml:"namespace"`
		EnvName       string `yaml:"env_name"`
		ContainerName string `yaml:"container_name"`
		Host          string `yaml:"host"`
	}
	Server struct {
		Filters  []string         `yaml:"filters"`
		Services []*ServiceConfig `yaml:"services"`
	}
	Client struct {
		Filters []string        `yaml:"filters"`
		Clients []*ClientConfig `yaml:"clients"`
	}
	Plugins plugin.Config
}

type ServiceConfig struct {
	Name        string   `yaml:"name"`
	ServiceName string   `yaml:"service_name"`
	Protocol    string   `yaml:"protocol"`
	Port        uint16   `yaml:"port"`
	Target      string   `yaml:"target"`
	Registry    string   `yaml:"registry"`
	Filters     []string `yaml:"filters"`
}

type ClientConfig struct {
	Name           string   `yaml:"name"`
	Namespace      string   `yaml:"namespace"`
	ServiceName    string   `yaml:"service_name"`
	SrcServiceName string   `yaml:"src_service_name"`
	Protocol       string   `yaml:"protocol"`
	Port           uint16   `yaml:"port"`
	Address        string   `yaml:"address"`
	Discovery      string   `yaml:"discovery"`
	Filters        []string `yaml:"filters"`
}

// RepairServerConfig 修复配置数据，填充默认值
func repairServerConfig(cfg *Config) error {
	return nil
}

func repairClientConfig(cfg *Config) error {

	return nil
}

func parseConfigFromFile(configPath string) (*Config, error) {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	// 支持环境变量
	buf = []byte(config.ExpandEnv(string(buf)))

	cfg := defaultConfig()
	if err := yaml.Unmarshal(buf, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func defaultConfig() *Config {
	cfg := &Config{}

	return cfg
}
