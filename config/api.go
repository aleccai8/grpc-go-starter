package config

import (
	"flag"
)

// SetGlobalConfig 设置全局配置对象
func SetGlobalConfig(cfg *Config) {
	gm.Store(cfg)
}

// ServerConfigPath 获取服务启动配置文件路径
//	最高优先级：服务主动修改ServerConfigPath变量
//	第二优先级：服务通过--conf或者-conf传入配置文件路径
//	第三优先级：默认路径./grpc_go.yaml
func ServerConfigPath() string {
	if Path == defaultConfigPath {
		flag.StringVar(&Path, "conf", defaultConfigPath, "server config path")
		flag.Parse()
	}
	return Path
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
