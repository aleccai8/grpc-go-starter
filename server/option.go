package server

import (
	"github.com/zhengheng7913/grpc-go-starter/filter"
	"github.com/zhengheng7913/grpc-go-starter/naming/registry"
	"google.golang.org/grpc"
)

type Options struct {
	Namespace string // 当前服务命名空间 正式环境 Production 测试环境 Development

	ServiceName string

	IP string

	Port uint16

	Address string

	// GrpcGateway 代理的目标地址
	Target string

	Registry registry.Registry

	Filters []filter.Filter

	Others []interface{}
}

type Option func(*Options)

func WithServiceRegisterAdapter(srv Service) grpc.ServiceRegistrar {
	return newServiceRegisterAdapter(srv)
}

func WithNamespace(namespace string) Option {
	return func(opt *Options) {
		opt.Namespace = namespace
	}
}

func WithServiceName(name string) Option {
	return func(opt *Options) {
		opt.ServiceName = name
	}
}

func WithOther(option interface{}) Option {
	return func(opt *Options) {
		opt.Others = append(opt.Others, option)
	}
}

func WithTarget(target string) Option {
	return func(opt *Options) {
		opt.Target = target
	}
}

func WithPort(port uint16) Option {
	return func(opt *Options) {
		opt.Port = port
	}
}

// WithRegistry 指定server服务注册中心, 一个服务只能支持一个registry
func WithRegistry(r registry.Registry) Option {
	return func(opt *Options) {
		opt.Registry = r
	}
}

func WithFilters(fs filter.Chain) Option {
	return func(opt *Options) {
		opt.Filters = append(opt.Filters, fs...)
	}
}

func WithIP(ip string) Option {
	return func(options *Options) {
		options.IP = ip
	}
}

func WithAddress(address string) Option {
	return func(options *Options) {
		options.Address = address
	}
}
