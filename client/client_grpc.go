package client

import (
	"context"
	"fmt"
	"github.com/zhengheng7913/grpc-go-starter/naming/discovery"
	"google.golang.org/grpc"
	"reflect"
)

const GrpcProtocol = "grpc"

var (
	ErrNotGrpcClient = fmt.Errorf("not a valid grpc client")
)

func NewGrpcClient(opts ...Option) Client {
	return &GrpcClient{
		options:    applyOption(opts...),
		realClient: nil,
	}
}

type GrpcClient struct {
	options    *Options
	realClient any
}

func (g *GrpcClient) RealClient() any {
	return g.realClient
}

func (g *GrpcClient) Register(realClient interface{}, opts ...Option) {
	g.applyOption(opts...)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rr := reflect.ValueOf(realClient)
	if g.isGrpcCons(rr.Type()) {
		panic(ErrNotGrpcClient)
	}
	target, err := g.options.Discovery.Target(g.options.ServiceName, discovery.WithNamespace(g.options.Namespace))
	if err != nil {
		panic(fmt.Errorf("get target error: %v", err))
	}
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	g.realClient = rr.Call([]reflect.Value{reflect.ValueOf(conn)})[0].Interface()
}
func (g *GrpcClient) isGrpcCons(t reflect.Type) bool {
	return !(t.NumIn() == 1 && t.NumOut() == 1)
}

func (g *GrpcClient) applyOption(opts ...Option) {
	for _, opt := range opts {
		opt(g.options)
	}
}
