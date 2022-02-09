package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"reflect"
)

const GrpcProtocol = "grpc"

var (
	ErrNotGrpcMethod = fmt.Errorf("not a valid grpc method")
	ErrNotGrpcClient = fmt.Errorf("not a valid grpc client")
)

func NewGrpcClient[T](opts ...Option) Client[T] {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	return &GrpcClient[]{
		options:    options,
		realClient: nil,
	}
}

type GrpcClient[T interface{}] struct {
	options    *Options
	realClient T
}

func (g *GrpcClient[T]) RealClient() T {
	return nil
}

func (g *GrpcClient[T]) Register(realClient interface{}, opts ...Option) {
	for _, opt := range opts {
		opt(g.options)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cons, ok := realClient.(func(cc grpc.ClientConnInterface) interface{})
	if !ok {
		panic(ErrNotGrpcClient)
	}
	target, err := g.options.Discovery.Target(fmt.Sprintf("%v://%v", g.options.Discovery, g.options.ServiceName))
	if err != nil {
		panic(fmt.Errorf("get target error: %v", err))
	}
	conn, err := grpc.DialContext(ctx, target)
	if err != nil {
		panic(err)
	}
	g.realClient = cons(conn)
}

func (g *GrpcClient[T]) isGrpcMethod(t reflect.Type) bool {
	return t.NumIn() != 3 || t.NumOut() != 2
}

func (g *GrpcClient[T]) Invoke(context context.Context, method string, req T1, opts ...Option) (T2, error) {

	handle := reflect.ValueOf(g.realClient).MethodByName(method)
	if !g.isGrpcMethod(handle.Type()) {
		return nil, ErrNotGrpcMethod
	}

	values := handle.Call([]reflect.Value{reflect.ValueOf(context), reflect.ValueOf(req)})
	reply := values[0]
	err := values[1]
	return reply.Interface(), err.Interface().(error)
}
