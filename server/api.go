package server

import (
	"github.com/zhengheng7913/grpc-config/config"
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

func NewHttpService(cfg *config.ServiceConfig, opts ...Options) Service {
	return &ServiceHTTP{
		cfg: cfg,
	}
}
func NewHttpServiceDesc(registrar RegistrarHTTP) *ServiceDescHTTP {
	return &ServiceDescHTTP{registrar: registrar}
}

func NewGrpcService(cfg *config.ServiceConfig, opts ...Options) Service {
	// 初始化GrpcServer
	gOption := &OptionGRPC{}
	for _, f := range opts {
		f(gOption)
	}
	return &ServiceGRPC{
		opt: gOption,
		cfg: cfg,
	}
}

func WithServiceRegisterAdapter(srv Service) grpc.ServiceRegistrar {
	return newServiceRegisterAdapter(srv)
}

func WithGrpcOptions(serviceName string, option ...grpc.ServerOption) Options {
	return func(opt Option) {
		opt.Apply(dessertGrpcOptions(option...))
	}
}
