package client

import (
	"context"
	"fmt"
	"github.com/zhengheng7913/grpc-go-starter/naming/discovery"
	"github.com/zhengheng7913/grpc-go-starter/server"
	"google.golang.org/grpc"
	"reflect"
)

const GrpcProtocol = "grpc"

var (
	ErrNotGrpcClient = fmt.Errorf("not a valid grpc client")
)

func NewGrpcClient(opts ...Option) Client {
	return &GrpcClient{
		options:           applyOption(opts...),
		realClientFactory: nil,
	}
}

type GrpcClient struct {
	options           *Options
	realClientFactory any
}

func (g *GrpcClient) RealClient(ctx context.Context) any {
	srcService, ok := ctx.Value(server.KeyService).(string)
	if !ok {
		srcService = g.options.SrcServiceName
	}
	rr := reflect.ValueOf(g.realClientFactory)
	if g.isGrpcFactory(rr.Type()) {
		panic(ErrNotGrpcClient)
	}
	target, err := g.options.Discovery.Target(
		g.options.ServiceName,
		discovery.WithNamespace(g.options.Namespace),
		discovery.WithSrcService(srcService),
	)
	if err != nil {
		panic(fmt.Errorf("get target error: %v", err))
	}
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return rr.Call([]reflect.Value{reflect.ValueOf(conn)})[0].Interface()
}

func (g *GrpcClient) Register(realClientFactory interface{}, opts ...Option) {
	g.apply(opts...)
	g.realClientFactory = realClientFactory

}
func (g *GrpcClient) isGrpcFactory(t reflect.Type) bool {
	return !(t.NumIn() == 1 && t.NumOut() == 1)
}

func (g *GrpcClient) apply(opts ...Option) {
	for _, opt := range opts {
		opt(g.options)
	}
}
