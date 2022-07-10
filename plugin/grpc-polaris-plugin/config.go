package grpc_polaris_plugin

type FactoryConfig struct {
	HeartbeatInterval  int       `yaml:"heartbeat_interval"`
	AddressList        string    `yaml:"address_list"`
	Services           []Service `yaml:"services"`
	Clients            []Client  `yaml:"clients"`
	DisableHealthCheck bool      `yaml:"disable_health_check"`
}

type Client struct {
	Name     string `yaml:"name"`
	MetaData map[string]string
}

type Service struct {
	Name           string            `yaml:"name"`
	Namespace      string            `yaml:"namespace"`
	ServiceName    string            `yaml:"service_name"`
	Token          string            `yaml:"token"`
	Weight         int               `yaml:"weight"`
	MetaData       map[string]string `yaml:"metadata"`
	Protocol       string            `yaml:"protocol"`
	EnableRegister bool              `yaml:"register_self"`
	Version        string            `yaml:"version"`
}

type DiscoveryConfig struct {
	Name     string
	MetaData map[string]string
}

type RegistryConfig struct {
	// ServiceToken 服务访问Token
	ServiceToken string
	// HeartBeat 上报心跳时间间隔，默认为建议 为TTL/2
	HeartBeat int
	// EnableRegister 默认只上报心跳，不注册服务，为 true 则启动注册
	EnableRegister bool
	// Weight
	Weight int
	// TTL 单位s，服务端检查周期实例是否健康的周期
	TTL int

	// DisableHealthCheck 禁用健康检查
	DisableHealthCheck bool

	Version string

	Services map[string]*RegistryService
}

type RegistryService struct {
	// InstanceID 实例名
	InstanceID string
	// Namespace 命名空间
	Namespace string
	// ServiceName 服务名
	ServiceName string
	// Protocol 服务端访问方式，支持 http grpc，默认 grpc
	Protocol string

	Port int

	Host string
	// Metadata 用户自定义 metadata 信息
	Metadata map[string]string
}
