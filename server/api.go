package server

import (
	"github.com/zhengheng7913/grpc-go-starter/config"
	"github.com/zhengheng7913/grpc-go-starter/filter"
	"github.com/zhengheng7913/grpc-go-starter/naming/registry"
	"google.golang.org/grpc"
)

var (
	implementMap = make(map[string]ServiceConstructor)
)

const (
	ProtocolNameGrpc = "grpc"
	ProtocolNameHTTP = "http"
)

// Register 非线程安全
func Register(name string, constructor ServiceConstructor) {
	implementMap[name] = constructor
}

func Get(name string) ServiceConstructor {
	return implementMap[name]
}

func init() {
	Register(ProtocolNameGrpc, NewGrpcService)
	Register(ProtocolNameHTTP, NewHttpService)
}

func NewHttpService(cfg *config.ServiceConfig, opts ...Option) Service {
	gOption := &Options{
		ServiceOptions: &HttpOptions{},
	}
	for _, f := range opts {
		f(gOption)
	}
	return &HttpService{
		opt: gOption,
		cfg: cfg,
	}
}

func NewHttpServiceDesc(registrar HttpRegistrar) *ServiceDescHTTP {
	return &ServiceDescHTTP{registrar: registrar}
}

func NewGrpcService(cfg *config.ServiceConfig, opts ...Option) Service {
	// 初始化GrpcServer
	gOption := &Options{
		ServiceOptions: &GrpcOptions{},
	}
	for _, f := range opts {
		f(gOption)
	}
	return &GrpcService{
		opt: gOption,
		cfg: cfg,
	}
}

func WithServiceRegisterAdapter(srv Service) grpc.ServiceRegistrar {
	return newServiceRegisterAdapter(srv)
}

func WithGrpcOptions(serviceName string, option ...grpc.ServerOption) Option {
	return func(opt *Options) {
		opt.ServiceOptions.Apply(dessertGrpcOptions(option...))
	}
}

// WithRegistry 指定server服务注册中心, 一个服务只能支持一个registry
func WithRegistry(r registry.Registry) Option {
	return func(opt *Options) {
		opt.Registry = r
	}
}

func WithFilters(fs filter.Chain) Option {
	return func(o *Options) {
		o.ServiceOptions.Apply(fs...)
	}
}
