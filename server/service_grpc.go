package server

import (
	"fmt"
	"github.com/zhengheng7913/grpc-go-starter/naming/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
	"reflect"
)

func NewGrpcService(opts ...Option) Service {
	// 初始化GrpcServer
	gOption := &Options{}
	for _, f := range opts {
		f(gOption)
	}
	return &GrpcService{
		opts: gOption,
	}
}

type GrpcService struct {
	server *grpc.Server
	opts   *Options
}

func (g *GrpcService) Register(factory interface{}, impl interface{}) {
	var opts []grpc.ServerOption
	opts = append(opts, arrayConvert[grpc.ServerOption](g.opts.Filters)...)
	opts = append(opts, arrayConvert[grpc.ServerOption](g.opts.Others)...)
	opts = append(opts, grpc.ChainUnaryInterceptor(GetContextValueInterceptor(g.opts)))
	g.server = grpc.NewServer(opts...)
	reflect.ValueOf(factory).Call([]reflect.Value{reflect.ValueOf(g), reflect.ValueOf(impl)})
}

func (g *GrpcService) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	g.server.RegisterService(desc, impl)
}

func (g *GrpcService) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", g.opts.Host, g.opts.Port))
	if err != nil {
		return fmt.Errorf("Failed to listen: %v ", err)
	}
	go func() {
		defer g.opts.Registry.Deregister(g.opts.ServiceName)

		err = g.opts.Registry.Register(
			g.opts.ServiceName,
			registry.WithNamespace(g.opts.Namespace),
			registry.WithHost(g.opts.Host),
			registry.WithProtocol(ProtocolNameGrpc),
			registry.WithServiceName(g.opts.ServiceName),
			registry.WithPort(g.opts.Port),
		)
		if err != nil {
			grpclog.Errorln(err)
		}
		err := g.server.Serve(lis)
		if err != nil {
			grpclog.Fatalln(err)
		}
	}()

	return nil
}

func (g *GrpcService) Close(c chan struct{}) error {
	g.opts.Registry.Deregister(g.opts.ServiceName)
	g.server.Stop()
	return nil
}
