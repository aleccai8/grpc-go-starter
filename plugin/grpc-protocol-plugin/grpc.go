package grpc_protocol_plugin

import (
	"context"
	"fmt"
	"github.com/zhengheng7913/grpc-go-starter/pkg/client"
	"github.com/zhengheng7913/grpc-go-starter/pkg/naming/discovery"
	"github.com/zhengheng7913/grpc-go-starter/pkg/naming/registry"
	"github.com/zhengheng7913/grpc-go-starter/pkg/server"
	"github.com/zhengheng7913/grpc-go-starter/plugin/grpc-protocol-plugin/utils"
	"google.golang.org/grpc"
	"log"
	"net"
	"reflect"
)

func init() {
	client.Register("grpc", NewGrpcClient[interface{}])
	server.Register("grpc", NewService)
	server.Register("http", NewGatewayService)
}

func NewService(opts ...server.Option) server.Service {
	// 初始化GrpcServer
	gOption := &server.Options{}
	for _, f := range opts {
		f(gOption)
	}
	return &Service{
		opts: gOption,
	}
}

type Service struct {
	server *grpc.Server
	opts   *server.Options
}

func (g *Service) Register(factory interface{}, impl interface{}) {
	var opts []grpc.ServerOption
	opts = append(opts, utils.ArrayConvert[grpc.ServerOption](g.opts.Filters)...)
	opts = append(opts, utils.ArrayConvert[grpc.ServerOption](g.opts.Others)...)
	opts = append(opts, grpc.ChainUnaryInterceptor(utils.GetContextValueInterceptor(g.opts)))
	g.server = grpc.NewServer(opts...)
	reflect.ValueOf(factory).Call([]reflect.Value{reflect.ValueOf(g), reflect.ValueOf(impl)})
}

func (g *Service) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	g.server.RegisterService(desc, impl)
}

func (g *Service) Serve() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", g.opts.Host, g.opts.Port))
	if err != nil {
		return fmt.Errorf("Failed to listen: %v ", err)
	}
	go func() {
		err = g.opts.Registry.Register(
			g.opts.Name,
			registry.WithNamespace(g.opts.Namespace),
			registry.WithHost(g.opts.Host),
			registry.WithProtocol("grpc"),
			registry.WithServiceName(g.opts.ServiceName),
			registry.WithPort(g.opts.Port),
		)
		if err != nil {
			log.Println(err)
		}
		err := g.server.Serve(lis)
		if err != nil {
			log.Println(err)
		}
	}()

	return nil
}

func (g *Service) Close() error {
	g.opts.Registry.Deregister(g.opts.Name)
	g.server.Stop()
	return nil
}

var (
	ErrNotGrpcClient = fmt.Errorf("not a valid grpc client")
)

func NewGrpcClient[T interface{}](opts ...client.Option) client.Client[T] {
	return &GrpcClient[T]{
		options:           utils.ApplyOption(opts...),
		realClientFactory: nil,
	}
}

type GrpcClient[T interface{}] struct {
	options           *client.Options
	realClientFactory any
}

func (g *GrpcClient[T]) RealClient(ctx context.Context) T {
	var err error
	var target string
	srcService, ok := ctx.Value("service").(string)
	if !ok {
		srcService = g.options.SrcServiceName
	}
	rr := reflect.ValueOf(g.realClientFactory)
	if isGrpcFactory(rr.Type()) {
		panic(ErrNotGrpcClient)
	}
	if g.options.Target == "" {
		if target, err = g.options.Discovery.Target(
			g.options.ServiceName,
			discovery.WithNamespace(g.options.Namespace),
			discovery.WithSrcService(srcService),
			discovery.WithContext(ctx),
			discovery.WithMetadata(g.options.Metadata),
			discovery.WithProtocol(g.options.Protocol),
		); err != nil {
			log.Printf("get target error: %v\n", err)
		}
	} else {
		target = g.options.Target
	}
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure())
	if err != nil {
		log.Printf("grpc dial error: %v\n", err.Error())
	}
	return rr.Call([]reflect.Value{reflect.ValueOf(conn)})[0].Interface().(T)
}

func (g *GrpcClient[T]) Register(realClientFactory interface{}, opts ...client.Option) {
	g.apply(opts...)
	g.realClientFactory = realClientFactory

}

func (g *GrpcClient[T]) apply(opts ...client.Option) {
	for _, opt := range opts {
		opt(g.options)
	}
}

func isGrpcFactory(t reflect.Type) bool {
	return !(t.NumIn() == 1 && t.NumOut() == 1)
}
